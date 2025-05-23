package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type pengeluaranServiceImpl struct {
	PengeluaranRepo repository.PengeluaranRepo
	DB              *sql.DB
}

func NewPengeluaranService(pengeluaranRepo repository.PengeluaranRepo, db *sql.DB) PengeluaranService {
	return &pengeluaranServiceImpl{
		PengeluaranRepo: pengeluaranRepo,
		DB:              db,
	}
}

// AddPengeluaran implements PengeluaranService.
func (s *pengeluaranServiceImpl) AddPengeluaran(ctx context.Context, r *http.Request, pengeluaranRequest dto.PengeluaranRequest) (dto.PengeluaranResponse, error) {
	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal menguraikan formulir: %v", err)
	}

	// Ambil nilai tanggal dan waktu dari form dan parse ke time.Time
	tanggalWaktuStr := r.FormValue("tanggal") // Contoh format: "2006-01-02T15:04"
	tanggalWaktu, err := time.Parse("2006-01-02 15:04", tanggalWaktuStr)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("format tanggal atau waktu tidak valid: %v", err)
	}

	// Isi pengeluaranRequest dengan data yang diambil dari form
	pengeluaranRequest = dto.PengeluaranRequest{
		Tanggal:    tanggalWaktuStr, // Sekarang bertipe string
		Keterangan: r.FormValue("keterangan"),
	}

	// Ambil file dari form
	file, handler, err := r.FormFile("nota")
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal membaca file: %v", err)
	}
	defer file.Close()

	// Format tanggal dan waktu untuk nama file
	formattedDateTime := tanggalWaktu.Format("2006-01-02-15-04") // Format: YYYY-MM-DD-HH-MM

	// Buat nama file dengan format: tanggal-waktu-uuid
	fileName := fmt.Sprintf("%s-%s.jpeg", formattedDateTime, uuid.New().String())
	handler.Filename = fileName

	// Buat direktori upload jika belum ada
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Simpan file ke direktori upload
	filePath := filepath.Join(uploadDir, handler.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal membuat file: %v", err)
	}
	defer out.Close()

	// Salin file yang diunggah ke file yang baru dibuat
	_, err = io.Copy(out, file)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal menyalin file: %v", err)
	}

	// Simpan nama file ke dalam request
	pengeluaranRequest.Nota = handler.Filename

	// Konversi nominal dari string ke integer
	nominal, err := (strconv.Atoi(r.FormValue("nominal")))
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("nominal harus berupa angka yang valid: %v", err)
	}
	pengeluaranRequest.Nominal = uint64(nominal)

	// Validasi input
	if pengeluaranRequest.Tanggal == "" || pengeluaranRequest.Nota == "" || pengeluaranRequest.Nominal == 0 {
		return dto.PengeluaranResponse{}, fmt.Errorf("tanggal, catatan, atau nominal tidak boleh kosong")
	}

	// Parsing tanggal dari string ke time.Time
	tanggal, err := time.Parse("2006-01-02 15:04", pengeluaranRequest.Tanggal)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("format tanggal tidak valid, yang diharapkan 'YYYY-MM-DD HH:MM:SS'")
	}

	// Mulai transaksi database
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	// Buat objek Pengeluaran
	pengeluaran := model.Pengeluaran{
		Id:         uuid.New().String(),
		Tanggal:    tanggal, // Gunakan tanggal yang sudah di-parse
		Nota:       pengeluaranRequest.Nota,
		Nominal:    pengeluaranRequest.Nominal,
		Keterangan: pengeluaranRequest.Keterangan,
	}

	// Tambahkan pengeluaran ke database
	addPengeluaran, err := s.PengeluaranRepo.AddPengeluaran(ctx, tx, pengeluaran)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal menambahkan pengeluaran: %v", err)
	}

	// Kembalikan respons
	return util.ConvertPengeluaranToResponseDTO(addPengeluaran), nil
}

