package server

import (
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"google.golang.org/grpc"
)

type server struct {
	l            *logrus.Logger
	config       *nebula.Config
	initialized  bool
	store        *store.Store
	ipManager    *store.IPManager
	unixGrpc     *grpc.Server
	agentService *grpc.Server
}

func Main(config *nebula.Config, buildVersion string, logger *logrus.Logger) (*Control, error) {
	l := logger
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	server := server{l, config, false, nil, nil, nil, nil}

	return &Control{l, server.start, server.stop, make(chan interface{})}, nil
}

func (s *server) start() error {
	unsealed := make(chan interface{})

	st, err := store.NewStore(s.l, s.config, unsealed)
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
		err = s.startAgentService(st)
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *server) stop() {
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
