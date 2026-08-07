package customers

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
	rg.GET("/customers", h.list)
	rg.GET("/customers/:id", h.get)

	write := rg.Group("")
	write.Use(middleware.RequireRole("admin", "manager"))
	write.POST("/customers", h.create)
	write.PUT("/customers/:id", h.update)
	write.DELETE("/customers/:id", h.delete)
}

func (h *Handler) list(c *gin.Context) {
	items, err := h.svc.List(c.Query("search"))
	if err != nil {
		httpresp.Err(c, 500, "failed to load customers")
		return
	}
	httpresp.OK(c, 200, items)
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Get(uint(id))
	if err != nil {
		httpresp.Err(c, 404, "customer not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) create(c *gin.Context) {
	var req Input
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "name is required")
		return
	}
	item, err := h.svc.Create(req)
	if err != nil {
		httpresp.Err(c, 400, "could not create customer")
		return
	}
	httpresp.OK(c, 201, item)
}

func (h *Handler) update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req Input
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "name is required")
		return
	}
	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		httpresp.Err(c, 404, "customer not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(uint(id)); err != nil {
		httpresp.Err(c, 400, "could not delete customer")
		return
	}
	httpresp.OK(c, 200, gin.H{"deleted": true})
}
