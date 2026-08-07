package dashboard

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
	group := rg.Group("")
	group.Use(middleware.RequireRole("admin", "manager"))
	group.GET("/dashboard/summary", h.summary)
	group.GET("/dashboard/sales-over-time", h.salesOverTime)
	group.GET("/dashboard/profit-trend", h.profitTrend)
	group.GET("/dashboard/top-products", h.topProducts)
}

func (h *Handler) summary(c *gin.Context) {
	s, err := h.svc.Summary()
	if err != nil {
		httpresp.Err(c, 500, "failed to load dashboard summary")
		return
	}
	httpresp.OK(c, 200, s)
}

func (h *Handler) salesOverTime(c *gin.Context) {
	days, _ := strconv.Atoi(c.Query("days"))
	points, err := h.svc.SalesOverTime(days)
	if err != nil {
		httpresp.Err(c, 500, "failed to load sales over time")
		return
	}
	httpresp.OK(c, 200, points)
}

func (h *Handler) profitTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.Query("days"))
	points, err := h.svc.ProfitTrend(days)
	if err != nil {
		httpresp.Err(c, 500, "failed to load profit trend")
		return
	}
	httpresp.OK(c, 200, points)
}

func (h *Handler) topProducts(c *gin.Context) {
	days, _ := strconv.Atoi(c.Query("days"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit <= 0 {
		limit = 10
	}
	items, err := h.svc.TopSellingProducts(days, limit)
	if err != nil {
		httpresp.Err(c, 500, "failed to load top products")
		return
	}
	httpresp.OK(c, 200, items)
}
