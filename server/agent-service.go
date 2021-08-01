package server

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
)

type agentService struct {
	protocol.UnimplementedAgentServiceServer

	l *logrus.Logger
}

func getClientCertFingerprint(ctx context.Context) ([32]byte, error) {
	pe, ok := peer.FromContext(ctx)
	if ok {
		switch v := pe.AuthInfo.(type) {
		case credentials.TLSInfo:
			securityValue := v.GetSecurityValue()
			switch v := securityValue.(type) {
			case *credentials.TLSChannelzSecurityValue:
				return sha256.Sum256(v.RemoteCertificate), nil
			}
		}
	}
	var d [32]byte
	return d, fmt.Errorf("failed to get fingerprint")
}

func (p *agentService) Enroll(ctx context.Context, request *protocol.EnrollRequest) (*protocol.EnrollResponse, error) {
	fingerprint, err := getClientCertFingerprint(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to initialize store")
	}
	if request.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "Token is required")
	}
	p.l.Printf("Client fingerprint: %x\n", fingerprint)

	return &protocol.EnrollResponse{}, nil
}

func (s *server) startProvisioner(store *store.Store) error {
	s.l.Println("Starting http agentService server")

	srvcert, err := tls.LoadX509KeyPair("examples/server.pem", "examples/server.key")
	if err != nil {
		s.l.WithError(err).Errorf("SERVER: unable to read server key pair: %v", err)
		return err
	}
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{srvcert},
		ClientAuth:   tls.RequireAnyClientCert,
	})

	lis, err := net.Listen("tcp", ":51150")
	if err != nil {
		s.l.WithError(err).Errorf("SERVER: unable to listen: %v", err)
		return err
	}
	s.provisionerGrpc = grpc.NewServer(grpc.Creds(ta))
	protocol.RegisterAgentServiceServer(s.provisionerGrpc, &agentService{l: s.l})
	if err := s.provisionerGrpc.Serve(lis); err != nil {
		s.l.WithError(err).Errorf("SERVER: failed to serve: %v", err)
		return err
	}
	return nil
}

func (s *server) stopProvisioner() error {
	s.l.Println("Stopping http agentService server")
	if s.provisionerGrpc != nil {
		s.provisionerGrpc.GracefulStop()
	}
	return nil
}
