package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize server with new shared secret",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		check, err := cmd.Flags().GetBool("check")
		if err != nil {
			fmt.Printf("failed to get check flag: %s\n", err)
			return
		}

		if check {
			if err = isInit(sc); err != nil {
				fmt.Printf("failed to check is server is initialized: %s\n", err)
				return
			}
		}

		if err = initServer(sc); err != nil {
			fmt.Printf("failed to initialize server: %s\n", err)
			return
		}
	},
}

func init() {
	initCmd.Flags().Bool("check", false, "Check if server is initialized")
}

func isInit(c *serverClient) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.client.IsInit(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	if response.IsInitialized {
		fmt.Println("Server is initialized")
	} else {
		fmt.Println("Server is not yet initialized")
	}
	return nil
}

func initServer(c *serverClient) error {

	ctx, cancelCheck := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCheck()
	isInitRes, err := c.client.IsInit(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	if isInitRes.IsInitialized {
		fmt.Println("Server is already initialized")
		return nil
	}

	fmt.Print("How many parts do you want to split the key in? ")
	var nParts uint32
	_, err = fmt.Scanf("%d", &nParts)
	if err != nil {
		return err
	}
	fmt.Print("How many parts will you require to unseal the server? ")
	var threshold uint32
	_, err = fmt.Scanf("%d", &threshold)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := c.client.Init(ctx, &protocol.InitRequest{
		KeyParts:     nParts,
		KeyThreshold: threshold,
	})
	if err != nil {
		return err
	}

	fmt.Println("Keep these parts stored secure and seperated!!!")
	fmt.Printf("You require %d parts to unseal the server.\n", threshold)
	for _, keyPart := range response.KeyParts {
		fmt.Println(keyPart)
	}

	fmt.Println()
	fmt.Println("Server is now initialized and ready to be unseal")

	return nil
}
