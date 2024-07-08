package xgo

//go:generate goversioninfo --platform-specific

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

// Compiler is a struct containing relevant data for cross-compiling
// Go.
type Compiler struct {
	Debug bool
	Zig   bool
}

func (x *Compiler) debugRun(enviro []string, args []string) string {
	var relevant []string
	var tmp string

	for _, v := range enviro {
		if strings.HasPrefix(v, "CC=") ||
			strings.HasPrefix(v, "CGO_ENABLED=") ||
			strings.HasPrefix(v, "CXX=") ||
			strings.HasPrefix(v, "GOARCH=") ||
			strings.HasPrefix(v, "GOOS=") {

			switch runtime.GOOS {
			case "windows":
				relevant = append(relevant, "$env:"+quote(v)+";")
			default:
				relevant = append(relevant, quote(v)+" \\")
			}
		}
	}

	for i := range args {
		args[i] = quote(args[i])
	}

	tmp = fmt.Sprintf(
		"%s\ngo %s",
		strings.Join(relevant, "\n"),
		strings.Join(args, " "),
	)

	return tmp
}

func (x *Compiler) defaultEnv(
	goos string, goarch string, cgo string,
) (map[string]string, error) {
	var debug bool
	var e error
	var env map[string]string = map[string]string{}
	var stdout string
	var tmp map[string]string

	// Get default env
	for _, line := range os.Environ() {
		if k, v, ok := strings.Cut(line, "="); ok {
			v = strings.TrimPrefix(v, "'")
			v = strings.TrimSuffix(v, "'")
			env[k] = v
		}
	}

	// Enable CGO if cross-compiling
	env["CGO_ENABLED"] = cgo

	// Modify env for target GOOS/GOARCH
	env["GOARCH"] = goarch
	env["GOOS"] = goos

	debug = x.Debug
	x.Debug = false
	defer func() { x.Debug = debug }()

	// Get default Go env vars for target GOOS/GOARCH
	if stdout, e = x.Run(env, "env", "--json"); e != nil {
		return nil, e
	}

	if e = json.Unmarshal([]byte(stdout), &tmp); e != nil {
		return nil, e
	}

	// Get default env
	for k, v := range tmp {
		v = strings.TrimPrefix(v, "'")
		v = strings.TrimSuffix(v, "'")
		env[k] = v
	}

	return env, nil
}

// Run will run the go command.
func (x *Compiler) Run(
	env map[string]string, args ...string,
) (string, error) {
	var b []byte
	var cmd *exec.Cmd
	var e error
	var enviro []string
	var tmp string

	for k, v := range env {
		enviro = append(enviro, k+"="+v)
	}

	slices.Sort(enviro)

	if x.Debug {
		return x.debugRun(enviro, args), nil
	}

	cmd = exec.Command("go", args...)
	cmd.Env = enviro

	if b, e = cmd.Output(); e != nil {
		switch e := e.(type) {
		case *exec.ExitError:
			tmp = strings.TrimSuffix(string(e.Stderr), "\n")
			if tmp != "" {
				return "", fmt.Errorf(tmp)
			}
		default:
			return "", e
		}
	}

	return strings.TrimSuffix(string(b), "\n"), nil
}

// SetupEnv will set the following ENV vars:
// - CC
// - CGO_ENABLED
// - CXX
// - GOARCH
// - GOOS
func (x *Compiler) SetupEnv(
	goos string, goarch string,
) (map[string]string, error) {
	var cc string
	var cgo string = "0"
	var cxx string
	var e error
	var env map[string]string

	// Get configured cross-compiler
	if x.Zig {
		cc, cxx = setupZig(goos, goarch)
	} else {
		cc, cxx = setupCC(goos, goarch)
	}

	// Enable CGO if we hare compiling for host OS
	// Enable CGO if we have cross-compilers
	if (runtime.GOOS == goos) || ((cc != "") && (cxx != "")) {
		cgo = "1"
	}

	if env, e = x.defaultEnv(goos, goarch, cgo); e != nil {
		return nil, e
	}

	if (runtime.GOOS == goos) || ((cc != "") && (cxx != "")) {
		// Set cross-compilers in env
		env["CC"] = cc
		env["CGO_ENABLED"] = "1" // Redundant
		env["CXX"] = cxx
	}

	return env, nil
}
