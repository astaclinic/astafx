package info

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	BuildDate string
)

func GetInfo() (string, error) {
	var info strings.Builder
	fmt.Fprintf(&info, "       Platform: %v %v\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&info, "        Runtime: %v\n", runtime.Version())
	hostname, err := os.Hostname()
	if err != nil {
		hostname = fmt.Sprintf("Fail to get hostname: %v", err)
	}
	fmt.Fprintf(&info, "       Hostname: %v\n", hostname)
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("failed to read build info")
	}
	buildCommit := os.Getenv("BUILD_COMMIT")
	for _, buildSetting := range buildInfo.Settings {
		if buildSetting.Key == "vcs.revision" {
			buildCommit = buildSetting.Value
		}
	}
	fmt.Fprintf(&info, "   Build Commit: %v\n", buildCommit)
	fmt.Fprintf(&info, "     Build Date: %v", BuildDate)
	return info.String(), nil
}
