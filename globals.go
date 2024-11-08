package xgo

// crossCC is a mapping of GOHOSTOS/GOOS/GOARCH to CC and CXX.
var crossCC = map[string]map[string]map[string][]string{
	"darwin": {
		"linux": {
			// brew install musl-cross --with-i468 --without-aarch64
			"386": {
				"i486-linux-musl-gcc --static",
				"i486-linux-musl-g++ --static",
			},
			"amd64": {
				"x86_64-linux-musl-gcc --static",
				"x86_64-linux-musl-g++ --static",
			},
		},
		"windows": {
			// brew install mingw-w64
			"386": {"i686-w64-mingw32-gcc", "i686-w64-mingw32-g++"},
			"amd64": {
				"x86_64-w64-mingw32-gcc",
				"x86_64-w64-mingw32-g++",
			},
		},
	},
	"linux": {
		"darwin": {
			// https://github.com/tpoechtrager/osxcross
			"amd64": {"o64-clang", "o64-clang++"},
			"arm64": {"oa64-clang", "oa64-clang++"},
		},
		"windows": {
			// mingw-w64
			"386": {"i686-w64-mingw32-gcc", "i686-w64-mingw32-g++"},
			"amd64": {
				"x86_64-w64-mingw32-gcc",
				"x86_64-w64-mingw32-g++",
			},
		},
	},
	"windows": {
		"darwin": {
			// How to install?
			"amd64": {"o64-clang", "o64-clang++"},
			"arm64": {"oa64-clang", "oa64-clang++"},
		},
		"linux": {
			// choco install mingw
			"386": {"i686-w64-mingw32-gcc", "i686-w64-mingw32-g++"},
			"amd64": {
				"x86_64-w64-mingw32-gcc",
				"x86_64-w64-mingw32-g++",
			},
		},
	},
}

// Version is the package version.
const Version = "0.2.0"
