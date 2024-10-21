package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cego/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
)

var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "Managing certificate authorities",
}

var caListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing certificate authorities",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		if err = listCA(sc); err != nil {
			fmt.Printf("failed certificate authorities: %s\n", err)
			return
		}
	},
}

func init() {
	caListCmd.Flags().StringSliceP("networks", "n", nil, "comma separated list of network names, to filter CA by")

	caCmd.AddCommand(caListCmd)
}

func listCA(c *serverClient) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.client.ListCertificateAuthorityByNetwork(ctx, &protocol.ListCertificateAuthorityByNetworkRequest{})
	if err != nil {
		return err
	}

	fmt.Println("CA's")
	fmt.Println("Network\tSha256sum")
	for _, ca := range res.CertificateAuthorities {
		fmt.Printf("%s\t%s\n", ca.NetworkName, ca.Sha256Sum)
	}

	return nil
}
