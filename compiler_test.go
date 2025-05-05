package xgo_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/mjwhitta/xgo"
	assert "github.com/stretchr/testify/require"
)

type compileTest struct {
	os   string
	arch string
}

func (test compileTest) Bin(fn string, garble bool, zig bool) string {
	var tmp string = fmt.Sprintf(
		"%s.%s.%s",
		strings.TrimSuffix(fn, filepath.Ext(fn)),
		test.os,
		test.arch,
	)

	if garble {
		tmp += ".garble"
	}

	if zig {
		tmp += ".zig"
	}

	if test.os == "windows" {
		tmp += ".exe"
	}

	return tmp
}

var tests = map[string][]compileTest{
	"cgoSupported": {
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"linux", "386"},
		{"linux", "amd64"},
		{"windows", "386"},
		{"windows", "amd64"},
	},
	"cgoUnsupported": {
		{"aix", "ppc64"},
		{"android", "386"},
		{"android", "amd64"},
		{"android", "arm"},
		{"android", "arm64"},
		{"dragonfly", "amd64"},
		{"freebsd", "386"},
		{"freebsd", "amd64"},
		{"freebsd", "arm"},
		{"freebsd", "arm64"},
		{"freebsd", "riscv64"},
		{"illumos", "amd64"},
		{"ios", "amd64"},
		{"ios", "arm64"},
		{"js", "wasm"},
		{"linux", "arm"},
		{"linux", "arm64"},
		{"linux", "loong64"},
		{"linux", "mips"},
		{"linux", "mips64"},
		{"linux", "mips64le"},
		{"linux", "mipsle"},
		{"linux", "ppc64"},
		{"linux", "ppc64le"},
		{"linux", "riscv64"},
		{"linux", "s390x"},
		{"netbsd", "386"},
		{"netbsd", "amd64"},
		{"netbsd", "arm"},
		{"netbsd", "arm64"},
		{"openbsd", "386"},
		{"openbsd", "amd64"},
		{"openbsd", "arm"},
		{"openbsd", "arm64"},
		{"openbsd", "ppc64"},
		{"plan9", "386"},
		{"plan9", "amd64"},
		{"plan9", "arm"},
		{"solaris", "amd64"},
		{"wasip1", "wasm"},
		{"windows", "arm"},
		{"windows", "arm64"},
	},
	"garbleUnsupported": {
		{"wasip1", "wasm"},
	},
	"supported": {
		{"aix", "ppc64"},
		{"android", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"dragonfly", "amd64"},
		{"freebsd", "386"},
		{"freebsd", "amd64"},
		{"freebsd", "arm"},
		{"freebsd", "arm64"},
		{"freebsd", "riscv64"},
		{"illumos", "amd64"},
		{"js", "wasm"},
		{"linux", "386"},
		{"linux", "amd64"},
		{"linux", "arm"},
		{"linux", "arm64"},
		{"linux", "loong64"},
		{"linux", "mips"},
		{"linux", "mips64"},
		{"linux", "mips64le"},
		{"linux", "mipsle"},
		{"linux", "ppc64"},
		{"linux", "ppc64le"},
		{"linux", "riscv64"},
		{"linux", "s390x"},
		{"netbsd", "386"},
		{"netbsd", "amd64"},
		{"netbsd", "arm"},
		{"netbsd", "arm64"},
		{"openbsd", "386"},
		{"openbsd", "amd64"},
		{"openbsd", "arm"},
		{"openbsd", "arm64"},
		{"openbsd", "ppc64"},
		{"plan9", "386"},
		{"plan9", "amd64"},
		{"plan9", "arm"},
		{"solaris", "amd64"},
		{"wasip1", "wasm"},
		{"windows", "386"},
		{"windows", "amd64"},
		{"windows", "arm"},
		{"windows", "arm64"},
	},
	"unsupported": {
		{"android", "386"},
		{"android", "amd64"},
		{"android", "arm"},
		{"ios", "amd64"},
		{"ios", "arm64"},
	},
}

