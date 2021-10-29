package store

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"
	"github.com/slackhq/nebula/cert"
	"github.com/slyngdk/nebula-provisioner/protocol"
)

func (s *Store) ListNetworks() ([]*protocol.Network, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.Prefix = prefix_network
	it := txn.NewIterator(opts)
	defer it.Close()

	var networks []*protocol.Network

	for it.Seek(prefix_network); it.ValidForPrefix(prefix_network); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			n := &protocol.Network{}
			if err := proto.Unmarshal(v, n); err != nil {
				s.l.WithError(err).Error("Failed to parse network")
			}
			networks = append(networks, n)

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return networks, nil
}

func (s *Store) getNetwork(txn *badger.Txn, name string) (*protocol.Network, error) {
	name = strings.ToLower(name)

	if !exists(txn, prefix_network, []byte(name)) {
		return nil, fmt.Errorf("network was not found: %s", name)
	}

	t, err := txn.Get(append(prefix_network, name...))
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %s", err)
	}

	n := &protocol.Network{}
	err = t.Value(func(val []byte) error {
		return proto.Unmarshal(val, n)
	})
	if err != nil {
		s.l.WithError(err).Error("Failed to parse network")
		return nil, fmt.Errorf("failed to parse network: %s", err)
	}

	return n, nil
}

func (s *Store) CreateNetwork(req *protocol.CreateNetworkRequest) (*protocol.Network, error) {
	name := strings.ToLower(req.Name)

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	if exists(txn, prefix_network, []byte(name)) {
		return nil, fmt.Errorf("network already exists")
	}

	bytes, err := proto.Marshal(&protocol.Network{
		Name:     name,
		Duration: req.Duration,
		Groups:   req.Groups,
		Ips:      req.Ips,
		Subnets:  req.Subnets,
		IpPools:  req.IpPools,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal network: %s", err)
	}

	err = txn.Set(append(prefix_network, name...), bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add network: %s", err)
	}

	// Generate first CA for network

	var ips []*net.IPNet
	if len(req.Ips) != 0 {
		for _, rs := range req.Ips {
			rs := strings.Trim(rs, " ")
			if rs != "" {
				ip, ipNet, err := net.ParseCIDR(rs)
				if err != nil {
					return nil, fmt.Errorf("invalid ip definition: %s", err)
				}

				ipNet.IP = ip
				ips = append(ips, ipNet)
			}
		}
	}
	var subnets []*net.IPNet
	if len(req.Subnets) != 0 {
		for _, rs := range req.Subnets {
			rs := strings.Trim(rs, " ")
			if rs != "" {
				_, s, err := net.ParseCIDR(rs)
				if err != nil {
					return nil, fmt.Errorf("invalid subnet definition: %s", err)
				}
				subnets = append(subnets, s)
			}
		}
	}

	pub, rawPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error while generating ed25519 keys: %s", err)
	}

	duration := time.Hour * 24 * 365

	nc := cert.NebulaCertificate{
		Details: cert.NebulaCertificateDetails{
			Name:      name,
			Groups:    req.Groups,
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

	ca := &CA{NetworkName: name, PrivateKey: key, PublicKey: crt, Sha256Sum: sum}

	caBytes, err := proto.Marshal(ca)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal network: %s", err)
	}

	err = txn.Set(append(prefix_ca, sum...), caBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add ca: %s", err)
	}

	// DONE Generate first CA for network

	_, err = s.generateEnrollmentToken(txn, name)
	if err != nil {
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to add network: %s", err)
	}

	return &protocol.Network{Name: name}, nil
}

func containsIgnoreCase(s []string, e string) bool {
	e = strings.ToLower(e)
	for _, a := range s {
		if strings.ToLower(a) == e {
			return true
		}
	}
	return false
}
