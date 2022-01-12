// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Agent struct {
	Fingerprint string    `json:"fingerprint"`
	Created     string    `json:"created"`
	NetworkName string    `json:"networkName"`
	Groups      []*string `json:"groups"`
	AssignedIP  *string   `json:"assignedIP"`
	IssuedAt    *string   `json:"issuedAt"`
	ExpiresAt   *string   `json:"expiresAt"`
	Name        *string   `json:"name"`
}

type Ca struct {
	Fingerprint string   `json:"fingerprint"`
	Status      CAStatus `json:"status"`
	IssuedAt    string   `json:"issuedAt"`
	ExpiresAt   string   `json:"expiresAt"`
}

type EnrollmentRequest struct {
	Fingerprint string    `json:"fingerprint"`
	Created     string    `json:"created"`
	NetworkName string    `json:"networkName"`
	ClientIP    *string   `json:"clientIP"`
	Name        *string   `json:"name"`
	RequestedIP *string   `json:"requestedIP"`
	Groups      []*string `json:"groups"`
}

type Network struct {
	Name               string               `json:"name"`
	Duration           *string              `json:"duration"`
	Groups             []*string            `json:"groups"`
	Ips                []*string            `json:"ips"`
	Subnets            []*string            `json:"subnets"`
	IPPools            []*string            `json:"ipPools"`
	Agents             []*Agent             `json:"agents"`
	EnrollmentToken    *string              `json:"enrollmentToken"`
	EnrollmentRequests []*EnrollmentRequest `json:"enrollmentRequests"`
	Cas                []*Ca                `json:"cas"`
}

type User struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	UserApprove *UserApprove `json:"userApprove"`
	Disabled    bool         `json:"disabled"`
}

type UserApprove struct {
	Approved       bool   `json:"approved"`
	ApprovedBy     string `json:"approvedBy"`
	ApprovedByUser *User  `json:"approvedByUser"`
	ApprovedAt     string `json:"approvedAt"`
}

type CAStatus string

const (
	CAStatusActive   CAStatus = "active"
	CAStatusExpired  CAStatus = "expired"
	CAStatusInactive CAStatus = "inactive"
	CAStatusNext     CAStatus = "next"
)

var AllCAStatus = []CAStatus{
	CAStatusActive,
	CAStatusExpired,
	CAStatusInactive,
	CAStatusNext,
}

func (e CAStatus) IsValid() bool {
	switch e {
	case CAStatusActive, CAStatusExpired, CAStatusInactive, CAStatusNext:
		return true
	}
	return false
}

func (e CAStatus) String() string {
	return string(e)
}

func (e *CAStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CAStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CAStatus", str)
	}
	return nil
}

func (e CAStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
