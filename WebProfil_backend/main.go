package main

import (
	"fmt"
	"net/http"
	"strings"
	"godesaapps/config"
	"godesaapps/controller"
	"godesaapps/repository"
	"godesaapps/service"
	"godesaapps/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"

	"gorm.io/driver/mysql"
    "gorm.io/gorm"

)

func main() {
	fmt.Println("DesaApps Runn...")

	// Koneksi ke database
	mysql, err := config.ConnectToDatabase()
	util.SentPanicIfError(err)

	// Inisialisasi repository
	userRepository := repository.NewUserRepositoryImpl(mysql)
	wargaRepository := repository.NewWargaRepository(mysql)
	contactMessageRepository := repository.NewContactMessageRepository(mysql) // Perbaikan: ganti db dengan mysql

	// Inisialisasi service
	userService := service.NewUserServiceImpl(userRepository, mysql)
	wargaService := service.NewWargaService(wargaRepository)
	contactMessageService := service.NewContactMessageService(contactMessageRepository)

	// Inisialisasi controller
	userController := controller.NewUserControllerImpl(userService)
	wargaController := controller.NewWargaController(wargaService)
	contactMessageController := controller.NewContactMessageController(contactMessageService)

	// Inisialisasi router
	router := httprouter.New()

	// User routes
	router.POST("/api/user/create", userController.CreateUser)
	router.GET("/api/user", userController.ReadUser)
	router.POST("/api/user/login", userController.LoginUser)
	router.GET("/api/user/me", VerifyJWT(userController.GetUserInfo))
	router.POST("/api/user/forgot-password", userController.ForgotPassword)
	router.POST("/api/user/reset-password", userController.ResetPassword)

	// Warga routes
	router.POST("/api/warga/register", wargaController.RegisterWarga)

	// Contact message route
	router.POST("/api/contact/message", contactMessageController.CreateMessage)

	// Terapkan middleware CORS
	handler := corsMiddleware(router)

	// Inisialisasi server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.AppHost, config.AppPort),
		Handler: handler,
	}

	// Jalankan server
	errServer := server.ListenAndServe()
	util.SentPanicIfError(errServer)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VerifyJWT(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Token Format", http.StatusUnauthorized)
			return
		}

		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret_key"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or Expired Token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("User-Email", claims.Email)
		next(w, r, ps)
	}
}

// Fungsi ConnectToDatabase (sebagai cadangan jika tidak ada di package config)
func ConnectToDatabase() (*gorm.DB, error) {
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
