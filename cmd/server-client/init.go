package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/grpc"
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

	if *flags.check {
		return isInit(out)
	}

	return fmt.Errorf("NotImplemented")
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

func isInit(out io.Writer) error {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("unix:///tmp/nebula-provisioner.socket", opts...) // TODO change socket path
	if err != nil {
		return err
	}
	defer conn.Close()

	client := protocol.NewServerCommandClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.IsInit(ctx, &empty.Empty{})
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
