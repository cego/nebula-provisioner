package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/slackhq/nebula/cert"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/curve25519"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"io/ioutil"
	"os"
	"time"
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

		token, err := cmd.Flags().GetString("token")
		if err != nil {
			l.WithError(err).Fatalln("failed to get token")
		}

		if err = enroll(agent, token); err != nil {
			l.WithError(err).Fatalln("failed to enroll to server")
		}
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

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		status := getStatus()
		if status == 2 {
			return
		}

		for {
			select {
			case _ = <-ticker.C:
				status := getStatus()
				if status == 2 {
					return
				}
			}
		}
	},
}

func getStatus() int8 {
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
		return 2
	} else if res.IsEnrollmentRequested {
		l.Info("Agent has requested to be enrolled")
		return 1
	} else {
		l.Info("Agent enrollment not started")
		return 0
	}
}

func init() {
	enrollCmd.Flags().StringP("token", "t", "", "Enrollment token")
	enrollCmd.MarkFlagRequired("token")

	enrollCmd.AddCommand(enrollStatusCmd)
	enrollCmd.AddCommand(enrollWaitCmd)
}

func enroll(c *agentClient, enrollmentToken string) error {
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

	hostname, err := os.Hostname()
	if err == nil {
		enrollRequest.Name = hostname
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.client.Enroll(ctx, enrollRequest)
	if err != nil {
		return err
	}

	c.l.Println(res)

	return nil
}

func generateNebulaKeyPair() ([]byte, error) {

	var csr []byte
	var err error

	exists, info := fileExists("agent-nebula.csr")
	if exists && info.IsDir() {
		return nil, fmt.Errorf("expected agent-nebula.csr to be a file")
	} else if exists {
		csr, err = ioutil.ReadFile("agent-nebula.csr")
		if err != nil {
			return nil, fmt.Errorf("error while reading csr: %s", err)
		}
	} else {
		var pubkey, privkey [32]byte
		if _, err := io.ReadFull(rand.Reader, privkey[:]); err != nil {
			panic(err)
		}
		curve25519.ScalarBaseMult(&pubkey, &privkey)

		err := ioutil.WriteFile("agent-nebula.key", cert.MarshalX25519PrivateKey(privkey[:]), 0600)
		if err != nil {
			return nil, fmt.Errorf("error while writing key: %s", err)
		}

		csr = cert.MarshalX25519PublicKey(pubkey[:])
		err = ioutil.WriteFile("agent-nebula.csr", csr, 0600)
		if err != nil {
			return nil, fmt.Errorf("error while writing csr: %s", err)
		}
	}

	return csr, nil
}
