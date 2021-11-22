package store

import (
	"fmt"
	"net"
	"sync"

	"github.com/slyngdk/nebula-provisioner/protocol"

	"github.com/sirupsen/logrus"
)

type IPManager struct {
	l     *logrus.Logger
	store *Store

	lock     sync.Mutex
	networks map[string]*NetworkIPManager
}

type NetworkIPManager struct {
	l     *logrus.Logger
	store *Store

	networkName string
	lock        sync.Mutex
	pools       []*IPPool
	ipsInUse    []net.IP
	ranges      []*net.IPNet
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

	netManagers := make(map[string]*NetworkIPManager)

	for _, network := range networks {
		manager, err := newNetworkIPManager(i.l, i.store, network)
		if err != nil {
			i.l.WithError(err).Errorf("failed to load network: %s", network.Name)
			continue
		}
		netManagers[network.Name] = manager
	}

	i.networks = netManagers

	return nil
}

func newNetworkIPManager(l *logrus.Logger, store *Store, network *protocol.Network) (*NetworkIPManager, error) {

	var pools = make([]*IPPool, 0)
	for _, pool := range network.IpPools {

		_, ipNet, err := net.ParseCIDR(pool)
		if err != nil {
			l.WithError(err).Errorf("Failed to load pool: %s", pool)
			continue
		}
		ipPool, err := store.NewIPPool(network.Name, ipNet)
		if err != nil {
			l.WithError(err).Errorf("Failed to load pool: %s", pool)
			continue
		}
		pools = append(pools, ipPool)
	}
	var ranges = make([]*net.IPNet, 0)
	for _, ip := range network.Ips {

		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			l.WithError(err).Errorf("Failed to load ip range: %s", ip)
			continue
		}

		ranges = append(ranges, ipNet)
	}

	l.Infof("Finding IP's used on network : %s", network.Name)
	txn := store.db.NewTransaction(false)
	defer txn.Discard()

	agentsByNet, err := store.listAgentByNetwork(txn, network.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent for network: %s %s", network.Name, err)
	}

	var usedIps []net.IP

	for _, agent := range agentsByNet {
		if len(agent.AssignedIP) > 0 {
			ip, _, err := net.ParseCIDR(agent.AssignedIP)
			if err != nil {
				l.WithError(err).Errorf("Failed to parse assigned ip: %s", agent.AssignedIP)
				continue
			}
			usedIps = append(usedIps, ip)
		}
	}

	return &NetworkIPManager{networkName: network.Name, store: store, l: l, pools: pools, ranges: ranges, ipsInUse: usedIps}, nil
}

func (i *IPManager) Next(networkName string) *net.IPNet {
	i.lock.Lock()
	defer i.lock.Unlock()

	if n, ok := i.networks[networkName]; ok {
		return n.Next()
	}

	return nil
}

func (i *IPManager) RequestForAgent(networkName string, agentFingerprint []byte, ip net.IP) (*net.IPNet, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if n, ok := i.networks[networkName]; ok {
		return n.RequestForAgent(agentFingerprint, ip)
	}

	return nil, fmt.Errorf("network not found: %s", networkName)
}

func (i *NetworkIPManager) Next() *net.IPNet {
	i.lock.Lock()
	defer i.lock.Unlock()

	var ip *net.IPNet

	for _, pool := range i.pools {
		ip = pool.Next(i.ipsInUse)
		if ip != nil {
			i.ipsInUse = append(i.ipsInUse, ip.IP)
			break
		}
	}
	return ip
}

func (i *NetworkIPManager) RequestForAgent(agentFingerprint []byte, ip net.IP) (*net.IPNet, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	txn := i.store.db.NewTransaction(false)
	defer txn.Discard()

	validOnNetwork := false
	var ipNet *net.IPNet

	for _, ipNet = range i.ranges {

		networkIP := networkIP(ipNet)
		broadcastIP := broadcastIP(ipNet)
		if ip.Equal(networkIP) || ip.Equal(broadcastIP) {
			break
		}

		if ipNet.Contains(ip) {
			validOnNetwork = true
			break
		}
	}

	if !validOnNetwork {
		return nil, fmt.Errorf("IP %s is not valid for network %s", ip.String(), i.networkName)
	}

	if containsIP(i.ipsInUse, ip) {
		if i.store.isAgentEnrolled(txn, agentFingerprint) {
			agent, err := i.store.getAgentByFingerprint(txn, agentFingerprint)
			if err == nil {
				assigned := net.ParseIP(agent.AssignedIP)
				if ip.Equal(assigned) {
					return &net.IPNet{IP: ip, Mask: ipNet.Mask}, nil
				}
			}
			return nil, fmt.Errorf("IP %s is already in use on network %s", ip.String(), i.networkName)
		}
	}

	i.ipsInUse = append(i.ipsInUse, ip)
	return &net.IPNet{IP: ip, Mask: ipNet.Mask}, nil
}