func build(
	t *testing.T,
	test compileTest,
	file string,
	garble bool,
	zig bool,
	pass bool,
	canSkip bool,
) {
	var e error
	var env map[string]string
	var x *xgo.Compiler = &xgo.Compiler{Garble: garble, Zig: zig}

	// XGo entry
	env, e = x.SetupEnv(test.os, test.arch)
	assert.Nil(t, e)
	assert.NotNil(t, env)

	if canSkip && (env["CC"] != "") {
		_, e = exec.LookPath(strings.Fields(env["CC"])[0])
		if e != nil {
			t.Skipf(
				"%s is not installed",
				strings.Fields(env["CC"])[0],
			)
		}
	}

	t.Cleanup(
		func() {
			os.Remove(
				filepath.Join(
					"testdata",
					test.Bin(file, garble, zig),
				),
			)
		},
	)

	// Compile
	_, e = x.Run(
		env,
		"build",
		"-o",
		filepath.Join("testdata", test.Bin(file, garble, zig)),
		filepath.Join("testdata", file),
	)
	if pass {
		assert.Nil(t, e)
	} else {
		assert.NotNil(t, e)
	}
}

func TestCompileCGOSupported(t *testing.T) {
	var src string = "main_cgo.go"

	t.Parallel()

	for _, test := range tests["cgoSupported"] {
		t.Run(
			"Target("+test.os+"/"+test.arch+")",
			func(t *testing.T) {
				build(t, test, src, false, false, true, true)
			},
		)
	}
}

func TestCompileCGOUnsupported(t *testing.T) {
	var src string = "main_cgo.go"

	t.Parallel()

	for _, test := range tests["cgoUnsupported"] {
		t.Run(
			"Target("+test.os+"/"+test.arch+")",
			func(t *testing.T) {
				build(t, test, src, false, false, false, false)
			},
		)
	}
}

func TestCompileCGOZig(t *testing.T) {
	var src string = "main_cgo.go"

	t.Parallel()

	for _, test := range tests["cgoSupported"] {
		t.Run(
			"Target("+test.os+"/"+test.arch+")",
			func(t *testing.T) {
				build(t, test, src, false, true, true, true)
			},
		)
	}
}

func TestCompileSupported(t *testing.T) {
	var src string = "main.go"

	t.Parallel()

	for _, test := range tests["supported"] {
		t.Run(
			"Target("+test.os+"/"+test.arch+")",
			func(t *testing.T) {
				build(t, test, src, false, false, true, false)
			},
		)
	}
}

func TestCompileSupportedWithGarble(t *testing.T) {
	var src string = "main.go"

	t.Parallel()

	for _, test := range tests["supported"] {
		if slices.Contains(tests["garbleUnsupported"], test) {
			continue
		}

		t.Run(
			"GarbleTarget("+test.os+"/"+test.arch+")",
			func(t *testing.T) {
				build(t, test, src, true, false, true, false)
			},
		)
	}
}

func TestCompileUnsupported(t *testing.T) {
	var src string = "main.go"

	t.Parallel()

	for _, test := range tests["unsupported"] {
		t.Run(
			"Target("+test.os+"/"+test.arch+")",
			func(t *testing.T) {
				build(t, test, src, false, false, false, false)
			},
		)
	}
}

func TestDebug(t *testing.T) {
	var e error
	var env map[string]string
	var stdout string
	var x *xgo.Compiler = &xgo.Compiler{Debug: true}

	t.Parallel()

	env, e = x.SetupEnv(runtime.GOOS, runtime.GOARCH)
	assert.Nil(t, e)
	assert.NotNil(t, env)

	stdout, e = x.Run(env, "vet", ".")
	assert.Nil(t, e)
	assert.NotEqual(t, "", stdout)
}
