package handler

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx" // Or any other DB library you are using, like GORM
)

// ProductServiceHandler implements the ProductServiceServer interface
type ProductServiceHandler struct {
	db *sqlx.DB // DB connection
}

// NewProductServiceHandler creates a new ProductServiceHandler with DB connection
func NewProductServiceHandler(db *sqlx.DB) *ProductServiceHandler {
	return &ProductServiceHandler{
		db: db,
	}
}

// GetProductDetails retrieves the details of a specific product by its ID
func (h *ProductServiceHandler) GetProductDetails(ctx context.Context, req *product.ProductRequest) (*product.ProductDetails, error) {
	// Use the DB connection to fetch product details
	var productDetails product.ProductDetails
	query := "SELECT product_id, name, description, price, stock FROM products WHERE product_id = ?"
	err := h.db.Get(&productDetails, query, req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("Product with ID %d not found: %v", req.ProductId, err)
	}

	return &productDetails, nil
}

// ListProducts retrieves a list of all available products
func (h *ProductServiceHandler) ListProducts(ctx context.Context, req *product.Empty) (*product.ProductList, error) {
	// Use the DB connection to fetch all products
	var products []product.ProductDetails
	query := "SELECT product_id, name, description, price, stock FROM products"
	err := h.db.Select(&products, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch products: %v", err)
	}

	return &product.ProductList{
		Products: products,
	}, nil
}
