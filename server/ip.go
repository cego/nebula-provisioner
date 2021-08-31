package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"net"
	"sync"
)

type IPManager struct {
	l     *logrus.Logger
	store *store.Store

	lock *sync.Mutex
}

func NewIPManager(l *logrus.Logger, s *store.Store) (*IPManager, error) {
	return &IPManager{l: l, store: s, lock: &sync.Mutex{}}, nil
}

func (i *IPManager) reload() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	networks, err := i.store.ListNetworks()
	if err != nil {
		return fmt.Errorf("IPManager failed to get networks: %s", err)
	}

	var netPools = make(map[string][]*store.IPPool)

	for _, network := range networks {
		var pools = make([]*store.IPPool, len(network.IpPools))
		for _, pool := range network.IpPools {

			_, ipNet, err := net.ParseCIDR(pool)
			if err != nil {
				i.l.WithError(err).Errorf("Failed to load pool: %s", pool)
				continue
			}
			ipRange, err := i.store.NewIPRange(network.Name, ipNet)
			if err != nil {
				i.l.WithError(err).Errorf("Failed to load pool: %s", pool)
				continue
			}
			pools = append(pools, ipRange)
		}
		netPools[network.Name] = pools
	}

	fmt.Println(netPools)

	return nil
}
