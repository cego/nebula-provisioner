package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var Build string

type helpError struct {
	s string
}

func (he *helpError) Error() string {
	return he.s
}

func newHelpErrorf(s string, v ...interface{}) error {
	return &helpError{s: fmt.Sprintf(s, v...)}
}
func main() {
	flag.Usage = func() {
		help("", os.Stderr)
		os.Exit(1)
	}

	printVersion := flag.Bool("version", false, "Print version")
	flagHelp := flag.Bool("help", false, "Print command line usage")
	flagH := flag.Bool("h", false, "Print command line usage")
	printUsage := false

	flag.Parse()

	if *flagH || *flagHelp {
		printUsage = true
	}

	args := flag.Args()

	if *printVersion {
		fmt.Printf("Version: %v\n", Build)
		os.Exit(0)
	}

	if len(args) < 1 {
		if printUsage {
			help("", os.Stderr)
			os.Exit(0)
		}

		help("No mode was provided", os.Stderr)
		os.Exit(1)
	} else if printUsage {
		handleError(args[0], &helpError{}, os.Stderr)
		os.Exit(0)
	}

	var err error
	switch args[0] {
	case "init":
		err = runInit(args[1:], os.Stdout, os.Stderr)
	default:
		err = fmt.Errorf("unknown mode: %s", args[0])
	}

	if err != nil {
		os.Exit(handleError(args[0], err, os.Stderr))
	}

}

func help(err string, out io.Writer) {
	if err != "" {
		fmt.Fprintln(out, "Error:", err)
		fmt.Fprintln(out, "")
	}

	fmt.Fprintf(out, "Usage of %s <global flags> <mode>:\n", os.Args[0])
	fmt.Fprintln(out, "  Global flags:")
	fmt.Fprintln(out, "    -version: Prints the version")
	fmt.Fprintln(out, "    -h, -help: Prints this help message")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "  Modes:")
	fmt.Fprintln(out, "    "+initSummary())
}

func handleError(mode string, e error, out io.Writer) int {
	code := 1

	// Handle -help, -h flags properly
	if e == flag.ErrHelp {
		code = 0
		e = &helpError{}
	} else if e != nil && e.Error() != "" {
		fmt.Fprintln(out, "Error:", e)
	}

	switch e.(type) {
	case *helpError:
		switch mode {
		case "init":
			initHelp(out)
		}
	}

	return code
}
