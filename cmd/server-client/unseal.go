package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"regexp"
	"time"
)

type unsealFlags struct {
	set *flag.FlagSet
}

func newUnsealFlags() *unsealFlags {
	flags := unsealFlags{set: flag.NewFlagSet("unseal", flag.ContinueOnError)}
	flags.set.Usage = func() {}
	return &flags
}

func runUnseal(args []string, out io.Writer, errOut io.Writer) error {
	flags := newUnsealFlags()
	err := flags.set.Parse(args)
	if err != nil {
		return err
	}

	sc, err := NewClient()
	if err != nil {
		return err
	}
	defer sc.Close()

	return sc.unseal(out)
}

func unsealSummary() string {
	return "unseal <flags>: unseal server with secret part"
}

func unsealHelp(out io.Writer) {
	cf := newUnsealFlags()
	out.Write([]byte("Usage of " + os.Args[0] + " " + unsealSummary() + "\n"))
	cf.set.SetOutput(out)
	cf.set.PrintDefaults()
}

func (c serverClient) unseal(out io.Writer) error {

	fmt.Fprint(out, "Secret part to use for unsealing: ")
	part, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "")

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

	fmt.Fprintln(out,"Successfully unsealed server")

	return nil
}
