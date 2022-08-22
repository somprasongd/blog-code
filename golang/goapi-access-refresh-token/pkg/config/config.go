package config

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	App     appConfig     `validate:"dive"`
	Db      dbConfig      `validate:"dive"`
	Redis   redisConfig   `validate:"dive"`
	Server  serverConfig  `validate:"dive"`
	Gateway gatewayConfig `validate:"dive"`
	Token   tokenConfig   `validate:"dive"`
}

type appConfig struct {
	BaseUrl      string `env:"APP_BASEURL" validate:"required"`
	Mode         string `env:"APP_MODE" validate:"required,oneof=development test production"`
	Name         string `env:"APP_NAME"`
	Version      string `env:"APP_VERSION"`
	LivenessFile string `env:"APP_LIVENESSFILE"`
}

func (c *appConfig) IsDevMode() bool {
	return c.Mode == "development"
}

func (c *appConfig) IsTestMode() bool {
	return c.Mode == "test"
}

func (c *appConfig) IsProdMode() bool {
	return c.Mode == "production"
}

type dbConfig struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	Port     uint   `env:"DB_PORT"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Database string `env:"DB_DATABASE"`
	Sslmode  string `env:"DB_SSLMODE"`
}

type redisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     uint   `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	Database int    `env:"REDIS_DATABASE"`
}

type serverConfig struct {
	Port         uint          `env:"SERVER_PORT"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE"`
}

type gatewayConfig struct {
	Host    string `env:"GATEWAY_HOST"`
	BaseURL string `env:"GATEWAY_BASEURL"`
}

type tokenConfig struct {
	AccessSecretKey  string        `env:"TOKEN_ACCESS_SECRET" validate:"required"`
	AccessExpires    time.Duration `env:"TOKEN_ACCESS_EXPIRES" validate:"required"`
	RefreshSecretKey string        `env:"TOKEN_REFRESH_SECRET" validate:"required"`
	RefreshExpires   time.Duration `env:"TOKEN_REFRESH_EXPIRES" validate:"required"`
}

func LoadConfig() *Config {
	setDefault()
	viper.AutomaticEnv()                                   // ให้อ่านค่าจาก env มา replace ในไฟล์ config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // แปลง _ underscore ใน env เป็น . dot notation ใน viper

	viper.SetConfigName("config") // กำหนดชื่อไฟล์ config (without extension)
	viper.SetConfigType("yaml")   // ระบุประเภทของไฟล์ config
	viper.AddConfigPath(".")      // ระบุตำแหน่งของไฟล์ config อยู่ที่ working directory
	// อ่านไฟล์ config
	err := viper.ReadInConfig()
	if err != nil { // ถ้าอ่านไฟล์ config ไม่ได้ให้ข้ามไปเพราะให้เอาค่าจาก env มาแทนได้
		log.Println("please consider environment variables", err.Error())
	}
	config := setConfig()

	// ตรวจสอบว่ากำหนดค่ามาครบหรือไม่
	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		panic(errors.New("load config error: " + err.Error()))
	}
	return config
}

func setDefault() {
	// กำหนด Default Value
	viper.SetDefault("app.baseurl", "/api/v1")
	viper.SetDefault("app.mode", "development")
	viper.SetDefault("app.version", "1.0")
	viper.SetDefault("app.livenessfile", "./tmp-live")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.timeout.read", "15s")
	viper.SetDefault("server.timteout.write", "15s")
	viper.SetDefault("server.timeout.idle", "60s")

	viper.SetDefault("db.sslmode", "disable")

	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.database", 0)
}

func setConfig() *Config {
	return &Config{
		App: appConfig{
			BaseUrl:      viper.GetString("app.baseurl"),
			Mode:         viper.GetString("app.mode"),
			Name:         viper.GetString("app.name"),
			Version:      viper.GetString("app.version"),
			LivenessFile: viper.GetString("app.livenessfile"),
		},
		Db: dbConfig{
			Driver:   viper.GetString("db.driver"),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetUint("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
			Database: viper.GetString("db.database"),
			Sslmode:  viper.GetString("db.sslmode"),
		},
		Redis: redisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetUint("redis.port"),
			Password: viper.GetString("redis.password"),
			Database: viper.GetInt("redis.database"),
		},
		Server: serverConfig{
			Port:         viper.GetUint("server.port"),
			TimeoutRead:  parseDuration(viper.GetString("server.timeout.read")),
			TimeoutWrite: parseDuration(viper.GetString("server.timteout.write")),
			TimeoutIdle:  parseDuration(viper.GetString("server.timeout.idle")),
		},
		Gateway: gatewayConfig{
			Host:    viper.GetString("gateway.host"),
			BaseURL: viper.GetString("gateway.baseurl"),
		},
		Token: tokenConfig{
			AccessSecretKey:  viper.GetString("token.access.secret"),
			AccessExpires:    parseDuration(viper.GetString("token.access.expires")),
			RefreshSecretKey: viper.GetString("token.refresh.secret"),
			RefreshExpires:   parseDuration(viper.GetString("token.refresh.expires")),
		},
	}
}

func parseDuration(t string) time.Duration {
	d, _ := time.ParseDuration(t)
	return d
}
