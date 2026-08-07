// Package models contains the GORM entities mapped onto the schema created
// by the SQL files in backend/migrations. AutoMigrate is never used —
// migrations own the schema, GORM only queries it.
package models

import "time"

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	RoleID       uint      `json:"role_id"`
	Role         Role      `json:"role"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Supplier struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SKU           string    `json:"sku"`
	Barcode       *string   `json:"barcode"`
	Name          string    `json:"name"`
	CategoryID    *uint     `json:"category_id"`
	Category      *Category `json:"category,omitempty"`
	Brand         string    `json:"brand"`
	Description   string    `json:"description"`
	UnitType      string    `json:"unit_type"`
	PurchasePrice float64   `json:"purchase_price"`
	SellingPrice  float64   `json:"selling_price"`
	MinimumStock  float64   `json:"minimum_stock"`
	CurrentStock  float64   `json:"current_stock"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const (
	MovementPurchase   = "PURCHASE"
	MovementSale       = "SALE"
	MovementReturn     = "RETURN"
	MovementAdjustment = "ADJUSTMENT"
)

type StockMovement struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `json:"product_id"`
	Product        *Product  `json:"product,omitempty"`
	Quantity       float64   `json:"quantity"`
	Type           string    `json:"type"`
	ReferenceType  string    `json:"reference_type"`
	ReferenceID    *uint     `json:"reference_id"`
	Note           string    `json:"note"`
	CreatedBy      *uint     `json:"created_by"`
	CreatedByUser  *User     `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type Purchase struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ReferenceNumber string         `json:"reference_number"`
	SupplierID      uint           `json:"supplier_id"`
	Supplier        *Supplier      `json:"supplier,omitempty"`
	TotalAmount     float64        `json:"total_amount"`
	CreatedBy       *uint          `json:"created_by"`
	CreatedByUser   *User          `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	Items           []PurchaseItem `gorm:"foreignKey:PurchaseID" json:"items,omitempty"`
}

type PurchaseItem struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	PurchaseID uint     `json:"purchase_id"`
	ProductID  uint     `json:"product_id"`
	Product    *Product `json:"product,omitempty"`
	Quantity   float64  `json:"quantity"`
	Price      float64  `json:"price"`
}

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	SaleStatusCompleted        = "COMPLETED"
	SaleStatusPartiallyReturned = "PARTIALLY_RETURNED"
	SaleStatusReturned         = "RETURNED"
)

type Sale struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	InvoiceNumber string     `json:"invoice_number"`
	CustomerID    uint       `json:"customer_id"`
	Customer      *Customer  `json:"customer,omitempty"`
	Subtotal      float64    `json:"subtotal"`
	Discount      float64    `json:"discount"`
	Total         float64    `json:"total"`
	PaymentMethod string     `json:"payment_method"`
	Status        string     `json:"status"`
	CreatedBy     *uint      `json:"created_by"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	Items         []SaleItem `gorm:"foreignKey:SaleID" json:"items,omitempty"`
}

type SaleItem struct {
	ID                uint     `gorm:"primaryKey" json:"id"`
	SaleID            uint     `json:"sale_id"`
	ProductID         uint     `json:"product_id"`
	Product           *Product `json:"product,omitempty"`
	Quantity          float64  `json:"quantity"`
	ReturnedQuantity  float64  `json:"returned_quantity"`
	SellingPrice      float64  `json:"selling_price"`
	PurchasePrice     float64  `json:"purchase_price"`
	Profit            float64  `json:"profit"`
}

type Return struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	ReturnNumber  string       `json:"return_number"`
	SaleID        uint         `json:"sale_id"`
	Sale          *Sale        `json:"sale,omitempty"`
	TotalRefund   float64      `json:"total_refund"`
	CreatedBy     *uint        `json:"created_by"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	Items         []ReturnItem `gorm:"foreignKey:ReturnID" json:"items,omitempty"`
}

type ReturnItem struct {
	ID          uint     `gorm:"primaryKey" json:"id"`
	ReturnID    uint     `json:"return_id"`
	SaleItemID  uint     `json:"sale_item_id"`
	ProductID   uint     `json:"product_id"`
	Product     *Product `json:"product,omitempty"`
	Quantity    float64  `json:"quantity"`
	RefundAmount float64 `json:"refund_amount"`
}

type Settings struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ShopName       string    `json:"shop_name"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	ReceiptFooter  string    `json:"receipt_footer"`
	TaxPercentage  float64   `json:"tax_percentage"`
	UpdatedAt      time.Time `json:"updated_at"`
}
