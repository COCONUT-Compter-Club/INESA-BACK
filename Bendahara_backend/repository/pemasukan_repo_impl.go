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

type pemasukanRepoImpl struct {
}

func NewPemasukanRepo() PemasukanRepo {
	return &pemasukanRepoImpl{}
}

// AddPemasukan implements PemasukanRepo.
func (s *pemasukanRepoImpl) AddPemasukan(ctx context.Context, tx *sql.Tx, pemasukan model.Pemasukan) (model.Pemasukan, error) {
	idTransaksi := uuid.New().String()

	// Validasi tanggal
	if pemasukan.Tanggal.IsZero() {
		return pemasukan, fmt.Errorf("tanggal tidak boleh kosong")
	}

	// Insert ke history_transaksi
	queryTransaksi := `
		INSERT INTO history_transaksi (id_transaksi, tanggal, keterangan, jenis_transaksi, nominal) 
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := tx.ExecContext(ctx, queryTransaksi, idTransaksi, pemasukan.Tanggal, pemasukan.Keterangan, "Pemasukan", pemasukan.Nominal)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal menyisipkan data ke history_transaksi: %v", err)
	}

	// Insert ke tabel pemasukan
	queryPemasukan := `
		INSERT INTO pemasukan (id_pemasukan, tanggal, kategori, keterangan, nominal, nota, id_transaksi) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryPemasukan, pemasukan.Id, pemasukan.Tanggal, pemasukan.Kategori, pemasukan.Keterangan, pemasukan.Nominal, pemasukan.Nota, idTransaksi)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal menyisipkan data ke pemasukan: %v", err)
	}

	// Ambil saldo terakhir sebelum tanggal pemasukan
	var saldoSebelumnya uint64
	querySaldo := `
		SELECT saldo FROM laporan_keuangan 
		WHERE tanggal <= ?
		ORDER BY tanggal DESC
		LIMIT 1
	`
	err = tx.QueryRowContext(ctx, querySaldo, pemasukan.Tanggal).Scan(&saldoSebelumnya)
	if err != nil && err != sql.ErrNoRows {
		return pemasukan, fmt.Errorf("gagal mengambil saldo sebelumnya: %v", err)
	}

	// Hitung saldo baru
	saldoBaru := saldoSebelumnya + pemasukan.Nominal

	// Insert laporan keuangan baru
	idLaporan := uuid.New().String()
	queryLaporan := `
		INSERT INTO laporan_keuangan 
		(id_laporan, tanggal, keterangan, pemasukan, pengeluaran, saldo, id_transaksi, nota)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryLaporan, idLaporan, pemasukan.Tanggal, pemasukan.Keterangan, pemasukan.Nominal, 0, saldoBaru, idTransaksi, pemasukan.Nota)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal menyisipkan data ke laporan_keuangan: %v", err)
	}

	// Update saldo semua entri setelah tanggal pemasukan
	queryUpdate := `
		UPDATE laporan_keuangan
		SET saldo = saldo + ?
		WHERE tanggal > ?
	`
	_, err = tx.ExecContext(ctx, queryUpdate, pemasukan.Nominal, pemasukan.Tanggal)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memperbarui saldo di masa depan: %v", err)
	}

	pemasukan.IdTransaksi = idTransaksi
	return pemasukan, nil
}

// UpdatePemasukan implements PemasukanRepo.
func (s *pemasukanRepoImpl) UpdatePemasukan(ctx context.Context, tx *sql.Tx, pemasukan model.Pemasukan, id string) (model.Pemasukan, error) {
	// Pastikan tanggal sudah dalam format time.Time
	if pemasukan.Tanggal.IsZero() {
		return pemasukan, fmt.Errorf("tanggal tidak boleh kosong")
	}

	// Ambil data pemasukan sebelumnya untuk mendapatkan nominal lama, tanggal lama, dan id_transaksi
	var oldNominal uint64
	var tanggalRaw []byte
	var idTransaksi string
	queryFetch := `
		SELECT nominal, tanggal, id_transaksi 
		FROM pemasukan 
		WHERE id_pemasukan = ?
	`
	err := tx.QueryRowContext(ctx, queryFetch, id).Scan(&oldNominal, &tanggalRaw, &idTransaksi)
	if err != nil {
		if err == sql.ErrNoRows {
			return pemasukan, fmt.Errorf("pemasukan dengan id %s tidak ditemukan", id)
		}
		return pemasukan, fmt.Errorf("gagal mengambil data pemasukan sebelumnya: %v", err)
	}

	// Konversi tanggalRaw ke time.Time
	tanggalStr := string(tanggalRaw)
	oldTanggal, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal mem-parsing tanggal lama: %v", err)
	}

	// Hitung selisih nominal
	nominalDiff := int64(pemasukan.Nominal) - int64(oldNominal)

	// Perbarui tabel pemasukan
	queryPemasukan := `
		UPDATE pemasukan 
		SET tanggal = ?, kategori = ?, keterangan = ?, nominal = ?, nota = ? 
		WHERE id_pemasukan = ?
	`
	_, err = tx.ExecContext(ctx, queryPemasukan, pemasukan.Tanggal, pemasukan.Kategori, pemasukan.Keterangan, pemasukan.Nominal, pemasukan.Nota, id)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memperbarui pemasukan: %v", err)
	}

	// Perbarui tabel history_transaksi
	queryHistory := `
		UPDATE history_transaksi 
		SET tanggal = ?, keterangan = ?, nominal = ? 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryHistory, pemasukan.Tanggal, pemasukan.Keterangan, pemasukan.Nominal, idTransaksi)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memperbarui history_transaksi: %v", err)
	}

	// Perbarui entri laporan_keuangan yang terkait dengan transaksi ini
	queryLaporan := `
		UPDATE laporan_keuangan 
		SET tanggal = ?, keterangan = ?, pemasukan = ?, saldo = saldo + ?, nota = ?
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryLaporan, pemasukan.Tanggal, pemasukan.Keterangan, pemasukan.Nominal, nominalDiff, pemasukan.Nota, idTransaksi)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memperbarui laporan_keuangan: %v", err)
	}

	// Perbarui saldo dan total pemasukan untuk semua entri laporan_keuangan setelah tanggal baru
	queryUpdateFuture := `
		UPDATE laporan_keuangan 
		SET saldo = saldo + ?, pemasukan = pemasukan + ? 
		WHERE tanggal > ?
	`
	_, err = tx.ExecContext(ctx, queryUpdateFuture, nominalDiff, nominalDiff, pemasukan.Tanggal)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memperbarui laporan_keuangan di masa depan: %v", err)
	}

	// Jika tanggal berubah, perbarui saldo dan pemasukan untuk entri antara tanggal lama dan baru
	if !oldTanggal.Equal(pemasukan.Tanggal) {
		// Kurangi pengaruh nominal lama pada entri setelah tanggal lama
		queryAdjustOld := `
			UPDATE laporan_keuangan 
			SET saldo = saldo - ?, pemasukan = pemasukan - ? 
			WHERE tanggal > ? AND tanggal <= ?
		`
		_, err = tx.ExecContext(ctx, queryAdjustOld, oldNominal, oldNominal, oldTanggal, pemasukan.Tanggal)
		if err != nil {
			return pemasukan, fmt.Errorf("gagal menyesuaikan laporan_keuangan untuk tanggal lama: %v", err)
		}

		// Tambahkan pengaruh nominal baru pada entri setelah tanggal baru
		queryAdjustNew := `
			UPDATE laporan_keuangan 
			SET saldo = saldo + ?, pemasukan = pemasukan + ? 
			WHERE tanggal > ?
		`
		_, err = tx.ExecContext(ctx, queryAdjustNew, pemasukan.Nominal, pemasukan.Nominal, pemasukan.Tanggal)
		if err != nil {
			return pemasukan, fmt.Errorf("gagal menyesuaikan laporan_keuangan untuk tanggal baru: %v", err)
		}
	}

	return pemasukan, nil
}

// GetPemasukan implements PemasukanRepo.
func (s *pemasukanRepoImpl) GetPemasukan(ctx context.Context, tx *sql.Tx, page int, pageSize int) ([]model.Pemasukan, int, error) {
	// Hitung offset
	offset := (page - 1) * pageSize

	// Query untuk mendapatkan data dengan pagination
	query := "SELECT id_pemasukan, tanggal, kategori, keterangan, nominal, nota FROM pemasukan ORDER BY tanggal DESC LIMIT ? OFFSET ?"
	rows, err := tx.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data pemasukan: %v", err)
	}
	defer rows.Close()

	var pemasukanSlice []model.Pemasukan
	for rows.Next() {
		pemasukan := model.Pemasukan{}
		var tanggalRaw []byte

		err := rows.Scan(&pemasukan.Id, &tanggalRaw, &pemasukan.Kategori, &pemasukan.Keterangan, &pemasukan.Nominal, &pemasukan.Nota)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal memindai data pemasukan: %v", err)
		}

		// Konversi tanggalRaw ke time.Time
		tanggalStr := string(tanggalRaw)
		parsedTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
		}
		pemasukan.Tanggal = parsedTime

		pemasukanSlice = append(pemasukanSlice, pemasukan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("kesalahan setelah mengiterasi baris: %v", err)
	}

	// Query untuk mendapatkan total data
	var total int
	countQuery := "SELECT COUNT(*) FROM pemasukan"
	err = tx.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total pemasukan: %v", err)
	}

	return pemasukanSlice, total, nil
}

// FindById implements PemasukanRepo.
func (s *pemasukanRepoImpl) FindById(ctx context.Context, tx *sql.Tx, id string) (model.Pemasukan, error) {
	query := "SELECT id_pemasukan, tanggal, kategori, keterangan, nominal, nota FROM pemasukan WHERE id_pemasukan = ?"
	row := tx.QueryRowContext(ctx, query, id)

	pemasukan := model.Pemasukan{}
	var tanggalRaw []byte

	err := row.Scan(&pemasukan.Id, &tanggalRaw, &pemasukan.Kategori, &pemasukan.Keterangan, &pemasukan.Nominal, &pemasukan.Nota)
	if err != nil {
		if err == sql.ErrNoRows {
			return pemasukan, fmt.Errorf("pemasukan tidak ditemukan")
		}
		return pemasukan, fmt.Errorf("gagal memindai data pemasukan: %v", err)
	}

	// Konversi tanggalRaw ke time.Time
	tanggalStr := string(tanggalRaw)
	parsedTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
	}
	pemasukan.Tanggal = parsedTime

	return pemasukan, nil
}

// DeletePemasukan implements PemasukanRepo.
func (s *pemasukanRepoImpl) DeletePemasukan(ctx context.Context, tx *sql.Tx, pemasukan model.Pemasukan) (model.Pemasukan, error) {
	// Validate input
	if pemasukan.Id == "" {
		return pemasukan, fmt.Errorf("id_pemasukan tidak boleh kosong")
	}

	// Fetch id_transaksi, nominal, and tanggal from pemasukan
	var idTransaksi string
	var nominal int
	var tanggalRaw []byte
	queryFetch := `
		SELECT id_transaksi, nominal, tanggal 
		FROM pemasukan 
		WHERE id_pemasukan = ?
	`
	err := tx.QueryRowContext(ctx, queryFetch, pemasukan.Id).Scan(&idTransaksi, &nominal, &tanggalRaw)
	if err != nil {
		if err == sql.ErrNoRows {
			return pemasukan, fmt.Errorf("pemasukan dengan id %s tidak ditemukan", pemasukan.Id)
		}
		return pemasukan, fmt.Errorf("gagal mengambil data pemasukan: %v", err)
	}

	// Konversi tanggalRaw ke time.Time
	tanggalStr := string(tanggalRaw)
	tanggalTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
	}

	// Validate nominal
	if nominal <= 0 {
		return pemasukan, fmt.Errorf("nilai nominal tidak valid: %d", nominal)
	}

	log.Printf("Berhasil mengambil pemasukan: id=%s, id_transaksi=%s, nominal=%d, tanggal=%v", pemasukan.Id, idTransaksi, nominal, tanggalTime)

	// Delete from laporan_keuangan
	queryLaporan := `
		DELETE FROM laporan_keuangan 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryLaporan, idTransaksi)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal menghapus dari laporan_keuangan: %v", err)
	}

	// Delete from history_transaksi
	queryHistory := `
		DELETE FROM history_transaksi 
		WHERE id_transaksi = ?
	`
	_, err = tx.ExecContext(ctx, queryHistory, idTransaksi)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal menghapus dari history_transaksi: %v", err)
	}

	// Delete from pemasukan
	queryPemasukan := `
		DELETE FROM pemasukan 
		WHERE id_pemasukan = ?
	`
	_, err = tx.ExecContext(ctx, queryPemasukan, pemasukan.Id)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal menghapus dari pemasukan: %v", err)
	}

	// Update saldo for all records after the deleted pemasukan's tanggal
	queryUpdateSaldo := `
		UPDATE laporan_keuangan
		SET saldo = saldo - ?
		WHERE tanggal > ?
	`
	result, err := tx.ExecContext(ctx, queryUpdateSaldo, nominal, tanggalTime)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memperbarui saldo di masa depan: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memeriksa jumlah baris yang diperbarui untuk saldo: %v", err)
	}
	log.Printf("Berhasil memperbarui %d baris di laporan_keuangan untuk saldo dengan id_pemasukan %s, nominal %d, tanggal %v", rowsAffected, pemasukan.Id, nominal, tanggalTime)

	// Update total pemasukan for all records after the deleted pemasukan's tanggal
	queryUpdatePemasukan := `
		UPDATE laporan_keuangan
		SET pemasukan = pemasukan - ?
		WHERE tanggal > ?
	`
	result, err = tx.ExecContext(ctx, queryUpdatePemasukan, nominal, tanggalTime)
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memperbarui pemasukan di masa depan: %v", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return pemasukan, fmt.Errorf("gagal memeriksa jumlah baris yang diperbarui untuk pemasukan: %v", err)
	}
	log.Printf("Berhasil memperbarui %d baris di laporan_keuangan untuk pemasukan dengan id_pemasukan %s, nominal %d, tanggal %v", rowsAffected, pemasukan.Id, nominal, tanggalTime)

	return pemasukan, nil
}