// UpdatePengeluaran implements PengeluaranService.
func (s *pengeluaranServiceImpl) UpdatePengeluaran(ctx context.Context, r *http.Request, pengeluaranRequest dto.PengeluaranRequest, id string) (dto.PengeluaranResponse, error) {
	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal menguraikan formulir: %v", err)
	}

	// Ambil nilai tanggal dan waktu dari form dan parse ke time.Time
	tanggalWaktuStr := r.FormValue("tanggal") // Contoh format: "2006-01-02T15:04"
	tanggalWaktu, err := time.Parse("2006-01-02 15:04", tanggalWaktuStr)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("format tanggal atau waktu tidak valid: %v", err)
	}

	// Isi pengeluaranRequest dengan data yang diambil dari form
	pengeluaranRequest = dto.PengeluaranRequest{
		Tanggal:    tanggalWaktuStr, // Sekarang bertipe string
		Keterangan: r.FormValue("keterangan"),
	}

	// Mulai transaksi database
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	// Cari pengeluaran berdasarkan ID
	pengeluaran, err := s.PengeluaranRepo.FindById(ctx, tx, id)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal menemukan pengeluaran: %v", err)
	}

	// Jika file nota diunggah, simpan file baru dan hapus file lama
	file, handler, err := r.FormFile("nota")
	if err == nil {
		defer file.Close()

		// Format tanggal dan waktu untuk nama file
		formattedDateTime := tanggalWaktu.Format("2006-01-02-15-04") // Format: YYYY-MM-DD-HH-MM

		// Buat nama file dengan format: tanggal-waktu-uuid
		fileName := fmt.Sprintf("%s-%s.jpeg", formattedDateTime, uuid.New().String())
		handler.Filename = fileName

		// Buat direktori upload jika belum ada
		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		// Simpan file ke direktori upload
		filePath := filepath.Join(uploadDir, handler.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			return dto.PengeluaranResponse{}, fmt.Errorf("gagal membuat file: %v", err)
		}
		defer out.Close()

		// Salin file yang diunggah ke file yang baru dibuat
		_, err = io.Copy(out, file)
		if err != nil {
			return dto.PengeluaranResponse{}, fmt.Errorf("gagal menyalin file: %v", err)
		}

		// Hapus file nota lama jika ada
		if pengeluaran.Nota != "" {
			oldFilePath := filepath.Join(uploadDir, pengeluaran.Nota)
			if _, err := os.Stat(oldFilePath); err == nil {
				os.Remove(oldFilePath)
			}
		}

		// Simpan nama file baru ke dalam request
		pengeluaranRequest.Nota = handler.Filename
	} else {
		// Jika tidak ada file baru, gunakan nota yang sudah ada
		pengeluaranRequest.Nota = pengeluaran.Nota
	}

	// Konversi nominal dari string ke integer
	nominal, err := strconv.Atoi(r.FormValue("nominal"))
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("nominal harus berupa angka yang valid: %v", err)
	}
	pengeluaranRequest.Nominal = uint64(nominal)

	// Validasi input
	if pengeluaranRequest.Tanggal == "" || pengeluaranRequest.Nota == "" || pengeluaranRequest.Nominal == 0 {
		return dto.PengeluaranResponse{}, fmt.Errorf("tanggal, catatan, atau nominal tidak boleh kosong")
	}

	// Parsing tanggal dari string ke time.Time
	tanggal, err := time.Parse("2006-01-02 15:04", pengeluaranRequest.Tanggal)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("format tanggal tidak valid, yang diharapkan 'YYYY-MM-DD HH:MM:SS'")
	}

	// Update data pengeluaran
	pengeluaran.Tanggal = tanggal // Gunakan tanggal yang sudah di-parse
	pengeluaran.Nota = pengeluaranRequest.Nota
	pengeluaran.Nominal = pengeluaranRequest.Nominal
	pengeluaran.Keterangan = pengeluaranRequest.Keterangan

	// Simpan perubahan ke database
	updatePengeluaran, err := s.PengeluaranRepo.UpdatePengeluaran(ctx, tx, pengeluaran, id)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal memperbarui pengeluaran: %v", err)
	}

	// Kembalikan respons
	return util.ConvertPengeluaranToResponseDTO(updatePengeluaran), nil
}

// GetPengeluaran implements PengeluaranService.
func (s *pengeluaranServiceImpl) GetPengeluaran(ctx context.Context, page int, pageSize int) (dto.PengeluaranPaginationResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PengeluaranPaginationResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	// Validasi parameter pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15 // Default page size
	}

	pengeluaran, total, err := s.PengeluaranRepo.GetPengeluaran(ctx, tx, page, pageSize)
	if err != nil {
		return dto.PengeluaranPaginationResponse{}, fmt.Errorf("gagal mendapatkan semua pengeluaran: %v", err)
	}

	// Hitung total halaman
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	response := dto.PengeluaranPaginationResponse{
		Items:      util.ConvertPengeluaranToListResponseDTO(pengeluaran),
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return response, nil
}

// DeletePengeluaran implements PengeluaranService.
func (s *pengeluaranServiceImpl) DeletePengeluaran(ctx context.Context, id string) (dto.PengeluaranResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Commit()

	pengeluaran, err := s.PengeluaranRepo.FindById(ctx, tx, id)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("pengeluaran tidak ditemukan: %v", err)
	}

	pengeluaran, err = s.PengeluaranRepo.DeletePengeluaran(ctx, tx, pengeluaran)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal menghapus pengeluaran: %v", err)
	}

	return util.ConvertPengeluaranToResponseDTO(pengeluaran), nil
}

// GetById implements PengeluaranService.
func (s *pengeluaranServiceImpl) GetById(ctx context.Context, id string) (dto.PengeluaranResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Commit()

	pengeluaran, err := s.PengeluaranRepo.FindById(ctx, tx, id)
	if err != nil {
		return dto.PengeluaranResponse{}, fmt.Errorf("pengeluaran tidak ditemukan: %v", err)
	}

	return util.ConvertPengeluaranToResponseDTO(pengeluaran), nil
}

// GetPengeluaranByDateRange implements PengeluaranService.
func (s *pengeluaranServiceImpl) GetPengeluaranByDateRange(ctx context.Context, startDate, endDate string, page int, pageSize int) (dto.PengeluaranPaginationResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.PengeluaranPaginationResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	// Validasi parameter pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15 // Default page size
	}

	// Validasi tanggal
	if startDate == "" || endDate == "" {
		return dto.PengeluaranPaginationResponse{}, fmt.Errorf("tanggal mulai dan tanggal akhir diperlukan")
	}

	// Parse tanggal untuk memastikan format valid
	_, err = time.Parse("2006-01-02", startDate)
	if err != nil {
		return dto.PengeluaranPaginationResponse{}, fmt.Errorf("format tanggal mulai tidak valid, diharapkan 'YYYY-MM-DD': %v", err)
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return dto.PengeluaranPaginationResponse{}, fmt.Errorf("format tanggal akhir tidak valid, diharapkan 'YYYY-MM-DD': %v", err)
	}

	pengeluaran, total, err := s.PengeluaranRepo.GetPengeluaranByDateRange(ctx, tx, startDate, endDate, page, pageSize)
	if err != nil {
		return dto.PengeluaranPaginationResponse{}, fmt.Errorf("gagal mendapatkan pengeluaran berdasarkan rentang tanggal: %v", err)
	}

	// Hitung total halaman
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	response := dto.PengeluaranPaginationResponse{
		Items:      util.ConvertPengeluaranToListResponseDTO(pengeluaran),
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return response, nil
}