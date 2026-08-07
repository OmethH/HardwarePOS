package sales

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
	rg.POST("/sales", h.checkout)
	rg.GET("/sales", h.list)
	rg.GET("/sales/:id", h.get)
	rg.GET("/sales/invoice/:invoice", h.getByInvoice)
}

func (h *Handler) checkout(c *gin.Context) {
	var req CreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "at least one item is required")
		return
	}

	userIDVal, _ := c.Get(middleware.CtxUserID)
	userID, _ := userIDVal.(uint)

	sale, err := h.svc.Checkout(req, userID)
	if err != nil {
		httpresp.Err(c, 400, err.Error())
		return
	}
	httpresp.OK(c, 201, sale)
}

func (h *Handler) list(c *gin.Context) {
	var customerID uint
	if v := c.Query("customer_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		customerID = uint(id)
	}

	f := ListFilter{
		Search:     c.Query("search"),
		From:       c.Query("from"),
		To:         c.Query("to"),
		CustomerID: customerID,
	}

	// Cashiers may only see their own transactions; admins/managers see all.
	role, _ := c.Get(middleware.CtxRole)
	if roleStr, _ := role.(string); roleStr == "cashier" {
		userIDVal, _ := c.Get(middleware.CtxUserID)
		f.OnlyUserID, _ = userIDVal.(uint)
	}

	items, err := h.svc.List(f)
	if err != nil {
		httpresp.Err(c, 500, "failed to load sales")
		return
	}
	httpresp.OK(c, 200, items)
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Get(uint(id))
	if err != nil {
		httpresp.Err(c, 404, "sale not found")
		return
	}

	if !canView(c, item.CreatedBy) {
		httpresp.Err(c, 403, "you may only view your own transactions")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) getByInvoice(c *gin.Context) {
	item, err := h.svc.GetByInvoice(c.Param("invoice"))
	if err != nil {
		httpresp.Err(c, 404, "sale not found")
		return
	}

	if !canView(c, item.CreatedBy) {
		httpresp.Err(c, 403, "you may only view your own transactions")
		return
	}
	httpresp.OK(c, 200, item)
}

func canView(c *gin.Context, createdBy *uint) bool {
	role, _ := c.Get(middleware.CtxRole)
	roleStr, _ := role.(string)
	if roleStr != "cashier" {
		return true
	}
	userIDVal, _ := c.Get(middleware.CtxUserID)
	userID, _ := userIDVal.(uint)
	return createdBy != nil && *createdBy == userID
}
