package main

import (
	"context"
	"fmt"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"time"
)

var enrollCmd = &cobra.Command{
	Use:   "enroll",
	Short: "Managing enrollment tokens",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		network, _ := cmd.Flags().GetString("network")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := sc.client.GetEnrollmentTokenForNetwork(ctx, &protocol.GetEnrollmentTokenForNetworkRequest{
			Network: network,
		})
		if err != nil {
			fmt.Printf("failed to get enrollment token: %s\n", err)
			return
		}

		fmt.Println(res)
	},
}

var enrollmentPendingCmd = &cobra.Command{
	Use:   "pending",
	Short: "Lists pending enrollment requests",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := sc.client.ListEnrollmentRequests(ctx, &emptypb.Empty{})
		if err != nil {
			fmt.Printf("failed to get enrollment tokens: %s\n", err)
			return
		}

		for _, request := range res.EnrollmentRequests {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", request.NetworkName, request.Name, request.Created.AsTime().Format(time.RFC3339), request.ClientFingerprint, request.ClientIP)
		}
	},
}

var enrollmentApproveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve pending enrollment requests",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		fingerprint, _ := cmd.Flags().GetString("fingerprint")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = sc.client.ApproveEnrollmentRequest(ctx, &protocol.ApproveEnrollmentRequestRequest{
			ClientFingerprint: fingerprint,
		})
		if err != nil {
			fmt.Printf("failed to approve enrollment request: %s\n", err)
			return
		}
	},
}

func init() {
	enrollCmd.Flags().StringP("network", "n", "", "")
	enrollCmd.MarkFlagRequired("network")
	enrollmentApproveCmd.Flags().StringP("fingerprint", "f", "", "")
	enrollmentApproveCmd.MarkFlagRequired("fingerprint")

	enrollCmd.AddCommand(enrollmentPendingCmd, enrollmentApproveCmd)
}
