package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port uint
	Dsn  string
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath(".")                               // ระบุตำแหน่งของไฟล์ config อยู่ที่ working directory
	viper.SetConfigName("config")                          // กำหนดชื่อไฟล์ config (without extension)
	viper.SetConfigType("yaml")                            // ระบุประเภทของไฟล์ config
	viper.AutomaticEnv()                                   // ให้อ่านค่าจาก env มา replace ในไฟล์ config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // แปลง _ underscore ใน env เป็น . dot notation ใน viper

	// อ่านไฟล์ config
	e := viper.ReadInConfig()
	if e != nil { // ถ้าอ่านไฟล์ config ไม่ได้ให้ข้ามไปเพราะให้เอาค่าจาก env มาแทนได้
		fmt.Println("please consider environment variables", e.Error())
	}

	// กำหนด Default Value
	viper.SetDefault("port", 3000)
	// err = viper.Unmarshal(&c)
	c = Config{
		Port: viper.GetUint("port"),
		Dsn:  viper.GetString("dsn"),
	}
	err = nil
	if len(c.Dsn) == 0 {
		err = errors.New("env DSN is required")
	}

	return
}
