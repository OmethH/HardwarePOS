package reports

import (
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type SalesReportRow struct {
	Date     string  `json:"date"`
	Invoices int64   `json:"invoices"`
	Revenue  float64 `json:"revenue"`
	Discount float64 `json:"discount"`
	Profit   float64 `json:"profit"`
}

// SalesReport aggregates revenue/discount/profit per day within [from, to]
// (inclusive, "YYYY-MM-DD"). Empty bounds mean unbounded.
func (s *Service) SalesReport(from, to string) ([]SalesReportRow, error) {
	q := s.db.Table("sales").
		Select(`date(sales.created_at) as date,
			COUNT(DISTINCT sales.id) as invoices,
			COALESCE(SUM(sales.total),0) as revenue,
			COALESCE(SUM(sales.discount),0) as discount,
			COALESCE((SELECT SUM(si.profit) FROM sale_items si WHERE si.sale_id IN (
				SELECT id FROM sales s2 WHERE date(s2.created_at) = date(sales.created_at)
			)),0) as profit`).
		Group("date(sales.created_at)").
		Order("date")

	if from != "" {
		q = q.Where("date(sales.created_at) >= ?", from)
	}
	if to != "" {
		q = q.Where("date(sales.created_at) <= ?", to)
	}

	var rows []SalesReportRow
	err := q.Scan(&rows).Error
	return rows, err
}

type InventoryReportRow struct {
	Product      string  `json:"product"`
	SKU          string  `json:"sku"`
	CurrentStock float64 `json:"current_stock"`
	StockValue   float64 `json:"stock_value"`
	Status       string  `json:"status"`
}

// InventoryReport shows current stock valuation for every product.
// Status reflects the stock level relative to minimum_stock.
func (s *Service) InventoryReport() ([]InventoryReportRow, error) {
	var rows []InventoryReportRow
	err := s.db.Table("products").
		Select(`name as product, sku as sku,
			current_stock as current_stock,
			current_stock * purchase_price as stock_value,
			CASE
				WHEN current_stock <= 0 THEN 'out_of_stock'
				WHEN current_stock <= minimum_stock THEN 'low_stock'
				ELSE 'in_stock'
			END as status`).
		Order("name").
		Scan(&rows).Error
	return rows, err
}

type ProductPerformanceRow struct {
	Product      string  `json:"product"`
	QuantitySold float64 `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
	Profit       float64 `json:"profit"`
}

// ProductPerformance ranks products by revenue within [from, to].
func (s *Service) ProductPerformance(from, to string) ([]ProductPerformanceRow, error) {
	q := s.db.Table("sale_items").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Joins("JOIN products ON products.id = sale_items.product_id").
		Select(`products.name as product,
			COALESCE(SUM(sale_items.quantity),0) as quantity_sold,
			COALESCE(SUM(sale_items.quantity * sale_items.selling_price),0) as revenue,
			COALESCE(SUM(sale_items.profit),0) as profit`).
		Group("products.id, products.name").
		Order("revenue desc")

	if from != "" {
		q = q.Where("date(sales.created_at) >= ?", from)
	}
	if to != "" {
		q = q.Where("date(sales.created_at) <= ?", to)
	}

	var rows []ProductPerformanceRow
	err := q.Scan(&rows).Error
	return rows, err
}
