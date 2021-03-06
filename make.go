// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	commandName     = "eribo"
	commandLocation = "./cmd/eribo"
	deployFolder    = "/home/kusubooru/code/go/bin"
)

type platform struct {
	os   string
	arch string
}

type binary struct {
	name    string
	version string
	targets []platform
}

func (bin binary) Name(os, arch string) string {
	s := fmt.Sprintf("%s_%s-%s_%s", bin.name, os, arch, bin.version)
	if os == "windows" {
		s = s + ".exe"
	}
	return s
}

func (bin binary) Names() []string {
	names := make([]string, len(bin.targets))
	for i, t := range bin.targets {
		names[i] = bin.Name(t.os, t.arch)
	}
	return names
}

var (
	release   = flag.Bool("release", false, "Build binaries for all target platforms.")
	clean     = flag.Bool("clean", false, "Remove all created binaries from current directory.")
	deploy    = flag.Bool("deploy", false, "Deploy binary to server.")
	buildARCH = flag.String("arch", runtime.GOARCH, "Architecture to build for.")
	buildOS   = flag.String("os", runtime.GOOS, "Operating system to build for.")
	password  = flag.String("pass", "", "Password of the machine used for deployment.")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: go run make.go [OPTIONS]\n\n")
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	bin := binary{
		name: commandName,
		targets: []platform{
			{os: "linux", arch: "amd64"},
			{os: runtime.GOOS, arch: runtime.GOARCH},
		},
	}
	bin.version = getVersion()

	if *release {
		start := time.Now()
		forEachBinary(bin, buildBinary)
		fmt.Println("Time elapsed:", time.Since(start))
		os.Exit(0)
	}

	if *clean {
		forEachBinary(bin, rmBinary)
		os.Exit(0)
	}

	if *deploy {
		deployBin(bin, *buildOS, *buildARCH)
		os.Exit(0)
	}

	buildBinary(bin, *buildOS, *buildARCH)
}

func deployBin(bin binary, OS, arch string) {
	OS = "linux"
	arch = "amd64"
	buildBinary(bin, OS, arch)
	name := bin.Name(OS, arch)
	defer func() {
		if err := os.Remove(name); err != nil {
			fmt.Fprintf(os.Stderr, "cleaning up %s: %v", name, err)
		}
	}()

	server := "kusubooru@kusubooru.com"
	fmt.Println("Deploying to server", server)
	args := []string{name, server + ":" + deployFolder}
	cmdName := "scp"
	if runtime.GOOS == "windows" {
		cmdName = "pscp"
		args = append([]string{"-scp", "-pw", *password}, args...)
	}
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error deploying to server: %v", err)
	}
}

func getVersion() string {
	cmd := exec.Command("git", "describe", "--tags")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git describe: %v", err)
	}
	return strings.TrimPrefix(strings.TrimSpace(string(out)), "v")
}

type binaryFunc func(bin binary, OS, arch string)

func forEachBinary(bin binary, fn binaryFunc) {
	var wg sync.WaitGroup
	for _, t := range bin.targets {
		wg.Add(1)
		go func(bin binary, os, arch string) {
			defer wg.Done()
			fn(bin, os, arch)
		}(bin, t.os, t.arch)
	}
	wg.Wait()
}

func buildBinary(bin binary, OS, arch string) {
	ldflags := fmt.Sprintf("--ldflags=-s -w -X main.theVersion=%s", bin.version)
	tags := fmt.Sprint("--tags=prod")
	cmd := exec.Command("go", "build", tags, ldflags, "-o", bin.Name(OS, arch), commandLocation)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = copyGoEnv()
	cmd.Env = setEnv(cmd.Env, "GOOS", OS)
	cmd.Env = setEnv(cmd.Env, "GOARCH", arch)
	fmt.Println("Building binary:", bin.Name(OS, arch))
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running go build:", err)
	}
}

func rmBinary(bin binary, OS, arch string) {
	err := os.Remove(bin.Name(OS, arch))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Error removing binary:", err)
		}
	}
}

func copyGoEnv() (environ []string) {
	for _, env := range os.Environ() {
		environ = append(environ, env)
	}
	return
}

func setEnv(env []string, key, value string) []string {
	for i, e := range env {
		if strings.HasPrefix(e, fmt.Sprintf("%s=", key)) {
			env[i] = fmt.Sprintf("%s=%s", key, value)
			return env
		}
	}
	env = append(env, fmt.Sprintf("%s=%s", key, value))
	return env
}
