package xgo

import (
	"runtime"
	"strings"
)

func quote(env string) string {
	if before, after, ok := strings.Cut(env, "="); ok {
		// Shouldn't have spaces before equal, must be a value
		if strings.Contains(before, " ") {
			return "\"" + strings.ReplaceAll(env, "\"", "\\\"") + "\""
		}

		// No need to wrap booleans
		switch after {
		case "false", "true":
			return env
		}

		env = before + "=\"" + after + "\""
	}

	return env
}

func setupCC(goos string, goarch string) (string, string) {
	var cc string
	var cxx string

	if (goarch == runtime.GOARCH) && (goos == runtime.GOOS) {
		return "", ""
	}

	if _, ok := crossCC[runtime.GOOS]; !ok {
		return "", ""
	} else if _, ok := crossCC[runtime.GOOS][goos]; !ok {
		return "", ""
	} else if _, ok := crossCC[runtime.GOOS][goos][goarch]; !ok {
		return "", ""
	}

	cc = crossCC[runtime.GOOS][goos][goarch][0]
	cxx = crossCC[runtime.GOOS][goos][goarch][1]

	return cc, cxx
}

func setupZig(goos string, goarch string) (string, string) {
	var cc string = "zig cc --target="
	var cxx string = "zig c++ --target="
	var translate map[string]string = map[string]string{
		"386":     "x86",
		"amd64":   "x86_64",
		"arm64":   "aarch64",
		"darwin":  "macos",
		"linux":   "linux",
		"windows": "windows",
	}

	if (goarch == runtime.GOARCH) && (goos == runtime.GOOS) {
		return "", ""
	}

	if (goarch == "arm64") && (goos != "darwin") {
		return "", ""
	}

	if _, ok := translate[goarch]; ok {
		cc += translate[goarch]
		cxx += translate[goarch]
	} else {
		return "", ""
	}

	cc += "-"
	cxx += "-"

	if _, ok := translate[goos]; ok {
		cc += translate[goos]
		cxx += translate[goos]
	} else {
		return "", ""
	}

	return cc, cxx
}
