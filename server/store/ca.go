package store

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"
	"github.com/slackhq/nebula/cert"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var DEFAULT_CA_DURATION = time.Hour * 24 * 365

func (s *Store) RevokeAgent(fingerprint []byte) error {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	agent, err := s.getAgentByFingerprint(txn, fingerprint)
	if err != nil {
		return err
	}

	nebulaFingerprints := make([]string, 0)

	nebulaFingerprint, err := NebulaFingerprintFromPEM(agent.SignedPEM)
	if err != nil {
		return err
	}
	nebulaFingerprints = append(nebulaFingerprints, nebulaFingerprint)

	for _, f := range agent.OldSignedPEMs {
		nebulaFingerprint, err = NebulaFingerprintFromPEM(f)
		if err != nil {
			return err
		}
		nebulaFingerprints = append(nebulaFingerprints, nebulaFingerprint)
	}

	err = s.addRevokedForNetwork(txn, agent.NetworkName, nebulaFingerprints)
	if err != nil {
		return fmt.Errorf("failed to add revoked fingerprint for network: %s", err)
	}

	err = s.deleteAgent(txn, fingerprint)
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit revoke of agent")
	}

	return nil
}

func (s *Store) ListCAByNetwork(networks []string) ([]*CA, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	cas, err := s.listCAByNetwork(txn, networks)
	if err != nil {
		return nil, err
	}

	return cas, nil
}

func (s *Store) ListCRLByNetwork(networks []string) ([]*protocol.NetworkCRL, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.listCRLByNetwork(txn, networks)
}

func (s *Store) PrepareCARollover(networkName string) error {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	network, err := s.getNetwork(txn, networkName)
	if err != nil {
		return fmt.Errorf("failed to get network: %s %s", networkName, err)
	}

	cas, err := s.listCAByNetwork(txn, []string{networkName})
	if err != nil {
		return fmt.Errorf("failed to get CA` for network: %s %s", networkName, err)
	}

	_, err = s.prepareCARollover(txn, cas, network)
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit prepare of CA rollover")
	}

	return nil
}

func (s *Store) prepareCARollover(txn *badger.Txn, cas []*CA, network *protocol.Network) (*CA, error) {

	var next *CA

	for _, ca := range cas {
		if ca.Status == CA_Next {
			if next != nil {
				return nil, fmt.Errorf("multiple CA`s with next status found for network: %s", network.Name)
			}
			next = ca
		}
	}

	if next != nil {
		expired, err := s.expireCA(txn, next)
		if err != nil {
			return nil, fmt.Errorf("failed to enure next CA is not expired for network: %s %s", network.Name, err)
		}
		if expired {
			next = nil
		}

	}

	if next == nil {
		s.l.Infof("Creating next CA for %s", network.Name)
		ips, err := stringsToIPNet(network.Ips)
		if err != nil {
			return nil, fmt.Errorf("invalid ip definition: %s", err)
		}

		subnets, err := stringsToIPNet(network.Subnets)
		if err != nil {
			return nil, fmt.Errorf("invalid subnet definition: %s", err)
		}

		ca, err := generateCA(network.Name, network.Groups, ips, subnets, network.Duration.AsDuration())
		if err != nil {
			return nil, fmt.Errorf("failed to generated next CA: %s", err)
		}

		ca.Status = CA_Next

		err = s.saveCA(txn, ca)
		if err != nil {
			return nil, err
		}
	}

	return next, nil
}

func (s *Store) SwitchActiveCA(networkName string) error {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	network, err := s.getNetwork(txn, networkName)
	if err != nil {
		return fmt.Errorf("failed to get network: %s %s", networkName, err)
	}

	cas, err := s.listCAByNetwork(txn, []string{networkName})
	if err != nil {
		return fmt.Errorf("failed to get CA` for network: %s", network)
	}

	err = s.switchActiveCA(txn, cas, networkName)
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit switch of active CA for network: %s %s", networkName, err)
	}

	return nil
}

