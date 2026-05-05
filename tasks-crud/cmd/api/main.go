package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "tasks-crud/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"tasks-crud/internal/config"
	"tasks-crud/internal/domain"
	"tasks-crud/internal/handler"
	"tasks-crud/internal/middleware"
	"tasks-crud/internal/repository"
	"tasks-crud/internal/service"
)

// @title artemydottech API with JWT Authentication
// @version 1.0.0
// @description REST API for task management with JWT authentication

// @contact.name Artemij Zverev
// @contact.url https://github.com/artemydottech
// @contact.email artemiy.zverev@bk.ru

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @type apiKey
// @description JWT Authorization header. Введите ТОЛЬКО токен.

// @security BearerAuth
// HealthCheck проверка работоспособности сервиса
// @Summary Health check
// @Description Проверка работоспособности сервиса
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "todo-api",
		"version":   "1.0.0",
		"auth":      "jwt-enabled",
	})
}

// HealthResponse структура ответа health check
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Auth      string `json:"auth"`
}

// sendError отправка ошибки в JSON формате
func sendError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := domain.ErrorResponse{
		Error:  message,
		Status: statusCode,
		Time:   time.Now().Format(time.RFC3339),
	}

	if err != nil {
		errorResponse.Details = err.Error()
	}

	json.NewEncoder(w).Encode(errorResponse)
}

// NotFoundHandler обработчик для 404 ошибок
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	sendError(w, http.StatusNotFound, "Endpoint not found", nil)
}

// MethodNotAllowedHandler обработчик для 405 ошибок
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	sendError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequestLogger middleware для логирования запросов
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем обертку для ResponseWriter чтобы перехватить статус код
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf("[%s] %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.statusCode,
			duration,
		)
	})
}

// responseWriter кастомный ResponseWriter для перехвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	fmt.Println("🚀 Запуск Todo API с JWT аутентификацией...")
	fmt.Println("=============================================")

	cfg := config.Load()

	fmt.Printf("📋 Конфигурация:\n")
	fmt.Printf("Порт: %d\n", cfg.Port)
	fmt.Printf("Окружение: %s\n", cfg.Env)
	fmt.Printf("JWT Expiry: %v\n", cfg.JWTExpiry)
	fmt.Printf("Bcrypt Cost: %d\n", cfg.BcryptCost)

	if cfg.Env == "development" && cfg.JWTSecret == "your-secret-key-change-in-production" {
		fmt.Println("⚠️ВНИМАНИЕ: Используется дефолтный JWT секрет. В продакшене установите JWT_SECRET!")
	}

	fmt.Println("=============================================")

	fmt.Println("📦 Инициализация репозиториев...")
	taskRepo := repository.NewInMemoryTaskRepository()
	userRepo := repository.NewInMemoryUserRepository()
	fmt.Println("✅ Репозитории инициализированы")

	fmt.Println("⚙️  Инициализация сервисов...")
	taskService := service.NewTaskService(taskRepo)
	authService := service.NewAuthService(userRepo, cfg)
	fmt.Println("✅ Сервисы инициализированы")

	fmt.Println("🔄 Инициализация хендлеров...")
	taskHandler := handler.NewTaskHandler(taskService)
	authHandler := handler.NewAuthHandler(authService)
	fmt.Println("✅ Хендлеры инициализированы")

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	fmt.Println("🛣️ Настройка маршрутов...")
	public := router.PathPrefix("/api/v1").Subrouter()

	public.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	public.HandleFunc("/health", HealthCheck).Methods("GET")

	public.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/api/v1/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
		httpSwagger.UIConfig(map[string]string{
			"defaultModelsExpandDepth": "3",
		}),
	))

	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware(authService))

	protected.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	protected.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	protected.HandleFunc("/tasks/{id}", taskHandler.GetTaskByID).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	protected.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	swaggerURL := "/api/v1/swagger/index.html"
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, swaggerURL, http.StatusTemporaryRedirect)
	})

	router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, swaggerURL, http.StatusPermanentRedirect)
	})

	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, swaggerURL, http.StatusTemporaryRedirect)
	})

	handlerChain := enableCORS(
		RequestLogger(
			middleware.Logger(
				middleware.JSONContentType(
					router,
				),
			),
		),
	)

	// Настройка HTTP сервера
	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handlerChain,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("\n=============================================")
	fmt.Println("🌐 Сервер запущен:")
	fmt.Printf("   Основной URL: http://localhost:%d\n", cfg.Port)
	fmt.Printf("   API Base URL: http://localhost:%d/api/v1\n", cfg.Port)
	fmt.Println("\n📚 Документация:")
	fmt.Printf("   Swagger UI: http://localhost:%d%s\n", cfg.Port, swaggerURL)
	fmt.Printf("   OpenAPI JSON: http://localhost:%d/api/v1/swagger/doc.json\n", cfg.Port)
	fmt.Println("\n🔐 Аутентификация:")
	fmt.Printf("   Регистрация: POST http://localhost:%d/api/v1/auth/register\n", cfg.Port)
	fmt.Printf("   Вход: POST http://localhost:%d/api/v1/auth/login\n", cfg.Port)
	fmt.Println("\n📋 Примеры запросов:")
	fmt.Println("   curl -X POST http://localhost:8080/api/v1/auth/register \\")
	fmt.Println("     -H \"Content-Type: application/json\" \\")
	fmt.Println("     -d '{\"username\":\"test\",\"email\":\"test@example.com\",\"password\":\"password123\"}'")
	fmt.Println("\n   curl -X POST http://localhost:8080/api/v1/auth/login \\")
	fmt.Println("     -H \"Content-Type: application/json\" \\")
	fmt.Println("     -d '{\"email\":\"test@example.com\",\"password\":\"password123\"}'")
	fmt.Println("\n   curl -X GET http://localhost:8080/api/v1/tasks \\")
	fmt.Println("     -H \"Authorization: Bearer YOUR_JWT_TOKEN\"")
	fmt.Println("\n=============================================")
	fmt.Println("🛑 Для остановки нажмите Ctrl+C")
	fmt.Println("=============================================")

	// Запуск сервера
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Ошибка запуска сервера: %v", err)
	}
}
