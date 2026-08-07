package auth

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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authed *gin.RouterGroup) {
	rg.POST("/auth/login", h.login)
	authed.GET("/auth/me", h.me)
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.Err(c, 400, "username and password are required")
		return
	}

	result, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		httpresp.Err(c, 401, "invalid username or password")
		return
	}

	httpresp.OK(c, 200, result)
}

func (h *Handler) me(c *gin.Context) {
	userID, _ := c.Get(middleware.CtxUserID)
	uid, _ := userID.(uint)

	user, err := h.svc.Me(uid)
	if err != nil {
		httpresp.Err(c, 404, "user not found")
		return
	}

	httpresp.OK(c, 200, user)
}
