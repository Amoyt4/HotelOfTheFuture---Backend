package controllers

import (
	"context"
	"diplom_back/internal/entity"
	"diplom_back/internal/repository"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllCleaningHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cleanings, err := repository.GetAllCleanings(ctx, db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cleanings); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error encoding response"))
			return
		}
	}
}

func PostNewCleaningHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}
		var cleaning entity.NewClean
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&cleaning); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("JSON не соответствует entity на бэке"))
			return
		}
		defer r.Body.Close()

		//проверю что поля пришли не пустыми
		if cleaning.RoomNum == 0 || cleaning.StartTime.IsZero() || cleaning.EndTime.IsZero() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("все поля должны быть заполнены"))
			return
		}

		id, err := repository.PostNewCleaning(ctx, db, &cleaning)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка добавления в бд: " + err.Error()))
			return
		}
		slog.Info("Успешно добавлено сущность в cleaning")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%d - успешно созданная запись убоки", id)))
		slog.Info("Уборка успешно добавлена")
	}

}

func DeleteCleaningByIdHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request entity.CleanDel

		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		//проверка что пришёл не пустой Id
		if request.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пришёл пустой id"))
			return
		}

		err := repository.DeleteCleaning(ctx, db, request.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cleaning deleted successfully"))

	}
}
