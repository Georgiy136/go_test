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
//	@Tags			Goodss
//	@Description	create Goods
//	@ID				create-Goods
//	@Accept			json
//	@Produce		json
//	@Param			input	body		PostGoodsRequest	true	"Goods info"
//	@Success		201		{object}	models.Goods
//	@Router			/Goods [post]
func (h *GoodsHandler) PostGoods(c *gin.Context) {
	type PostGoodsRequest struct {
		Id         uuid.UUID `json:"id"`
		FirstName  string    `json:"firstName" binding:"required"`
		LastName   string    `json:"lastName" binding:"required"`
		Patronymic string    `json:"patronymic" binding:"required"`
		City       string    `json:"city" binding:"required"`
		Phone      string    `json:"phone" binding:"required"`
		Email      string    `json:"email" binding:"required"`
		Password   string    `json:"password"`
	}

	postGoodsRequest := &PostGoodsRequest{}

	if err := c.Bind(postGoodsRequest); err != nil {
		httpresponse.SendFailBadRequest(c, nil)
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

	Goods, err := h.us.AddGoods(c.Request.Context(), *Goods)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	httpresponse.SendSuccess(c, http.StatusCreated, Goods)
}

// GetOneGoods godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves Goods based on given ID
//	@Tags			Goodss
//	@Description	get Goods by id
//	@ID				get-Goods-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Goods ID"
//	@Success		202	{object}	models.Goods
//	@Router			/Goods/{id} [get]
func (h *GoodsHandler) GetOneGood(c *gin.Context) {
	projects, err := h.us.GetOneGood(c.Request.Context(), "1")
	if err != nil {
		logrus.Error(err)
		httpresponse.SendFailBadRequest(c, nil)
		return
	}
	httpresponse.SendSuccessOK(c, projects)
}

// GetAllGoodss godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Retrieves All Goodss
//	@Tags			Goodss
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
		httpresponse.SendFailBadRequest(c, nil)
		return
	}
	httpresponse.SendSuccessOK(c, Goodss)
}

// DeleteGoods godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Delete Goods based on given ID
//	@Tags			Goodss
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
		httpresponse.SendFailBadRequest(c, nil)
		return
	}
	httpresponse.SendNoContent(c)
}

// UpdateGoods godoc
//
//	@Security		ApiKeyAuth
//	@Summary		Update Goods based on given ID
//	@Tags			Goodss
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
		httpresponse.SendFailBadRequest(c, nil)
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
		httpresponse.SendFailBadRequest(c, nil)
		return
	}
	httpresponse.SendNoContent(c)
}