func (s *Store) switchActiveCA(txn *badger.Txn, cas []*CA, networkName string) error {
	var err error
	var next *CA
	var active *CA

	for _, ca := range cas {
		if ca.Status == CA_Next {
			if next != nil {
				return fmt.Errorf("multiple CA`s with next status found for network: %s", networkName)
			}
			next = ca
		}
		if ca.Status == CA_Active {
			if active != nil {
				return fmt.Errorf("multiple CA`s with active status found for network: %s", networkName)
			}
			active = ca
		}
	}

	if next == nil {
		return fmt.Errorf("unable to switch active CA, without having the next CA created for network: %s", networkName)
	}

	s.l.Infof("Switching CA for %s", networkName)

	if active != nil {
		active.Status = CA_Inactive
		err = s.saveCA(txn, active)
		if err != nil {
			return fmt.Errorf("failed to save active CA as inactive for network: %s %s", networkName, err)
		}
	}

	next.Status = CA_Active
	err = s.saveCA(txn, next)
	if err != nil {
		return fmt.Errorf("failed to save next CA as active for network: %s %s", networkName, err)
	}

	return nil
}

func (s *Store) RenewCAs() error {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	networks, err := s.listNetworks(txn)
	if err != nil {
		return fmt.Errorf("failed to get networks %s", err)
	}

	for _, network := range networks {
		cas, err := s.listCAByNetwork(txn, []string{network.Name})
		if err != nil {
			return fmt.Errorf("failed to get CA` for network: %s", network)
		}

		var next *CA
		var active *CA

		for _, ca := range cas {
			if ca.Status == CA_Next {
				if next != nil {
					return fmt.Errorf("multiple CA`s with next status found for network: %s", network.Name)
				}
				next = ca
			}
			if ca.Status == CA_Active {
				if active != nil {
					return fmt.Errorf("multiple CA`s with active status found for network: %s", network.Name)
				}
				active = ca
			}
		}

		if active == nil {
			s.l.Infof("No active CA for %s to renew", network.Name)
			continue
		}

		activePublicKey, _, err := cert.UnmarshalNebulaCertificateFromPEM(active.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to parse active CA %s", err)
		}

		// Ensure next is created if active expires in less than 60 days
		if time.Now().Add(60*24*time.Hour).After(activePublicKey.Details.NotAfter) && next == nil {
			next, err = s.prepareCARollover(txn, cas, network)
			if err != nil {
				return err
			}
		}

		if time.Now().Add(45*24*time.Hour).After(activePublicKey.Details.NotAfter) && next != nil {
			err = s.switchActiveCA(txn, cas, network.Name)
			if err != nil {
				return err
			}
		}

	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit RenewCAs: %v", err)
	}

	return nil
}

func (s *Store) expireCA(txn *badger.Txn, ca *CA) (bool, error) {
	publicKey, _, err := cert.UnmarshalNebulaCertificateFromPEM(ca.PublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse CA %s", err)
	}
	if publicKey.Details.NotAfter.Before(time.Now()) {
		ca.Status = CA_Expired
		err = s.saveCA(txn, ca)
		if err != nil {
			return false, fmt.Errorf("failed to save expired CA %s", err)
		}
	}
	return false, nil
}

