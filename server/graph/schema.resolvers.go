package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/slackhq/nebula/cert"
	"github.com/slyngdk/nebula-provisioner/server/graph/generated"
	"github.com/slyngdk/nebula-provisioner/server/graph/model"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ApproveUser is the resolver for the approveUser field.
func (r *mutationResolver) ApproveUser(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, gqlerror.Errorf("userId is required")
	}
	currentUser := UserFormContext(ctx)

	user, err := r.store.ApproveUserAccess(userID, &store.UserApprove{Approved: true, ApprovedBy: currentUser.ID})
	if err != nil {
		return nil, gqlerror.Errorf("Failed to approve user access: %s", err)
	}

	return userToModel(user), nil
}

// DisableUser is the resolver for the disableUser field.
func (r *mutationResolver) DisableUser(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, gqlerror.Errorf("userId is required")
	}

	user, err := r.store.DisableUserAccess(userID)
	if err != nil {
		return nil, gqlerror.Errorf("Failed to disable user access: %s", err)
	}

	return userToModel(user), nil
}

// ApproveEnrollmentRequest is the resolver for the approveEnrollmentRequest field.
func (r *mutationResolver) ApproveEnrollmentRequest(ctx context.Context, fingerprint string) (*model.Agent, error) {
	if fingerprint == "" {
		return nil, gqlerror.Errorf("fingerprint is required")
	}

	bytes, err := hex.DecodeString(fingerprint)
	if err != nil {
		return nil, gqlerror.Errorf("failed to decode fingerprint: %s", fingerprint)
	}

	agent, err := r.store.ApproveEnrollmentRequest(r.ipManager, bytes)
	if err != nil {
		r.l.WithError(err).Errorf("Failed to approve enrollment request: %s", fingerprint)
		return nil, gqlerror.Errorf("Failed to approve enrollment request: %s", fingerprint)
	}

	return agentToModel(agent), nil
}

// DeleteEnrollmentRequest is the resolver for the deleteEnrollmentRequest field.
func (r *mutationResolver) DeleteEnrollmentRequest(ctx context.Context, fingerprint string) (*bool, error) {
	if fingerprint == "" {
		return nil, gqlerror.Errorf("fingerprint is required")
	}

	bytes, err := hex.DecodeString(fingerprint)
	if err != nil {
		return nil, gqlerror.Errorf("failed to decode fingerprint: %s", fingerprint)
	}

	err = r.store.DeleteEnrollmentRequest(bytes)
	if err != nil {
		r.l.WithError(err).Errorf("failed to delete enrollment request: %s", fingerprint)
		return nil, gqlerror.Errorf("failed to delete enrollment request: %s", fingerprint)
	}

	return nil, nil
}

