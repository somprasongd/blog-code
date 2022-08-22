package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type configuration struct {
	App appConfig
	Db  dbConfig
}

type appConfig struct {
	Port uint
	Env  string
}

type dbConfig struct {
	Driver   string
	Host     string
	Port     uint
	Username string
	Password string
	Database string
}

var Config *configuration

func LoadConfig() {
	viper.SetConfigName("config")                          // กำหนดชื่อไฟล์ config (without extension)
	viper.SetConfigType("yaml")                            // ระบุประเภทของไฟล์ config
	viper.AddConfigPath(".")                               // ระบุตำแหน่งของไฟล์ config อยู่ที่ working directory
	viper.AutomaticEnv()                                   // ให้อ่านค่าจาก env มา replace ในไฟล์ config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // แปลง _ underscore ใน env เป็น . dot notation ใน viper

	err := viper.ReadInConfig() // อ่านไฟล์ config
	if err != nil {             // ถ้าอ่านไฟล์ config ไม่ได้ให้ข้ามไปเพราะให้เอาค่าจาก env มาแทนได้
		fmt.Println("please consider environment variables", err.Error())
	}

	// กำหนด Default Value
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.env", "development")

	Config = &configuration{
		App: appConfig{
			Port: viper.GetUint("app.port"),
			Env:  viper.GetString("app.env"),
		},
		Db: dbConfig{
			Driver:   viper.GetString("db.driver"),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetUint("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
			Database: viper.GetString("db.database"),
		},
	}

	// ตรวจสอบว่ากำหนดค่ามาครบหรือไม่
	validate := validator.New()

	err = validate.Struct(Config)
	if err != nil {
		panic(errors.New("load config error: " + err.Error()))
	}
}
