package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbarros/golang-ecommerce-api/types"
	"github.com/joshbarros/golang-ecommerce-api/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Clear ID to prevent conflicts, assuming ID is auto-incremented in the database
	product.ID = 0

	// Validate the product input
	if err := validateProduct(product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Create the product in the store
	if err := h.store.CreateProduct(&product); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return the created product, including the generated ID
	utils.WriteJSON(w, http.StatusCreated, product)
}

// validateProduct checks if the product fields are valid
func validateProduct(product types.Product) error {
	if product.Name == "" {
		return fmt.Errorf("Product name is required")
	}
	if product.Price <= 0 {
		return fmt.Errorf("Product price must be greater than zero")
	}
	if product.Quantity < 0 {
		return fmt.Errorf("Product quantity cannot be negative")
	}
	return nil
}
