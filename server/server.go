package server

import (
	context "context"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type server struct {
	l        *logrus.Logger
	config   *nebula.Config
	unixGrpc *grpc.Server
	initialized bool
	store *store
}

func Main(config *nebula.Config, buildVersion string, logger *logrus.Logger) (*Control, error) {
	l := logger
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	server := server{l, config, nil, false, nil}

	return &Control{l, server.start, server.stop, make(chan interface{})}, nil
}

func (s *server) start() error {

	store, err := NewStore(s.config)
	if err != nil {
		return err
	}
	s.store = store

	if err := s.startUnixSocket(); err != nil {
		return err
	}

	return nil
}

func (s *server) stop() {
	if err := s.stopUnixSocket(); err != nil {
		s.l.WithError(err).Error("Failed to stop unix socket server")
	}

	if s.store != nil {
		if err := s.store.Close(); err != nil {
			s.l.WithError(err).Error("Failed to stop store")
		}
	}
}

func (s *server) startUnixSocket() error {
	s.l.Println("Starting http unix socket")
	lis, err := net.Listen("unix", "/tmp/nebula-provisioner.socket") // TODO add address to config
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	s.unixGrpc = grpc.NewServer(opts...)
	protocol.RegisterServerCommandServer(s.unixGrpc, &commandServer{})
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

type commandServer struct {
	protocol.UnimplementedServerCommandServer
}

func (_ commandServer) Init(ctx context.Context, in *emptypb.Empty) (*protocol.InitResponse, error) {
	return &protocol.InitResponse{Message: "Hi from server"}, nil
}

func (_ commandServer) IsInit(context.Context, *emptypb.Empty) (*protocol.IsInitResponse, error){
	return &protocol.IsInitResponse{IsInitialized: false}, nil
}
