package store

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) GetAgentByFingerprint(clientFingerprint []byte) (*Agent, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.getAgentByFingerprint(txn, clientFingerprint)
}

func (s *Store) getAgentByFingerprint(txn *badger.Txn, clientFingerprint []byte) (*Agent, error) {
	if !s.isAgentEnrolled(txn, clientFingerprint) {
		return nil, fmt.Errorf("agent was not found: %s", clientFingerprint)
	}

	t, err := txn.Get(append(prefix_agent, clientFingerprint...))
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

func (s *Store) IsAgentEnrolled(clientFingerprint []byte) bool {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.isAgentEnrolled(txn, clientFingerprint)
}

func (s *Store) isAgentEnrolled(txn *badger.Txn, clientFingerprint []byte) bool {
	return exists(txn, prefix_agent, clientFingerprint)
}

func (s *Store) addAgent(txn *badger.Txn, agent *Agent) (*Agent, error) {

	if exists(txn, prefix_agent, agent.ClientFingerprint) {
		return nil, fmt.Errorf("agent already exists")
	}

	agent.Created = timestamppb.Now()

	bytes, err := proto.Marshal(agent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent: %s", err)
	}

	err = txn.Set(append(prefix_agent, agent.ClientFingerprint...), bytes)
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
