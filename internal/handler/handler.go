package handler

import (
	"context"
	"diplom_back/config"
	"diplom_back/internal/auth"
	"diplom_back/internal/handler/controllers"
	"diplom_back/internal/handler/controllers/storeHandler"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func Setup(cfg *config.Config, ctx context.Context) http.Handler {
	muxR := mux.NewRouter()
	db := cfg.Client

	// ✅ CORS на весь роутер
	muxR.Use(corsMiddleware)

	// ✅ Логи на весь роутер
	muxR.Use(loggingMiddleware)

	// ✅ catch-all для preflight OPTIONS (иначе mux может отдать 404/405)
	muxR.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Auth
	muxR.HandleFunc("/api/login", controllers.LoginHandler(ctx, db)).Methods("POST")

	// Cleaning
	muxR.HandleFunc(GetAllCleaning, auth.RequireAuth(controllers.GetAllCleaningHandler(ctx, db))).Methods("GET")
	muxR.HandleFunc(DeleteCleaningById, auth.RequireAuth(controllers.DeleteCleaningByIdHandler(ctx, db))).Methods("DELETE")
	muxR.HandleFunc(PostNewCleaning, controllers.PostNewCleaningHandler(ctx, db)).Methods("POST")

	// Employee
	muxR.HandleFunc(GetAllEmployee, auth.RequireAuth(controllers.GetAllEmployeeHandler(ctx, db))).Methods("GET")
	muxR.HandleFunc(DeleteEmployeeById, auth.RequireAuth(controllers.DeleteEmployeeByIdHandler(ctx, db))).Methods("DELETE")
	muxR.HandleFunc(PostNewEmployee, auth.RequireAuth(controllers.PostNewEmployeeHandler(ctx, db))).Methods("POST")

	// Contact Employee
	muxR.HandleFunc(GetAllEmployeeContacts, auth.RequireAuth(controllers.GetAllEmployeeContactsHandler(ctx, db))).Methods("GET")
	muxR.HandleFunc(PostNewEmployeeContacts, controllers.PostNewEmployeeContactsHandler(ctx, db)).Methods("POST")
	muxR.HandleFunc(DeleteEmployeeContactsById, auth.RequireAuth(controllers.DeleteEmployeeContactsByIdHandler(ctx, db))).Methods("DELETE")

	// STORE - Dishes
	muxR.HandleFunc(GetAllDishes, storeHandler.GetAllDishesHandler(ctx, db)).Methods("GET")
	muxR.HandleFunc(PostNewDishes, storeHandler.PostNewDishesHandler(ctx, db)).Methods("POST")
	muxR.HandleFunc(DeleteDishesById, storeHandler.DeleteDishByIdHandler(ctx, db)).Methods("DELETE")

	// STORE - Orders
	muxR.HandleFunc(GetAllOrders, storeHandler.GetAllOrdersHandler(ctx, db)).Methods("GET")
	muxR.HandleFunc(GetOrderByID, storeHandler.GetOrderByIDHandler(ctx, db)).Methods("GET")
	muxR.HandleFunc(PostNewOrder, storeHandler.PostNewOrderHandler(ctx, db)).Methods("POST")
	muxR.HandleFunc(DeleteOrder, storeHandler.DeleteOrderHandler(ctx, db)).Methods("DELETE")

	// STORE - OrderItem
	muxR.HandleFunc(GetOrderItems, storeHandler.GetOrderItemsByOrderIDHandler(ctx, db)).Methods("GET")
	muxR.HandleFunc(UpdateOrderItem, storeHandler.UpdateOrderItemQuantityHandler(ctx, db)).Methods("PUT")
	muxR.HandleFunc(DeleteOrderItem, storeHandler.DeleteOrderItemHandler(ctx, db)).Methods("DELETE")

	return muxR
}

// ========================
// Middlewares
// ========================

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// не логируем OPTIONS (preflight)
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		ip := r.Header.Get("X-Forwarded-For")
		userAgent := r.Header.Get("User-Agent")

		slog.Info(fmt.Sprintf(
			"IP: %s, Method: %s, Route: %s, Query: %s, UserAgent: %s, AuthHeader: %s",
			ip, r.Method, r.URL.Path, r.URL.Query(), userAgent, r.Header.Get("Authorization"),
		))

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowed := map[string]bool{
			"http://localhost:5173":  true, // Vite
			"http://localhost:63342": true, // WebStorm preview (если нужно)
		}

		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Если не используешь cookies/sessions — credentials не нужно
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// ✅ preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
