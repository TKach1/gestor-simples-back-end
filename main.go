package main

import (
	"database/sql"
	"encoding/json"
	"gestor-simples-ecs/internal/database"
	"gestor-simples-ecs/internal/models"
	"gestor-simples-ecs/pkg/auth"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// --- Main Application Setup ---

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	
	// Initialize packages
	database.Connect()
	auth.Initialize()

	// Set up router
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// Authentication routes
	api.HandleFunc("/auth/login", loginHandler).Methods("POST")
	api.HandleFunc("/auth/register", registerHandler).Methods("POST")

	// User routes
	userRouter := api.PathPrefix("/users").Subrouter()
	userRouter.Use(auth.AuthMiddleware) // Protect all user routes
	userRouter.HandleFunc("", adminOnly(getUsersHandler)).Methods("GET")
	userRouter.HandleFunc("", adminOnly(createUserHandler)).Methods("POST")
	userRouter.HandleFunc("/{id}", getUserHandler).Methods("GET")
	userRouter.HandleFunc("/{id}", updateUserHandler).Methods("PUT")
	userRouter.HandleFunc("/{id}", adminOnly(deleteUserHandler)).Methods("DELETE")
	
	// Product routes
	productRouter := api.PathPrefix("/products").Subrouter()
	productRouter.Use(auth.AuthMiddleware)
	productRouter.HandleFunc("", getProductsHandler).Methods("GET")
	productRouter.HandleFunc("", adminOnly(createProductHandler)).Methods("POST")
	productRouter.HandleFunc("/{id}", getProductHandler).Methods("GET")
	productRouter.HandleFunc("/{id}", adminOnly(updateProductHandler)).Methods("PUT")
	productRouter.HandleFunc("/{id}", adminOnly(deleteProductHandler)).Methods("DELETE")

	// Sales routes
	salesRouter := api.PathPrefix("/sales").Subrouter()
	salesRouter.Use(auth.AuthMiddleware)
	salesRouter.HandleFunc("", getSalesHandler).Methods("GET")
	salesRouter.HandleFunc("", createSaleHandler).Methods("POST")

	// Dashboard routes
	dashboardRouter := api.PathPrefix("/dashboard").Subrouter()
	dashboardRouter.Use(auth.AuthMiddleware)
	dashboardRouter.HandleFunc("/summary", getDashboardSummaryHandler).Methods("GET")

	// Start server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

// --- Helper Functions ---

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// adminOnly is a convenience function to chain the AdminMiddleware.
func adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth.AdminMiddleware(h).ServeHTTP(w, r)
	})
}


// --- Handlers ---

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user models.User
	err := database.DB.QueryRow("SELECT id, name, username, password_hash, role FROM users WHERE username = $1", req.Username).Scan(&user.ID, &user.Name, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"token": token, "user": user})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	var userID int64
	// Default role is 'vendedor'
	err = database.DB.QueryRow(
		"INSERT INTO users (name, username, password_hash, role) VALUES ($1, $2, $3, 'vendedor') RETURNING id",
		req.Name, req.Username, hashedPassword,
	).Scan(&userID)

	if err != nil {
		// You might want to check for specific errors, like duplicate username
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	user := models.User{
		ID:       userID,
		Name:     req.Name,
		Username: req.Username,
		Role:     "vendedor",
	}

	respondWithJSON(w, http.StatusCreated, user)
}


