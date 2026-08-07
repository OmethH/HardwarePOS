package reports

import (
	"github.com/gin-gonic/gin"

	"hardwarepos/backend/internal/httpresp"
	"hardwarepos/backend/internal/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("")
	group.Use(middleware.RequireRole("admin", "manager"))
	group.GET("/reports/sales", h.salesReport)
	group.GET("/reports/inventory", h.inventoryReport)
	group.GET("/reports/product-performance", h.productPerformance)
}

var salesHeaders = []string{"Date", "Invoices", "Revenue", "Discount", "Profit"}
var inventoryHeaders = []string{"Product", "SKU", "Current Stock", "Stock Value", "Status"}
var performanceHeaders = []string{"Product", "Quantity Sold", "Revenue", "Profit"}

func (h *Handler) salesReport(c *gin.Context) {
	rows, err := h.svc.SalesReport(c.Query("from"), c.Query("to"))
	if err != nil {
		httpresp.Err(c, 500, "failed to generate sales report")
		return
	}

	cell := func(r, col int) string {
		row := rows[r]
		switch col {
		case 0:
			return row.Date
		case 1:
			return fmtQty(float64(row.Invoices))
		case 2:
			return fmtMoney(row.Revenue)
		case 3:
			return fmtMoney(row.Discount)
		default:
			return fmtMoney(row.Profit)
		}
	}

	switch c.Query("format") {
	case "csv":
		data, err := ToCSV(salesHeaders, len(rows), cell)
		if err != nil {
			httpresp.Err(c, 500, "failed to export csv")
			return
		}
		serveFile(c, "sales-report.csv", "text/csv", data)
	case "pdf":
		data, err := ToPDF("Sales Report", salesHeaders, len(rows), cell)
		if err != nil {
			httpresp.Err(c, 500, "failed to export pdf")
			return
		}
		serveFile(c, "sales-report.pdf", "application/pdf", data)
	default:
		httpresp.OK(c, 200, rows)
	}
}

func (h *Handler) inventoryReport(c *gin.Context) {
	rows, err := h.svc.InventoryReport()
	if err != nil {
		httpresp.Err(c, 500, "failed to generate inventory report")
		return
	}

	cell := func(r, col int) string {
		row := rows[r]
		switch col {
		case 0:
			return row.Product
		case 1:
			return row.SKU
		case 2:
			return fmtQty(row.CurrentStock)
		case 3:
			return fmtMoney(row.StockValue)
		default:
			return row.Status
		}
	}

	switch c.Query("format") {
	case "csv":
		data, err := ToCSV(inventoryHeaders, len(rows), cell)
		if err != nil {
			httpresp.Err(c, 500, "failed to export csv")
			return
		}
		serveFile(c, "inventory-report.csv", "text/csv", data)
	case "pdf":
		data, err := ToPDF("Inventory Report", inventoryHeaders, len(rows), cell)
		if err != nil {
			httpresp.Err(c, 500, "failed to export pdf")
			return
		}
		serveFile(c, "inventory-report.pdf", "application/pdf", data)
	default:
		httpresp.OK(c, 200, rows)
	}
}

func (h *Handler) productPerformance(c *gin.Context) {
	rows, err := h.svc.ProductPerformance(c.Query("from"), c.Query("to"))
	if err != nil {
		httpresp.Err(c, 500, "failed to generate product performance report")
		return
	}

	cell := func(r, col int) string {
		row := rows[r]
		switch col {
		case 0:
			return row.Product
		case 1:
			return fmtQty(row.QuantitySold)
		case 2:
			return fmtMoney(row.Revenue)
		default:
			return fmtMoney(row.Profit)
		}
	}

	switch c.Query("format") {
	case "csv":
		data, err := ToCSV(performanceHeaders, len(rows), cell)
		if err != nil {
			httpresp.Err(c, 500, "failed to export csv")
			return
		}
		serveFile(c, "product-performance.csv", "text/csv", data)
	case "pdf":
		data, err := ToPDF("Product Performance", performanceHeaders, len(rows), cell)
		if err != nil {
			httpresp.Err(c, 500, "failed to export pdf")
			return
		}
		serveFile(c, "product-performance.pdf", "application/pdf", data)
	default:
		httpresp.OK(c, 200, rows)
	}
}

func serveFile(c *gin.Context, filename, contentType string, data []byte) {
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(200, contentType, data)
}
