package inventory

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
	rg.GET("/inventory/movements", h.list)
	rg.GET("/inventory/movements/product/:id", h.listByProduct)

	adjust := rg.Group("")
	adjust.Use(middleware.RequireRole("admin", "manager"))
	adjust.POST("/inventory/adjustments", h.adjust)
}

func (h *Handler) list(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, err := h.svc.List(limit)
	if err != nil {
		httpresp.Err(c, 500, "failed to load stock movements")
		return
	}
	httpresp.OK(c, 200, items)
}

func (h *Handler) listByProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	items, err := h.svc.ListByProduct(uint(id))
	if err != nil {
		httpresp.Err(c, 500, "failed to load stock movements")
		return
	}
	httpresp.OK(c, 200, items)
}

type adjustRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	Note      string  `json:"note" binding:"required"`
}

func (h *Handler) adjust(c *gin.Context) {
	var req adjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "product_id, quantity and note are required")
		return
	}

	userIDVal, _ := c.Get(middleware.CtxUserID)
	userID, _ := userIDVal.(uint)

	movement, err := h.svc.ManualAdjustment(req.ProductID, req.Quantity, req.Note, userID)
	if err != nil {
		httpresp.Err(c, 400, err.Error())
		return
	}
	httpresp.OK(c, 201, movement)
}
