package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	InputDir        string `yaml:"input_dir"`
	OutputDir       string `yaml:"output_dir"`
	GoogleAnalytics string `yaml:"google_analytics"`
	BaseURL         string `yaml:"base_url"`
	PostFolder      string `yaml:"post_folder"`

	Method string `yaml:"method"`

	RemoteAddr string `yaml:"remote_addr"`
	User       string `yaml:"user"`
	KeyPath    string `yaml:"key_path"`
	Password   string `yaml:"password"`
}

var config Config

func ParseConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}
