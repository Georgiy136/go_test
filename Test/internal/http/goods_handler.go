package http

import (
	"myapp/internal/http/httpresponse"
	"myapp/internal/models"
	"myapp/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GoodsHandler struct {
	us usecase.GoodsUseCases
}

// PostGoods godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Add Goods to database
//	@Tags			Goods
//	@Description	create Goods
//	@ID				create-Goods
//	@Accept			json
//	@Produce		json
//	@Param			input	body		postGoodsRequest	true	"Goods info"
//	@Success		201		{object}	models.Goods
//	@Router			/Goods [post]
func (h *GoodsHandler) PostGoods(c *gin.Context) {
	type postGoodsParamsRequest struct {
		ProjectID int `form:"projectID" binding:"required,gt=0"`
	}
	var req postGoodsParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	type postGoodsRequest struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description" binding:"omitempty"`
		Priority    *int    `json:"priority" binding:"required,gt=0"`
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
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}
	httpresponse.SendSuccess(c, http.StatusCreated, goods)
}

// UpdateGoods godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Update Goods based on given ID
//	@Tags			Goods
//	@Description	update Goods by id
//	@ID				update-Goods-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Goods ID"
//	@Param			input	body		PutGoodsRequest	true	"Goods info"
//	@Success		201		{object}	models.Goods
//	@Router			/Goods/{id} [put]
func (h *GoodsHandler) UpdateGood(c *gin.Context) {
	type updateGoodParamsRequest struct {
		ID        int `form:"id" binding:"required,gt=0"`
		ProjectID int `form:"projectID" binding:"required,gt=0"`
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
		ID:          req.ID,
		ProjectID:   req.ProjectID,
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}
	httpresponse.SendSuccessOK(c, goods)
}

// DeleteGoods godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Delete Goods based on given ID
//	@Tags			Goods
//	@Description	delete Goods by id
//	@ID				delete-Goods-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Goods ID"
//	@Success		200
//	@Router			/Goods/{id} [delete]
func (h *GoodsHandler) DeleteGood(c *gin.Context) {
	type deleteGoodParamsRequest struct {
		ID        int `form:"id" binding:"required,gt=0"`
		ProjectID int `form:"projectID" binding:"required,gt=0"`
	}
	var req deleteGoodParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.DeleteGood(c.Request.Context(), models.DataFromRequestGoodsDelete{
		ID:        req.ID,
		ProjectID: req.ProjectID,
	})

	if err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}
	httpresponse.SendSuccessOK(c, goods)
}

// GetAllGoodss godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves All Goodss
//	@Tags			Goods
//	@Description	get all Goodss
//	@ID				get-all-Goodss
//	@Accept			json
//	@Produce		json
//	@Success		202	{array}	[]models.Goods
//	@Router			/Goods [get]
func (h *GoodsHandler) ListGoods(c *gin.Context) {
	type listGoodParamsRequest struct {
		ID        int `form:"id" binding:"omitempty,gt=0"`
		ProjectID int `form:"projectID" binding:"omitempty,gt=0"`
		Limit     int `form:"limit" binding:"omitempty,gt=0"`
		Offset    int `form:"offset" binding:"omitempty,gt=0"`
	}
	var req listGoodParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.ListGoods(c.Request.Context(), models.DataFromRequestGoodsList{
		ID:        req.ID,
		ProjectID: req.ProjectID,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}
	httpresponse.SendSuccessOK(c, goods)
}

// GetGoods godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves All Goodss
//	@Tags			Goods
//	@Description	get all Goodss
//	@ID				get-all-Goodss
//	@Accept			json
//	@Produce		json
//	@Success		202	{array}	[]models.Goods
//	@Router			/Goods [get]
func (h *GoodsHandler) ReprioritizeGood(c *gin.Context) {
	type reprioritizeGoodParamsRequest struct {
		ID int `form:"id" binding:"required,gt=0"`
	}
	var req reprioritizeGoodParamsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	type reprioritizeGoodsRequest struct {
		NewPriority int `json:"NewPriority" binding:"required"`
	}
	var body reprioritizeGoodsRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.GetGoods(c.Request.Context())
	if err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}
	httpresponse.SendSuccessOK(c, goods)
}
