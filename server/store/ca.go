package store

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"
	"github.com/slackhq/nebula/cert"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"time"
)

func (s *Store) signCSR(txn *badger.Txn, agent *Agent, ip *net.IPNet) (*Agent, error) {
	if agent.NetworkName == "" {
		return nil, fmt.Errorf("missing network name for agent")
	}

	cas, err := s.listCAByNetwork(txn, []string{agent.NetworkName})
	if err != nil {
		return nil, fmt.Errorf("failed to get ca for network: %s", err)
	}

	if len(cas) == 0 {
		return nil, fmt.Errorf("no ca found for network: %s", agent.NetworkName)
	}
	// TODO Check CA is valid
	// TODO Support to use the active CA
	// TODO Add parameters for agent
	ca := cas[0]

	caKey, _, err := cert.UnmarshalEd25519PrivateKey(ca.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error while parsing ca-key: %s", err)
	}

	caCert, _, err := cert.UnmarshalNebulaCertificateFromPEM(ca.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("error while parsing ca-crt: %s", err)
	}

	issuer, err := caCert.Sha256Sum()
	if err != nil {
		return nil, fmt.Errorf("error while getting ca-crt fingerprint: %s", err)
	}

	if caCert.Expired(time.Now()) {
		return nil, fmt.Errorf("ca certificate is expired")
	}

	duration := time.Hour * 24 * 30
	pub, _, err := cert.UnmarshalX25519PublicKey([]byte(agent.CsrPEM))
	if err != nil {
		return nil, fmt.Errorf("error while parsing in-pub: %s", err)
	}

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:   hex.EncodeToString(agent.ClientFingerprint), // TODO friendly name
			Ips:    []*net.IPNet{ip},
			Groups: agent.Groups,
			//Subnets:   subnets,
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(duration), // TODO load default duration from config
			PublicKey: pub,
			IsCA:      false,
			Issuer:    issuer,
		},
	}

	if err := nc.CheckRootConstrains(caCert); err != nil {
		return nil, fmt.Errorf("refusing to sign, root certificate constraints violated: %s", err)
	}

	err = nc.Sign(caKey)
	if err != nil {
		return nil, fmt.Errorf("error while signing: %s", err)
	}

	b, err := nc.MarshalToPEM()
	if err != nil {
		return nil, fmt.Errorf("error while marshalling certificate: %s", err)
	}

	agent.SignedPEM = string(b)
	agent.IssuedAt = timestamppb.New(nc.Details.NotBefore)
	agent.ExpiresAt = timestamppb.New(nc.Details.NotAfter)
	agent.AssignedIP = ip.String()

	return agent, nil
}

func (s *Store) ListCAByNetwork(networks []string) ([]*protocol.CertificateAuthority, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	cas, err := s.listCAByNetwork(txn, networks)
	if err != nil {
		return nil, err
	}

	var mCas []*protocol.CertificateAuthority

	for _, ca := range cas {
		c := &protocol.CertificateAuthority{
			NetworkName:  ca.NetworkName,
			Sha256Sum:    ca.Sha256Sum,
			PublicKeyPEM: string(ca.PublicKey),
		}
		mCas = append(mCas, c)
	}

	return mCas, nil
}

func (s *Store) listCAByNetwork(txn *badger.Txn, networks []string) ([]*CA, error) {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.Prefix = prefix_ca
	it := txn.NewIterator(opts)
	defer it.Close()

	var cas []*CA

	for it.Seek(prefix_ca); it.ValidForPrefix(prefix_ca); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			ca := &CA{}
			if err := proto.Unmarshal(v, ca); err != nil {
				s.l.WithError(err).Error("Failed to parse network")
				return nil
			}
			if len(networks) == 0 || containsIgnoreCase(networks, ca.NetworkName) {
				cas = append(cas, ca)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return cas, nil
}
