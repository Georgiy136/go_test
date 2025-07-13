package http

import (
	"github.com/Georgiy136/go_test/auth_service/internal/http/httpresponse"
	"github.com/Georgiy136/go_test/auth_service/internal/models"
	"github.com/Georgiy136/go_test/auth_service/internal/service"
	"github.com/Georgiy136/go_test/auth_service/internal/service/app_errors"
	"github.com/gin-gonic/gin"
	"github.com/go-faster/errors"
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

	tokens, err := h.us.GetTokens(c.Request.Context(), models.DataFromRequestGetTokens{
		UserID:    req.UserID,
		UserAgent: c.Request.UserAgent(),
		IpAddress: c.ClientIP(),
	})
	if err != nil {
		httpresponse.HandleError(c, err, nil)
		return
	}
	httpresponse.SendSuccessOK(c, tokens)
}

func (h *AuthHandler) UpdateTokens(c *gin.Context) {
	type updateTokensRequest struct {
		AccessToken  string `json:"access_token" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var body updateTokensRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	tokens, err := h.us.UpdateTokens(c.Request.Context(), models.DataFromRequestUpdateTokens{
		AccessToken:  body.AccessToken,
		RefreshToken: body.RefreshToken,
		UserAgent:    c.Request.UserAgent(),
		IpAddress:    c.ClientIP(),
	})
	if err != nil {
		switch {
		case errors.Is(err, app_errors.SessionUserNotFoundError):
			httpresponse.SendFailUnauthorized(c, err.Error(), nil)
			return
		case errors.Is(err, app_errors.UserAgentNotMatchInDB):
			httpresponse.SendFailUnauthorized(c, err.Error(), nil)
			return
		default:
			httpresponse.HandleError(c, err, nil)
			return
		}
	}

	httpresponse.SendSuccessOK(c, tokens)
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	type getUserRequest struct {
		AccessToken  string `json:"access_token" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var body getUserRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	user, err := h.us.GetUser(c.Request.Context(), models.DataFromRequestGetUser{
		AccessToken:  body.AccessToken,
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		switch {
		case errors.Is(err, app_errors.TokenIsExpiredError):
			httpresponse.SendFailUnauthorized(c, err.Error(), nil)
			return
		case errors.Is(err, app_errors.SessionUserNotFoundError):
			httpresponse.SendFailUnauthorized(c, err.Error(), nil)
			return
		default:
			httpresponse.HandleError(c, err, nil)
			return
		}
	}

	httpresponse.SendSuccessOK(c, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	type logoutRequest struct {
		AccessToken  string `json:"access_token" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var body logoutRequest
	if err := c.BindJSON(&body); err != nil {
		httpresponse.SendFailBadRequest(c, err.Error(), nil)
		return
	}

	err := h.us.Logout(c.Request.Context(), models.DataFromRequestLogout{
		AccessToken:  body.AccessToken,
		RefreshToken: body.RefreshToken,
		UserAgent:    c.Request.UserAgent(),
	})
	if err != nil {
		switch {
		case errors.Is(err, app_errors.SessionUserNotFoundError):
			httpresponse.SendFailUnauthorized(c, err.Error(), nil)
			return
		case errors.Is(err, app_errors.UserAgentNotMatchInDB):
			httpresponse.SendFailUnauthorized(c, err.Error(), nil)
			return
		default:
			httpresponse.HandleError(c, err, nil)
			return
		}
	}

	httpresponse.SendNoContent(c)
}
