package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"os"
	"time"
)

// Build A version string that can be set with
//
//     -ldflags "-X main.Build=SOMEVERSION"
//
// at compile-time.
var Build string

func main() {
	configPath := flag.String("config", "", "Path to either a file or directory to load configuration from")
	printVersion := flag.Bool("version", false, "Print version")
	printUsage := flag.Bool("help", false, "Print command line usage")

	flag.Parse()

	l := logrus.New()
	l.Out = os.Stdout

	if *printVersion {
		fmt.Printf("Version: %s\n", Build)
		os.Exit(0)
	}

	if *printUsage {
		flag.Usage()
		os.Exit(0)
	}

	if *configPath == "" {
		l.Error("-config flag must be set")
		flag.Usage()
		os.Exit(1)
	}

	config := nebula.NewConfig(l)
	err := config.Load(*configPath)
	if err != nil {
		l.WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	client, err := NewClient()
	if err != nil {
		l.WithError(err).Error("Failed to create client")
		os.Exit(1)
	}

	err = client.join(os.Stdout)
	if err != nil {
		l.WithError(err).Error("Failed to call say hello")
		os.Exit(1)
	}
}

type serverClient struct {
	conn   *grpc.ClientConn
	client protocol.AgentServiceClient
}

func NewClient() (*serverClient, error) {

	var opts []grpc.DialOption

	clcert, err := tls.LoadX509KeyPair("examples/client.pem", "examples/client.key")
	if err != nil {
		return nil, fmt.Errorf("unable to load keypair: %v", err)
	}

	// TODO make an option in config to trust self signed or else default trust public CA
	srvcert, err := ioutil.ReadFile("examples/server.pem")
	if err != nil {
		return nil, fmt.Errorf("unable to load server cert pool: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(srvcert); !ok {
		return nil, fmt.Errorf("unable to append server cert to ca pool")
	}

	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clcert},
		RootCAs:      caCertPool,
	})

	opts = append(opts, grpc.WithTransportCredentials(ta))

	conn, err := grpc.Dial("localhost:51150", opts...) // TODO change address load from config
	if err != nil {
		return nil, err
	}

	client := protocol.NewAgentServiceClient(conn)

	return &serverClient{conn, client}, nil
}

func (c serverClient) Close() error {
	return c.conn.Close()
}

func (c serverClient) join(out io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.client.Enroll(ctx, &protocol.EnrollRequest{Token: "ksjdafkljadsfkl;j"})
	if err != nil {
		return err
	}

	fmt.Fprintln(out, res)

	return nil
}
