package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal server with secret part",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		if err = unseal(sc); err != nil {
			fmt.Printf("failed to to unseal server: %s\n", err)
			return
		}
	},
}

func unseal(c *serverClient) error {

	fmt.Print("Secret part to use for unsealing: ")
	part, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Println("")

	partPattern := "^[a-f0-9]*$"
	match, _ := regexp.Match(partPattern, part)
	if !match {
		return fmt.Errorf("secret part does not match %s", partPattern)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = c.client.Unseal(ctx, &protocol.UnsealRequest{
		KeyPart: string(part),
	})
	if err != nil {
		return err
	}

	fmt.Println("Successfully unsealed server")

	return nil
}
