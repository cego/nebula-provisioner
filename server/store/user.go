package store

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dgraph-io/badger/v3"
	"google.golang.org/protobuf/proto"
)

func (s *Store) GetUserByID(id string) (*User, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	return s.getUserByID(txn, id)
}

func (s *Store) AddUser(user *User) (*User, error) {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	user, err := s.addUser(txn, user)
	if err != nil {
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %s", err)
	}
	s.l.Infof("User is created: %s %s", user.Id, user.Email)
	return user, nil
}

func (s *Store) IsUserApproved(id string) (*User, bool) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	user, err := s.getUserByID(txn, id)
	if err != nil {
		return user, false
	}

	if user.Approve != nil && user.Approve.Approved {
		return user, true
	}

	return user, false
}

func (s *Store) ListUsersWaitingForApproval() ([]*User, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.Prefix = prefix_user
	it := txn.NewIterator(opts)
	defer it.Close()

	var users []*User

	for it.Seek(prefix_user); it.ValidForPrefix(prefix_user); it.Next() {
		item := it.Item()
		err := item.Value(func(v []byte) error {
			u := &User{}
			if err := proto.Unmarshal(v, u); err != nil {
				s.l.WithError(err).Error("Failed to parse user")
			}
			if u.Approve == nil || !u.Approve.Approved {
				users = append(users, u)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (s *Store) ApproveUserAccess(userId string, approve *UserApprove) (*User, error) {
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	user, err := s.approveUserAccess(txn, userId, approve)
	if err != nil {
		return nil, err
	}
	err = txn.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %s", err)
	}
	s.l.Infof("User is approved: %s %s", user.Id, user.Email)
	return user, nil
}

func (s *Store) getUserByID(txn *badger.Txn, id string) (*User, error) {
	t, err := txn.Get(append(prefix_user, id...))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %s", err)
	}

	u := &User{}
	err = t.Value(func(val []byte) error {
		return proto.Unmarshal(val, u)
	})
	if err != nil {
		s.l.WithError(err).Error("Failed to parse user")
		return nil, fmt.Errorf("failed to parse user: %s", err)
	}
	return u, nil
}

func (s *Store) addUser(txn *badger.Txn, user *User) (*User, error) {

	if exists(txn, prefix_user, []byte(user.Id)) {
		return nil, fmt.Errorf("user already exists")
	}

	if user.Id == "" {
		return nil, fmt.Errorf("missing id on user")
	}
	if user.Name == "" {
		return nil, fmt.Errorf("missing name on user")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("missing email on user")
	}

	user.Approve = nil
	user.Created = timestamppb.Now()

	bytes, err := proto.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %s", err)
	}

	err = txn.Set(append(prefix_user, user.Id...), bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %s", err)
	}

	return user, nil
}

func (s *Store) approveUserAccess(txn *badger.Txn, id string, approve *UserApprove) (*User, error) {

	if !exists(txn, prefix_user, []byte(id)) {
		return nil, fmt.Errorf("user don't exists")
	}

	if approve == nil {
		return nil, fmt.Errorf("user approve is nil")
	}
	if approve.ApprovedBy == "" {
		return nil, fmt.Errorf("missing approvedBy on user approve")
	}

	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user.Approve != nil && user.Approve.Approved {
		return nil, fmt.Errorf("user is already approved")
	}

	user.Approve = approve
	user.Approve.ApprovedAt = timestamppb.Now()

	bytes, err := proto.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %s", err)
	}

	err = txn.Set(append(prefix_user, user.Id...), bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %s", err)
	}

	return user, nil
}
