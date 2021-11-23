package store

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetAgentByFingerprint(clientFingerprint []byte) (*Agent, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.getAgentByFingerprint(txn, clientFingerprint)
}

func (s *Store) getAgentByFingerprint(txn *badger.Txn, fingerprint []byte) (*Agent, error) {
	if !s.isAgentEnrolled(txn, fingerprint) {
		return nil, fmt.Errorf("agent was not found: %x", fingerprint)
	}

	t, err := txn.Get(append(prefix_agent, fingerprint...))
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %s", err)
	}

	a := &Agent{}
	err = t.Value(func(val []byte) error {
		return proto.Unmarshal(val, a)
	})
	if err != nil {
		s.l.WithError(err).Error("Failed to parse agent")
		return nil, fmt.Errorf("failed to parse agent: %s", err)
	}
	return a, err
}

func (s *Store) IsAgentEnrolled(fingerprint []byte) bool {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.isAgentEnrolled(txn, fingerprint)
}

func (s *Store) ListAgentByNetwork(networkName string) ([]*Agent, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.listAgentByNetwork(txn, networkName)
}

func (s *Store) isAgentEnrolled(txn *badger.Txn, fingerprint []byte) bool {
	return exists(txn, prefix_agent, fingerprint)
}

func (s *Store) addAgent(txn *badger.Txn, agent *Agent) (*Agent, error) {

	if exists(txn, prefix_agent, agent.Fingerprint) {
		return nil, fmt.Errorf("agent already exists")
	}

	agent.Created = timestamppb.Now()

	bytes, err := proto.Marshal(agent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent: %s", err)
	}

	err = txn.Set(append(prefix_agent, agent.Fingerprint...), bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add agent: %s", err)
	}

	return agent, nil
}

func (s *Store) listAgentByNetwork(txn *badger.Txn, networkName string) ([]*Agent, error) {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.Prefix = prefix_agent
	it := txn.NewIterator(opts)
	defer it.Close()

	var agents []*Agent

	for it.Seek(prefix_agent); it.ValidForPrefix(prefix_agent); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			agent := &Agent{}
			if err := proto.Unmarshal(v, agent); err != nil {
				s.l.WithError(err).Error("Failed to parse agent")
				return nil
			}
			if networkName == agent.NetworkName {
				agents = append(agents, agent)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return agents, nil
}

func (s *Store) addRevokedForNetwork(txn *badger.Txn, networkName string, fingerprint string) error {
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

	if !containsIgnoreCase(crl.Fingerprints, fingerprint) {
		crl.Fingerprints = append(crl.Fingerprints, fingerprint)
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

func (s *Store) deleteAgent(txn *badger.Txn, fingerprint []byte) error {
	if !exists(txn, prefix_agent, fingerprint) {
		return fmt.Errorf("agent was not found: %x", fingerprint)
	}

	err := txn.Delete(append(prefix_agent, fingerprint...))
	if err != nil {
		return fmt.Errorf("failed to remove agent: %s", err)
	}
	return nil
}
