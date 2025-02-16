package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	v1 "simple-service/internal/api/v1"
	"simple-service/internal/device"
	"simple-service/internal/notification"

	_ "simple-service/docs"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           FCM Device Management API
// @version         1.0
// @description     API для управления устройствами и отправки push-уведомлений

// @host      4272517-lw36995.twc1.net
// @schemes   https
// @BasePath  /

type Response struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

var (
	deviceHandler       *v1.DeviceHandler
	notificationHandler *v1.NotificationHandler
)

func init() {
	// Инициализация базы данных
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "devices_db")
	dbPort := getEnvAsInt("DB_PORT", 5432)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Инициализация репозитория и обработчиков
	deviceRepo := device.NewPostgresRepository(db)
	deviceHandler = v1.NewDeviceHandler(deviceRepo)

	notificationService, err := notification.NewService("/app/fcm-config.json", deviceRepo)
	if err != nil {
		log.Fatalf("Failed to initialize notification service: %v", err)
	}
	notificationHandler = v1.NewNotificationHandler(notificationService)
}

func main() {
	// Регистрируем обработчики маршрутов
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/devices", deviceHandler.HandleDevices)
	http.HandleFunc("/devices/", deviceHandler.HandleDevice)
	http.HandleFunc("/notify/device/", notificationHandler.SendToDevice)
	http.HandleFunc("/notify/all", notificationHandler.SendToAll)

	// Добавляем Swagger UI
	http.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	certFile := os.Getenv("SSL_CERT")
	keyFile := os.Getenv("SSL_KEY")

	if certFile == "" || keyFile == "" {
		log.Fatal("SSL_CERT and SSL_KEY environment variables are required")
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":443",
		Handler:   http.DefaultServeMux,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting HTTPS server on https://%s", os.Getenv("HOST"))
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(Response{
		Message: "Добро пожаловать в наш API!",
		Status:  true,
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	json.NewEncoder(w).Encode(Response{
		Message: "Сервис работает нормально",
		Status:  true,
	})
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}
