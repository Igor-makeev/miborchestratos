package handler

import (
	"miborchestrator/internal/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initTransfer(c *gin.Context) {
	var input entities.TransactionRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := entities.Request{
		TransactionRequest: input,
		Event:              entities.InitTx,
	}
	h.Service.AddtoTxManagerQueue(req)
	c.Status(http.StatusAccepted)
}

func (h *Handler) sendToTxManager(c *gin.Context) {
	var input entities.TxDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req := entities.Request{
		TxDTO: input,
		Event: input.DTOEvent,
	}
	h.Service.AddtoTxManagerQueue(req)
	c.Status(http.StatusAccepted)
}

func (h *Handler) sendToCreateQueue(c *gin.Context) {
	userID, err := c.Get(userCtx)
	if !err {
		return
	}

	var input entities.CreateWalletRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := entities.WcTDO{UserID: userID.(int), WalletID: input.WalletID}
	c.Status(http.StatusAccepted)
	h.Service.AddtoCreateWalletQueue(&data)

}
