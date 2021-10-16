package server

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strings"
)

type agentService struct {
	protocol.UnimplementedAgentServiceServer

	l     *logrus.Logger
	store *store.Store
}

func getClientCertFingerprint(ctx context.Context) ([]byte, error) {
	pe, ok := peer.FromContext(ctx)
	if ok {
		switch v := pe.AuthInfo.(type) {
		case credentials.TLSInfo:
			securityValue := v.GetSecurityValue()
			switch v := securityValue.(type) {
			case *credentials.TLSChannelzSecurityValue:
				sum := sha256.Sum256(v.RemoteCertificate)
				return sum[:], nil
			}
		}
	}
	var d []byte
	return d, fmt.Errorf("failed to get fingerprint")
}

func (a *agentService) Enroll(ctx context.Context, request *protocol.EnrollRequest) (*protocol.EnrollResponse, error) {
	fingerprint, err := getClientCertFingerprint(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to get certificate fingerprint")
	}
	if request.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "Token is required")
	}
	if request.CsrPEM == "" {
		return nil, status.Error(codes.InvalidArgument, "CsrPEM is required")
	}

	_, _, err = cert.UnmarshalX25519PublicKey([]byte(request.CsrPEM))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "CsrPEM is invalid")
	}

	p, _ := peer.FromContext(ctx)
	addr := p.Addr.String()
	ip := addr[0:strings.LastIndex(addr, ":")]

	_, err = a.store.CreateEnrollmentRequest(fingerprint, request.Token, request.CsrPEM, ip, request.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	return &protocol.EnrollResponse{}, nil
}

func (a *agentService) GetEnrollStatus(ctx context.Context, _ *emptypb.Empty) (*protocol.GetEnrollStatusResponse, error) {
	fingerprint, err := getClientCertFingerprint(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to get certificate fingerprint")
	}

	res := &protocol.GetEnrollStatusResponse{
		IsEnrolled:            a.store.IsAgentEnrolled(fingerprint),
		IsEnrollmentRequested: a.store.EnrollmentRequestExists(fingerprint),
	}

	if res.IsEnrolled {
		agent, err := a.store.GetAgentByFingerprint(fingerprint)
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to get agent")
		}
		res.SignedPEM = agent.SignedPEM
		res.IssuedAt = agent.IssuedAt
		res.ExpiresAt = agent.ExpiresAt

		cas, err := a.store.ListCAByNetwork([]string{agent.NetworkName})
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to certificateAuthorities for network")
		}

		res.CertificateAuthorities = cas
	}

	return res, nil
}

func (a *agentService) GetCertificateAuthorityByNetwork(ctx context.Context, request *protocol.GetCertificateAuthorityByNetworkRequest) (*protocol.GetCertificateAuthorityByNetworkResponse, error) {
	fingerprint, err := getClientCertFingerprint(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to get certificate fingerprint")
	}
	if !a.store.IsAgentEnrolled(fingerprint) {
		return nil, status.Error(codes.PermissionDenied, "")
	}

	cas, err := a.store.ListCAByNetwork(request.NetworkNames)
	if err != nil {
		a.l.WithError(err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	return &protocol.GetCertificateAuthorityByNetworkResponse{CertificateAuthorities: cas}, nil
}

func (srv *server) startAgentService(s *store.Store) error {
	srv.l.Println("Starting http agentService server")

	cert := srv.config.GetString("pki.cert", "server.crt")
	key := srv.config.GetString("pki.key", "server.key")

	keyPair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		srv.l.WithError(err).Errorf("SERVER: unable to read server key pair: %v", err)
		return err
	}
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{keyPair},
		ClientAuth:   tls.RequireAnyClientCert,
	})

	lis, err := net.Listen("tcp", ":51150") // TODO address from config
	if err != nil {
		srv.l.WithError(err).Errorf("SERVER: unable to listen: %v", err)
		return err
	}
	srv.agentService = grpc.NewServer(grpc.Creds(ta))
	svc := &agentService{
		l:     srv.l,
		store: s,
	}
	protocol.RegisterAgentServiceServer(srv.agentService, svc)

	go func() {
		if err := srv.agentService.Serve(lis); err != nil {
			srv.l.WithError(err).Errorf("Agent Service: failed to serve: %v", err)
		}
	}()
	return nil
}

func (s *server) stopAgentService() error {
	s.l.Println("Stopping http agentService server")
	if s.agentService != nil {
		s.agentService.GracefulStop()
	}
	return nil
}
