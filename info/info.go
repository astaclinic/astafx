package info

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

var (
	BuildDate   string
	ProgramName string
)

type InfoDisplay struct {
	Name        string
	Platform    string
	Runtime     string
	HostName    string
	BuildCommit string
	BuildDate   string
}

func GetInfo() (*InfoDisplay, error) {
	display := &InfoDisplay{
		Name:     fmt.Sprintf("        Program: %v\n", ProgramName),
		Platform: fmt.Sprintf("       Platform: %v %v\n", runtime.GOOS, runtime.GOARCH),
		Runtime:  fmt.Sprintf("        Runtime: %v\n", runtime.Version()),
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = fmt.Sprintf("Fail to get hostname: %v", err)
	}
	display.HostName = fmt.Sprintf("       Hostname: %v\n", hostname)
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, errors.New("failed to read build info")
	}
	buildCommit := os.Getenv("BUILD_COMMIT")
	for _, buildSetting := range buildInfo.Settings {
		if buildSetting.Key == "vcs.revision" {
			buildCommit = buildSetting.Value
		}
	}
	display.BuildCommit = fmt.Sprintf("   Build Commit: %v\n", buildCommit)
	display.BuildDate = fmt.Sprintf("     Build Date: %v", BuildDate)
	return display, nil
}
