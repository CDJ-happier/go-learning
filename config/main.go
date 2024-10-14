package main

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"strings"
)

func main() {
	yamlConfig, err := LoadYAML("./config/config.yaml")
	if err != nil {
		log.Println(err)
		return
	}
	//printStructFields(reflect.ValueOf(*yamlConfig), 0)
	printYaml(yamlConfig)
	tomlConfig, err := LoadTOML("./config/config.toml")
	if err != nil {
		log.Println(err)
		return
	}
	// we can print all settings according to viper.Viper instance
	//fmt.Println("all settings: ", tomlConfig.AllSettings())
	appName := tomlConfig.GetString("appname")
	ip := tomlConfig.GetString("servers.alpha.ip")
	ports := tomlConfig.GetIntSlice("database.ports")
	fmt.Println("appname: ", appName)
	fmt.Println("ip: ", ip)
	fmt.Println("ports: ", ports)
}

func LoadTOML(fileName string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(fileName)

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("Config file not found")
		}
		return nil, err
	}

	return v, nil
}

func LoadYAML(fileName string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(fileName)

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("Config file not found")
		}
		return nil, err
	}

	// parse and unmarshal into config struct
	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func printYaml(cfg *Config) {
	fmt.Printf("%+v\n", cfg.Email.Host)
	fmt.Printf("%+v\n", cfg.Email.Port)
}

// 递归输出结构体成员
func printStructFields(v reflect.Value, indent int) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("json")
		if tag == "" {
			tag = t.Field(i).Name
		}
		if field.Kind() == reflect.Struct {
			fmt.Printf("%s%s:\n", strings.Repeat("  ", indent), tag)
		} else {
			fmt.Printf("%s%s: %v\n", strings.Repeat("  ", indent), tag, field.Interface())
		}

		// 如果字段是结构体类型，则递归调用
		if field.Kind() == reflect.Struct {
			printStructFields(field, indent+1)
		}
	}
}