func (s *Store) signCSR(txn *badger.Txn, agent *Agent, ip *net.IPNet) (*Agent, error) {
	if agent.NetworkName == "" {
		return nil, fmt.Errorf("missing network name for agent")
	}

	cas, err := s.listCAByNetwork(txn, []string{agent.NetworkName})
	if err != nil {
		return nil, fmt.Errorf("failed to get ca for network: %s", err)
	}

	if len(cas) == 0 {
		return nil, fmt.Errorf("no CA`s found for network: %s", agent.NetworkName)
	}
	// TODO Check CA is valid
	// TODO Add parameters for agent

	var ca *CA
	for _, c := range cas {
		if c.Status == CA_Active {
			ca = c
			break
		}
	}

	if ca == nil {
		return nil, fmt.Errorf("no active CA found for network: %s", agent.NetworkName)
	}

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

	name := hex.EncodeToString(agent.Fingerprint)

	if agent.Name != "" {
		name = agent.Name + "-" + name
	}

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:   name,
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

func (s *Store) listCRLByNetwork(txn *badger.Txn, networks []string) ([]*protocol.NetworkCRL, error) {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.Prefix = prefix_network_crl
	it := txn.NewIterator(opts)
	defer it.Close()

	var crls []*protocol.NetworkCRL

	for it.Seek(prefix_ca); it.ValidForPrefix(prefix_network_crl); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			crl := &protocol.NetworkCRL{}
			if err := proto.Unmarshal(v, crl); err != nil {
				s.l.WithError(err).Error("Failed to parse network")
				return nil
			}
			if len(networks) == 0 || containsIgnoreCase(networks, crl.NetworkName) {
				crls = append(crls, crl)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return crls, nil
}

func (s *Store) addRevokedForNetwork(txn *badger.Txn, networkName string, fingerprints []string) error {
	crl := &protocol.NetworkCRL{NetworkName: networkName}

	key := append(prefix_network_crl, networkName...)
	if exists(txn, prefix_network_crl, []byte(networkName)) {
		item, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf("failed to get NetworkCRL: %s", err)
		}

		err = item.Value(func(v []byte) error {
			if err := proto.Unmarshal(v, crl); err != nil {
				return fmt.Errorf("failed to unmarhal NetworkCRL: %s", err)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	for _, f := range fingerprints {
		if !containsIgnoreCase(crl.Fingerprints, f) {
			crl.Fingerprints = append(crl.Fingerprints, f)
		}
	}

	bytes, err := proto.Marshal(crl)
	if err != nil {
		return fmt.Errorf("failed to marshal NetworkCRL: %s", err)
	}

	err = txn.Set(key, bytes)
	if err != nil {
		return fmt.Errorf("failed to set NetworkCRL: %s", err)
	}

	return nil
}

func (s Store) saveCA(txn *badger.Txn, ca *CA) error {
	caBytes, err := proto.Marshal(ca)
	if err != nil {
		return fmt.Errorf("failed to marshal CA: %s", err)
	}

	err = txn.Set(append(prefix_ca, ca.Sha256Sum...), caBytes)
	if err != nil {
		return fmt.Errorf("failed to save ca: %s", err)
	}

	return nil
}

func NebulaFingerprintFromPEM(pem string) (string, error) {
	publicKey, _, err := cert.UnmarshalNebulaCertificateFromPEM([]byte(pem))
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %s", err.Error())
	}

	nebulaFingerprint, err := publicKey.Sha256Sum()
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %s", err.Error())
	}

	return nebulaFingerprint, nil
}

func generateCA(networkName string, groups []string, ips, subnets []*net.IPNet, duration time.Duration) (*CA, error) {
	pub, rawPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error while generating ed25519 keys: %s", err)
	}

	if duration <= (time.Hour * 24 * 25) {
		duration = DEFAULT_CA_DURATION
	}

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      networkName,
			Groups:    groups,
			Ips:       ips,
			Subnets:   subnets,
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(duration),
			PublicKey: pub,
			IsCA:      true,
		},
	}

	err = nc.Sign(rawPriv)
	if err != nil {
		return nil, fmt.Errorf("error while signing: %s", err)
	}

	sum, err := nc.Sha256Sum()
	if err != nil {
		return nil, err
	}

	key := cert.MarshalEd25519PrivateKey(rawPriv)
	crt, err := nc.MarshalToPEM()
	if err != nil {
		return nil, err
	}

	ca := &CA{NetworkName: networkName, PrivateKey: key, PublicKey: crt, Sha256Sum: sum}

	return ca, nil
}
