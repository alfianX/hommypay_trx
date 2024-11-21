package handlers

import (
	"net/http"

	"github.com/alfianX/hommypay_trx/helper"
	"github.com/gin-gonic/gin"
)

func (s service) Health(c *gin.Context) {
	helper.Respond(c, gin.H{"Message":"App OK"}, http.StatusOK)
}