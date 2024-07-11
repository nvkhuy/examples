package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	JwtSecret        string `mapstructure:"JWT_SECRET" json:"jwt_secret"`
	AwsStorageUrl    string `mapstructure:"AWS_STORAGE_URL" json:"aws_storage_url"`
	AwsStorageBucket string `mapstructure:"AWS_STORAGE_BUCKET" json:"aws_storage_bucket"`
	AwsCdnUrl        string `mapstructure:"AWS_CDN_URL" json:"aws_cdn_url"`
	AwsCdnBucket     string `mapstructure:"AWS_CDN_BUCKET" json:"aws_cdn_bucket"`
	AwsRegion        string `mapstructure:"AWS_REGION" json:"aws_region"`
}

var config *Config

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	var err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config = &Config{}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	PrintJSON(config)
	return config, nil
}

func PrintJSON(i interface{}) {
	data, _ := json.MarshalIndent(i, "", "   ")
	fmt.Println(string(data))
}
