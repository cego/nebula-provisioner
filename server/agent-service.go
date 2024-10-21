package server

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"net"
	"strings"

	"github.com/cego/nebula-provisioner/protocol"
	"github.com/cego/nebula-provisioner/server/store"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/cert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
				if len(v.RemoteCertificate) != 0 {
					_, err := x509.ParseCertificate(v.RemoteCertificate)
					if err == nil {
						sum := sha256.Sum256(v.RemoteCertificate)
						return sum[:], nil
					}
				}
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
	if request.RequestedIP != "" {
		ip := net.ParseIP(request.RequestedIP)
		if ip == nil {
			return nil, status.Error(codes.InvalidArgument, "RequestedIP is not valid")
		}
	}

	_, _, err = cert.UnmarshalX25519PublicKey([]byte(request.CsrPEM))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "CsrPEM is invalid")
	}

	p, _ := peer.FromContext(ctx)
	addr := p.Addr.String()
	ip := addr[0:strings.LastIndex(addr, ":")]

	_, err = a.store.CreateEnrollmentRequest(fingerprint, request.Token, request.CsrPEM, ip, request.Name, request.RequestedIP, request.Groups)
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

	if res.IsEnrollmentRequested {
		er, err := a.store.GetEnrollmentRequest(fingerprint)
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to get enrollment request")
		}
		res.EnrollmentRequest = &protocol.EnrollRequest{
			CsrPEM:      er.CsrPEM,
			Groups:      er.Groups,
			Name:        er.Name,
			RequestedIP: er.RequestedIP,
		}
	}

	if res.IsEnrolled {
		agent, err := a.store.GetAgentByFingerprint(fingerprint)
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to get agent")
		}
		res.SignedPEM = agent.SignedPEM
		res.IssuedAt = agent.IssuedAt
		res.ExpiresAt = agent.ExpiresAt
		res.Name = agent.Name
		res.Groups = agent.Groups

		nebulaFingerprint, err := store.NebulaFingerprintFromPEM(agent.SignedPEM)
		if err == nil {
			res.SignedPEMFingerprint = nebulaFingerprint
		}

		ip, _, err := net.ParseCIDR(agent.AssignedIP)
		if err == nil {
			res.AssignedIP = ip.String()
		}

		cas, err := a.store.ListCAByNetwork([]string{agent.NetworkName})
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to get certificate authorities for network")
		}
		res.CertificateAuthorities = caToProtocol(cas)

		crl, err := a.store.ListCRLByNetwork([]string{agent.NetworkName})
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to get certificate revoke list for network")
		}
		res.CertificateRevocationList = crl
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

	return &protocol.GetCertificateAuthorityByNetworkResponse{CertificateAuthorities: caToProtocol(cas)}, nil
}

func (a *agentService) GetCRLByNetwork(ctx context.Context, request *protocol.GetCRLByNetworkRequest) (*protocol.GetCRLByNetworkResponse, error) {
	fingerprint, err := getClientCertFingerprint(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Failed to get certificate fingerprint")
	}
	if !a.store.IsAgentEnrolled(fingerprint) {
		return nil, status.Error(codes.PermissionDenied, "")
	}

	crls, err := a.store.ListCRLByNetwork(request.NetworkNames)
	if err != nil {
		a.l.WithError(err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}

	return &protocol.GetCRLByNetworkResponse{Crls: crls}, nil
}

func (s *server) stopAgentService() error {
	s.l.Println("Stopping http agentService server")
	if s.agentService != nil {
		s.agentService.GracefulStop()
	}
	return nil
}

func caToProtocol(cas []*store.CA) []*protocol.CertificateAuthority {
	var mCas []*protocol.CertificateAuthority

	for _, ca := range cas {
		c := &protocol.CertificateAuthority{
			NetworkName:  ca.NetworkName,
			Sha256Sum:    ca.Sha256Sum,
			PublicKeyPEM: string(ca.PublicKey),
		}
		mCas = append(mCas, c)
	}

	return mCas
}
