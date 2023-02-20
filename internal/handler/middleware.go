package handler

import (
	customerrors "miborchestrator/internal/entities/custom_errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func (h *Handler) useridentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		customerrors.NewErorResponse(c, http.StatusUnauthorized, "error:empty auth header")
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		customerrors.NewErorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	userID, err := h.Service.ParseToken(headerParts[1])
	if err != nil {
		customerrors.NewErorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userID)

}
