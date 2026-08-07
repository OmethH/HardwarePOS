package categories

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
	rg.GET("/categories", h.list)
	rg.GET("/categories/:id", h.get)

	write := rg.Group("")
	write.Use(middleware.RequireRole("admin", "manager"))
	write.POST("/categories", h.create)
	write.PUT("/categories/:id", h.update)
	write.DELETE("/categories/:id", h.delete)
}

type categoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) list(c *gin.Context) {
	items, err := h.svc.List()
	if err != nil {
		httpresp.Err(c, 500, "failed to load categories")
		return
	}
	httpresp.OK(c, 200, items)
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Get(uint(id))
	if err != nil {
		httpresp.Err(c, 404, "category not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) create(c *gin.Context) {
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "name is required")
		return
	}
	item, err := h.svc.Create(req.Name)
	if err != nil {
		httpresp.Err(c, 400, "could not create category, name may already exist")
		return
	}
	httpresp.OK(c, 201, item)
}

func (h *Handler) update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "name is required")
		return
	}
	item, err := h.svc.Update(uint(id), req.Name)
	if err != nil {
		httpresp.Err(c, 404, "category not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(uint(id)); err != nil {
		httpresp.Err(c, 400, "could not delete category")
		return
	}
	httpresp.OK(c, 200, gin.H{"deleted": true})
}
