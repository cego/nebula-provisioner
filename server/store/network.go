package store

import (
	"fmt"
	"net"
	"strings"

	"github.com/cego/nebula-provisioner/protocol"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/protobuf/proto"
)

func (s *Store) ListNetworks() ([]*protocol.Network, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.listNetworks(txn)
}

func (s *Store) listNetworks(txn *badger.Txn) ([]*protocol.Network, error) {
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

func (s *Store) GetNetworkByName(name string) (*protocol.Network, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	n, err := s.getNetwork(txn, name)
	if err != nil {
		return nil, err
	}

	return n, nil
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

	ips, err := stringsToIPNet(req.Ips)
	if err != nil {
		return nil, fmt.Errorf("invalid ip definition: %s", err)
	}

	subnets, err := stringsToIPNet(req.Subnets)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet definition: %s", err)
	}

	ca, err := generateCA(name, req.Groups, ips, subnets, req.Duration.AsDuration())
	if err != nil {
		return nil, fmt.Errorf("failed to generate CA: %s", err)
	}

	err = s.saveCA(txn, ca)
	if err != nil {
		return nil, err
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

func stringsToIPNet(s []string) ([]*net.IPNet, error) {
	var nets []*net.IPNet
	if len(s) != 0 {
		for _, rs := range s {
			rs := strings.Trim(rs, " ")
			if rs != "" {
				_, s, err := net.ParseCIDR(rs)
				if err != nil {
					return nil, fmt.Errorf("invalid cidr definition: %s", err)
				}
				nets = append(nets, s)
			}
		}
	}
	return nets, nil
}
