package config_test

import (
	"os"
	"testing"

	"github.com/spf13/viper"

	"github.com/astaclinic/astafx/config"
)

func TestInitConfig(t *testing.T) {
	t.Run("Test config file path", func(t *testing.T) {
		f, err := os.CreateTemp("", "astafx-")
		if err != nil {
			t.Errorf("failed to create temp file: %v", err)
		}
		f.Write([]byte(`
foo:
  bar: "123"
`))
		if err := f.Close(); err != nil {
			t.Errorf("failed to close temp file %s: %v", f.Name(), err)
		}
		config.InitConfig(f.Name())
		gotPath := viper.ConfigFileUsed()
		if gotPath != f.Name() {
			t.Errorf("wrong config file read, got %s, expected %s", gotPath, f.Name())
		}
		expectedVal := "123"
		gotVal := viper.Get("foo.bar")
		if gotVal != expectedVal {
			t.Errorf("unexpected config value, got %s, expected %s", gotVal, expectedVal)
		}
	})
	t.Run("Test env variable input", func(t *testing.T) {
		t.Setenv("FOO_BAR", "123")
		viper.SetDefault("foo.bar", "default_value")
		config.InitConfig("")
		expectedVal := "123"
		gotVal := viper.Get("foo.bar")
		if gotVal != expectedVal {
			t.Errorf("unexpected config value, got %s, expected %s", gotVal, expectedVal)
		}
	})
}
