package httpapi

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"pos/backend/internal/model"
	"pos/backend/internal/store"
)

type Server struct {
	products      store.ProductStore
	adminUsername string
	adminPassword string
}

type productPayload struct {
	Name              string `json:"name"`
	Category          string `json:"category"`
	Description       string `json:"description"`
	Price             int64  `json:"price"`
	ImageURL          string `json:"imageUrl"`
	LowStockThreshold int    `json:"lowStockThreshold"`
	IsActive          bool   `json:"isActive"`
}

type stockPayload struct {
	Delta  int    `json:"delta"`
	Reason string `json:"reason"`
}

func NewHandler(products store.ProductStore, adminUsername, adminPassword string) http.Handler {
	server := &Server{products: products, adminUsername: adminUsername, adminPassword: adminPassword}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", server.health)
	mux.HandleFunc("/api/products", server.listCashierProducts)
	mux.Handle("/api/admin/products", server.requireAdmin(http.HandlerFunc(server.adminProducts)))
	mux.Handle("/api/admin/products/", server.requireAdmin(http.HandlerFunc(server.adminProductAction)))
	return mux
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) listCashierProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	products, err := s.products.ListProducts(r.Context(), false)
	if err != nil {
		log.Printf("list cashier products: %v", err)
		writeError(w, http.StatusInternalServerError, "Gagal memuat menu")
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (s *Server) adminProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		products, err := s.products.ListProducts(r.Context(), true)
		if err != nil {
			log.Printf("list admin products: %v", err)
			writeError(w, http.StatusInternalServerError, "Gagal memuat daftar produk")
			return
		}
		writeJSON(w, http.StatusOK, products)
	case http.MethodPost:
		payload, ok := decodeProductPayload(w, r)
		if !ok {
			return
		}
		product, err := payload.product()
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		product.ID, err = newProductID()
		if err != nil {
			log.Printf("generate product id: %v", err)
			writeError(w, http.StatusInternalServerError, "Gagal membuat produk")
			return
		}
		if err := s.products.CreateProduct(r.Context(), product); err != nil {
			log.Printf("create product: %v", err)
			writeError(w, http.StatusInternalServerError, "Gagal menyimpan produk")
			return
		}
		created, err := s.products.GetProduct(r.Context(), product.ID)
		if err != nil {
			log.Printf("read created product: %v", err)
			writeError(w, http.StatusInternalServerError, "Produk tersimpan tetapi gagal dimuat ulang")
			return
		}
		writeJSON(w, http.StatusCreated, created)
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) adminProductAction(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/products/"), "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	id := parts[0]

	if len(parts) == 2 && parts[1] == "stock" {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var payload stockPayload
		if err := decodeJSON(w, r, &payload); err != nil {
			writeError(w, http.StatusBadRequest, "Format penyesuaian stok tidak valid")
			return
		}
		payload.Reason = strings.TrimSpace(payload.Reason)
		if payload.Delta == 0 || payload.Reason == "" || utf8.RuneCountInString(payload.Reason) > 200 {
			writeError(w, http.StatusBadRequest, "Jumlah perubahan harus bukan nol dan alasan maksimal 200 karakter")
			return
		}
		product, err := s.products.AdjustStock(r.Context(), id, payload.Delta, payload.Reason, r.Header.Get("X-Admin-User"))
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, product)
		return
	}
	if len(parts) != 1 {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodPut:
		payload, ok := decodeProductPayload(w, r)
		if !ok {
			return
		}
		product, err := payload.product()
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		product.ID = id
		if err := s.products.UpdateProduct(r.Context(), product); err != nil {
			writeStoreError(w, err)
			return
		}
		updated, err := s.products.GetProduct(r.Context(), id)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		if err := s.products.SetProductActive(r.Context(), id, false); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "archived"})
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.adminUsername == "" || s.adminPassword == "" {
			writeError(w, http.StatusServiceUnavailable, "ADMIN_USERNAME dan ADMIN_PASSWORD belum dikonfigurasi")
			return
		}
		username, password, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(username), []byte(s.adminUsername)) != 1 || subtle.ConstantTimeCompare([]byte(password), []byte(s.adminPassword)) != 1 {
			writeError(w, http.StatusUnauthorized, "Username atau password admin salah")
			return
		}
		r.Header.Set("X-Admin-User", username)
		next.ServeHTTP(w, r)
	})
}

func decodeProductPayload(w http.ResponseWriter, r *http.Request) (productPayload, bool) {
	var payload productPayload
	if err := decodeJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "Format data produk tidak valid")
		return productPayload{}, false
	}
	return payload, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, value interface{}) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("request must contain one JSON value")
	}
	return nil
}

func (p productPayload) product() (model.Product, error) {
	p.Name = strings.TrimSpace(p.Name)
	p.Category = strings.TrimSpace(p.Category)
	p.Description = strings.TrimSpace(p.Description)
	p.ImageURL = strings.TrimSpace(p.ImageURL)
	if p.Name == "" || utf8.RuneCountInString(p.Name) > 120 {
		return model.Product{}, errors.New("Nama wajib diisi dan maksimal 120 karakter")
	}
	if p.Category == "" || utf8.RuneCountInString(p.Category) > 60 {
		return model.Product{}, errors.New("Kategori wajib diisi dan maksimal 60 karakter")
	}
	if utf8.RuneCountInString(p.Description) > 300 {
		return model.Product{}, errors.New("Deskripsi maksimal 300 karakter")
	}
	if p.Price <= 0 {
		return model.Product{}, errors.New("Harga harus lebih besar dari nol")
	}
	if p.LowStockThreshold < 0 {
		return model.Product{}, errors.New("Batas stok minimum tidak boleh negatif")
	}
	if p.ImageURL != "" {
		imageURL, err := url.ParseRequestURI(p.ImageURL)
		if err != nil || imageURL.Scheme != "https" || imageURL.Host == "" {
			return model.Product{}, errors.New("URL gambar harus menggunakan HTTPS")
		}
	}
	return model.Product{
		Name: p.Name, Category: p.Category, Description: p.Description,
		Price: p.Price, ImageURL: p.ImageURL,
		LowStockThreshold: p.LowStockThreshold, IsActive: p.IsActive,
	}, nil
}

func newProductID() (string, error) {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return "product-" + hex.EncodeToString(bytes[:]), nil
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrProductNotFound):
		writeError(w, http.StatusNotFound, "Produk tidak ditemukan")
	case errors.Is(err, store.ErrInsufficientStock):
		writeError(w, http.StatusConflict, "Stok tidak mencukupi untuk pengurangan ini")
	case errors.Is(err, store.ErrInvalidStockChange):
		writeError(w, http.StatusBadRequest, "Perubahan stok tidak valid")
	default:
		log.Printf("admin product operation: %v", err)
		writeError(w, http.StatusInternalServerError, "Operasi produk gagal")
	}
}

func methodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode JSON response: %v", err)
	}
}
