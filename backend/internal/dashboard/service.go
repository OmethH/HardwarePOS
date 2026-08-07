package dashboard

import (
	"strconv"

	"gorm.io/gorm"

	"hardwarepos/backend/internal/products"
	"hardwarepos/backend/internal/sales"
)

type Service struct {
	db       *gorm.DB
	products *products.Service
	sales    *sales.Service
}

func NewService(db *gorm.DB, productsSvc *products.Service, salesSvc *sales.Service) *Service {
	return &Service{db: db, products: productsSvc, sales: salesSvc}
}

type InventorySummary struct {
	TotalProducts    int64 `json:"total_products"`
	LowStockProducts int64 `json:"low_stock_products"`
	OutOfStock       int64 `json:"out_of_stock_products"`
}

type Summary struct {
	Today     sales.TodaySummary `json:"today"`
	Inventory InventorySummary   `json:"inventory"`
}

func (s *Service) Summary() (*Summary, error) {
	today, err := s.sales.TodaySummary()
	if err != nil {
		return nil, err
	}

	total, err := s.products.TotalCount()
	if err != nil {
		return nil, err
	}
	low, err := s.products.LowStockCount()
	if err != nil {
		return nil, err
	}
	out, err := s.products.OutOfStockCount()
	if err != nil {
		return nil, err
	}

	return &Summary{
		Today: *today,
		Inventory: InventorySummary{
			TotalProducts:    total,
			LowStockProducts: low,
			OutOfStock:       out,
		},
	}, nil
}

type SalesPoint struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

// SalesOverTime returns daily revenue totals for the last N days.
func (s *Service) SalesOverTime(days int) ([]SalesPoint, error) {
	var points []SalesPoint
	err := s.db.Table("sales").
		Select("date(created_at) as date, COALESCE(SUM(total),0) as total").
		Where("created_at >= date('now', ?)", "-"+strconv.Itoa(normalizeDays(days))+" days").
		Group("date(created_at)").
		Order("date").
		Scan(&points).Error
	return points, err
}

type ProfitPoint struct {
	Date   string  `json:"date"`
	Profit float64 `json:"profit"`
}

// ProfitTrend returns daily profit totals for the last N days.
func (s *Service) ProfitTrend(days int) ([]ProfitPoint, error) {
	var points []ProfitPoint
	err := s.db.Table("sale_items").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Select("date(sales.created_at) as date, COALESCE(SUM(sale_items.profit),0) as profit").
		Where("sales.created_at >= date('now', ?)", "-"+strconv.Itoa(normalizeDays(days))+" days").
		Group("date(sales.created_at)").
		Order("date").
		Scan(&points).Error
	return points, err
}

type TopProduct struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	QuantitySold float64 `json:"quantity_sold"`
	Revenue     float64 `json:"revenue"`
}

// TopSellingProducts returns the best sellers by quantity within the last N days.
func (s *Service) TopSellingProducts(days, limit int) ([]TopProduct, error) {
	var items []TopProduct
	err := s.db.Table("sale_items").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Joins("JOIN products ON products.id = sale_items.product_id").
		Select("products.id as product_id, products.name as product_name, "+
			"COALESCE(SUM(sale_items.quantity),0) as quantity_sold, "+
			"COALESCE(SUM(sale_items.quantity * sale_items.selling_price),0) as revenue").
		Where("sales.created_at >= date('now', ?)", "-"+strconv.Itoa(normalizeDays(days))+" days").
		Group("products.id, products.name").
		Order("quantity_sold desc").
		Limit(limit).
		Scan(&items).Error
	return items, err
}

func normalizeDays(days int) int {
	if days <= 0 {
		return 30
	}
	return days
}
