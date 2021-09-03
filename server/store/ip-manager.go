package store

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

type IPManager struct {
	l     *logrus.Logger
	store *Store

	lock     sync.Mutex
	netPools map[string][]*IPPool
}

func NewIPManager(l *logrus.Logger, s *Store) (*IPManager, error) {
	return &IPManager{l: l, store: s}, nil
}

func (i *IPManager) Reload() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	networks, err := i.store.ListNetworks()
	if err != nil {
		return fmt.Errorf("IPManager failed to get networks: %s", err)
	}

	var netPools = make(map[string][]*IPPool)

	for _, network := range networks {
		var pools = make([]*IPPool, 0)
		for _, pool := range network.IpPools {

			_, ipNet, err := net.ParseCIDR(pool)
			if err != nil {
				i.l.WithError(err).Errorf("Failed to load pool: %s", pool)
				continue
			}
			ipPool, err := i.store.NewIPPool(network.Name, ipNet)
			if err != nil {
				i.l.WithError(err).Errorf("Failed to load pool: %s", pool)
				continue
			}
			pools = append(pools, ipPool)
		}
		netPools[network.Name] = pools
	}

	i.netPools = netPools

	return nil
}

func (i *IPManager) Next(networkName string) *net.IPNet {
	i.lock.Lock()
	defer i.lock.Unlock()

	var ip *net.IPNet

	pools := i.netPools[networkName]
	for _, pool := range pools {
		ip = pool.Next()
		if ip != nil {
			break
		}
	}
	return ip
}
