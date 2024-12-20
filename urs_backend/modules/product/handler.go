package product

import (
	"context"
	"fmt"

	pb "urs_backend/proto/protobuf"
)

// ProductService implements the ProductServiceServer interface from the generated proto
type ProductService struct {
	Repo *ProductRepository
	pb.UnimplementedProductServiceServer
}

func NewProductService(repo *ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

// GetProductDetails retrieves a product by its ID
func (s *ProductService) GetProductDetails(ctx context.Context, req *pb.ProductDetailsRequest) (*pb.ProductDetailsResponse, error) {
	fmt.Printf("Comes inside GetProductDetails for product ID: %d\n", req.Id)

	result, err := s.Repo.GetProductDetails(int(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.ProductDetailsResponse{
		ProductId:   int32(result.ID),
		Name:        result.Name,
		Description: result.Description,
		Price:       float32(result.Price),
		Stock:       int32(result.Stock),
	}, nil
}

// ListProducts returns a list of products
func (s *ProductService) ListProducts(ctx context.Context, req *pb.Empty) (*pb.ProductListResponse, error) {
	result, err := s.Repo.ListProducts()
	if err != nil {
		return nil, err
	}

	var productList []*pb.ProductDetailsResponse

	for _, product := range result {
		productList = append(productList, &pb.ProductDetailsResponse{
			ProductId:   int32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			Stock:       int32(product.Stock),
		})
	}

	return &pb.ProductListResponse{
		Products: productList,
	}, nil
}

// AddProduct add product to db
func (s *ProductService) AddProduct(ctx context.Context, req *pb.AddProductRequest,
) (*pb.AddProductResponse, error) {
	fmt.Printf("Comes inside AddProduct for product ID: %s\n", req.Name)

	err := s.Repo.AddProduct(req)
	if err != nil {
		return nil, err
	}

	return &pb.AddProductResponse{
		Status: 200,
		Error:  "",
	}, nil
}
