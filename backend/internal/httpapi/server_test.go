package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pos/backend/internal/model"
	"pos/backend/internal/store"
)

type memoryProductStore struct {
	products map[string]model.Product
}

func newMemoryProductStore() *memoryProductStore {
	return &memoryProductStore{products: map[string]model.Product{
		"food-001": {ID: "food-001", Name: "Nasi Goreng", Category: "Makanan", Price: 28000, StockQuantity: 4, LowStockThreshold: 5, IsActive: true, Available: true},
	}}
}

func (m *memoryProductStore) ListProducts(_ context.Context, includeInactive bool) ([]model.Product, error) {
	products := make([]model.Product, 0, len(m.products))
	for _, product := range m.products {
		if includeInactive || product.IsActive {
			products = append(products, product)
		}
	}
	return products, nil
}

func (m *memoryProductStore) GetProduct(_ context.Context, id string) (model.Product, error) {
	product, ok := m.products[id]
	if !ok {
		return model.Product{}, store.ErrProductNotFound
	}
	return product, nil
}

func (m *memoryProductStore) CreateProduct(_ context.Context, product model.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *memoryProductStore) UpdateProduct(_ context.Context, product model.Product) error {
	current, ok := m.products[product.ID]
	if !ok {
		return store.ErrProductNotFound
	}
	product.StockQuantity = current.StockQuantity
	product.Available = product.IsActive && product.StockQuantity > 0
	m.products[product.ID] = product
	return nil
}

func (m *memoryProductStore) SetProductActive(_ context.Context, id string, active bool) error {
	product, ok := m.products[id]
	if !ok {
		return store.ErrProductNotFound
	}
	product.IsActive = active
	product.Available = active && product.StockQuantity > 0
	m.products[id] = product
	return nil
}

func (m *memoryProductStore) AdjustStock(_ context.Context, id string, delta int, _ string, _ string) (model.Product, error) {
	product, ok := m.products[id]
	if !ok {
		return model.Product{}, store.ErrProductNotFound
	}
	if delta == 0 {
		return model.Product{}, store.ErrInvalidStockChange
	}
	if product.StockQuantity+delta < 0 {
		return model.Product{}, store.ErrInsufficientStock
	}
	product.StockQuantity += delta
	product.Available = product.IsActive && product.StockQuantity > 0
	m.products[id] = product
	return product, nil
}

func TestAdminProductsRequiresCredentials(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/admin/products", nil)

	NewHandler(newMemoryProductStore(), "admin", "secret").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestAdminCanCreateProductAndAdjustStock(t *testing.T) {
	products := newMemoryProductStore()
	handler := NewHandler(products, "admin", "secret")

	createRecorder := httptest.NewRecorder()
	createRequest := httptest.NewRequest(http.MethodPost, "/api/admin/products", strings.NewReader(`{"name":"Es Teh","category":"Minuman","description":"Teh dingin","price":8000,"imageUrl":"https://images.example.com/tea.jpg","lowStockThreshold":3,"isActive":true}`))
	createRequest.SetBasicAuth("admin", "secret")
	handler.ServeHTTP(createRecorder, createRequest)

	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, createRecorder.Code, createRecorder.Body.String())
	}
	var created model.Product
	if err := json.NewDecoder(createRecorder.Body).Decode(&created); err != nil {
		t.Fatalf("decode created product: %v", err)
	}
	if created.StockQuantity != 0 {
		t.Fatalf("new products should start with zero stock, got %d", created.StockQuantity)
	}

	stockRecorder := httptest.NewRecorder()
	stockRequest := httptest.NewRequest(http.MethodPost, "/api/admin/products/"+created.ID+"/stock", strings.NewReader(`{"delta":12,"reason":"Stok awal"}`))
	stockRequest.SetBasicAuth("admin", "secret")
	handler.ServeHTTP(stockRecorder, stockRequest)

	if stockRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, stockRecorder.Code, stockRecorder.Body.String())
	}
	var adjusted model.Product
	if err := json.NewDecoder(stockRecorder.Body).Decode(&adjusted); err != nil {
		t.Fatalf("decode adjusted product: %v", err)
	}
	if adjusted.StockQuantity != 12 || !adjusted.Available {
		t.Fatalf("expected 12 available units, got stock=%d available=%t", adjusted.StockQuantity, adjusted.Available)
	}
}

func TestAdminStockCannotBecomeNegative(t *testing.T) {
	handler := NewHandler(newMemoryProductStore(), "admin", "secret")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/admin/products/food-001/stock", strings.NewReader(`{"delta":-5,"reason":"Koreksi"}`))
	request.SetBasicAuth("admin", "secret")

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
	}
}
