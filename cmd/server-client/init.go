package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"io"
	"os"
	"time"
)

type initFlags struct {
	set   *flag.FlagSet
	check *bool
}

func newInitFlags() *initFlags {
	flags := initFlags{set: flag.NewFlagSet("init", flag.ContinueOnError)}
	flags.set.Usage = func() {}
	flags.check = flags.set.Bool("check", false, "Check if server is initialized")
	return &flags
}

func runInit(args []string, out io.Writer, errOut io.Writer) error {
	flags := newInitFlags()
	err := flags.set.Parse(args)
	if err != nil {
		return err
	}

	sc, err := NewClient()
	if err != nil {
		return err
	}
	defer sc.Close()

	if *flags.check {
		return sc.isInit(out)
	}

	return sc.init(out)
}

func initSummary() string {
	return "init <flags>: init provision server with new shared secret"
}

func initHelp(out io.Writer) {
	cf := newInitFlags()
	out.Write([]byte("Usage of " + os.Args[0] + " " + initSummary() + "\n"))
	cf.set.SetOutput(out)
	cf.set.PrintDefaults()
}

func (c serverClient) isInit(out io.Writer) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.client.IsInit(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	if response.IsInitialized {
		fmt.Fprintln(out, "Server is initialized")
	} else {
		fmt.Fprintln(out, "Server is not yet initialized")
	}
	return nil
}

func (c serverClient) init(out io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isInitRes, err := c.client.IsInit(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	fmt.Fprint(out,"How many parts do you want to split the key in? ")
	var nParts uint32
	_, err = fmt.Scanf("%d", &nParts)
	if err != nil {
		return err
	}
	fmt.Fprint(out,"How many parts will you require to unseal the server? ")
	var threshold uint32
	_, err = fmt.Scanf("%d", &threshold)
	if err != nil {
		return err
	}

	if isInitRes.IsInitialized {
		fmt.Fprintln(out, "Server is already initialized")
		return nil
	}

	response, err := c.client.Init(ctx, &protocol.InitRequest{
		KeyParts:     nParts,
		KeyThreshold: threshold,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "Keep these parts stored secure and seperated!!!")
	fmt.Fprintf(out, "You require %d parts to unseal the server.\n", threshold)
	for _, keyPart := range response.KeyParts {
		fmt.Fprintln(out, keyPart)
	}

	return nil
}
