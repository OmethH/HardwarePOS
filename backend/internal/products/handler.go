package products

import (
	"strconv"

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
	rg.GET("/products", h.list)
	rg.GET("/products/:id", h.get)
	rg.GET("/products/barcode/:barcode", h.getByBarcode)

	write := rg.Group("")
	write.Use(middleware.RequireRole("admin", "manager"))
	write.POST("/products", h.create)
	write.PUT("/products/:id", h.update)
	write.DELETE("/products/:id", h.delete)
}

func (h *Handler) list(c *gin.Context) {
	var categoryID uint
	if v := c.Query("category_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		categoryID = uint(id)
	}
	f := ListFilter{
		Search:     c.Query("search"),
		CategoryID: categoryID,
		LowStock:   c.Query("low_stock") == "true",
		Status:     c.Query("status"),
	}
	items, err := h.svc.List(f)
	if err != nil {
		httpresp.Err(c, 500, "failed to load products")
		return
	}
	httpresp.OK(c, 200, items)
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Get(uint(id))
	if err != nil {
		httpresp.Err(c, 404, "product not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) getByBarcode(c *gin.Context) {
	item, err := h.svc.GetByBarcode(c.Param("barcode"))
	if err != nil {
		httpresp.Err(c, 404, "product not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) create(c *gin.Context) {
	var req Input
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "sku, name and unit_type are required")
		return
	}
	item, err := h.svc.Create(req)
	if err != nil {
		httpresp.Err(c, 400, "could not create product, sku/barcode may already exist")
		return
	}
	httpresp.OK(c, 201, item)
}

func (h *Handler) update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req Input
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "sku, name and unit_type are required")
		return
	}

	role, _ := c.Get(middleware.CtxRole)
	roleStr, _ := role.(string)
	allowPriceChange := roleStr == "admin" || roleStr == "manager"

	item, err := h.svc.Update(uint(id), req, allowPriceChange)
	if err != nil {
		if err == ErrPurchasePriceLocked {
			httpresp.Err(c, 403, err.Error())
			return
		}
		httpresp.Err(c, 404, "product not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(uint(id)); err != nil {
		httpresp.Err(c, 400, "could not delete product")
		return
	}
	httpresp.OK(c, 200, gin.H{"deleted": true})
}
