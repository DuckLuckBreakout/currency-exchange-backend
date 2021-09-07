package configer

import (
	"log"

	"github.com/spf13/viper"
)

type ServerData struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type RedisData struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
}

type PostgresqlData struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db-name"`
	Port     string `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	Sslmode  string `mapstructure:"sslmode"`
}

type SecretData struct {
	AccessSecret  string `mapstructure:"access-secret"`
	RefreshSecret string `mapstructure:"refresh-secret"`
}

type S3Data struct {
	AccessKeyId     string `mapstructure:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key"`
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	Acl             string `mapstructure:"acl"`
	Endpoint        string `mapstructure:"endpoint"`
}

type Config struct {
	Server     ServerData     `mapstructure:"server-data"`
	Redis      RedisData      `mapstructure:"redis-data"`
	Postgresql PostgresqlData `mapstructure:"postgresql-data"`
	Secret     SecretData     `mapstructure:"secrets-data"`
	S3         S3Data         `mapstructure:"s3-data"`
}

var AppConfig Config

func InitConfig(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatal(err)
	}
}
