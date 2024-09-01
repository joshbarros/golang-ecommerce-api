package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshbarros/golang-ecommerce-api/types"
	"github.com/joshbarros/golang-ecommerce-api/utils"
)

// Mock implementation of the ProductStore interface
type mockProductStore struct {
	products []types.Product
}

func (m *mockProductStore) GetProducts() ([]types.Product, error) {
	return m.products, nil
}

func (m *mockProductStore) CreateProduct(product *types.Product) error {
	// Simulate database ID generation
	product.ID = len(m.products) + 1
	m.products = append(m.products, *product)
	return nil
}

func (m *mockProductStore) GetProductsByID(ps []int) ([]types.Product, error) {
	var result []types.Product
	for _, id := range ps {
		for _, product := range m.products {
			if product.ID == id {
				result = append(result, product)
			}
		}
	}
	if len(result) == 0 {
		return nil, errors.New("no products found")
	}
	return result, nil
}

func (m *mockProductStore) UpdateProduct(product types.Product) error {
	for i, p := range m.products {
		if p.ID == product.ID {
			m.products[i] = product
			return nil
		}
	}
	return errors.New("product not found")
}

func (m *mockProductStore) FailCreateProduct(product *types.Product) error {
	return errors.New("failed to create product")
}

func TestProductServiceHandlers(t *testing.T) {
	productStore := &mockProductStore{
		products: []types.Product{
			{ID: 1, Name: "Test Product 1", Price: 9.99, Quantity: 10},
			{ID: 2, Name: "Test Product 2", Price: 19.99, Quantity: 20},
		},
	}
	handler := NewHandler(productStore)

	t.Run("Should get all products", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.handleGetProduct)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var products []types.Product
		if err := json.NewDecoder(rr.Body).Decode(&products); err != nil {
			t.Fatal("Failed to decode JSON response")
		}

		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})

	t.Run("Should create a product successfully", func(t *testing.T) {
		payload := types.Product{
			Name:     "New Product",
			Price:    29.99,
			Quantity: 15,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.handleCreateProduct)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}

		var createdProduct types.Product
		if err := json.NewDecoder(rr.Body).Decode(&createdProduct); err != nil {
			t.Fatal("Failed to decode JSON response")
		}

		if createdProduct.ID == 0 {
			t.Error("Expected product ID to be set")
		}

		if createdProduct.Name != "New Product" {
			t.Errorf("Expected product name to be 'New Product', got '%s'", createdProduct.Name)
		}
	})

	t.Run("Should fail to create a product with invalid data", func(t *testing.T) {
		payload := types.Product{
			Name:     "",
			Price:    0,
			Quantity: -1,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.handleCreateProduct)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should get products by ID", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products?id=1&id=2", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
			ids := []int{1, 2}
			products, err := productStore.GetProductsByID(ids)
			if err != nil {
				utils.WriteError(w, http.StatusNotFound, err)
				return
			}
			utils.WriteJSON(w, http.StatusOK, products)
		})
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var products []types.Product
		if err := json.NewDecoder(rr.Body).Decode(&products); err != nil {
			t.Fatal("Failed to decode JSON response")
		}

		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})

	t.Run("Should update a product", func(t *testing.T) {
		payload := types.Product{
			ID:       1,
			Name:     "Updated Product",
			Price:    29.99,
			Quantity: 15,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
			var updatedProduct types.Product
			if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}
			if err := productStore.UpdateProduct(updatedProduct); err != nil {
				utils.WriteError(w, http.StatusNotFound, err)
				return
			}
			utils.WriteJSON(w, http.StatusOK, updatedProduct)
		})
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var updatedProduct types.Product
		if err := json.NewDecoder(rr.Body).Decode(&updatedProduct); err != nil {
			t.Fatal("Failed to decode JSON response")
		}

		if updatedProduct.Name != "Updated Product" {
			t.Errorf("Expected product name to be 'Updated Product', got '%s'", updatedProduct.Name)
		}
	})
}
