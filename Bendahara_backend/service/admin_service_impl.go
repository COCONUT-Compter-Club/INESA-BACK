package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
	"golang.org/x/crypto/bcrypt"
)

type adminServiceImpl struct {
	AdminRepo repository.AdminRepo
	DB        *sql.DB
}

func NewAdminService(adminRepo repository.AdminRepo, db *sql.DB) AdminService {
	return adminServiceImpl{
		AdminRepo: adminRepo,
		DB:        db,
	}
}

// Function for hash password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Function to verify password
func verifyPassword(storedHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

// SignUp implements AdminService.
func (a adminServiceImpl) SignUp(ctx context.Context, adminRequest dto.AdminRequest) (dto.AdminResponse, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Commit()

	hassedPass, err := hashPassword(adminRequest.Password)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("gagal membuat kata sandi hash: %v", err)
	}

	admin := model.Admin{
		Id:       uuid.New().String(),
		Nik:      adminRequest.Nik,
		Username: adminRequest.Username,
		Password: hassedPass,
		Role:     adminRequest.Role,
	}

	createAdmin, err := a.AdminRepo.SignUp(ctx, tx, admin)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("gagal melakukan pendaftaran admin: %v", err)
	}

	return util.ConvertAdminToResponseDTO(createAdmin), nil
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (a adminServiceImpl) GenerateJWT(username string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	jwtKey := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtKey))
}

// SignIn implements AdminService.
func (a adminServiceImpl) SignIn(ctx context.Context, loginRequest dto.LoginRequest) (string, error) {
	if loginRequest.Nik == "" || loginRequest.Password == "" {
		return "", fmt.Errorf("nik atau kata sandi tidak boleh kosong")
	}
	tx, err := a.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Commit()

	admin, err := a.AdminRepo.FindByNik(ctx, tx, loginRequest.Nik)
	if err != nil {
		return "", fmt.Errorf("nik atau kata sandi tidak valid: %v", err)
	}

	if verifyPassword(admin.Password, loginRequest.Password) {
		fmt.Println("Login berhasil!")
	} else {
		return "", fmt.Errorf("nik atau kata sandi tidak valid")
	}

	token, err := a.GenerateJWT(loginRequest.Nik)
	if err != nil {
		return "", fmt.Errorf("gagal menghasilkan token: %v", err)
		// panic(err)
	}

	return token, nil
}

// GetAdminByNik implements AdminService.
func (a adminServiceImpl) GetAdminByNik(ctx context.Context, nik string) (dto.AdminResponse, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer util.CommitOrRollBack(tx)

	admin, err := a.AdminRepo.FindByNik(ctx, tx, nik)
	if err != nil {
		return dto.AdminResponse{}, fmt.Errorf("admin dengan nik %s tidak ditemukan: %v", nik, err)
	}

	return util.ConvertAdminToResponseDTO(admin), nil
}
