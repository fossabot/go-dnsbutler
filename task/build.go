package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var flags = struct {
	Name        string
	UseNetgo    bool
	SrcFilePath string
	DistDirPath string
	UpxEnabled  bool
	UpxLevel    string // -7, --best, --brute
	ShowDebug   bool
}{
	"dnsbutler",
	true,
	"./cmd/dnsbutler/main.go",
	"./dist",
	false,
	"-7",
	true,
}

var upxBinary = struct {
	Available bool
	Version   string
	Enabled   bool
}{
	false, "", false,
}

var buildTargets = []buildConfig{
	{"linux", "amd64", ""},
	{"linux", "arm", "6"},  // (Raspberry Pi A, A+, B, B+, Zero)
	{"linux", "arm", "7"},  // (Raspberry Pi 2, 3)
	{"linux", "arm64", ""}, // GOARM is not available!
	{"windows", "amd64", ""},
}

// All available GOOS and GOARCH are listed in go src/go/build/syslist.go
type buildConfig struct {
	GOOS   string
	GOARCH string
	GOARM  string
}

func (b buildConfig) ArchFileString() string {
	arch := b.GOARCH
	if b.GOARCH == "arm" {
		arch += b.GOARM
	}
	return arch
}

func (b buildConfig) CreateDistFilePath(name, distPath string) string {
	path := fmt.Sprintf(distPath+"/%s-%s-%s", name, b.ArchFileString(), b.GOOS)
	if b.GOOS == "windows" {
		path += ".exe"
	}

	return path
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func deleteAndCreateDistDir(path string) {
	err := os.RemoveAll(path)
	must(err)

	err = os.MkdirAll(path, 0755)
	must(err)
}

func execBuild(name, srcFile, distPath string, c buildConfig) error {
	args := []string{
		"build",
		"-a",
	}

	if flags.UseNetgo {
		args = append(args, "-tags", "'netgo'")
	}

	args = append(args, "-ldflags",
		"-s -w",
		"-o",
		c.CreateDistFilePath(name, distPath),
		srcFile,
	)

	cmd := exec.Command("go", args...)

	cmd.Env = append(
		os.Environ(),
		"CGO_ENABLED=0",
		fmt.Sprintf("GOOS=%s", c.GOOS),
		fmt.Sprintf("GOARCH=%s", c.GOARCH),
	)

	if c.GOARM != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOARM=%s", c.GOARM))
	}

	if flags.ShowDebug {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}

func upx(filePath string) error {
	basename := filepath.Base(filePath)
	ext := filepath.Ext(basename)
	name := strings.TrimSuffix(basename, ext) + "-upx" + ext

	target := filepath.Clean(filepath.Dir(filePath) + "/" + name)

	cmd := exec.Command("upx", []string{
		"--no-color",
		"--no-progress",
		flags.UpxLevel,
		"-o",
		target,
		filePath,
	}...)

	return cmd.Run()
}

func path(path string) string {
	tmp, err := filepath.Abs(path)
	must(err)
	return filepath.Clean(tmp)
}

func name(nameFlag, srcFile string) string {
	if nameFlag != "" {
		return nameFlag
	}
	basename := filepath.Base(srcFile)
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

func build(name, srcFile, distPath string, buildTargets []buildConfig) {
	buildStarted := time.Now()

	for _, c := range buildTargets {
		now := time.Now()
		fmt.Printf("Start building for arch=%s os=%s...\n", c.ArchFileString(), c.GOOS)
		err := execBuild(name, srcFile, distPath, c)

		dur := time.Since(now).Round(time.Second)
		if err != nil {
			fmt.Printf("build failed after %s with '%s'", dur, err)
			panic(err)
		} else {
			fmt.Printf("Done after %s\n", dur)
		}

		if upxBinary.Enabled && upxBinary.Available {
			fmt.Println("Start upx...")
			now = time.Now()

			err := upx(c.CreateDistFilePath(name, distPath))
			must(err)

			fmt.Printf("Done after %s\n", time.Since(now).Round(time.Second))
		}
	}

	fmt.Printf("All builds done after %s", time.Since(buildStarted).Round(time.Second))
}

func detectUpx() {
	upxBinary.Enabled = flags.UpxEnabled

	outBuf := new(bytes.Buffer)

	cmd := exec.Command("upx", []string{
		"--no-color",
		"--no-progress",
		"--version",
	}...)
	cmd.Stdout = outBuf

	err := cmd.Run()
	if err != nil {
		upxBinary.Available = false
		return
	}

	reg := regexp.MustCompile(".*upx\\s(.*?)\\s.*")
	res := reg.FindStringSubmatch(outBuf.String())

	if len(res) >= 1 {
		upxBinary.Version = res[1]
	} else {
		upxBinary.Version = "unknown"
	}
	upxBinary.Available = true
}

func main() {
	srcFile := path(flags.SrcFilePath)
	distPath := path(flags.DistDirPath)
	name := name(flags.Name, srcFile)

	detectUpx()
	if upxBinary.Available {
		fmt.Printf("Detected upx %s\n", upxBinary.Version)
	} else if upxBinary.Enabled {
		panic("Using upx enabled but no upx found in PATH")
	} else {
		fmt.Println("No upx in PATH detected")
	}

	fmt.Printf("Creating dist dir '%s'\n", distPath)
	deleteAndCreateDistDir(distPath)

	fmt.Println("Start the build process")
	build(name, srcFile, distPath, buildTargets)

	os.Exit(0)
}
