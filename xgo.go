package xgo

import (
	"os/exec"
	"runtime"
	"strings"
)

// BuildArgsSanityCheck will ensure the build args include some sane
// defaults. It will not alter existing args.
func BuildArgsSanityCheck(args []string) []string {
	var add []string
	var found bool
	var opts [][]string = [][]string{
		{"--buildvcs", "false"},
		{"--ldflags", "-s -w"},
		{"--trimpath", ""},
	}

	if len(args) == 0 {
		return nil
	}

	// Only for compilation commands
	switch args[0] {
	case "build", "get", "install":
	default:
		return args
	}

	add = []string{args[0]}

	for _, o := range opts {
		found = false
		for _, arg := range args {
			if found = strings.HasPrefix(arg, o[0]); found {
				break
			}
		}

		if found {
			continue
		}

		add = append(add, strings.TrimSuffix(o[0]+"="+o[1], "="))
	}

	return append(add, args[1:]...)
}

// MissingToolchains returns a list of toolchains that are not
// installed.
func MissingToolchains() map[string][]string {
	var e error
	var missing map[string][]string = map[string][]string{}
	var tmp []string

	for goos, target := range crossCC[runtime.GOOS] {
		for goarch, cccxx := range target {
			tmp = []string{}

			_, e = exec.LookPath(strings.Fields(cccxx[0])[0])
			if e != nil {
				tmp = append(tmp, cccxx[0])
			}

			_, e = exec.LookPath(strings.Fields(cccxx[1])[0])
			if e != nil {
				tmp = append(tmp, cccxx[1])
			}

			if len(tmp) > 0 {
				missing[goos+"/"+goarch] = tmp
			}
		}
	}

	if _, e = exec.LookPath("zig"); e != nil {
		missing["all targets"] = []string{"zig"}
	}

	return missing
}
