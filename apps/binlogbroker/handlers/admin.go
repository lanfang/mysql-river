package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lanfang/go-lib/log"
	"net/http"
)

type GetMyRoleResp struct {
	Role string `json:"role"`
}

func GetMyRole(c *gin.Context) {
	log.Info("Start GetMyRole")
	resp := GetMyRoleResp{Role: "only test, not work"}
	c.JSON(http.StatusOK, resp)
	log.Info("End GetMyRole, %+v, %+v", resp)
}
