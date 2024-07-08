package xgo_test

import (
	"testing"

	"github.com/mjwhitta/xgo"
	assert "github.com/stretchr/testify/require"
)

type buildArgsTest struct {
	name string
	in   []string
	out  []string
}

func TestBuildArgsSanityCheck(t *testing.T) {
	var bld string = "--ldflags=-s -w"
	var btrim string = "--trimpath"
	var bvcs string = "--buildvcs=false"
	var tests []buildArgsTest = []buildArgsTest{
		{"Nothing", nil, nil},
		{"Wrong command", []string{"vet", "."}, []string{"vet", "."}},
		{
			"Missing all",
			[]string{"build"},
			[]string{"build", bvcs, bld, btrim},
		},
		{
			"Missing buildvcs",
			[]string{"build", bld, btrim},
			[]string{"build", bvcs, bld, btrim},
		},
		{
			"Missing ldflags",
			[]string{"build", bvcs, btrim},
			[]string{"build", bld, bvcs, btrim},
		},
		{
			"Missing trimpath",
			[]string{"build", bvcs, bld},
			[]string{"build", btrim, bvcs, bld},
		},
		{
			"Existing",
			[]string{"build", "--buildvcs=a", "--ldflags=-s", btrim},
			[]string{"build", "--buildvcs=a", "--ldflags=-s", btrim},
		},
	}

	t.Parallel()

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var args []string = xgo.BuildArgsSanityCheck(test.in)
				assert.Equal(t, test.out, args)
			},
		)
	}
}

func TestMissingToolchains(t *testing.T) {
	t.Parallel()

	missing := xgo.MissingToolchains()
	assert.NotNil(t, missing)
}
