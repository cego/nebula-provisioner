package server

import (
	"fmt"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"sync"
)

type IPManager struct {
	store *store.Store

	lock *sync.Mutex
}

func NewIPManager(s *store.Store) (*IPManager, error) {
	return &IPManager{store: s, lock: &sync.Mutex{}}, nil
}

func (i *IPManager) reload() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	networks, err := i.store.ListNetworks()
	if err != nil {
		return fmt.Errorf("IPManager failed to get networks: %s", err)
	}

	for _, network := range networks {
		fmt.Println(network.IpPools)
	}

	return nil
}
