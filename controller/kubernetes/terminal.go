package kubernetes

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ops-api/kubernetes"
	service "ops-api/service/kubernetes"
)

var PodTerminal podTerminal

type podTerminal struct{}

func wsInit(c *gin.Context) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return upgrader.Upgrade(c.Writer, c.Request, nil)
}

func (p *podTerminal) Init(c *gin.Context) {

	params := new(struct {
		UUID          string `form:"uuid" binding:"required"`
		PodName       string `form:"podName" binding:"required"`
		ContainerName string `form:"containerName" binding:"required"`
		Namespace     string `form:"namespace" binding:"required"`
		Cols          int    `form:"cols,default=120"`
		Rows          int    `form:"rows,default=30"`
	})

	if err := c.Bind(params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	ws, err := wsInit(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)

	client := c.MustGet("kc").(*kubernetes.ClientList)
	if err := service.PodTerminal.Init(params.Namespace, params.PodName, params.ContainerName, params.Cols, params.Rows, ws, client); err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
		return
	}
}
