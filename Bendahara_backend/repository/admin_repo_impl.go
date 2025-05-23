package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type adminRepoImpl struct {
}

func NewAdminRepo() AdminRepo {
	return &adminRepoImpl{}
}

// SignUp implements AdminRepo.
func (a adminRepoImpl) SignUp(ctx context.Context, tx *sql.Tx, admin model.Admin) (model.Admin, error) {
	query := "INSERT INTO Admin (id, nik, username, password, role) VALUES (?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, admin.Id, admin.Nik, admin.Username, admin.Password, admin.Role)
	if err != nil {
		return admin, err
	}

	return admin, nil
}

// FindById implements AdminRepo.
func (a adminRepoImpl) FindByNik(ctx context.Context, tx *sql.Tx, nik string) (model.Admin, error) {
	query := "SELECT id, nik, username, password, role FROM Admin WHERE nik = ?"

	rows, err := tx.QueryContext(ctx, query, nik)
	if err != nil {
		return model.Admin{}, err
	}

	defer rows.Close()
	admin := model.Admin{}
	if rows.Next() {
		err := rows.Scan(&admin.Id, &admin.Nik, &admin.Username, &admin.Password, &admin.Role)
		if err != nil {
			return model.Admin{}, err
		}
		return admin, nil
	} else {
		return admin, errors.New("nik user tidak ditemukan")
	}
}
