package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

func (s *server) startUnixSocket(store *store.Store) error {
	s.l.Println("Starting http unix socket")
	lis, err := net.Listen("unix", "/tmp/nebula-provisioner.socket") // TODO add address to config
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	s.unixGrpc = grpc.NewServer(opts...)

	c := commandServer{l: s.l, store: store}
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

type commandServer struct {
	protocol.UnimplementedServerCommandServer

	l     *logrus.Logger
	store *store.Store
}

func (c *commandServer) Init(ctx context.Context, in *protocol.InitRequest) (*protocol.InitResponse, error) {

	if c.store.IsInitialized() {
		return nil, status.Error(codes.FailedPrecondition, "Server is already initialized")
	}

	keyParts, err := c.store.Initialize(in.KeyParts, in.KeyThreshold)
	if err != nil {
		c.l.WithError(err).Println("Failed to initialize store")
		return nil, status.Error(codes.Internal, "Failed to initialize store")
	}

	return &protocol.InitResponse{KeyParts: keyParts}, nil
}

func (c *commandServer) IsInit(context.Context, *emptypb.Empty) (*protocol.IsInitResponse, error) {
	return &protocol.IsInitResponse{IsInitialized: c.store.IsInitialized()}, nil
}

func (c *commandServer) Unseal(ctx context.Context, in *protocol.UnsealRequest) (*protocol.UnsealResponse, error) {
	c.l.Infof("Recived unseal request")

	if len(in.KeyPart) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Missing KeyPart")
	}

	if !c.store.IsInitialized() {
		return nil, status.Error(codes.FailedPrecondition, "Server not initialized")
	}

	if c.store.IsOpen() {
		return nil, status.Error(codes.FailedPrecondition, "Server is already unsealed")
	}

	err := c.store.Unseal(in.KeyPart, in.RemoveExistingParts)
	if err != nil {
		c.l.WithError(err)
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s", err))
	}

	return &protocol.UnsealResponse{}, nil
}
