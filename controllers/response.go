package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var res response

type response struct {
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *response) json(c *gin.Context, msg interface{}, data interface{}) {
	c.JSON(http.StatusOK, response{
		Msg:  msg,
		Data: data,
	})

}

func (r *response) text(c *gin.Context, data string) {
	c.String(http.StatusOK, data)
}
