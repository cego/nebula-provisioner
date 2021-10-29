package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) generateEnrollmentToken(txn *badger.Txn, network string) (*EnrollmentToken, error) {
	network = strings.ToLower(network)

	if !exists(txn, prefix_network, []byte(network)) {
		return nil, fmt.Errorf("network was not found: %s", network)
	}

	b := make([]byte, 100)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed generate random for token: %s", err)
	}
	sum := sha256.Sum256(b)

	nt := &EnrollmentToken{
		Token:       hex.EncodeToString(sum[:]),
		NetworkName: network,
	}

	if exists(txn, prefix_enrollment_token, []byte(nt.Token)) {
		return nil, fmt.Errorf("enrollment token already exists")
	}

	bytes, err := proto.Marshal(nt)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enrollment token: %s", err)
	}

	err = txn.Set(append(prefix_enrollment_token, nt.Token...), bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add enrollment token: %s", err)
	}

	return nil, nil
}

func (s *Store) getEnrollmentToken(txn *badger.Txn, token string) (*EnrollmentToken, error) {
	if !exists(txn, prefix_enrollment_token, []byte(token)) {
		return nil, fmt.Errorf("enrollment token was not found: %s", token)
	}

	t, err := txn.Get(append(prefix_enrollment_token, token...))
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment token: %s", err)
	}

	nt := &EnrollmentToken{}
	err = t.Value(func(val []byte) error {
		return proto.Unmarshal(val, nt)
	})
	if err != nil {
		s.l.WithError(err).Error("Failed to parse enrollment token")
		return nil, fmt.Errorf("failed to parse enrollment token: %s", err)
	}

	return nt, nil
}

func (s *Store) GetNetworkByEnrollmentToken(token string) (*protocol.Network, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	nt, err := s.getEnrollmentToken(txn, token)
	if err != nil {
		return nil, err
	}

	n, err := s.getNetwork(txn, nt.NetworkName)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (s *Store) GetEnrollmentTokenByNetwork(network string) (*EnrollmentToken, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.getEnrollmentTokenByNetwork(txn, network)
}

func (s *Store) getEnrollmentTokenByNetwork(txn *badger.Txn, network string) (*EnrollmentToken, error) {
	network = strings.ToLower(network)

	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = true
	opt.Prefix = prefix_enrollment_token
	it := txn.NewIterator(opt)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()

		nt := &EnrollmentToken{}
		err := item.Value(func(val []byte) error {
			return proto.Unmarshal(val, nt)
		})
		// Just continue if err
		if err == nil {
			if nt.NetworkName == network {
				return nt, nil
			}
		}
	}

	return nil, fmt.Errorf("no token found for: %s", network)
}

func (s *Store) CreateEnrollmentRequest(clientFingerprint []byte, token string, csrPEM string, clientIP string, name string) (*EnrollmentRequest, error) {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	er, err := s.createEnrollmentRequest(txn, clientFingerprint, token, csrPEM, clientIP, name)
	if err != nil {
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to add enrollment token: %s", err)
	}
	s.l.Infof("Enrollment request created by: %x", clientFingerprint)
	return er, err
}

func (s *Store) createEnrollmentRequest(txn *badger.Txn, clientFingerprint []byte, token string, csrPEM string, clientIP string, name string) (*EnrollmentRequest, error) {
	if exists(txn, prefix_enrollment_req, clientFingerprint) {
		return nil, fmt.Errorf("enrollement request already exists")
	}

	t, err := s.getEnrollmentToken(txn, token)
	if err != nil {
		return nil, err
	}

	e := &EnrollmentRequest{
		ClientFingerprint: clientFingerprint,
		Created:           timestamppb.Now(),
		Token:             token,
		NetworkName:       t.NetworkName,
		CsrPEM:            csrPEM,
		ClientIP:          clientIP,
		Name:              name,
	}

	bytes, err := proto.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enrollment request: %s", err)
	}

	err = txn.Set(append(prefix_enrollment_req, clientFingerprint...), bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add enrollment request: %s", err)
	}

	return e, nil
}

func (s *Store) EnrollmentRequestExists(clientFingerprint []byte) bool {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return exists(txn, prefix_enrollment_req, clientFingerprint)
}

func (s *Store) ListEnrollmentRequests() ([]*EnrollmentRequest, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.Prefix = prefix_enrollment_req
	it := txn.NewIterator(opts)
	defer it.Close()

	var requests []*EnrollmentRequest

	for it.Seek(prefix_enrollment_req); it.ValidForPrefix(prefix_enrollment_req); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			n := &EnrollmentRequest{}
			if err := proto.Unmarshal(v, n); err != nil {
				s.l.WithError(err).Error("Failed to parse enrollment request")
			}
			requests = append(requests, n)

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return requests, nil
}

func (s *Store) ApproveEnrollmentRequest(ipManager *IPManager, clientFingerprint []byte) error {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	err := s.approveEnrollmentRequest(txn, ipManager, clientFingerprint)
	if err != nil {
		return fmt.Errorf("failed to approve enrollment token: %s", err)
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("failed to approve enrollment token: %s", err)
	}

	return nil
}

func (s *Store) approveEnrollmentRequest(txn *badger.Txn, ipManager *IPManager, clientFingerprint []byte) error {

	er, err := s.getEnrollmentRequest(txn, clientFingerprint)
	if err != nil {
		return err
	}

	agent := &Agent{
		ClientFingerprint: clientFingerprint,
		NetworkName:       er.NetworkName,
		CsrPEM:            er.CsrPEM,
		Groups:            er.Groups,
	}

	ip := ipManager.Next(er.NetworkName)
	if ip == nil {
		return fmt.Errorf("failed to get ip for agent")
	}

	agent, err = s.signCSR(txn, agent, ip)
	if err != nil {
		return fmt.Errorf("failed to sign agent csr: %s", err)
	}

	agent, err = s.addAgent(txn, agent)
	if err != nil {
		return fmt.Errorf("failed to add agent as part of approving: %s", err)
	}

	if err = s.deleteEnrollmentRequest(txn, clientFingerprint); err != nil {
		return err
	}

	return nil
}

func (s *Store) getEnrollmentRequest(txn *badger.Txn, clientFingerprint []byte) (*EnrollmentRequest, error) {
	if !exists(txn, prefix_enrollment_req, clientFingerprint) {
		return nil, fmt.Errorf("enrollement token was not found: %s", clientFingerprint)
	}

	t, err := txn.Get(append(prefix_enrollment_req, clientFingerprint...))
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment token: %s", err)
	}

	er := &EnrollmentRequest{}
	err = t.Value(func(val []byte) error {
		return proto.Unmarshal(val, er)
	})
	if err != nil {
		s.l.WithError(err).Error("Failed to parse enrollment request")
		return nil, fmt.Errorf("failed to parse enrollment request: %s", err)
	}

	return er, nil
}

func (s *Store) deleteEnrollmentRequest(txn *badger.Txn, clientFingerprint []byte) error {
	if !exists(txn, prefix_enrollment_req, clientFingerprint) {
		return fmt.Errorf("enrollement token was not found: %s", clientFingerprint)
	}

	err := txn.Delete(append(prefix_enrollment_req, clientFingerprint...))
	if err != nil {
		return fmt.Errorf("failed to remove enrollment token: %s", err)
	}
	return nil
}
