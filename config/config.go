package config

import (
	//	"fmt"
	"os"
)

type Config struct {
	DBUser     string //`json:"db_user"`
	DBPassword string //`json:"db_password"`
	DBHost     string //`json:"db_host"`
	DBPort     string //`json:"db_port"`
	DBName     string //`json:"db_name"`
}

// 获取数据库配置信息
/*
func GetDatabaseInfo() (*Config, error) {
	/*
	cfg := &Config{}
	cfg.DBUser = getEnv("DB_USER", "root")
	cfg.DBPassword = getEnv("DB_PASSWORD", "123456")
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	port, err := strconv.ParseInt(getEnv("DB_PORT", "3306"), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Invalid DB_PORT: %w", err)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("DB_PORT must be between 1 and 65535")
	}
	这里本来想写检查，但是strconv.ParseInt已经检查了，所以就不写了

	cfg.DBName = getEnv("DB_NAME", "usersystem")
	cfg.DBPort = uint16(port) //强转，但安全
	//	检查是否为空
	if err := cfg.CheckValid(); err != nil {
		return nil, fmt.Errorf("Invalid database configuration: %w", err)
	}
	// 这里本来想写检查，但数据库连接时自己会检查的，直接懒得写
	return
}
*/
func GetDatabaseInfo() *Config {
	return &Config{
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "123456"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "usersystem"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

/*
func (c *Config) CheckValid() error {
	if c.DBUser == "" || c.DBPassword == "" || c.DBHost == "" || c.DBName == "" {
		return fmt.Errorf("Database configuration is incomplete")
	}
	return nil
}
*/
