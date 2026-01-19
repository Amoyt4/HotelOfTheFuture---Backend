package controllers

import (
	"context"
	"diplom_back/config"
	"diplom_back/internal/entity"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var env config.Env
var secret = env.JWT_SECRET

func LoginHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		var employee entity.Employee
		err := db.QueryRow(ctx,
			"SELECT id, login, password, name FROM employee WHERE login=$1",
			req.Login,
		).Scan(&employee.ID, &employee.Login, &employee.Password, &employee.Name)
		if err != nil || employee.Password != req.Password { // для начала без хэша
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{
			"employee_id": employee.ID,
			"exp":         time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(secret)
		if err != nil {
			http.Error(w, "token generation error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"token": tokenStr,
			"name":  employee.Name,
		})
	}
}
