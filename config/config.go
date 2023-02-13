package config

import (
	"path"
	"runtime/debug"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"

	"github.com/astaclinic/astafx/logger"
)

func InitConfig(cfgFile string) {
	packageName := GetPackageName()
	logger.Infof("Loading config for package %s", packageName)

	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		logger.Infof("Loading config from %s", cfgFile)
	} else {
		viper.SetConfigName("config")

		viper.AddConfigPath(path.Join("/etc", packageName))
		for _, configDir := range xdg.ConfigDirs {
			viper.AddConfigPath(path.Join(configDir, packageName))
		}
		viper.AddConfigPath(path.Join(xdg.ConfigHome, packageName))
		viper.AddConfigPath("./config")
		logger.Infof("Searching config from default paths")
	}

	// support reading from environmental variables
	// all env variables are capitalized, dot (levels) and dashes are replaced with underscores
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	err := viper.ReadInConfig()

	if err != nil {
		logger.Warnf("Error in reading config. %v", err)
	}
}

func GetPackageName() string {
	buildInfo, ok := debug.ReadBuildInfo()
	packageName := "asta"
	if ok {
		packageName = path.Base(buildInfo.Path)
	} else {
		logger.Warnf("Fail to read package info")
	}
	return packageName
}
