package store

import (
	"fmt"
	"net"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var usableIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		usableIPBlocks = append(usableIPBlocks, block)
	}
}

type IPPool struct {
	l *logrus.Logger

	cidr *net.IPNet

	db   *badger.DB
	lock sync.Mutex
	key  []byte
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (i *IPPool) Next(ipsInUse []net.IP) *net.IPNet {
	i.lock.Lock()
	defer i.lock.Unlock()

	networkIP := networkIP(i.cidr)
	broadcastIP := broadcastIP(i.cidr)

	for ip := i.cidr.IP.Mask(i.cidr.Mask); i.cidr.Contains(ip); inc(ip) {
		if ip.Equal(networkIP) || ip.Equal(broadcastIP) || containsIP(ipsInUse, ip) {
			continue
		}
		return &net.IPNet{IP: ip, Mask: i.cidr.Mask}
	}
	return nil
}

func safeXORBytes(dst, a, b []byte, n int) {
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
}

func safeORBytes(dst, a, b []byte, n int) {
	for i := 0; i < n; i++ {
		dst[i] = a[i] | b[i]
	}
}

func safeANDBytes(dst, a, b []byte, n int) {
	for i := 0; i < n; i++ {
		dst[i] = a[i] & b[i]
	}
}

func networkIP(block *net.IPNet) net.IP {
	ip := []byte(block.IP)
	mask := []byte(block.Mask)

	network := make([]byte, len(mask))
	safeANDBytes(network, ip, mask, len(mask))

	return network
}

func broadcastIP(block *net.IPNet) net.IP {
	ip := []byte(block.IP)
	mask := []byte(block.Mask)

	fill := make([]byte, len(mask))
	for i := range fill {
		fill[i] = 0xFF
	}

	xor := make([]byte, len(mask))
	safeXORBytes(xor, fill, mask, len(mask))

	broadcast := make([]byte, len(mask))
	safeORBytes(broadcast, ip, xor, len(mask))
	return broadcast
}

func containsIP(ips []net.IP, ip net.IP) bool {
	for _, a := range ips {
		if ip.Equal(a) {
			return true
		}
	}
	return false
}

func IsUsableBlock(b *net.IPNet) bool {
	for _, block := range usableIPBlocks {
		if block.Contains(b.IP) && block.Contains(broadcastIP(b)) {
			return true
		}
	}
	return false
}

func (s *Store) NewIPPool(networkName string, cidr *net.IPNet) (*IPPool, error) {
	s.l.Infof("Creating IPPool for: %s", cidr.IP.String())

	if !IsUsableBlock(cidr) {
		return nil, fmt.Errorf("cidr %s is not in an uable private range", cidr)
	}

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	key := append([]byte(networkName), []byte("-")...)
	key = append(key, []byte(cidr.IP.String())...)

	i := &IPPool{l: s.l, cidr: cidr, db: s.db, key: key}

	var r *IPRange
	if !exists(txn, prefix_ip_range, key) {
		s.l.Infof("Adding IP Range : %s on network : %s\n", cidr.IP, networkName)
		r = &IPRange{Network: cidr.IP, Netmask: cidr.Mask}
		err := i.save(txn, r)
		if err != nil {
			return nil, fmt.Errorf("failed to create IP Range: %s %s", cidr.IP.String(), err)
		}
	}

	err := txn.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to construct IPPool: %s %s", cidr.IP.String(), err)
	}

	return i, nil
}

func (i *IPPool) save(txn *badger.Txn, ipRange *IPRange) error {

	bytes, err := proto.Marshal(ipRange)
	if err != nil {
		return fmt.Errorf("failed to marshal IPRange: %s", err)
	}

	err = txn.Set(append(prefix_ip_range, i.key...), bytes)
	if err != nil {
		return fmt.Errorf("failed to add IPRange: %s", err)
	}

	return nil
}

func (r *IPPool) get(txn *badger.Txn) (*IPRange, error) {

	t, err := txn.Get(append(prefix_agent, r.key...))
	if err != nil {
		return nil, fmt.Errorf("failed to get IPRange: %s", err)
	}

	i := &IPRange{}
	err = t.Value(func(val []byte) error {
		return proto.Unmarshal(val, i)
	})
	if err != nil {
		r.l.WithError(err).Error("Failed to parse IPRange")
		return nil, fmt.Errorf("failed to parse IPRange: %s", err)
	}

	return i, nil
}
