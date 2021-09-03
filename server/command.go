package server

import (
	"context"
	"encoding/hex"
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

func (srv *server) startUnixSocket(s *store.Store) error {
	srv.l.Println("Starting http unix socket")
	socketPath := srv.config.GetString("command.socket", "/tmp/nebula-provisioner.socket") // TODO Change default path
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	srv.unixGrpc = grpc.NewServer(opts...)

	c := commandServer{
		l:         srv.l,
		store:     s,
		ipManager: srv.ipManager,
	}
	protocol.RegisterServerCommandServer(srv.unixGrpc, &c)
	go func() {
		err := srv.unixGrpc.Serve(lis)
		if err != nil {
			srv.l.WithError(err).Error("Failed to start http unix socket")
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

	l         *logrus.Logger
	store     *store.Store
	ipManager *store.IPManager
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

func (c *commandServer) CreateNetwork(ctx context.Context, in *protocol.CreateNetworkRequest) (*protocol.CreateNetworkResponse, error) {
	if !c.store.IsOpen() {
		return nil, status.Error(codes.FailedPrecondition, "Server is not unsealed")
	}

	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required")
	}

	if len(in.IpPools) == 0 {
		return nil, status.Error(codes.InvalidArgument, "At least one ip pool is required")
	}

	for _, pool := range in.IpPools {
		_, block, err := net.ParseCIDR(pool)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Invalid ip pool format: %s", err))
		}
		if !store.IsUsableBlock(block) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Invalid ip pool %s not in allowed ranges", pool))
		}
	}

	n, err := c.store.CreateNetwork(in)
	if err != nil {
		c.l.WithError(err).Error("failed to create network")
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	go func() {
		err := c.ipManager.Reload()
		if err != nil {
			c.l.WithError(err).Errorf("failed to reload ip manager after network was created")
		}
	}()

	return &protocol.CreateNetworkResponse{Network: n}, nil
}

func (c *commandServer) ListNetwork(ctx context.Context, in *protocol.ListNetworkRequest) (*protocol.ListNetworkResponse, error) {
	if !c.store.IsOpen() {
		return nil, status.Error(codes.FailedPrecondition, "Server is not unsealed")
	}

	nets, err := c.store.ListNetworks()
	if err != nil {
		c.l.WithError(err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	return &protocol.ListNetworkResponse{Networks: nets}, nil
}

func (c *commandServer) ListCertificateAuthorityByNetwork(ctx context.Context, in *protocol.ListCertificateAuthorityByNetworkRequest) (*protocol.ListCertificateAuthorityByNetworkResponse, error) {
	if !c.store.IsOpen() {
		return nil, status.Error(codes.FailedPrecondition, "Server is not unsealed")
	}

	cas, err := c.store.ListCAByNetwork(in.NetworkNames)
	if err != nil {
		c.l.WithError(err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	return &protocol.ListCertificateAuthorityByNetworkResponse{CertificateAuthorities: cas}, nil
}

func (c *commandServer) GetEnrollmentTokenForNetwork(ctx context.Context, in *protocol.GetEnrollmentTokenForNetworkRequest) (*protocol.GetEnrollmentTokenForNetworkResponse, error) {
	if !c.store.IsOpen() {
		return nil, status.Error(codes.FailedPrecondition, "Server is not unsealed")
	}

	t, err := c.store.GetEnrollmentTokenByNetwork(in.Network)
	if err != nil {
		c.l.WithError(err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	return &protocol.GetEnrollmentTokenForNetworkResponse{EnrollmentToken: t.Token}, nil
}

func (c *commandServer) ListEnrollmentRequests(context.Context, *emptypb.Empty) (*protocol.ListEnrollmentRequestsResponse, error) {
	if !c.store.IsOpen() {
		return nil, status.Error(codes.FailedPrecondition, "Server is not unsealed")
	}

	reqs, err := c.store.ListEnrollmentRequests()
	if err != nil {
		c.l.WithError(err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	var mReqs []*protocol.EnrollmentRequest

	for _, req := range reqs {
		mReqs = append(mReqs, &protocol.EnrollmentRequest{
			ClientFingerprint: hex.EncodeToString(req.ClientFingerprint),
			Created:           req.Created,
			NetworkName:       req.NetworkName,
			ClientIP:          req.ClientIP,
		})
	}

	return &protocol.ListEnrollmentRequestsResponse{
		EnrollmentRequests: mReqs,
	}, nil
}

func (c *commandServer) ApproveEnrollmentRequest(ctx context.Context, req *protocol.ApproveEnrollmentRequestRequest) (*emptypb.Empty, error) {
	if !c.store.IsOpen() {
		return nil, status.Error(codes.FailedPrecondition, "Server is not unsealed")
	}

	if req.ClientFingerprint == "" {
		return nil, status.Error(codes.InvalidArgument, "ClientFingerprint is required")
	}

	bytes, err := hex.DecodeString(req.ClientFingerprint)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ClientFingerprint is invalid")
	}

	err = c.store.ApproveEnrollmentRequest(c.ipManager, bytes)
	if err != nil {
		c.l.WithError(err).Error("Failed to approve enrollment request")
		return nil, status.Error(codes.Internal, "Failed to approve enrollment request")
	}

	return &emptypb.Empty{}, nil
}
