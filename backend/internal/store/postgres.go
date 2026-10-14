package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sort"

	_ "github.com/lib/pq"

	"pos/backend/internal/model"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type PostgresStore struct {
	db *sql.DB
}

const productColumns = `id, name, category, description, price, image_url, stock_quantity, low_stock_threshold, is_active, (is_active AND stock_quantity > 0)`

func OpenPostgres(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	store := &PostgresStore{db: db}
	if err := store.migrate(ctx); err != nil {
		db.Close()
		return nil, err
	}
	if err := store.seedSampleProducts(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) migrate(ctx context.Context) error {
	entries, err := fs.Glob(migrationFiles, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(entries)
	for _, name := range entries {
		contents, err := migrationFiles.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := s.db.ExecContext(ctx, string(contents)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
	}
	return nil
}

func (s *PostgresStore) seedSampleProducts(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&count); err != nil {
		return fmt.Errorf("count products: %w", err)
	}
	if count > 0 {
		return nil
	}

	products := []model.Product{
		{ID: "food-001", Name: "Nasi Goreng Spesial", Category: "Makanan", Description: "Nasi goreng dengan telur dan ayam", Price: 28000, ImageURL: "https://images.unsplash.com/photo-1741243413074-f87bd31f4971?auto=format&fit=crop&w=720&h=520&q=85", StockQuantity: 20, LowStockThreshold: 5, IsActive: true},
		{ID: "food-002", Name: "Ayam Bakar Madu", Category: "Makanan", Description: "Ayam bakar dengan bumbu madu", Price: 32000, ImageURL: "https://images.unsplash.com/photo-1763219802762-1d34ee0907c5?auto=format&fit=crop&w=720&h=520&q=85", StockQuantity: 20, LowStockThreshold: 5, IsActive: true},
		{ID: "food-003", Name: "Mie Goreng Jawa", Category: "Makanan", Description: "Mie goreng dengan sayuran segar", Price: 25000, ImageURL: "https://images.unsplash.com/photo-1752924349474-d717ed874d98?auto=format&fit=crop&w=720&h=520&q=85", StockQuantity: 20, LowStockThreshold: 5, IsActive: true},
		{ID: "drink-001", Name: "Es Teh Manis", Category: "Minuman", Description: "Teh melati dingin", Price: 8000, ImageURL: "https://images.unsplash.com/photo-1676159435365-e4401cd933af?auto=format&fit=crop&w=720&h=520&q=85", StockQuantity: 20, LowStockThreshold: 5, IsActive: true},
		{ID: "drink-002", Name: "Kopi Susu Gula Aren", Category: "Minuman", Description: "Espresso, susu, dan gula aren", Price: 18000, ImageURL: "https://images.unsplash.com/photo-1764298490445-6add12301888?auto=format&fit=crop&w=720&h=520&q=85", StockQuantity: 20, LowStockThreshold: 5, IsActive: true},
		{ID: "snack-001", Name: "Kentang Goreng", Category: "Camilan", Description: "Kentang goreng dengan saus pilihan", Price: 16000, ImageURL: "https://images.unsplash.com/photo-1676566399758-51b0d3927d48?auto=format&fit=crop&w=720&h=520&q=85", StockQuantity: 20, LowStockThreshold: 5, IsActive: true},
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin sample data transaction: %w", err)
	}
	defer tx.Rollback()
	for _, product := range products {
		if _, err := tx.ExecContext(ctx, `INSERT INTO products
			(id, name, category, description, price, image_url, stock_quantity, low_stock_threshold, is_active)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT (id) DO NOTHING`,
			product.ID, product.Name, product.Category, product.Description, product.Price,
			product.ImageURL, product.StockQuantity, product.LowStockThreshold, product.IsActive); err != nil {
			return fmt.Errorf("seed sample product %s: %w", product.ID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit sample products: %w", err)
	}
	return nil
}

func (s *PostgresStore) ListProducts(ctx context.Context, includeInactive bool) ([]model.Product, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+productColumns+` FROM products WHERE ($1 OR is_active) ORDER BY category, name`, includeInactive)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read products: %w", err)
	}
	return products, nil
}

func (s *PostgresStore) GetProduct(ctx context.Context, id string) (model.Product, error) {
	product, err := scanProduct(s.db.QueryRowContext(ctx, `SELECT `+productColumns+` FROM products WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Product{}, ErrProductNotFound
	}
	if err != nil {
		return model.Product{}, fmt.Errorf("get product: %w", err)
	}
	return product, nil
}

func (s *PostgresStore) CreateProduct(ctx context.Context, product model.Product) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO products
		(id, name, category, description, price, image_url, stock_quantity, low_stock_threshold, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,0,$7,$8)`, product.ID, product.Name, product.Category,
		product.Description, product.Price, product.ImageURL, product.LowStockThreshold, product.IsActive)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (s *PostgresStore) UpdateProduct(ctx context.Context, product model.Product) error {
	result, err := s.db.ExecContext(ctx, `UPDATE products SET name=$1, category=$2, description=$3,
		price=$4, image_url=$5, low_stock_threshold=$6, is_active=$7, updated_at=NOW() WHERE id=$8`,
		product.Name, product.Category, product.Description, product.Price, product.ImageURL,
		product.LowStockThreshold, product.IsActive, product.ID)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated product: %w", err)
	}
	if rows == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (s *PostgresStore) SetProductActive(ctx context.Context, id string, active bool) error {
	result, err := s.db.ExecContext(ctx, `UPDATE products SET is_active=$1, updated_at=NOW() WHERE id=$2`, active, id)
	if err != nil {
		return fmt.Errorf("set product active: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check product status: %w", err)
	}
	if rows == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (s *PostgresStore) AdjustStock(ctx context.Context, id string, delta int, reason string, changedBy string) (model.Product, error) {
	if delta == 0 {
		return model.Product{}, ErrInvalidStockChange
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Product{}, fmt.Errorf("begin stock transaction: %w", err)
	}
	defer tx.Rollback()

	product, err := scanProduct(tx.QueryRowContext(ctx, `UPDATE products SET stock_quantity=stock_quantity+$1,
		updated_at=NOW() WHERE id=$2 AND stock_quantity+$1 >= 0 RETURNING `+productColumns, delta, id))
	if errors.Is(err, sql.ErrNoRows) {
		var exists bool
		if queryErr := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)", id).Scan(&exists); queryErr != nil {
			return model.Product{}, fmt.Errorf("check product for stock adjustment: %w", queryErr)
		}
		if !exists {
			return model.Product{}, ErrProductNotFound
		}
		return model.Product{}, ErrInsufficientStock
	}
	if err != nil {
		return model.Product{}, fmt.Errorf("update stock: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO stock_movements (product_id, quantity_delta, reason, changed_by)
		VALUES ($1,$2,$3,$4)`, id, delta, reason, changedBy); err != nil {
		return model.Product{}, fmt.Errorf("record stock movement: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return model.Product{}, fmt.Errorf("commit stock adjustment: %w", err)
	}
	return product, nil
}

type rowScanner interface {
	Scan(...interface{}) error
}

func scanProduct(row rowScanner) (model.Product, error) {
	var product model.Product
	err := row.Scan(&product.ID, &product.Name, &product.Category, &product.Description,
		&product.Price, &product.ImageURL, &product.StockQuantity, &product.LowStockThreshold,
		&product.IsActive, &product.Available)
	return product, err
}
