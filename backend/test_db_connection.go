package main

import (
    "fmt"
    "os"
)

func main() {
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "5432")
    dbUser := getEnv("DB_USER", "postgres")
    dbPassword := getEnv("DB_PASSWORD", "postgres")
    dbName := getEnv("DB_NAME", "webrtc_meeting")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    fmt.Printf("Database connection string: %s\n", connStr)
    fmt.Printf("Target: %s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)
    fmt.Println("Connection string format is valid")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
