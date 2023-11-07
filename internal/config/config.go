package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lemjoe/md-blog/internal/models"
)

func LookupAndParseEnvInt(envName string, defaultVal int) (int, bool) {
	env, exists := os.LookupEnv(envName)
	if !exists {
		return defaultVal, false
	}
	parsedInt, err := strconv.Atoi(env)
	if err != nil {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not valid: %w", envName, err))
		return defaultVal, false
	}
	return parsedInt, true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CreateDefaultConfig() error {
	log.Println("Config file is missing. Creating default config")
	input, err := os.ReadFile("./.env.default")
	if err != nil {
		return err
	}

	err = os.WriteFile("./.env", input, 0644)
	if err != nil {
		return err
	}
	return nil
}

func InitConfig(confPath string) (models.ConfigDB, error) {
	if confPath != "" {
		if fileExists(confPath) { //if from dot env
			if err := godotenv.Load(confPath); err != nil {
				return models.ConfigDB{}, fmt.Errorf("InitConfig: unable to read file '%s'", confPath)
			}

		} else {
			if err := CreateDefaultConfig(); err != nil {
				return models.ConfigDB{}, fmt.Errorf("unable to create default config. You should create it manually")
			}
			godotenv.Load(confPath)
		}

	}
	defaultConf := models.ConfigDB{
		DbType: "cloverdb",
		Path:   "./db",
	}
	DB_TYPE, exist := os.LookupEnv("DB_TYPE")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "DB_TYPE"))
	} else {
		defaultConf.DbType = DB_TYPE
	}
	DB_PATH, exist := os.LookupEnv("DB_PATH")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "DB_PATH"))
	} else {
		defaultConf.Path = DB_PATH
	}
	DB_PORT, exist := os.LookupEnv("DB_PORT")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "DB_PORT"))
	} else {
		defaultConf.Port = DB_PORT
	}
	DB_HOST, exist := os.LookupEnv("DB_HOST")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "DB_HOST"))
	} else {
		defaultConf.Host = DB_HOST
	}
	DB_NAME, exist := os.LookupEnv("DB_NAME")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "DB_NAME"))
	} else {
		defaultConf.DBName = DB_NAME
	}
	return defaultConf, nil
}
