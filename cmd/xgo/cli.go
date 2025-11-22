package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mjwhitta/cli"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/xgo"
)

// Exit status
const (
	Good = iota
	InvalidOption
	MissingOption
	InvalidArgument
	MissingArgument
	ExtraArgument
	Exception
)

// Flags
var flags struct {
	check   bool
	debug   bool
	garble  bool
	goarch  string
	goos    string
	nocolor bool
	verbose bool
	version bool
}

func init() {
	// Configure cli package
	cli.Align = true // Defaults to false
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = "" +
		filepath.Base(os.Args[0]) + " [OPTIONS] <gocommand> [args]"
	cli.BugEmail = "xgo.bugs@whitta.dev"

	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		fmt.Sprintf("  %d: Invalid option\n", InvalidOption),
		fmt.Sprintf("  %d: Missing option\n", MissingOption),
		fmt.Sprintf("  %d: Invalid argument\n", InvalidArgument),
		fmt.Sprintf("  %d: Missing argument\n", MissingArgument),
		fmt.Sprintf("  %d: Extra argument\n", ExtraArgument),
		fmt.Sprintf("  %d: Exception", Exception),
	)
	cli.Info(
		"This tool aims to simplify cross-compiling Go with or",
		"without CGO support.",
	)

	cli.SeeAlso = []string{"gcc", "go", "mingw", "osxcross-git"}
	cli.Title = "XGo"

	// Parse cli flags
	cli.Flag(
		&flags.check,
		"c",
		"check",
		false,
		"Check for missing toolchains.",
	)
	cli.Flag(&flags.debug, "d", "debug", false, "n/a", true)
	cli.Flag(&flags.garble, "g", "garble", false, "n/a", true)
	cli.Flag(
		&flags.goarch,
		"goarch",
		runtime.GOARCH,
		"Set the GOARCH env var (useful for Windows).",
	)
	cli.Flag(
		&flags.goos,
		"goos",
		runtime.GOOS,
		"Set the GOOS env var (useful for Windows).",
	)
	cli.Flag(
		&flags.nocolor,
		"no-color",
		false,
		"Disable colorized output.",
	)
	cli.Flag(
		&flags.verbose,
		"v",
		"verbose",
		false,
		"Show stacktrace, if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit, if version was requested
	if flags.version {
		fmt.Println(
			filepath.Base(os.Args[0]) + " version " + xgo.Version,
		)
		os.Exit(Good)
	}

	// Validate cli flags
	if flags.check {
		if cli.NArg() > 0 {
			cli.Usage(ExtraArgument)
		}
	} else if cli.NArg() == 0 {
		cli.Usage(MissingArgument)
	}
}
