package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/slyngdk/nebula-provisioner/protocol"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Managing networks",
}

var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listing networks",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		if err = listNetworks(sc); err != nil {
			fmt.Printf("failed list networks: %s\n", err)
			return
		}
	},
}

var networkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating new network",
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := NewClient(config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer sc.Close()

		name, _ := cmd.Flags().GetString("name")

		duration, _ := cmd.Flags().GetDuration("duration")
		if duration < 0 {
			fmt.Println("duration needs to be a positive number")
			return
		}
		var d *durationpb.Duration
		if duration > 0 {
			d = durationpb.New(duration)
		}

		groups, _ := cmd.Flags().GetStringSlice("groups")
		ips, _ := cmd.Flags().GetStringSlice("ips")
		subnets, _ := cmd.Flags().GetStringSlice("subnets")
		pools, _ := cmd.Flags().GetStringSlice("pool")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := sc.client.CreateNetwork(ctx, &protocol.CreateNetworkRequest{
			Name:     name,
			Duration: d,
			Groups:   groups,
			Ips:      ips,
			Subnets:  subnets,
			IpPools:  pools,
		})
		if err != nil {
			fmt.Printf("failed create network: %s\n", err)
			return
		}

		fmt.Println(res)
	},
}

func init() {

	networkCreateCmd.Flags().StringP("name", "n", "", "Name of the network and used also used as name of the certificate authority")
	networkCreateCmd.MarkFlagRequired("name")
	networkCreateCmd.Flags().DurationP("duration", "d", 0, "amount of time the certificate should be valid for. Valid time units are seconds: \"s\", minutes: \"m\", hours: \"h\" (default 8760h0m0s)")
	networkCreateCmd.Flags().StringSliceP("groups", "g", nil, "comma separated list of groups. This will limit which groups subordinate certs can use")
	networkCreateCmd.Flags().StringSliceP("ips", "i", nil, "comma separated list of ip and network in CIDR notation. This will limit which ip addresses and networks subordinate certs can use")
	networkCreateCmd.Flags().StringSliceP("subnets", "s", nil, "comma separated list of ip and network in CIDR notation. This will limit which subnet addresses and networks subordinate certs can use")
	networkCreateCmd.Flags().StringSliceP("pool", "p", nil, "comma separated list of ip and network in CIDR notation. This will be used to assign IP's")

	networkCmd.AddCommand(networkListCmd, networkCreateCmd)
}

func listNetworks(c *serverClient) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.client.ListNetwork(ctx, &protocol.ListNetworkRequest{})
	if err != nil {
		return err
	}

	fmt.Println("Networks")
	fmt.Println("Name\tDuration\tGroups\tIPs\tSubnets\tPools")
	for _, network := range res.Networks {
		fmt.Printf("%s\t%s\t\t%s\t%s\t%s\t%s\n", network.Name, network.Duration, network.Groups, network.Ips, network.Subnets, network.IpPools)
	}

	return nil
}
