package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/slyngdk/nebula-provisioner/protocol"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Managing users",
}

var userApprovePendingCmd = &cobra.Command{
	Use:   "pending",
	Short: "Lists users pending approval for login",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := sc.client.ListUsersWaitingForApproval(ctx, &emptypb.Empty{})
		if err != nil {
			fmt.Printf("failed to get users: %s\n", err)
			return
		}

		for _, request := range res.Users {
			fmt.Printf("%s\t%s\t%s\t\n", request.Id, request.Email, request.Name)
		}
	},
}

var userApproveCmd = &cobra.Command{
	Use:   "approve",
	Short: "Approve access to user",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		id, _ := cmd.Flags().GetString("id")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = sc.client.ApproveUserAccess(ctx, &protocol.ApproveUserAccessRequest{
			UserId: id,
		})
		if err != nil {
			fmt.Printf("failed to approve access to user: %s\n", err)
			return
		}
	},
}

func init() {
	userApproveCmd.Flags().StringP("id", "i", "", "Id of the user to approve")
	userApproveCmd.MarkFlagRequired("id")

	userCmd.AddCommand(userApprovePendingCmd, userApproveCmd)
}
