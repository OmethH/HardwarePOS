package sales

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"hardwarepos/backend/internal/inventory"
	"hardwarepos/backend/internal/models"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type ItemInput struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
}

type CreateInput struct {
	CustomerID    uint        `json:"customer_id"`
	Discount      float64     `json:"discount"`
	PaymentMethod string      `json:"payment_method"`
	Items         []ItemInput `json:"items" binding:"required,min=1"`
}

var (
	ErrNoItems           = errors.New("a sale must include at least one item")
	ErrInsufficientStock = errors.New("insufficient stock for one or more items")
)

// Checkout performs the full POS checkout sequence — create sale, create
// sale items, reduce stock, record stock movements, calculate profit — as
// one database transaction so a partial failure never leaves stock or
// financial records inconsistent.
func (s *Service) Checkout(in CreateInput, userID uint) (*models.Sale, error) {
	if len(in.Items) == 0 {
		return nil, ErrNoItems
	}

	customerID := in.CustomerID
	if customerID == 0 {
		customerID = 1 // Walk-in Customer
	}
	paymentMethod := in.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "cash"
	}

	var sale models.Sale
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var products []models.Product
		productIDs := make([]uint, 0, len(in.Items))
		for _, item := range in.Items {
			productIDs = append(productIDs, item.ProductID)
		}
		if err := tx.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
			return err
		}
		productByID := make(map[uint]models.Product, len(products))
		for _, p := range products {
			productByID[p.ID] = p
		}

		var subtotal float64
		for _, item := range in.Items {
			product, ok := productByID[item.ProductID]
			if !ok {
				return fmt.Errorf("product %d not found", item.ProductID)
			}
			if product.CurrentStock < item.Quantity {
				return fmt.Errorf("%w: %s (available %.2f, requested %.2f)", ErrInsufficientStock, product.Name, product.CurrentStock, item.Quantity)
			}
			subtotal += product.SellingPrice * item.Quantity
		}

		discount := in.Discount
		if discount > subtotal {
			discount = subtotal
		}
		total := subtotal - discount

		sale = models.Sale{
			InvoiceNumber: fmt.Sprintf("INV-PENDING-%d", time.Now().UnixNano()),
			CustomerID:    customerID,
			Subtotal:      subtotal,
			Discount:      discount,
			Total:         total,
			PaymentMethod: paymentMethod,
			Status:        models.SaleStatusCompleted,
			CreatedBy:     &userID,
		}
		if err := tx.Create(&sale).Error; err != nil {
			return err
		}
		sale.InvoiceNumber = fmt.Sprintf("INV%05d", sale.ID)
		if err := tx.Model(&sale).Update("invoice_number", sale.InvoiceNumber).Error; err != nil {
			return err
		}

		for _, item := range in.Items {
			product := productByID[item.ProductID]
			profit := (product.SellingPrice - product.PurchasePrice) * item.Quantity

			saleItem := models.SaleItem{
				SaleID:        sale.ID,
				ProductID:     product.ID,
				Quantity:      item.Quantity,
				SellingPrice:  product.SellingPrice,
				PurchasePrice: product.PurchasePrice,
				Profit:        profit,
			}
			if err := tx.Create(&saleItem).Error; err != nil {
				return err
			}

			refID := sale.ID
			if err := inventory.ApplyMovement(tx, product.ID, -item.Quantity, models.MovementSale, "sale", &refID,
				fmt.Sprintf("Sale %s", sale.InvoiceNumber), &userID); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Get(sale.ID)
}

type ListFilter struct {
	Search     string
	From       string
	To         string
	CustomerID uint
	OnlyUserID uint // when > 0, restricts to sales created by this user (cashiers see only their own)
}

func (s *Service) List(f ListFilter) ([]models.Sale, error) {
	var sales []models.Sale
	q := s.db.Preload("Customer").Preload("CreatedByUser.Role").Order("created_at desc")

	if f.Search != "" {
		q = q.Where("invoice_number LIKE ?", "%"+f.Search+"%")
	}
	if f.From != "" {
		q = q.Where("created_at >= ?", f.From)
	}
	if f.To != "" {
		q = q.Where("created_at <= ?", f.To)
	}
	if f.CustomerID != 0 {
		q = q.Where("customer_id = ?", f.CustomerID)
	}
	if f.OnlyUserID != 0 {
		q = q.Where("created_by = ?", f.OnlyUserID)
	}

	err := q.Find(&sales).Error
	return sales, err
}

func (s *Service) Get(id uint) (*models.Sale, error) {
	var sale models.Sale
	err := s.db.Preload("Customer").Preload("CreatedByUser.Role").Preload("Items").Preload("Items.Product").
		First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (s *Service) GetByInvoice(invoice string) (*models.Sale, error) {
	var sale models.Sale
	err := s.db.Preload("Customer").Preload("CreatedByUser.Role").Preload("Items").Preload("Items.Product").
		Where("invoice_number = ?", invoice).First(&sale).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

// TodaySummary powers the "Today" tile group on the dashboard.
type TodaySummary struct {
	SalesAmount      float64 `json:"sales_amount"`
	Profit           float64 `json:"profit"`
	TransactionCount int64   `json:"transaction_count"`
	ItemsSold        float64 `json:"items_sold"`
}

func (s *Service) TodaySummary() (*TodaySummary, error) {
	today := time.Now().Format("2006-01-02")
	summary := &TodaySummary{}

	if err := s.db.Model(&models.Sale{}).
		Where("date(created_at) = ?", today).
		Select("COALESCE(SUM(total),0)").Row().Scan(&summary.SalesAmount); err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Sale{}).Where("date(created_at) = ?", today).Count(&summary.TransactionCount).Error; err != nil {
		return nil, err
	}
	if err := s.db.Table("sale_items").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Where("date(sales.created_at) = ?", today).
		Select("COALESCE(SUM(sale_items.profit),0)").Row().Scan(&summary.Profit); err != nil {
		return nil, err
	}
	if err := s.db.Table("sale_items").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Where("date(sales.created_at) = ?", today).
		Select("COALESCE(SUM(sale_items.quantity),0)").Row().Scan(&summary.ItemsSold); err != nil {
		return nil, err
	}

	return summary, nil
}
