package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"github.com/slyngdk/nebula-provisioner/server"
)

// Build A version string that can be set with
//
//     -ldflags "-X main.Build=SOMEVERSION"
//
// at compile-time.
var Build string

func main() {
	configPath := flag.String("config", "", "Path to either a file or directory to load configuration from")
	printVersion := flag.Bool("version", false, "Print version")
	printUsage := flag.Bool("help", false, "Print command line usage")

	flag.Parse()

	if *printVersion {
		fmt.Printf("Version: %s\n", Build)
		os.Exit(0)
	}

	if *printUsage {
		flag.Usage()
		os.Exit(0)
	}

	if *configPath == "" {
		fmt.Println("-config flag must be set")
		flag.Usage()
		os.Exit(1)
	}

	l := logrus.New()
	l.Out = os.Stdout

	config := nebula.NewConfig(l)
	err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("failed to load config: %s", err)
		os.Exit(1)
	}

	c, err := server.Main(config, Build, l)
	if err != nil {
		l.WithError(err).Error("Failed to start")
		os.Exit(1)
	}
	defer c.Stop()

	go c.Start()
	c.ShutdownBlock()
}
