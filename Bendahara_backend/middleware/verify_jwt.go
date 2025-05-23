package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
)

func VerifyJWT(next httprouter.Handle) httprouter.Handle {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	jwtKey := os.Getenv("JWT_SECRET")
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.WriteJSONError(w, http.StatusUnauthorized, "header otorisasi tidak ditemukan")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			helper.WriteJSONError(w, http.StatusUnauthorized, "format token tidak valid")
			return
		}

		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			helper.WriteJSONError(w, http.StatusUnauthorized, "token tidak valid atau kadaluarsa")
			return
		}

		r.Header.Set("User-Name", claims.Username)
		next(w, r, ps)
	}
}