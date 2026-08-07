package users

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
	group.Use(middleware.RequireRole("admin"))
	group.GET("/users", h.list)
	group.GET("/users/:id", h.get)
	group.POST("/users", h.create)
	group.PUT("/users/:id", h.update)
	group.DELETE("/users/:id", h.deactivate)
	group.GET("/roles", h.listRoles)
}

func (h *Handler) list(c *gin.Context) {
	items, err := h.svc.List()
	if err != nil {
		httpresp.Err(c, 500, "failed to load users")
		return
	}
	httpresp.OK(c, 200, items)
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Get(uint(id))
	if err != nil {
		httpresp.Err(c, 404, "user not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) create(c *gin.Context) {
	var req CreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "name, username, password (min 6 chars) and role_id are required")
		return
	}
	item, err := h.svc.Create(req)
	if err != nil {
		httpresp.Err(c, 400, "could not create user, username may already exist")
		return
	}
	httpresp.OK(c, 201, item)
}

func (h *Handler) update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req UpdateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "name and role_id are required")
		return
	}
	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		httpresp.Err(c, 404, "user not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) deactivate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Deactivate(uint(id)); err != nil {
		httpresp.Err(c, 400, "could not deactivate user")
		return
	}
	httpresp.OK(c, 200, gin.H{"deactivated": true})
}

func (h *Handler) listRoles(c *gin.Context) {
	items, err := h.svc.ListRoles()
	if err != nil {
		httpresp.Err(c, 500, "failed to load roles")
		return
	}
	httpresp.OK(c, 200, items)
}
