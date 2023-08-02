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
	Auth       string `json:"auth" default:"local"`
	BCSHost    string `json:"bcshost" default:""`
	BCSToken   string `json:"bcstoken" default:""`
	BCSCluster string `json:"bcscluster" default:""`
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

// setDefaultConfig2File 设置默认配置
func setDefaultConfig2File(v *viper.Viper) {
	v.Set("auth", "local")
	if err := v.WriteConfig(); err != nil {
		errorx.CheckError(fmt.Errorf("fatal error writing config file: %v \n", err))
	}
}

// getConfigDir 获取配置文件所在目录
func getConfigDir() string {
	home, err := os.UserHomeDir()
	errorx.CheckError(err)
	return filepath.Join(home, defaultRelativeConfig)
}

// getConfigpath 获取配置文件完整路径
func getConfigpath() string {
	return filepath.Join(getConfigDir(), configname)
}

// setViper 生成viper对象
func getViper() viper.Viper {
	viper.SetConfigName(configname)
	viper.AddConfigPath(getConfigDir())
	viper.SetConfigType("yaml")
	return *viper.GetViper()
}

// mkCongfigfile 创建配置文件并赋予默认值
func mkCongfigfile(v *viper.Viper) {
	f, err := os.Create(getConfigpath())
	if _, ok := err.(*os.PathError); ok {
		errorx.CheckError(os.Mkdir(getConfigDir(), 0755))
		f, err := os.Create(getConfigpath())
		errorx.CheckError(err)
		f.Close()
	} else {
		errorx.CheckError(err)
	}
	f.Close()
	setDefaultConfig2File(v)
}

func InitConfig() {
	v := getViper()
	InitConfigWithViper(&v)
}

func InitConfigWithViper(v *viper.Viper) {
	err := v.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		mkCongfigfile(v)
		// insert default config
	} else {
		errorx.CheckErrorWithCode(err, errorx.ErrorConfigErr)
	}
	err = v.Unmarshal(&kconsoleConfig)
	errorx.CheckErrorWithCode(err, errorx.ErrorConfigErr)
	// validate the configuration
	kconsoleConfig.validate()
}

// GetKconsoleConfig get kconsole config instance
func GetKconsoleConfig() *KconsoleConfig {
	once.Do(InitConfig)
	return kconsoleConfig
}

// UpdateConfilefile set new struct to config file
func UpdateConfilefile(config map[string]string) {
	v := getViper()
	InitConfigWithViper(&v)
	for key, val := range config {
		v.Set(key, val)
	}
	if err := v.WriteConfig(); err != nil {
		errorx.CheckError(fmt.Errorf("fatal error writing config file: %v \n", err))
	}
}
