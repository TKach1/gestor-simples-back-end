package models

import "time"

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Never expose this
	Role         string `json:"role"`
}

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type Sale struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userId"`
	Date       time.Time `json:"date"`
	Items      []SaleItem `json:"items"`
	TotalPrice float64   `json:"totalPrice"`
}

type SaleItem struct {
	ProductID   int64   `json:"productId"`
	ProductName string  `json:"productName,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice,omitempty"`
}

// Payloads for requests

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateSaleRequest struct {
	UserID int64      `json:"userId"`
	Items  []SaleItem `json:"items"`
}
