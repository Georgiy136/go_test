package http

import (
	"github.com/Georgiy136/go_test/auth_service/internal/http/httpresponse"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	us service.AuthService
}

func (h *AuthHandler) GetTokens(c *gin.Context) {
	type getTokensParamsRequest struct {
		UserID int `form:"user_id" binding:"required,gt=0"`
	}
	var req getTokensParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	userAgent := c.Request.UserAgent()
	//ipAddress := c.GetHeader("X-Forwarded-For")
	clientIp := c.ClientIP()

	tokens, err := h.us.GetTokens(c.Request.Context(), models.DataFromRequestGetTokens{
		UserID:    req.UserID,
		UserAgent: userAgent,
		IpAddress: clientIp,
	})
	if err != nil {
		httpresponse.HandleError(c, err, nil)
		return
	}
	httpresponse.SendSuccess(c, http.StatusOK, tokens)
}

func (h *AuthHandler) UpdateGood(c *gin.Context) {
	type updateGoodParamsRequest struct {
		GoodID    int `form:"good_id" binding:"required,gt=0"`
		ProjectID int `form:"project_id" binding:"required,gt=0"`
	}
	var req updateGoodParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	type updateGoodsRequest struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description" binding:"omitempty"`
	}
	var body updateGoodsRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.UpdateGood(c.Request.Context(), models.DataFromRequestGoodsUpdate{
		GoodID:      req.GoodID,
		ProjectID:   req.ProjectID,
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		httpresponse.HandleError(c, err, nil)
		return
	}
	httpresponse.SendSuccessOK(c, goods)
}

func (h *AuthHandler) DeleteGood(c *gin.Context) {
	type deleteGoodParamsRequest struct {
		GoodsID   int `form:"goods_id" binding:"required,gt=0"`
		ProjectID int `form:"project_id" binding:"required,gt=0"`
	}
	var req deleteGoodParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.DeleteGood(c.Request.Context(), models.DataFromRequestGoodsDelete{
		GoodID:    req.GoodsID,
		ProjectID: req.ProjectID,
		DeletedAt: time.Now(),
	})
	if err != nil {
		httpresponse.HandleError(c, err, nil)
		return
	}
	httpresponse.SendSuccessOK(c, goods)
}
