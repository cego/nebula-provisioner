package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/slackhq/nebula/cert"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/curve25519"
	"google.golang.org/protobuf/types/known/emptypb"
)

var enrollCmd = &cobra.Command{
	Use:   "enroll",
	Short: "Enroll agent into server",
	Run: func(cmd *cobra.Command, args []string) {
		agent, err := NewClient(l, config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer agent.Close()

		updateEnrollmentRequest(agent, cmd)
	},
}

var enrollStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Enrollment status",
	Run: func(cmd *cobra.Command, args []string) {
		agent, err := NewClient(l, config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer agent.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := agent.client.GetEnrollStatus(ctx, &emptypb.Empty{})
		if err != nil {
			l.WithError(err).Fatalln("failed to get enrollment status")
		}

		if res.IsEnrolled {
			l.Info("Agent is enrolled")
			l.Infof("IssuedAt: %s, ExpiresAt: %s\n", res.IssuedAt.AsTime().Format(time.RFC3339), res.ExpiresAt.AsTime().Format(time.RFC3339))
		} else if res.IsEnrollmentRequested {
			l.Info("Agent has requested to be enrolled")
		} else {
			l.Info("Agent enrollment not started")
		}
	},
}

var enrollWaitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait for enrollment",
	Run: func(cmd *cobra.Command, args []string) {
		agent, err := NewClient(l, config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer agent.Close()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		updateEnrollmentRequest(agent, cmd)

		if isEnrollDone(agent) {
			return
		}

		for {
			select {
			case _ = <-ticker.C:
				if isEnrollDone(agent) {
					return
				}
			}
		}
	},
}

func init() {
	enrollCmd.Flags().StringP("token", "t", "", "Enrollment token")
	enrollCmd.Flags().StringSliceP("groups", "g", []string{}, "Comma separated list of groups")
	enrollCmd.Flags().StringP("ip", "i", "", "Requesting for this specific nebula ip")
	enrollCmd.MarkFlagRequired("token")
	enrollWaitCmd.Flags().StringP("token", "t", "", "Enrollment token")
	enrollWaitCmd.Flags().StringSliceP("groups", "g", []string{}, "Comma separated list of groups")
	enrollWaitCmd.Flags().StringP("ip", "i", "", "Requesting for this specific nebula ip")
	enrollWaitCmd.MarkFlagRequired("token")

	enrollCmd.AddCommand(enrollStatusCmd)
	enrollCmd.AddCommand(enrollWaitCmd)
}

func updateEnrollmentRequest(agent *agentClient, cmd *cobra.Command) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := agent.client.GetEnrollStatus(ctx, &emptypb.Empty{})
	if err != nil {
		l.WithError(err).Fatalln("failed to get enrollment status")
		os.Exit(1)
	}

	token, err := cmd.Flags().GetString("token")
	if err != nil {
		l.WithError(err).Fatalln("failed to get token")
	}
	if token == "" {
		l.Errorln("Require the enrollment has been started, you can provide enrollment token and doing it in one process.")
		os.Exit(1)
	}
	groups, err := cmd.Flags().GetStringSlice("groups")
	if err != nil {
		l.WithError(err).Fatalln("failed to get groups")
	}
	if len(groups) == 0 && agent.config.IsSet("enroll.groups") {
		groups = agent.config.GetStringSlice("enroll.groups", []string{})
	}

	ip, err := cmd.Flags().GetString("ip")
	if err != nil {
		l.WithError(err).Fatalln("failed to get ip")
	} else {
		parseIP := net.ParseIP(ip)
		if parseIP == nil && ip != "" {
			l.WithError(err).Fatalln("ip is of invalid format")
			os.Exit(1)
		}
	}

	if ip == "" && agent.config.IsSet("enroll.ip") {
		ip = agent.config.GetString("enroll.ip", "")
		parseIP := net.ParseIP(ip)
		if parseIP == nil && ip != "" {
			l.WithError(err).Fatalln("ip is of invalid format")
			os.Exit(1)
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		l.WithError(err).Errorln("error when getting hostname")
		os.Exit(1)
	}

	var diff = false

	if status.EnrollmentRequest != nil {
		l.Debug("comparing against existing enrollment request")
		l.Debugf("hostname compare %s <=> %s", hostname, status.EnrollmentRequest.Name)
		if hostname != status.EnrollmentRequest.Name {
			diff = true
			l.Debugf("diff on hostname")
		}

		l.Debugf("ip compare %s <=> %s", ip, status.EnrollmentRequest.RequestedIP)
		if ip != status.EnrollmentRequest.RequestedIP {
			diff = true
			l.Debugf("diff on ip")
		}

		l.Debugf("groups compare %s <=> %s", groups, status.EnrollmentRequest.Groups)
		if !stringSlicesEqual(groups, status.EnrollmentRequest.Groups) {
			diff = true
			l.Debugf("diff on groups")
		}
	} else if status.IsEnrolled {
		l.Debug("comparing against enrolled agent")
		l.Debugf("hostname compare %s <=> %s", hostname, status.Name)
		if hostname != status.Name {
			diff = true
			l.Debugf("diff on hostname")
		}

		l.Debugf("ip compare %s <=> %s", ip, status.AssignedIP)
		if ip != status.AssignedIP && ip != "" {
			diff = true
			l.Debugf("diff on ip")
		}

		l.Debugf("groups compare %s <=> %s", groups, status.Groups)
		if !stringSlicesEqual(groups, status.Groups) {
			diff = true
			l.Debugf("diff on groups")
		}
	} else {
		diff = true
	}

	if diff {
		l.Info("adding enrollment request")
		if err := enroll(agent, token, ip, groups); err != nil {
			l.WithError(err).Fatalln("failed to enroll agent")
			os.Exit(2)
		}
	}
}

func isEnrollDone(agent *agentClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	status, err := agent.client.GetEnrollStatus(ctx, &emptypb.Empty{})
	if err != nil {
		l.WithError(err).Fatalln("failed to get enrollment status")
		os.Exit(1)
	}

	if status.IsEnrolled && !status.IsEnrollmentRequested {
		l.Info("Agent is enrolled")
		return true
	}

	return false
}

func enroll(c *agentClient, enrollmentToken, ip string, groups []string) error {
	if enrollmentToken == "" {
		return fmt.Errorf("requires enrollmentToken")
	}

	csr, err := generateNebulaKeyPair()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	enrollRequest := &protocol.EnrollRequest{
		Token:  enrollmentToken,
		CsrPEM: string(csr),
	}

	if len(groups) > 0 {
		enrollRequest.Groups = groups
	}

	if ip != "" {
		enrollRequest.RequestedIP = ip
	}

	hostname, err := os.Hostname()
	if err == nil {
		enrollRequest.Name = hostname
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.client.Enroll(ctx, enrollRequest)
	if err != nil {
		return err
	}

	return nil
}

func generateNebulaKeyPair() ([]byte, error) {

	var csr []byte
	var err error

	exists, info := fileExists(resolvePath(AgentNebulaCsrPath))
	if exists && info.IsDir() {
		return nil, fmt.Errorf("expected agent-nebula.csr to be a file")
	} else if exists {
		csr, err = ioutil.ReadFile(resolvePath(AgentNebulaCsrPath))
		if err != nil {
			return nil, fmt.Errorf("error while reading csr: %s", err)
		}
	} else {
		var pubkey, privkey [32]byte
		if _, err := io.ReadFull(rand.Reader, privkey[:]); err != nil {
			panic(err)
		}
		curve25519.ScalarBaseMult(&pubkey, &privkey)

		err := ioutil.WriteFile(resolvePath(AgentNebulaKeyPath), cert.MarshalX25519PrivateKey(privkey[:]), 0600)
		if err != nil {
			return nil, fmt.Errorf("error while writing key: %s", err)
		}

		csr = cert.MarshalX25519PublicKey(pubkey[:])
		err = ioutil.WriteFile(resolvePath(AgentNebulaCsrPath), csr, 0600)
		if err != nil {
			return nil, fmt.Errorf("error while writing csr: %s", err)
		}
	}

	return csr, nil
}
