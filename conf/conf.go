package conf

import (
	"bytes"
	_ "embed"
	"strings"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"gopkg.in/validator.v2"
)

var (
	//go:embed conf.yaml
	configFile []byte
	conf       *Config
	once       sync.Once
)

type Config struct {
	Hertz struct {
		Address         string `mapstructure:"address"`
		Service         string `mapstructure:"service"`
		EnablePprof     bool   `mapstructure:"enable_pprof"`
		EnableGzip      bool   `mapstructure:"enable_gzip"`
		EnableAccessLog bool   `mapstructure:"enable_access_log"`
		LogLevel        string `mapstructure:"log_level"`
		OtlpAddr        string `mapstructure:"otlp_address"`
		EnableCaptcha   bool   `mapstructure:"enable_captcha"`
	} `mapstructure:"hertz"`

	Registry struct {
		RegistryAddress []string `mapstructure:"registry_address"`
		Username        string   `mapstructure:"username"`
		Password        string   `mapstructure:"password"`
	} `mapstructure:"registry"`

	MySQL struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	} `mapstructure:"mysql"`

	Redis struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
		Username string `mapstructure:"username"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	OSS struct {
		AccessKeyID         string `mapstructure:"access_key_id"`
		AccessKeySecret     string `mapstructure:"access_key_secret"`
		RoleArnForOssUpload string `mapstructure:"role_arn_for_oss_upload"`
		Bucket              string `mapstructure:"bucket"`
		RegionID            string `mapstructure:"region_id"`
	} `mapstructure:"oss"`
	WX struct {
		AppID     string `mapstructure:"app_id"`
		AppSecret string `mapstructure:"app_secret"`
	} `mapstructure:"wx"`
	JWT struct {
		Secret    string `mapstructure:"secret"`
		ValidDays int    `mapstructure:"valid_days"`
	} `mapstructure:"jwt"`
	LOG struct {
		Path       string `mapstructure:"path"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxAge     int    `mapstructure:"max_age"`
		Compress   bool   `mapstructure:"compress"`
	}
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		hlog.Warn("Error loading .env file")
	}

	conf = new(Config)
	viper.SetConfigType("yaml")
	// 参数打包入可执行文件,可用.env覆盖
	err = viper.ReadConfig(bytes.NewBuffer(configFile))
	if err != nil {
		panic(err)
	}

	// Enable automatic environment variable reading
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(conf)
	if err != nil {
		panic(err)
	}

	if err := validator.Validate(conf); err != nil {
		hlog.Error("validate config error - %v", err)
		panic(err)
	}

	pretty.Printf("%+v\n", conf)
}

func LogLevel() hlog.Level {
	level := GetConf().Hertz.LogLevel
	switch level {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "info":
		return hlog.LevelInfo
	case "notice":
		return hlog.LevelNotice
	case "warn":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}
