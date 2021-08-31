package store

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
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

	cidr        *net.IPNet
	ipsInUse    []*net.IPAddr
	ipsNotInUse []*net.IPAddr

	db   *badger.DB
	lock *sync.Mutex
	key  []byte
}

func (r *IPPool) next() {

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

func IsUsableBlock(b *net.IPNet) bool {
	for _, block := range usableIPBlocks {
		if block.Contains(b.IP) && block.Contains(broadcastIP(b)) {
			return true
		}
	}
	return false
}

func (s *Store) NewIPRange(networkName string, cidr *net.IPNet) (*IPPool, error) {
	s.l.Infof("Creating IPPool for: %s", cidr.IP.String())

	if !IsUsableBlock(cidr) {
		return nil, fmt.Errorf("cidr %s is not in an uable private range", cidr)
	}

	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	key := append([]byte(networkName), []byte("-")...)
	key = append(key, []byte(cidr.IP.String())...)

	i := &IPPool{l: s.l, cidr: cidr, db: s.db, lock: &sync.Mutex{}, key: key}

	var r *IPRange
	if !exists(txn, prefix_ip_range, key) {
		fmt.Printf("Adding IP Range : %s\n", cidr.IP)
		r = &IPRange{Network: cidr.IP, Netmask: cidr.Mask}
		err := i.save(txn, r)
		if err != nil {
			return nil, fmt.Errorf("failed to create IP Range: %s %s", cidr.IP.String(), err)
		}
	} else {

	}

	err := txn.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to construct IPPool: %s %s", cidr.IP.String(), err)
	}

	return i, nil
}

func (r *IPPool) save(txn *badger.Txn, ipRange *IPRange) error {

	bytes, err := proto.Marshal(ipRange)
	if err != nil {
		return fmt.Errorf("failed to marshal IPRange: %s", err)
	}

	err = txn.Set(append(prefix_ip_range, r.key...), bytes)
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

//type IPAllocator struct {
//	l *logrus.Logger
//}
//
//func NewIPAllocator(l *logrus.Logger, network *net.IPNet) *IPAllocator {
//	return &IPAllocator{l}
//}
