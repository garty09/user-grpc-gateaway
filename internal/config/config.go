package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	defaultHTTPPort  = 8080
	defaultGRPCPort  = 8090
	defaultRedisAddr = "localhost:6379"
)

type Config struct {
	HTTPPort int `yaml:"http_port"`
	GRPCPort int `yaml:"grpc_port"`
	// the data source name (DSN) for connecting to the database
	DSN       string `yaml:"dsn"`
	RedisAddr string
}

func (c Config) Validate() error {
	if c.DSN == "" {
		return fmt.Errorf("connect string in empty")
	}
	if c.RedisAddr == "" {
		return fmt.Errorf("redis connect string in empty")
	}
	if c.HTTPPort <= 0 && c.HTTPPort > 64000 {
		return fmt.Errorf("http not in range")
	}
	if c.GRPCPort <= 0 && c.GRPCPort > 64000 {
		return fmt.Errorf("grpc not in range")
	}
	if c.HTTPPort == c.GRPCPort {
		return fmt.Errorf("can not be equal")
	}
	return nil
}

func Load(file string) (*Config, error) {
	// default config
	c := Config{
		HTTPPort:  defaultHTTPPort,
		GRPCPort:  defaultGRPCPort,
		RedisAddr: defaultRedisAddr,
	}

	// load from YAML config file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	// validation
	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}
