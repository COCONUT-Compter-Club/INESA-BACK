package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/controller"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
)

func Routes(db *sql.DB, port string) {
	router := httprouter.New()
	
	//admin
	adminRepo := repository.NewAdminRepo()
	adminService := service.NewAdminService(adminRepo, db)
	adminController := controller.NewAdminController(adminService)

	router.POST("/api/bendahara/admin/daftar", adminController.SignUp)
	router.POST("/api/bendahara/admin/login", adminController.SignIn)
	router.GET("/api/bendahara/admin/:nik", adminController.FindByNik)

	//pemasukan
	pemasukanRepo := repository.NewPemasukanRepo()
	pemasukanService := service.NewPemasukanService(pemasukanRepo, db)
	pemasukanController := controller.NewPemasukanController(pemasukanService)

	router.POST("/api/bendahara/pemasukan/add", pemasukanController.AddPemasukan)
	router.PUT("/api/bendahara/pemasukan/update/:id", pemasukanController.UpdatePemasukan)
	router.GET("/api/bendahara/pemasukan/getall", pemasukanController.GetPemasukan)
	router.GET("/api/bendahara/pemasukan/get/:id", pemasukanController.GetById)
	router.DELETE("/api/bendahara/pemasukan/delete/:id", pemasukanController.DeletePemasukan)

	//pengeluaran
	pengeluaranRepo := repository.NewPengeluaranRepo()
	pengeluaranService := service.NewPengeluaranService(pengeluaranRepo, db)
	pengeluaranController := controller.NewPengeluaranController(pengeluaranService)

	router.POST("/api/bendahara/pengeluaran/add", pengeluaranController.AddPengeluaran)
	router.PUT("/api/bendahara/pengeluaran/update/:id", pengeluaranController.UpdatePengeluaran)
	router.GET("/api/bendahara/pengeluaran/getall", pengeluaranController.GetPengeluaran)
	router.GET("/api/bendahara/pengeluaran/get/:id", pengeluaranController.GetById)
	router.DELETE("/api/bendahara/pengeluaran/delete/:id", pengeluaranController.DeletePengeluaran)

	//transaksi
	transactionRepo := repository.NewTransactionRepo()
	transactionService := service.NewTransactionService(transactionRepo,db)
	transactionController := controller.NewTransactionController(transactionService)

	router.GET("/api/bendahara/transaksi/getall", transactionController.GetAllTransaction)
	router.GET("/api/bendahara/transaksi/getlast", transactionController.GetLastTransaction)

	//laporan keuangan
	laporanKeuanganRepo := repository.NewLaporanKeuanganRepo()
	laporanKeuanganService := service.NewLaporanKeuanganService(laporanKeuanganRepo, db)
	laporanKeuanganController := controller.NewLaporanKeuanganController(laporanKeuanganService)

	router.GET("/api/bendahara/laporan/getall", laporanKeuanganController.GetAllLaporan)
	router.GET("/api/bendahara/laporan/saldo", laporanKeuanganController.GetLastBalance)
	router.GET("/api/bendahara/laporan/pengeluaran", laporanKeuanganController.GetTotalExpenditure)
	router.GET("/api/bendahara/laporan/pemasukan", laporanKeuanganController.GetTotalIncome)
	router.GET("/api/bendahara/laporan/range", laporanKeuanganController.GetLaporanByDateRange)

	router.ServeFiles("/api/bendahara/uploads/*filepath", http.Dir("./uploads/"))

	handler := corsMiddleware(router)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	errServer := server.ListenAndServe()
	if errServer != nil {
		panic(errServer)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
