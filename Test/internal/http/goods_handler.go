package http

import (
	"github.com/sirupsen/logrus"
	"myapp/internal/http/httpresponse"
	"myapp/internal/models"
	"myapp/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
//	@Param			input	body		PostGoodsRequest	true	"Goods info"
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

	type PostGoodsRequest struct {
		Name string `json:"name" binding:"required"`
	}
	var body PostGoodsRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	goods, err := h.us.AddGoods(c.Request.Context(), req.ProjectID, body.Name)
	if err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}
	httpresponse.SendSuccess(c, http.StatusCreated, goods)
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
	Goodss, err := h.us.GetAllGoods(c.Request.Context())
	if err != nil {
		logrus.Error(err)
		httpresponse.SendFailBadRequest(c, "", nil)
		return
	}
	httpresponse.SendSuccessOK(c, Goodss)
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
	id := c.Param("id")
	err := h.us.DeleteGood(c.Request.Context(), id)
	if err != nil {
		logrus.Error(err)
		httpresponse.SendFailBadRequest(c, "", nil)
		return
	}
	httpresponse.SendNoContent(c)
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

	type PutGoodsRequest struct {
		Id         uuid.UUID `json:"id"`
		FirstName  string    `json:"firstName" binding:"required"`
		LastName   string    `json:"lastName" binding:"required"`
		Patronymic string    `json:"patronymic" binding:"required"`
		City       string    `json:"city" binding:"required"`
		Phone      string    `json:"phone" binding:"required"`
		Email      string    `json:"email" binding:"required"`
		Password   string    `json:"password"`
	}

	id := c.Param("id")

	postGoodsRequest := &PutGoodsRequest{}

	if err := c.Bind(postGoodsRequest); err != nil {
		httpresponse.SendFailBadRequest(c, "", nil)
		return
	}

	Goods := &models.Goods{
		Id:         postGoodsRequest.Id,
		FirstName:  postGoodsRequest.FirstName,
		LastName:   postGoodsRequest.LastName,
		Patronymic: postGoodsRequest.Patronymic,
		City:       postGoodsRequest.City,
		Phone:      postGoodsRequest.Phone,
		Email:      postGoodsRequest.Email,
		Password:   postGoodsRequest.Password,
	}

	Goods, err := h.us.UpdateGood(c.Request.Context(), id, *Goods)
	if err != nil {
		logrus.Error(err)
		httpresponse.SendFailBadRequest(c, "", nil)
		return
	}
	httpresponse.SendNoContent(c)
}
