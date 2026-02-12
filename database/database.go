package database

import (
	"database/sql"
	"fmt"
	"log"
	"user_system/config"

	_ "github.com/go-sql-driver/mysql" //空白导入
)

var DB *sql.DB

func InitDB() error {
	cfg := config.GetDatabaseInfo()
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	))
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(0)

	if err := DB.Ping(); err != nil { //如果连接失败，尝试创建数据库
		DB.Close() //关闭连接
		DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort,
		))
		if err != nil {
			return fmt.Errorf("Failed to connect to MySQL server: %w", err)
		}
		_, err = DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		DB.Close()                                                                                              //关闭连接
		DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", //重新连接到新创建的数据库
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		))
		if err != nil {
			return fmt.Errorf("Failed to connect to database after creating it: %w", err)
		}
	}

	if err := DB.Ping(); err != nil { //再次尝试连接
		return fmt.Errorf("Failed to ping database: %w", err)
	}
	log.Println("Database connection established successfully")
	return nil
}

func CloseDB() error {
	if DB != nil {
		if err := DB.Close(); err != nil {
			return fmt.Errorf("Failed to close database: %w", err)
		}
		log.Println("Database connection closed successfully")
	}
	return nil
}
