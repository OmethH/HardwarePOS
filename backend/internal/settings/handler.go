package settings

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
	rg.GET("/settings", h.get)

	write := rg.Group("")
	write.Use(middleware.RequireRole("admin"))
	write.PUT("/settings", h.update)
}

func (h *Handler) get(c *gin.Context) {
	item, err := h.svc.Get()
	if err != nil {
		httpresp.Err(c, 500, "failed to load settings")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) update(c *gin.Context) {
	var req Input
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "shop_name is required")
		return
	}
	item, err := h.svc.Update(req)
	if err != nil {
		httpresp.Err(c, 500, "failed to update settings")
		return
	}
	httpresp.OK(c, 200, item)
}
