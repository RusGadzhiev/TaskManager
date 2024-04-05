package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	ConfigPath    = "../../config.yaml"
	EnvMySQLPass  = "mysql_pass"
	EnvServerPass = "server_pass"
)

type Config struct {
	// Env        string     `yaml:"env" env-required:"true"` можно и указать
	HTTPServer HTTPServer `yaml:"http_server"`
	MySQLDb    MySQLDb    `yaml:"mysql_db"`
	MongoDb    MongoDb    `yaml:"mongo_db"`
	RedisDb    RedisDb    `yaml:"redis_db"`
	LogLevel   string     `yaml:"log_level"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"username" env-required:"true"`
	Password    string        `env:"server_pass"`
}

type MySQLDb struct {
	Name     string `yaml:"name" env-default:"mysql"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"3306"`
	User     string `yaml:"username"`
	Password string `env:"mysql_pass"`
}

type MongoDb struct {
	Name string `yaml:"name" env-default:"mongodb"`
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"27017"`
}

type RedisDb struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"6379"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(ConfigPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
