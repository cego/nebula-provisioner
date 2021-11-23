package graph

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/slyngdk/nebula-provisioner/server/store"

	"github.com/slyngdk/nebula-provisioner/server/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	store     *store.Store
	ipManager *store.IPManager
	l         *logrus.Logger
}

func NewResolver(store *store.Store, ipManager *store.IPManager, l *logrus.Logger) *Resolver {
	return &Resolver{store: store, ipManager: ipManager, l: l}
}

func UserFormContext(ctx context.Context) model.User {
	raw, _ := ctx.Value("currentUser").(model.User)
	return raw
}

func WithUser(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, "currentUser", user)
}

func toPointerSliceString(in []string) []*string {
	out := make([]*string, len(in))
	for i := range in {
		out[i] = &in[i]
	}
	return out
}

func agentToModel(a *store.Agent) *model.Agent {
	agent := &model.Agent{
		Fingerprint: hex.EncodeToString(a.Fingerprint),
		Created:     a.Created.AsTime().Format(time.RFC3339),
		NetworkName: a.NetworkName,
		Groups:      toPointerSliceString(a.Groups),
		AssignedIP:  &a.AssignedIP,
		Name:        &a.Name,
	}
	if a.ExpiresAt != nil {
		expiresAt := a.ExpiresAt.AsTime().Format(time.RFC3339)
		agent.ExpiresAt = &expiresAt
	}
	if a.IssuedAt != nil {
		issuedAt := a.IssuedAt.AsTime().Format(time.RFC3339)
		agent.IssuedAt = &issuedAt
	}
	return agent
}

func userToModel(u *store.User) *model.User {
	var ua *model.UserApprove
	if u.Approve != nil {
		ua = &model.UserApprove{
			Approved:   u.Approve.Approved,
			ApprovedBy: u.Approve.ApprovedBy,
			ApprovedAt: u.Approve.ApprovedAt.AsTime().Format(time.RFC3339),
		}
	}
	return &model.User{
		ID:          u.Id,
		Name:        u.Name,
		Email:       u.Email,
		Disabled:    u.Disabled,
		UserApprove: ua,
	}
}

func networkToModel(n *protocol.Network) *model.Network {
	return &model.Network{
		Name:    n.Name,
		Groups:  toPointerSliceString(n.Groups),
		Ips:     toPointerSliceString(n.Ips),
		Subnets: toPointerSliceString(n.Subnets),
		IPPools: toPointerSliceString(n.IpPools),
	}
}
