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
	SecretKey string `env:"TOKEN_SECRET" validate:"required"`
}

func LoadConfig() *Config {
	start := time.Now()
	viper.SetConfigName("config") // กำหนดชื่อไฟล์ config (without extension)
	log.Printf("SetConfigName %s", time.Since(start))
	viper.SetConfigType("yaml") // ระบุประเภทของไฟล์ config
	log.Printf("SetConfigType %s", time.Since(start))
	viper.AddConfigPath(".") // ระบุตำแหน่งของไฟล์ config อยู่ที่ working directory
	log.Printf("AddConfigPath %s", time.Since(start))
	viper.AutomaticEnv() // ให้อ่านค่าจาก env มา replace ในไฟล์ config
	log.Printf("AutomaticEnv %s", time.Since(start))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // แปลง _ underscore ใน env เป็น . dot notation ใน viper
	log.Printf("SetEnvKeyReplacer %s", time.Since(start))

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

	log.Printf("SetDefault %s", time.Since(start))

	// อ่านไฟล์ config
	err := viper.ReadInConfig()
	log.Printf("ReadInConfig %s", time.Since(start))
	if err != nil { // ถ้าอ่านไฟล์ config ไม่ได้ให้ข้ามไปเพราะให้เอาค่าจาก env มาแทนได้
		log.Println("please consider environment variables", err.Error())
	}

	config := &Config{
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
			SecretKey: viper.GetString("token.secret"),
		},
	}
	log.Printf("New Config %s", time.Since(start))

	// ตรวจสอบว่ากำหนดค่ามาครบหรือไม่
	validate := validator.New()

	err = validate.Struct(config)
	log.Printf("Validate Config %s", time.Since(start))
	if err != nil {
		panic(errors.New("load config error: " + err.Error()))
	}
	log.Printf("End %s", time.Since(start))
	return config
}

func parseDuration(t string) time.Duration {
	d, _ := time.ParseDuration(t)
	return d
}
