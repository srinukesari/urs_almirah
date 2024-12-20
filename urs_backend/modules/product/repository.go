package product

import (
	"database/sql"
	"log"
	pb "urs_backend/proto/protobuf"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetProductDetails(id int) (*Product, error) {
	query := "SELECT id, name, description, price, stock FROM products WHERE id = $1"
	row := r.DB.QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) ListProducts() ([]*Product, error) {
	query := "SELECT id, name, description, price, stock FROM products"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *ProductRepository) AddProduct(req *pb.AddProductRequest) error {
	query := `INSERT INTO products (name, price, description, stock) VALUES ($1, $2, $3, $4)`

	_, err := r.DB.Exec(query, req.Name, req.Price, req.Description, req.Stock)
	if err != nil {
		return err
	}

	return nil
}
