package config

import (
	"fmt"
	"kconsole/utils/errorx"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	configname            = "config.yaml"
	defaultRelativeConfig = ".kconsole"
	once                  sync.Once
	kconsoleConfig        *KconsoleConfig = &KconsoleConfig{}
	// auth configuration enums
	LocalConfigAuth string = "local"
	BcsAuth         string = "bcs"
)

type KconsoleConfig struct {
	Auth     string `json:"auth" default:"local"`
	BCSHost  string `json:"bcsHost" default:""`
	BCSToken string `json:"bcsToken" default:""`
}

func (c *KconsoleConfig) validate() {
	if c.Auth != LocalConfigAuth && c.Auth != BcsAuth {
		errorx.CheckErrorWithCode(fmt.Errorf("auth:%s is not vaild, must be the `bcs` or `local`", c.Auth), errorx.ErrorAuthConfigErr)
	}
	if c.Auth == BcsAuth {
		if c.BCSHost == "" || c.BCSToken == "" {
			errorx.CheckErrorWithCode(fmt.Errorf("when auth=%s, must be set `bcsHost` and `bcsToken` option in your config", c.Auth), errorx.ErrorBCSAuthConfigErr)
		}
	}
}

func InitConfig() {
	viper.SetConfigName(configname)
	home, err := os.UserHomeDir()
	errorx.CheckError(err)
	viper.AddConfigPath(filepath.Join(home, defaultRelativeConfig))
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	errorx.CheckErrorWithCode(err, errorx.ErrorConfigErr)
	err = viper.Unmarshal(kconsoleConfig)
	errorx.CheckErrorWithCode(err, errorx.ErrorConfigErr)
	// validate the configuration
	kconsoleConfig.validate()
}

// GetKconsoleConfig get kconsole config instance
func GetKconsoleConfig() *KconsoleConfig {
	once.Do(InitConfig)
	return kconsoleConfig
}
