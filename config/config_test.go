package config

import (
	"fmt"
	"testing"
)

func TestLoadConf(t *testing.T) {
	LoadConf("../conf/config.yml")
	if GlobalConfig == nil {
		t.Error("GlobalConfig is nil")
	}

	// Print all configurations
	allSet := GlobalConfig.AllSettings()
	for key, value := range allSet {
		fmt.Printf("%q: %q\n", key, value)
	}
}
