package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cego/nebula-provisioner/protocol"
	"github.com/cego/nebula-provisioner/server/store"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/config"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
)

type server struct {
	l              *logrus.Logger
	config         *config.C
	buildVersion   string
	initialized    bool
	store          *store.Store
	ipManager      *store.IPManager
	unixGrpc       *grpc.Server
	agentService   *grpc.Server
	tasks          *tasks
	tlsLock        sync.RWMutex
	tlsCertificate *tls.Certificate
}

func Main(config *config.C, buildVersion string, logger *logrus.Logger) (*Control, error) {
	l := logger
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	server := server{l: l, config: config, buildVersion: buildVersion}

	return &Control{l, server.start, server.stop, make(chan interface{})}, nil
}

func (s *server) start() error {
	unsealed := make(chan interface{}, 1)

	dataDir := s.config.GetString("path", "/tmp/nebula-provisioner")
	stat, err := os.Stat(dataDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", dataDir)
	}

	encryptionEnabled := !(s.buildVersion == "DEBUG" && !s.config.GetBool("db.encrypted", true))

	st, err := store.NewStore(s.l, dataDir, unsealed, encryptionEnabled)
	if err != nil {
		return err
	}
	s.store = st

	ipManager, err := store.NewIPManager(s.l, st)
	if err != nil {
		return err
	}
	s.ipManager = ipManager

	err = s.startUnixSocket(st)
	if err != nil {
		return err
	}

	s.l.Println("Use server-client to continue startup")

	if s.store.IsInitialized() {
		s.l.Println("Waiting on unsealing...")
	} else {
		s.l.Println("Waiting on initializing...")
	}

	select {
	case _ = <-unsealed:
		s.l.Infoln("Server is unsealed")

		err = ipManager.Reload()
		if err != nil {
			return err
		}

		// continue startup when unsealed
		err := s.startHttpsServer(dataDir)
		if err != nil {
			return err
		}
		s.tasks = NewTasks(s.l, s.config, s.store)
		s.tasks.Start()
	}

	return nil
}

func (s *server) stop() {
	if s.tasks != nil {
		s.tasks.Stop()
	}

	if err := s.stopAgentService(); err != nil {
		s.l.WithError(err).Error("Failed to stop agentService server")
	}

	if err := s.stopUnixSocket(); err != nil {
		s.l.WithError(err).Error("Failed to stop unix socket server")
	}

	if s.store != nil {
		if err := s.store.Close(); err != nil {
			s.l.WithError(err).Error("Failed to stop store")
		}
	}
}

func (s *server) startHttpsServer(dataDir string) error {
	var tlsConfig *tls.Config

	svc := &agentService{
		l:     s.l,
		store: s.store,
	}

	server := grpc.NewServer()
	protocol.RegisterAgentServiceServer(server, svc)

	frontend, err := NewFrontend(s.config, s.l, s.store, s.ipManager)
	if err != nil {
		return err
	}

	httpsSrv := &http.Server{
		Addr:    s.config.GetString("listen.https", ":51150"),
		Handler: grpcHandlerFunc(server, frontend.ServeHTTP()),
	}

	if s.config.GetBool("acme.enabled", false) {
		hosts := s.config.GetStringSlice("acme.hosts", []string{})
		if len(hosts) == 0 {
			return fmt.Errorf("acme is enabled but no hosts specified 'acme.hosts'")
		}

		manager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(filepath.Join(dataDir, "autocert")),
			HostPolicy: autocert.HostWhitelist(hosts...),
			Email:      s.config.GetString("acme.email", ""),
		}

		// Create server for redirecting HTTP to HTTPS
		httpSrv := &http.Server{
			Addr:    s.config.GetString("listen.http", ":51151"),
			Handler: manager.HTTPHandler(nil),
		}
		go func() {
			s.l.Fatal(httpSrv.ListenAndServe())
		}()

		tlsConfig = manager.TLSConfig()
	} else {
		s.config.RegisterReloadCallback(func(_ *config.C) {
			s.l.Info("Reloading tls cert")
			keyPair, err := s.getKeyPair()
			if err != nil {
				return
			}
			s.tlsLock.Lock()
			defer s.tlsLock.Unlock()
			s.tlsCertificate = keyPair
		})

		keyPair, err := s.getKeyPair()
		if err != nil {
			return err
		}
		s.tlsLock.Lock()
		defer s.tlsLock.Unlock()
		s.tlsCertificate = keyPair

		tlsConfig = &tls.Config{
			GetCertificate: func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
				s.tlsLock.RLock()
				defer s.tlsLock.RUnlock()
				return s.tlsCertificate, nil
			},
		}
	}

	httpsSrv.TLSConfig = tlsConfig
	httpsSrv.TLSConfig.ClientAuth = tls.RequestClientCert
	httpsSrv.TLSConfig.MinVersion = tls.VersionTLS12
	httpsSrv.TLSConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
	}

	go func() {
		if err := httpsSrv.ListenAndServeTLS("", ""); err != nil {
			s.l.WithError(err).Errorf("SERVER: failed to serve: %v", err)
		}
	}()
	return nil
}

func (s *server) getKeyPair() (*tls.Certificate, error) {
	cert := s.config.GetString("pki.cert", "server.crt")
	key := s.config.GetString("pki.key", "server.key")

	keyPair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		s.l.WithError(err).Errorf("SERVER: unable to read server key pair: %v", err)
		return nil, err
	}
	return &keyPair, nil
}

func (s *server) startUnixSocket(store *store.Store) error {
	s.l.Println("Starting http unix socket")
	socketPath := s.config.GetString("command.socket", "/tmp/nebula-provisioner.socket") // TODO Change default path
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	s.unixGrpc = grpc.NewServer(opts...)

	c := commandServer{
		l:         s.l,
		store:     store,
		ipManager: s.ipManager,
	}
	protocol.RegisterServerCommandServer(s.unixGrpc, &c)
	go func() {
		err := s.unixGrpc.Serve(lis)
		if err != nil {
			s.l.WithError(err).Error("Failed to start http unix socket")
		}
	}()
	return nil
}

func (s *server) stopUnixSocket() error {
	s.l.Println("Stopping http unix socket")
	if s.unixGrpc != nil {
		s.unixGrpc.GracefulStop()
	}
	return nil
}

func grpcHandlerFunc(g *grpc.Server, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if r.ProtoMajor == 2 && strings.Contains(ct, "application/grpc") {
			g.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
