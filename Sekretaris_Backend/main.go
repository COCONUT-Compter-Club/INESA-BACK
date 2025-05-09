package main

import (
	"Sekertaris/config"
	"Sekertaris/controller"
	"Sekertaris/repository"
	"Sekertaris/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("Running on port: %s\n", port)

	// Permohonan Surat
	permohonanSuratRepo := repository.NewPermohonanSuratRepository(db)
	permohonanSuratService := service.NewPermohonanSuratService(permohonanSuratRepo)
	permohonanSuratController := controller.NewPermohonanSuratController(permohonanSuratService)

	// Surat Masuk
	suratMasukRepo := repository.NewSuratMasukRepository(db)
	suratMasukService := service.NewSuratMasukService(suratMasukRepo)
	suratMasukController := controller.NewSuratMasukController(suratMasukService)

	// Surat Keluar
	suratKeluarRepo := repository.NewSuratKeluarRepository(db)
	suratKeluarService := service.NewSuratKeluarService(suratKeluarRepo)
	suratKeluarController := controller.NewSuratKeluarController(suratKeluarService)

	// Setup router
	router := httprouter.New()
	router.HandleOPTIONS = true 

	// Serve static files
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	// Permohonan Surat Routes
	router.POST("/api/sekretaris/permohonansurat", permohonanSuratController.AddPermohonanSurat)
	router.GET("/api/sekretaris/permohonansurat", permohonanSuratController.GetPermohonanSurat)
	router.GET("/api/sekretaris/permohonansurat/get/:id", permohonanSuratController.GetPermohonanSuratByID)
	router.PUT("/api/sekretaris/permohonansurat/update/:id", permohonanSuratController.UpdatePermohonanSuratByID)
	router.DELETE("/api/sekretaris/permohonansurat/delete/:id", permohonanSuratController.DeletePermohonanSurat)
	router.PATCH("/api/sekretaris/permohonansurat/patch/:id", permohonanSuratController.UpdateStatus)

	// Surat Masuk Routes
	router.POST("/api/sekretaris/suratmasuk", suratMasukController.AddSuratMasuk)
	router.GET("/api/sekretaris/suratmasuk", suratMasukController.GetSuratMasuk)
	router.GET("/api/sekretaris/suratmasuk/get/:id", suratMasukController.GetSuratById)
	router.PUT("/api/sekretaris/suratmasuk/update/:id", suratMasukController.UpdateSuratMasukByID)
	router.DELETE("/api/sekretaris/suratmasuk/delete/:id", suratMasukController.DeleteSuratMasuk)

	// Surat Keluar Routes
	router.POST("/api/sekretaris/suratkeluar", suratKeluarController.AddSuratKeluar)
	router.GET("/api/sekretaris/suratkeluar", suratKeluarController.GetAllSuratKeluar)
	router.GET("/api/sekretaris/suratkeluar/get/:id", suratKeluarController.GetSuratKeluarById)
	router.PUT("/api/sekretaris/suratkeluar/update/:id", suratKeluarController.UpdateSuratKeluarByID)
	router.GET("/api/sekretaris/suratkeluar/file/:filename", suratKeluarController.ServeFile)
	router.DELETE("/api/sekretaris/suratkeluar/delete/:id", suratKeluarController.DeleteSuratKeluar)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://bontomanai.inesa.id",
		},
		 // Tambah semua kemungkinan origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Content-Disposition",
			"ngrok-skip-browser-warning",
		},
		AllowCredentials: true,
		Debug:            true,
	})

	// Wrap router with CORS
	handler := c.Handler(router)

	// Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("Server running on http://:%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
	}
	fmt.Println("Server stopped gracefully")
}
