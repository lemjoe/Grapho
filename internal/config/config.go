package config

import (
	"fmt"
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
	PATH, exist := os.LookupEnv("DB_PATH")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "PATH"))
	} else {
		defaultConf.Path = PATH
	}
	PORT, exist := os.LookupEnv("DB_PORT")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "PORT"))
	} else {
		defaultConf.Port = PORT
	}
	HOST, exist := os.LookupEnv("DB_HOST")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "HOST"))
	} else {
		defaultConf.Host = HOST
	}
	DB_NAME, exist := os.LookupEnv("DB_NAME")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "DB_NAME"))
	} else {
		defaultConf.DBName = DB_NAME
	}
	return defaultConf, nil
}
