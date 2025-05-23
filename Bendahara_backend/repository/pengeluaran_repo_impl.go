package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type pengeluaranRepoImpl struct {
}

func NewPengeluaranRepo() PengeluaranRepo {
	return &pengeluaranRepoImpl{}
}

// AddPengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) AddPengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran) (model.Pengeluaran, error) {
	idTransaksi := uuid.New().String()

	// Validasi tanggal
	if pengeluaran.Tanggal.IsZero() {
		return pengeluaran, fmt.Errorf("tanggal tidak boleh kosong")
	}

	// Insert ke history_transaksi
	queryTransaksi := `
		INSERT INTO history_transaksi (id_transaksi, tanggal, keterangan, jenis_transaksi, nominal)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := tx.ExecContext(ctx, queryTransaksi, idTransaksi, pengeluaran.Tanggal, pengeluaran.Keterangan, "Pengeluaran", pengeluaran.Nominal)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal menyisipkan data ke history_transaksi: %v", err)
	}

	// Insert ke tabel pengeluaran
	queryPengeluaran := `
		INSERT INTO pengeluaran (id_pengeluaran, tanggal, nota, nominal, keterangan, id_transaksi)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryPengeluaran, pengeluaran.Id, pengeluaran.Tanggal, pengeluaran.Nota, pengeluaran.Nominal, pengeluaran.Keterangan, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal menyisipkan data ke pengeluaran: %v", err)
	}

	// Ambil saldo terakhir sebelum tanggal pengeluaran
	var saldoSebelumnya uint64
	querySaldo := `
		SELECT saldo FROM laporan_keuangan 
		WHERE tanggal <= ?
		ORDER BY tanggal DESC
		LIMIT 1
	`
	err = tx.QueryRowContext(ctx, querySaldo, pengeluaran.Tanggal).Scan(&saldoSebelumnya)
	if err != nil && err != sql.ErrNoRows {
		return pengeluaran, fmt.Errorf("gagal mengambil saldo sebelumnya: %v", err)
	}

	// Validasi saldo cukup untuk pengeluaran
	if pengeluaran.Nominal > saldoSebelumnya {
		return pengeluaran, fmt.Errorf("saldo tidak cukup: %d, dibutuhkan: %d", saldoSebelumnya, pengeluaran.Nominal)
	}

	// Hitung saldo baru
	saldoBaru := saldoSebelumnya - pengeluaran.Nominal

	// Insert laporan keuangan baru
	idLaporan := uuid.New().String()
	queryLaporan := `
		INSERT INTO laporan_keuangan 
		(id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo, id_transaksi, nota)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryLaporan, idLaporan, pengeluaran.Tanggal, pengeluaran.Keterangan, 0, pengeluaran.Nominal, saldoBaru, idTransaksi, pengeluaran.Nota)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal menyisipkan data ke laporan_keuangan: %v", err)
	}

	// Update saldo semua entri setelah tanggal pengeluaran
	queryUpdate := `
		UPDATE laporan_keuangan
		SET saldo = saldo - ?
		WHERE tanggal > ?
	`
	_, err = tx.ExecContext(ctx, queryUpdate, pengeluaran.Nominal, pengeluaran.Tanggal)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memperbarui saldo di masa depan: %v", err)
	}

	pengeluaran.IdTransaksi = idTransaksi
	return pengeluaran, nil
}

// UpdatePengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) UpdatePengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran, id string) (model.Pengeluaran, error) {
	// Pastikan tanggal sudah dalam format time.Time
	if pengeluaran.Tanggal.IsZero() {
		return pengeluaran, fmt.Errorf("tanggal tidak boleh kosong")
	}

	// Ambil data pengeluaran sebelumnya untuk mendapatkan nominal lama, tanggal lama, dan id_transaksi
	var oldNominal uint64
	var tanggalRaw []byte
	var idTransaksi string
	queryFetch := `
		SELECT nominal, tanggal, id_transaksi 
		FROM pengeluaran 
		WHERE id_pengeluaran = ?
	`
	err := tx.QueryRowContext(ctx, queryFetch, id).Scan(&oldNominal, &tanggalRaw, &idTransaksi)
	if err != nil {
		if err == sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("pengeluaran dengan id %s tidak ditemukan", id)
		}
		return pengeluaran, fmt.Errorf("gagal mengambil data pengeluaran sebelumnya: %v", err)
	}

	// Konversi tanggalRaw ke time.Time
	tanggalStr := string(tanggalRaw)
	oldTanggal, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal mem-parsing tanggal lama: %v", err)
	}

	// Hitung selisih nominal (positif jika nominal baru lebih kecil, negatif jika lebih besar)
	nominalDiff := int64(oldNominal) - int64(pengeluaran.Nominal)

	// Perbarui tabel pengeluaran
	queryPengeluaran := `
		UPDATE pengeluaran 
		SET tanggal = ?, nota = ?, nominal = ?, keterangan = ? 
		WHERE id_pengeluaran = ?
	`
	_, err = tx.ExecContext(ctx, queryPengeluaran, pengeluaran.Tanggal, pengeluaran.Nota, pengeluaran.Nominal, pengeluaran.Keterangan, id)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memperbarui pengeluaran: %v", err)
	}

	// Perbarui tabel history_transaksi
	queryHistory := `
		UPDATE history_transaksi 
		SET tanggal = ?, keterangan = ?, nominal = ? 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryHistory, pengeluaran.Tanggal, pengeluaran.Keterangan, pengeluaran.Nominal, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memperbarui history_transaksi: %v", err)
	}

	// Perbarui entri laporan_keuangan yang terkait dengan transaksi ini
	queryLaporan := `
		UPDATE laporan_keuangan 
		SET tanggal = ?, keterangan = ?, pengeluaran = ?, saldo = saldo + ?, nota = ?
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryLaporan, pengeluaran.Tanggal, pengeluaran.Keterangan, pengeluaran.Nominal, nominalDiff, pengeluaran.Nota, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memperbarui laporan_keuangan: %v", err)
	}

	// Perbarui saldo dan total pengeluaran untuk semua entri laporan_keuangan setelah tanggal baru
	queryUpdateFuture := `
		UPDATE laporan_keuangan 
		SET saldo = saldo + ?, pengeluaran = GREATEST(0, pengeluaran + ?) 
		WHERE tanggal > ?
	`
	_, err = tx.ExecContext(ctx, queryUpdateFuture, nominalDiff, -nominalDiff, pengeluaran.Tanggal)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memperbarui laporan_keuangan di masa depan: %v", err)
	}

	// Jika tanggal berubah, perbarui saldo dan pengeluaran untuk entri antara tanggal lama dan baru
	if !oldTanggal.Equal(pengeluaran.Tanggal) {
		// Tambahkan kembali nominal lama pada entri setelah tanggal lama hingga tanggal baru
		queryAdjustOld := `
			UPDATE laporan_keuangan 
			SET saldo = saldo + ?, pengeluaran = GREATEST(0, pengeluaran + ?) 
			WHERE tanggal > ? AND tanggal <= ?
		`
		_, err = tx.ExecContext(ctx, queryAdjustOld, oldNominal, oldNominal, oldTanggal, pengeluaran.Tanggal)
		if err != nil {
			return pengeluaran, fmt.Errorf("gagal menyesuaikan laporan_keuangan untuk tanggal lama: %v", err)
		}

		// Kurangi nominal baru pada entri setelah tanggal baru
		queryAdjustNew := `
			UPDATE laporan_keuangan 
			SET saldo = saldo - ?, pengeluaran = GREATEST(0, pengeluaran + ?) 
			WHERE tanggal > ?
		`
		_, err = tx.ExecContext(ctx, queryAdjustNew, pengeluaran.Nominal, pengeluaran.Nominal, pengeluaran.Tanggal)
		if err != nil {
			return pengeluaran, fmt.Errorf("gagal menyesuaikan laporan_keuangan untuk tanggal baru: %v", err)
		}
	}

	return pengeluaran, nil
}

// GetPengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) GetPengeluaran(ctx context.Context, tx *sql.Tx, page int, pageSize int) ([]model.Pengeluaran, int, error) {
	offset := (page - 1) * pageSize
	query := "SELECT id_pengeluaran, tanggal, nota, nominal, keterangan FROM pengeluaran ORDER BY tanggal DESC LIMIT ? OFFSET ?"
	rows, err := tx.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data pengeluaran: %v", err)
	}
	defer rows.Close()

	var pengeluaranSlice []model.Pengeluaran
	for rows.Next() {
		pengeluaran := model.Pengeluaran{}
		var tanggalRaw []byte
		err := rows.Scan(&pengeluaran.Id, &tanggalRaw, &pengeluaran.Nota, &pengeluaran.Nominal, &pengeluaran.Keterangan)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal memindai data pengeluaran: %v", err)
		}
		tanggalStr := string(tanggalRaw)
		parsedTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
		}
		pengeluaran.Tanggal = parsedTime
		pengeluaranSlice = append(pengeluaranSlice, pengeluaran)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("kesalahan setelah mengiterasi baris: %v", err)
	}
	var total int
	countQuery := "SELECT COUNT(*) FROM pengeluaran"
	err = tx.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total pengeluaran: %v", err)
	}
	return pengeluaranSlice, total, nil
}

// FindById implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (model.Pengeluaran, error) {
	query := "SELECT id_pengeluaran, tanggal, nota, nominal, keterangan FROM pengeluaran WHERE id_pengeluaran = ?"
	row := tx.QueryRowContext(ctx, query, id)
	pengeluaran := model.Pengeluaran{}
	var tanggalRaw []byte
	err := row.Scan(&pengeluaran.Id, &tanggalRaw, &pengeluaran.Nota, &pengeluaran.Nominal, &pengeluaran.Keterangan)
	if err != nil {
		if err == sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("pengeluaran tidak ditemukan")
		}
		return pengeluaran, fmt.Errorf("gagal memindai data pengeluaran: %v", err)
	}
	tanggalStr := string(tanggalRaw)
	parsedTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
	}
	pengeluaran.Tanggal = parsedTime
	return pengeluaran, nil
}

// DeletePengeluaran implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) DeletePengeluaran(ctx context.Context, tx *sql.Tx, pengeluaran model.Pengeluaran) (model.Pengeluaran, error) {
	if pengeluaran.Id == "" {
		return pengeluaran, fmt.Errorf("id_pengeluaran tidak boleh kosong")
	}
	var idTransaksi string
	var nominal uint64
	var tanggalRaw []byte
	queryFetch := `
		SELECT id_transaksi, nominal, tanggal 
		FROM pengeluaran 
		WHERE id_pengeluaran = ?
	`
	err := tx.QueryRowContext(ctx, queryFetch, pengeluaran.Id).Scan(&idTransaksi, &nominal, &tanggalRaw)
	if err != nil {
		if err == sql.ErrNoRows {
			return pengeluaran, fmt.Errorf("pengeluaran dengan id %s tidak ditemukan", pengeluaran.Id)
		}
		return pengeluaran, fmt.Errorf("gagal mengambil data pengeluaran: %v", err)
	}
	tanggalStr := string(tanggalRaw)
	tanggalTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
	}
	if nominal <= 0 {
		return pengeluaran, fmt.Errorf("nilai nominal tidak valid: %d", nominal)
	}
	log.Printf("Berhasil mengambil pengeluaran: id=%s, id_transaksi=%s, nominal=%d, tanggal=%v", pengeluaran.Id, idTransaksi, nominal, tanggalTime)
	queryLaporan := `
		DELETE FROM laporan_keuangan 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryLaporan, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal menghapus dari laporan_keuangan: %v", err)
	}
	queryHistory := `
		DELETE FROM history_transaksi 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryHistory, idTransaksi)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal menghapus dari history_transaksi: %v", err)
	}
	queryPengeluaran := `
		DELETE FROM pengeluaran 
		WHERE id_pengeluaran = ?
	`
	_, err = tx.ExecContext(ctx, queryPengeluaran, pengeluaran.Id)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal menghapus dari pengeluaran: %v", err)
	}
	queryUpdateSaldo := `
		UPDATE laporan_keuangan
		SET saldo = saldo + ?
		WHERE tanggal > ?
	`
	result, err := tx.ExecContext(ctx, queryUpdateSaldo, nominal, tanggalTime)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memperbarui saldo di masa depan: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memeriksa jumlah baris yang diperbarui untuk saldo: %v", err)
	}
	log.Printf("Berhasil memperbarui %d baris di laporan_keuangan untuk saldo dengan id_pengeluaran %s, nominal %d, tanggal %v", rowsAffected, pengeluaran.Id, nominal, tanggalTime)
	queryUpdatePengeluaran := `
		UPDATE laporan_keuangan
		SET pengeluaran = GREATEST(0, pengeluaran - ?)
		WHERE tanggal > ?
	`
	result, err = tx.ExecContext(ctx, queryUpdatePengeluaran, nominal, tanggalTime)
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memperbarui pengeluaran di masa depan: %v", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return pengeluaran, fmt.Errorf("gagal memeriksa jumlah baris yang diperbarui untuk pengeluaran: %v", err)
	}
	log.Printf("Berhasil memperbarui %d baris di laporan_keuangan untuk pengeluaran dengan id_pengeluaran %s, nominal %d, tanggal %v", rowsAffected, pengeluaran.Id, nominal, tanggalTime)
	return pengeluaran, nil
}

// GetPengeluaranByDateRange implements PengeluaranRepo.
func (s *pengeluaranRepoImpl) GetPengeluaranByDateRange(ctx context.Context, tx *sql.Tx, startDate, endDate string, page int, pageSize int) ([]model.Pengeluaran, int, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT id_pengeluaran, tanggal, nota, nominal, keterangan 
		FROM pengeluaran 
		WHERE tanggal BETWEEN ? AND ? 
		ORDER BY tanggal DESC 
		LIMIT ? OFFSET ?
	`
	rows, err := tx.QueryContext(ctx, query, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data pengeluaran berdasarkan rentang tanggal: %v", err)
	}
	defer rows.Close()
	var pengeluaranSlice []model.Pengeluaran
	for rows.Next() {
		pengeluaran := model.Pengeluaran{}
		var tanggalRaw []byte
		err := rows.Scan(&pengeluaran.Id, &tanggalRaw, &pengeluaran.Nota, &pengeluaran.Nominal, &pengeluaran.Keterangan)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal memindai data pengeluaran: %v", err)
		}
		tanggalStr := string(tanggalRaw)
		parsedTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
		}
		pengeluaran.Tanggal = parsedTime
		pengeluaranSlice = append(pengeluaranSlice, pengeluaran)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("kesalahan setelah mengiterasi baris: %v", err)
	}
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM pengeluaran 
		WHERE tanggal BETWEEN ? AND ?
	`
	err = tx.QueryRowContext(ctx, countQuery, startDate, endDate).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total pengeluaran dalam rentang tanggal: %v", err)
	}
	return pengeluaranSlice, total, nil
}