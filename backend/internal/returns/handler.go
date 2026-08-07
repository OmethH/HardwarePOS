package returns

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
	group.GET("/returns", h.list)
	group.GET("/returns/:id", h.get)
	group.POST("/returns", h.create)
}

func (h *Handler) list(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, err := h.svc.List(limit)
	if err != nil {
		httpresp.Err(c, 500, "failed to load returns")
		return
	}
	httpresp.OK(c, 200, items)
}

func (h *Handler) get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	item, err := h.svc.Get(uint(id))
	if err != nil {
		httpresp.Err(c, 404, "return not found")
		return
	}
	httpresp.OK(c, 200, item)
}

func (h *Handler) create(c *gin.Context) {
	var req CreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "sale_id and at least one item are required")
		return
	}

	userIDVal, _ := c.Get(middleware.CtxUserID)
	userID, _ := userIDVal.(uint)

	item, err := h.svc.Create(req, userID)
	if err != nil {
		httpresp.Err(c, 400, err.Error())
		return
	}
	httpresp.OK(c, 201, item)
}
