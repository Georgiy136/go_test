package http

import (
	"myapp/internal/http/httpresponse"
	"myapp/internal/models"
	"myapp/internal/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GoodsHandler struct {
	us usecase.GoodsUseCases
}

func (h *GoodsHandler) PostGoods(c *gin.Context) {
	type postGoodsParamsRequest struct {
		ProjectID int `form:"project_id" binding:"required,gt=0"`
	}
	var req postGoodsParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	type postGoodsRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"omitempty"`
		Priority    int    `json:"priority" binding:"omitempty,gt=0"`
	}
	var body postGoodsRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.AddGoods(c.Request.Context(), models.DataFromRequestGoodsAdd{
		ProjectID:   req.ProjectID,
		Name:        body.Name,
		Description: body.Description,
		Priority:    body.Priority,
	})
	if err != nil {
		httpresponse.HandleError(c, err, nil)
		return
	}
	httpresponse.SendSuccess(c, http.StatusCreated, goods)
}

func (h *GoodsHandler) UpdateGood(c *gin.Context) {
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

func (h *GoodsHandler) DeleteGood(c *gin.Context) {
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

func (h *GoodsHandler) ListGoods(c *gin.Context) {
	type listGoodParamsRequest struct {
		GoodsID   *int `form:"good_id" binding:"omitempty,gt=0"`
		ProjectID *int `form:"project_id" binding:"omitempty,gt=0"`
		Limit     *int `form:"limit" binding:"omitempty,gte=0"`
		Offset    *int `form:"offset" binding:"omitempty,gte=0"`
	}
	var req listGoodParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goodsList, err := h.us.ListGoods(c.Request.Context(), models.DataFromRequestGoodsList{
		GoodsID:   req.GoodsID,
		ProjectID: req.ProjectID,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		httpresponse.HandleError(c, err, nil)
		return
	}
	httpresponse.SendSuccessOK(c, goodsList)
}

func (h *GoodsHandler) ReprioritizeGood(c *gin.Context) {
	type reprioritizeGoodParamsRequest struct {
		GoodsID   int `form:"good_id" binding:"required,gt=0"`
		ProjectID int `form:"project_id" binding:"required,gt=0"`
	}
	var req reprioritizeGoodParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	type reprioritizeGoodsRequest struct {
		NewPriority int `json:"NewPriority" binding:"required,gte=0"`
	}
	var body reprioritizeGoodsRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.ReprioritizeGood(c.Request.Context(), models.DataFromRequestReprioritizeGood{
		GoodID:    req.GoodsID,
		ProjectID: req.ProjectID,
		Priority:  body.NewPriority,
	})
	if err != nil {
		httpresponse.HandleError(c, err, nil)
		return
	}
	httpresponse.SendSuccessOK(c, goods)
}
