package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	AppConfig *Config
)

type Config struct{
	AppPort string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBSSLMode string
	DBName string
	JWTSecret string
	// JWTExpireMinutes string
	JWTRefreshToken string
	JWTExpire string
	APPURL 	string
	CloudinaryURL string
	FrontendURL string
}
var envFiles = []string{".env.local", ".env"}

func LoadEnv() {
	loaded := false
	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("Config dimuat dari %s", file)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("Tidak ada file .env — memakai env var dari proses (normal di production)")
	}

	AppConfig =  &Config{
		AppPort: getEnv("PORT","3030"),
		DBHost: getEnv("DB_HOST","localhost"),
		DBPort: getEnv("DB_PORT","5432"),
		DBUser: getEnv("DB_USER","postgres"),
		DBPassword: getEnv("DB_PASSWORD","password"),
		DBName: getEnv("DB_NAME","project_management"),
		DBSSLMode: getEnv("DB_SSLMODE","disable"),
		JWTSecret: getEnv("JWT_SECRET","secret"),
		JWTExpire: getEnv("JWT_EXPIRED", getEnv("JWT_EXPIRY","6h")),
		// JWTExpireMinutes: getEnv("JWT_EXPIRY_MINUTES","60"),
		JWTRefreshToken: getEnv("REFRESH_TOKEN_EXPIRED","24h"),
		APPURL: getEnv("APP_URL","http://localhost:3030"),
		CloudinaryURL: getEnv("CLOUDINARY_URL",""),
		FrontendURL: getEnv("FRONTEND_URL","http://localhost:4321"),
	}
}

func getEnv(key string, fallback string)string {
valus, exist := os.LookupEnv(key)

if exist{
	return valus
}else{
	return fallback
}
}

func ConnectDB(){
	connect:= AppConfig
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		connect.DBHost, connect.DBPort, connect.DBUser, connect.DBPassword, connect.DBName, connect.DBSSLMode)

	db,err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal("Failed to connect Database", err)
	}
	log.Printf("Database Connected: %s:%s/%s (sslmode=%s)",
	connect.DBHost, connect.DBPort, connect.DBName, connect.DBSSLMode)
	sqlDB,err :=db.DB()
	if err != nil{
		log.Fatal("Failed to get Database instance",err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)



	DB = db
}