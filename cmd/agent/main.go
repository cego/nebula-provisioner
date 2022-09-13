package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	conf "github.com/slackhq/nebula/config"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Build A version string that can be set with
//
//	-ldflags "-X main.Build=SOMEVERSION"
//
// at compile-time.
var Build string

const AgentNebulaCsrPath = "agent-nebula.csr"
const AgentNebulaKeyPath = "agent-nebula.key"
const AgentNebulaCrtPath = "agent-nebula.crt"
const AgentNebulaCaPath = "agent-nebula-ca.crt"
const NebulaCrtPath = "nebula.crt"
const NebulaKeyPath = "nebula.key"
const NebulaCaPath = "ca.crt"

var (
	l          *logrus.Logger
	logLevel   string
	configPath string
	configDir  string
	config     *conf.C

	rootCmd = &cobra.Command{}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Use = os.Args[0]
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to either a file or directory to load configuration from")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")

	rootCmd.AddCommand(enrollCmd, exportCmd, serviceCmd, healthcheckCmd)
}

func initConfig() {
	l = logrus.New()
	l.Out = os.Stdout

	config = conf.NewC(l)

	if configPath == "" {
		configPath = getConfigPath()
	}

	if configPath != "" {

		configPathInfo, err := os.Stat(configPath)
		if err != nil {
			l.WithError(err).Errorln("failed to load config")
			os.Exit(1)
		}

		if configPathInfo.IsDir() {
			configDir = configPath
		} else {
			configDir = filepath.Dir(configPath)
		}

		err = config.Load(configPath)
		if err != nil {
			l.WithError(err).Errorln("failed to load config")
			os.Exit(1)
		}

		if config.IsSet("log.level") && !rootCmd.Flag("log-level").Changed {
			logLevel = config.GetString("log.level", "info")
		}

		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		l.SetLevel(level)
	} else {
		l.Errorf("failed to detect config path")
		os.Exit(1)
	}
	l.Tracef("using config: %s", configPath)
	l.Tracef("using config dir: %s", configDir)
}

func main() {
	Execute()
}

type agentClient struct {
	l      *logrus.Logger
	conn   *grpc.ClientConn
	client protocol.AgentServiceClient
	config *conf.C
}

func NewClient(l *logrus.Logger, config *conf.C) (*agentClient, error) {

	var opts []grpc.DialOption

	cert := resolvePath(config.GetString("pki.cert", "agent.crt"))
	key := resolvePath(config.GetString("pki.key", "agent.key"))

	certExists, _ := fileExists(cert)
	keyExists, _ := fileExists(key)

	var keyPair tls.Certificate
	var err error

	if certExists && keyExists {
		keyPair, err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("unable to load keypair: %v", err)
		}
	} else if !certExists && !keyExists {
		l.Info("Generating new keypair")
		if err = generateAgentKeyPair(cert, key); err != nil {
			return nil, err
		}
		keyPair, err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("unable to load keypair: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unable to load keypair: missing a part")
	}

	var caCertPool *x509.CertPool
	if config.IsSet("pki.ca") {
		ca := resolvePath(config.GetString("pki.ca", NebulaCaPath))
		srvcert, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("unable to load server cert pool: %v", err)
		}
		caCertPool = x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(srvcert); !ok {
			return nil, fmt.Errorf("unable to append server cert to ca pool")
		}
	} else {
		caCertPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("unable to load system cert pool: %v", err)
		}
	}

	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{keyPair},
		RootCAs:      caCertPool,
	})

	opts = append(opts, grpc.WithTransportCredentials(ta))

	addr := config.GetString("server", "127.0.0.1:51150")
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	client := protocol.NewAgentServiceClient(conn)

	return &agentClient{l, conn, client, config}, nil
}

func (c agentClient) Close() error {
	return c.conn.Close()
}

func generateAgentKeyPair(cert, key string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("cannot generate RSA key: %s", err)
	}

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	err = ioutil.WriteFile(key, pem.EncodeToMemory(privateKeyBlock), 0600)
	if err != nil {
		return fmt.Errorf("error while writing %s: %s", key, err)
	}
	serial, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	tml := x509.Certificate{
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(5, 0, 0),
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "New Name",           // FIXME
			Organization: []string{"New Org."}, // FIXME
		},
		BasicConstraintsValid: true,
	}
	publicKeyBytes, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("certificate cannot be created: %s", err)
	}

	publicKeyBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: publicKeyBytes,
	}

	err = ioutil.WriteFile(cert, pem.EncodeToMemory(publicKeyBlock), 0600)
	if err != nil {
		return fmt.Errorf("error while writing %s: %s", cert, err)
	}

	return nil
}