// GetPemasukanByDateRange implements PemasukanRepo.
func (s *pemasukanRepoImpl) GetPemasukanByDateRange(ctx context.Context, tx *sql.Tx, startDate, endDate string, page int, pageSize int) ([]model.Pemasukan, int, error) {
	// Hitung offset
	offset := (page - 1) * pageSize

	// Query untuk mendapatkan data dengan pagination dan date range
	query := `
		SELECT id_pemasukan, tanggal, kategori, keterangan, nominal, nota 
		FROM pemasukan 
		WHERE tanggal BETWEEN ? AND ? 
		ORDER BY tanggal DESC 
		LIMIT ? OFFSET ?
	`
	rows, err := tx.QueryContext(ctx, query, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data pemasukan berdasarkan rentang tanggal: %v", err)
	}
	defer rows.Close()

	var pemasukanSlice []model.Pemasukan
	for rows.Next() {
		pemasukan := model.Pemasukan{}
		var tanggalRaw []byte

		err := rows.Scan(&pemasukan.Id, &tanggalRaw, &pemasukan.Kategori, &pemasukan.Keterangan, &pemasukan.Nominal, &pemasukan.Nota)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal memindai data pemasukan: %v", err)
		}

		// Konversi tanggalRaw ke time.Time
		tanggalStr := string(tanggalRaw)
		parsedTime, err := time.Parse("2006-01-02 15:04:05", tanggalStr)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal mem-parsing tanggal: %v", err)
		}
		pemasukan.Tanggal = parsedTime

		pemasukanSlice = append(pemasukanSlice, pemasukan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("kesalahan setelah mengiterasi baris: %v", err)
	}

	// Query untuk mendapatkan total data dalam rentang tanggal
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM pemasukan 
		WHERE tanggal BETWEEN ? AND ?
	`
	err = tx.QueryRowContext(ctx, countQuery, startDate, endDate).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total pemasukan dalam rentang tanggal: %v", err)
	}

	return pemasukanSlice, total, nil
}