package ws

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// InitWS add ws handler for api
func InitWS(router *gin.Engine) {
	router.GET("/ws", wsHandler)
}

func wsHandler(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Warnf("Failed to set websocket upgrade %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	defer conn.Close()

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		logrus.Debugf("Got ws message type=%d %s", t, msg)
		conn.WriteMessage(t, msg)
	}
}
