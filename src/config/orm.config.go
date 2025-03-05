package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDatabase menghubungkan ke database
func ConnectDatabase() *gorm.DB {
	// Load .env file jika ada
	err := godotenv.Load()
	if err != nil {
		log.Println("Gagal memuat file .env, menggunakan environment variables default")
	}

	// Ambil konfigurasi dari environment variables
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	// Koneksi ke database MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	log.Println("Database Connected!")
	return db
}