// RevokeCertsForAgent is the resolver for the revokeCertsForAgent field.
func (r *mutationResolver) RevokeCertsForAgent(ctx context.Context, fingerprint string) (*bool, error) {
	if fingerprint == "" {
		return nil, gqlerror.Errorf("fingerprint is required")
	}

	bytes, err := hex.DecodeString(fingerprint)
	if err != nil {
		return nil, gqlerror.Errorf("failed to decode fingerprint: %s", fingerprint)
	}
	err = r.store.RevokeAgent(bytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// PrepareNextCa is the resolver for the prepareNextCA field.
func (r *mutationResolver) PrepareNextCa(ctx context.Context, networkName string) (*bool, error) {
	if networkName == "" {
		return nil, gqlerror.Errorf("networkName is required")
	}

	err := r.store.PrepareCARollover(networkName)
	if err != nil {
		r.l.WithError(err).Errorf("failed to prepare rollover of CA")
		return nil, gqlerror.Errorf("failed to prepare rollover of CA")
	}

	return nil, nil
}

// SwitchActiveCa is the resolver for the switchActiveCA field.
func (r *mutationResolver) SwitchActiveCa(ctx context.Context, networkName string) (*bool, error) {
	if networkName == "" {
		return nil, gqlerror.Errorf("networkName is required")
	}

	err := r.store.SwitchActiveCA(networkName)
	if err != nil {
		r.l.WithError(err).Errorf("failed to switch active CA")
		return nil, gqlerror.Errorf("failed to switch active CA")
	}

	return nil, nil
}

// Agents is the resolver for the agents field.
func (r *networkResolver) Agents(ctx context.Context, obj *model.Network) ([]*model.Agent, error) {
	if obj == nil {
		return []*model.Agent{}, nil
	}

	agents, err := r.store.ListAgentByNetwork(obj.Name)
	if err != nil {
		r.l.WithError(err).Errorf("failed to get agents for network: %s", obj.Name)
		return nil, gqlerror.Errorf("failed to get agents for network: %s", obj.Name)
	}

	gAgent := make([]*model.Agent, len(agents))
	for i, a := range agents {
		gAgent[i] = agentToModel(a)
	}
	return gAgent, nil
}

// EnrollmentToken is the resolver for the enrollmentToken field.
func (r *networkResolver) EnrollmentToken(ctx context.Context, obj *model.Network) (*string, error) {
	if obj == nil {
		return nil, nil
	}

	et, err := r.store.GetEnrollmentTokenByNetwork(obj.Name)
	if err != nil {
		r.l.WithError(err).Errorf("failed to get enrollment token for network: %s", obj.Name)
		return nil, gqlerror.Errorf("failed to get enrollment token for network: %s", obj.Name)
	}

	return &et.Token, nil
}

// EnrollmentRequests is the resolver for the enrollmentRequests field.
func (r *networkResolver) EnrollmentRequests(ctx context.Context, obj *model.Network) ([]*model.EnrollmentRequest, error) {
	if obj == nil {
		return []*model.EnrollmentRequest{}, nil
	}

	enrollmentRequests, err := r.store.ListEnrollmentRequestsByNetwork(obj.Name)
	if err != nil {
		r.l.WithError(err).Errorf("failed to get enrollement requests for network: %s", obj.Name)
		return nil, gqlerror.Errorf("failed to get enrollement requests for network: %s", obj.Name)
	}

	gEnrollmentRequests := make([]*model.EnrollmentRequest, len(enrollmentRequests))
	for i, a := range enrollmentRequests {
		gEnrollmentRequests[i] = &model.EnrollmentRequest{
			Fingerprint: hex.EncodeToString(a.Fingerprint),
			Created:     a.Created.AsTime().Format(time.RFC3339),
			NetworkName: a.NetworkName,
			ClientIP:    &a.ClientIP,
			Name:        &a.Name,
			Groups:      toPointerSliceString(a.Groups),
		}
		if a.RequestedIP != "" {
			gEnrollmentRequests[i].RequestedIP = &a.RequestedIP
		}
	}
	return gEnrollmentRequests, nil
}

// Cas is the resolver for the cas field.
func (r *networkResolver) Cas(ctx context.Context, obj *model.Network) ([]*model.Ca, error) {
	if obj == nil {
		return []*model.Ca{}, nil
	}

	cas, err := r.store.ListCAByNetwork([]string{obj.Name})
	if err != nil {
		r.l.WithError(err).Errorf("failed to get ca`s for network: %s", obj.Name)
		return nil, gqlerror.Errorf("failed to get ca`s for network: %s", obj.Name)
	}

	gCas := make([]*model.Ca, len(cas))
	for i, ca := range cas {
		publicKey, _, err := cert.UnmarshalNebulaCertificateFromPEM(ca.PublicKey)
		if err != nil {
			r.l.WithError(err).Error("failed to parse public key of ca")
			return nil, gqlerror.Errorf("failed to parse public key of ca")
		}

		fingerprint, err := publicKey.Sha256Sum()
		if err != nil {
			r.l.WithError(err).Error("failed to get fingerprint of ca for network")
			return nil, gqlerror.Errorf("failed to get fingerprint of ca")
		}

		gCas[i] = &model.Ca{
			Fingerprint: fingerprint,
			Status:      caStatusToModel(ca.Status),
			IssuedAt:    publicKey.Details.NotBefore.Format(time.RFC3339),
			ExpiresAt:   publicKey.Details.NotAfter.Format(time.RFC3339),
		}
	}

	return gCas, nil
}

// CurrentUser is the resolver for the currentUser field.
func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	user := UserFormContext(ctx)

	return &user, nil
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context) ([]*model.User, error) {
	users, err := r.store.ListUsers()
	if err != nil {
		r.l.WithError(err).Error("failed to get users")
		return nil, gqlerror.Errorf("failed to get users")
	}

	gUsers := make([]*model.User, len(users))

	for i, u := range users {
		gUsers[i] = userToModel(u)
	}

	return gUsers, nil
}

// GetNetworks is the resolver for the getNetworks field.
func (r *queryResolver) GetNetworks(ctx context.Context) ([]*model.Network, error) {
	networks, err := r.store.ListNetworks()
	if err != nil {
		r.l.WithError(err).Error("failed to get networks")
		return nil, gqlerror.Errorf("failed to get networks")
	}

	gNetworks := make([]*model.Network, 0)

	for _, n := range networks {
		gNetworks = append(gNetworks, networkToModel(n))
	}

	return gNetworks, nil
}

// GetNetwork is the resolver for the getNetwork field.
func (r *queryResolver) GetNetwork(ctx context.Context, name string) (*model.Network, error) {
	if name == "" {
		return nil, gqlerror.Errorf("name is required")
	}

	network, err := r.store.GetNetworkByName(name)
	if err != nil {
		r.l.WithError(err).Errorf("failed to get network: %s", name)
		return nil, gqlerror.Errorf("failed to get network: %s", name)
	}

	return networkToModel(network), nil
}

// GetAgent is the resolver for the getAgent field.
func (r *queryResolver) GetAgent(ctx context.Context, fingerprint string) (*model.Agent, error) {
	if fingerprint == "" {
		return nil, gqlerror.Errorf("fingerprint is required")
	}

	bytes, err := hex.DecodeString(fingerprint)
	if err != nil {
		return nil, gqlerror.Errorf("failed to decode fingerprint: %s", fingerprint)
	}

	isAgentEnrolled := r.store.IsAgentEnrolled(bytes)
	if !isAgentEnrolled {
		return nil, nil
	}

	agent, err := r.store.GetAgentByFingerprint(bytes)
	if err != nil {
		r.l.WithError(err).Errorf("failed to get agent: %s", fingerprint)
		return nil, gqlerror.Errorf("failed to get agent: %s", fingerprint)
	}

	return agentToModel(agent), nil
}

// ApprovedByUser is the resolver for the approvedByUser field.
func (r *userApproveResolver) ApprovedByUser(ctx context.Context, obj *model.UserApprove) (*model.User, error) {
	if obj == nil {
		return nil, nil
	}

	if obj.ApprovedBy == "server-client" {
		return &model.User{
			ID:   "server-client",
			Name: "server-client",
		}, nil
	}

	user, err := r.store.GetUserByID(obj.ApprovedBy)
	if err != nil {
		return nil, nil
	}
	return userToModel(user), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Network returns generated.NetworkResolver implementation.
func (r *Resolver) Network() generated.NetworkResolver { return &networkResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// UserApprove returns generated.UserApproveResolver implementation.
func (r *Resolver) UserApprove() generated.UserApproveResolver { return &userApproveResolver{r} }

type mutationResolver struct{ *Resolver }
type networkResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userApproveResolver struct{ *Resolver }
