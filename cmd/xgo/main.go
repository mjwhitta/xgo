package main

//go:generate goversioninfo --platform-specific

import (
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/xgo"
)

func booleanLike(name string) bool {
	switch strings.ToLower(os.Getenv(name)) {
	case "", "0", "disable", "f", "false", "no", "off":
		return false
	}

	return true
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r)
			}

			switch r := r.(type) {
			case error:
				log.ErrX(Exception, r.Error())
			case string:
				log.ErrX(Exception, r)
			}
		}
	}()

	var args []string
	var e error
	var env map[string]string
	var keys []string
	var missing map[string][]string
	var stdout string
	var x *xgo.Compiler

	validate()

	if flags.check {
		missing = xgo.MissingToolchains()

		for target := range missing {
			keys = append(keys, target)
		}

		slices.Sort(keys)

		for _, target := range keys {
			log.Warnf(
				"%s missing %s",
				target,
				strings.Join(missing[target], " and "),
			)
		}

		return
	}

	// Get env vars or default to runtime
	if flags.goarch = os.Getenv("GOARCH"); flags.goarch == "" {
		flags.goarch = runtime.GOARCH
	}

	if flags.goos = os.Getenv("GOOS"); flags.goos == "" {
		flags.goos = runtime.GOOS
	}

	// Enable debug, if requested
	flags.debug = flags.debug || booleanLike("XGODEBUG")
	flags.garble = flags.garble || booleanLike("XGOGARBLE")
	x = &xgo.Compiler{
		Debug:  flags.debug,
		Garble: flags.garble,
		Zig:    booleanLike("XGOZIG"),
	}

	// Preprocess cli args for some special cases
	args = xgo.BuildArgsSanityCheck(cli.Args())

	// Get env for specified GOOS/GOARCH
	if env, e = x.SetupEnv(flags.goos, flags.goarch); e != nil {
		panic(e)
	}

	// Run Go command
	if stdout, e = x.Run(env, args...); e != nil {
		panic(e)
	}

	if stdout != "" {
		fmt.Println(stdout)
	}
}