// --- User Handlers ---

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
        return
    }

    var userID int64
    err = database.DB.QueryRow(
        "INSERT INTO users (name, username, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id",
        req.Name, req.Username, hashedPassword, req.Role,
    ).Scan(&userID)

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create user")
        return
    }

    user := models.User{
        ID:       userID,
        Name:     req.Name,
        Username: req.Username,
        Role:     req.Role,
    }

    respondWithJSON(w, http.StatusCreated, user)
}
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, username, role FROM users")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to query users")
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Role); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to scan user")
			return
		}
		users = append(users, user)
	}

	respondWithJSON(w, http.StatusOK, users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	err := database.DB.QueryRow("SELECT id, name, username, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Username, &user.Role)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// For simplicity, we assume all fields are provided for update.
	_, err := database.DB.Exec("UPDATE users SET name = $1, username = $2, role = $3 WHERE id = $4", user.Name, user.Username, user.Role, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := database.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	count, err := res.RowsAffected()
	if err != nil || count == 0 {
		respondWithError(w, http.StatusNotFound, "User not found or already deleted")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Product Handlers ---

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, description, price, quantity FROM products")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to query products")
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to scan product")
			return
		}
		products = append(products, p)
	}

	respondWithJSON(w, http.StatusOK, products)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO products (name, description, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id",
		p.Name, p.Description, p.Price, p.Quantity,
	).Scan(&p.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p models.Product
	err := database.DB.QueryRow("SELECT id, name, description, price, quantity FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_, err := database.DB.Exec(
		"UPDATE products SET name = $1, description = $2, price = $3, quantity = $4 WHERE id = $5",
		p.Name, p.Description, p.Price, p.Quantity, id,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found or failed to delete")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


// --- Sales Handlers ---
func getSalesHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			s.id, s.user_id, s.date,
			si.product_id, si.quantity,
			p.name, p.price
		FROM sales s
		LEFT JOIN sales_items si ON s.id = si.sale_id
		LEFT JOIN products p ON si.product_id = p.id
		ORDER BY s.date DESC;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to query sales")
		return
	}
	defer rows.Close()

	salesMap := make(map[int64]*models.Sale)
	for rows.Next() {
		var (
			saleID       int64
			userID       int64
			saleDate     time.Time
			productID    sql.NullInt64 // Use sql.Null types for LEFT JOIN
			quantity     sql.NullInt32
			productName  sql.NullString
			productPrice sql.NullFloat64
		)

		if err := rows.Scan(&saleID, &userID, &saleDate, &productID, &quantity, &productName, &productPrice); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to scan sale data")
			return
		}

		// Check if the sale is already in the map
		sale, ok := salesMap[saleID]
		if !ok {
			sale = &models.Sale{
				ID:         saleID,
				UserID:     userID,
				Date:       saleDate,
				Items:      []models.SaleItem{},
				TotalPrice: 0,
			}
			salesMap[saleID] = sale
		}

		// Add item if it exists
		if productID.Valid {
			item := models.SaleItem{
				ProductID:   productID.Int64,
				ProductName: productName.String,
				Quantity:    int(quantity.Int32),
				UnitPrice:   productPrice.Float64,
			}
			sale.Items = append(sale.Items, item)
			sale.TotalPrice += float64(item.Quantity) * item.UnitPrice
		}
	}

	// Convert map to slice
	sales := make([]models.Sale, 0, len(salesMap))
	for _, sale := range salesMap {
		sales = append(sales, *sale)
	}
    // As salesMap doesn't keep the order, we should sort it again
    sort.Slice(sales, func(i, j int) bool {
        return sales[i].Date.After(sales[j].Date)
    })

	respondWithJSON(w, http.StatusOK, sales)
}

func createSaleHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	// Create the sale record
	var saleID int64
	err = tx.QueryRow("INSERT INTO sales (user_id, date) VALUES ($1, NOW()) RETURNING id", req.UserID).Scan(&saleID)
	if err != nil {
		tx.Rollback()
		respondWithError(w, http.StatusInternalServerError, "Failed to create sale record")
		return
	}

	// Loop through items, update stock, and insert into sales_items
	for _, item := range req.Items {
		// Decrease product quantity
		res, err := tx.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2 AND quantity >= $1", item.Quantity, item.ProductID)
		if err != nil {
			tx.Rollback()
			respondWithError(w, http.StatusInternalServerError, "Failed to update product stock")
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil || rowsAffected == 0 {
			tx.Rollback()
			respondWithError(w, http.StatusBadRequest, "Insufficient stock or product not found")
			return
		}
		// Insert into sales_items
		_, err = tx.Exec("INSERT INTO sales_items (sale_id, product_id, quantity) VALUES ($1, $2, $3)", saleID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			respondWithError(w, http.StatusInternalServerError, "Failed to record sale item")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int64{"saleId": saleID})
}

// --- Dashboard Handlers ---
func getDashboardSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simplified dashboard. A real one would have more complex logic
	// and would check the user's role from the JWT token.
	var totalSalesMonth float64
	database.DB.QueryRow("SELECT COALESCE(SUM(p.price * si.quantity), 0) FROM sales s JOIN sales_items si ON s.id = si.sale_id JOIN products p ON si.product_id = p.id WHERE s.date >= date_trunc('month', current_date)").Scan(&totalSalesMonth)

	var totalSellers int
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'vendedor'").Scan(&totalSellers)
	
	var lowStockProducts int
	database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE quantity < 10").Scan(&lowStockProducts)

	summary := map[string]interface{}{
		"totalSalesMonth":   totalSalesMonth,
		"totalSellers":      totalSellers,
		"lowStockProducts":  lowStockProducts,
	}

	respondWithJSON(w, http.StatusOK, summary)
}
