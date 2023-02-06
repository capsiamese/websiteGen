package main

import (
	"gopkg.in/yaml.v3"
	"log"
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

func newConfig(out string) {
	c := Config{
		InputDir:        "./input",
		OutputDir:       "./output",
		GoogleAnalytics: "",
		BaseURL:         "https://localhost",
		PostFolder:      "posts",
		Method:          "local",
		RemoteAddr:      "localhost:22",
		User:            "admin",
		KeyPath:         "~/.ssh/id_rsa",
		Password:        "stdin",
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalln("[error] marshal config", err)
	}
	err = os.WriteFile(out, data, 0644)
	if err != nil {
		log.Fatalln("[error] write config to", out, err)
	}
}
