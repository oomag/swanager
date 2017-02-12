package ws

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/swanager/core/auth"
	"github.com/da4nik/swanager/core/entities"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type authMessage struct {
	Token string `json:"token"`
}

type answer struct {
	AnswerType string `json:"type"`
	Data       string
	Service    *entities.Service `json:"service,omitempty"`
}

type clientConnection struct {
	State     string
	User      *entities.User
	Conn      *websocket.Conn
	AuthError error
	Incoming  chan entities.Service
}

const (
	stateWorking         = "working"
	stateUnauthenticated = "unauthenticated"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients = make(map[string]*clientConnection)

// InitWS add ws handler for api
func InitWS(router *gin.Engine) {
	router.GET("/ws", wsHandler)
}

func wsHandler(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log().Warnf("Failed to set websocket upgrade %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	defer conn.Close()

	context := clientConnection{
		State: stateUnauthenticated,
		Conn:  conn,
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			if context.State == stateWorking {
				delete(clients, context.User.ID)
			}
			break
		}

		log().Debugf("[%s] WS message type=%d %s", context.State, t, msg)

		switch context.State {
		case stateUnauthenticated:
			context.authenticate(msg)
			break
		case stateWorking:
			context.processMessage(msg)
		}
	}
}

func (c *clientConnection) listen() {
	for {
		select {
		case service := <-c.Incoming:
			log().WithField("UserID", c.User.ID).Debugf("Senging to client %s", service.NSName)
			if service.Name == "" {
				return
			}

			c.sendAnswer(answer{
				AnswerType: "data",
				Service:    &service,
			})
		}
	}
}

func (c *clientConnection) processMessage(msg []byte) {
	c.Conn.WriteMessage(1, msg)
}

func (c *clientConnection) authenticate(msg []byte) {
	var message authMessage
	c.AuthError = json.Unmarshal(msg, &message)
	if c.AuthError != nil {
		c.authError()
		return
	}

	c.User, c.AuthError = auth.WithToken(message.Token)
	if c.AuthError != nil {
		c.authError()
		return
	}

	log().Debugf("Authenticated (%s), proceeding with normal mode", c.User.Email)

	c.State = stateWorking
	incoming := make(chan entities.Service, 10)
	c.Incoming = incoming

	clients[c.User.ID] = c

	c.sendAnswer(answer{
		AnswerType: "authenticated",
		Data:       "Ok",
	})

	go c.listen()
}

func (c *clientConnection) authError() {
	log().Debugf("Auth error: %s", c.AuthError.Error())

	c.sendAnswer(answer{
		AnswerType: "error",
		Data:       c.AuthError.Error(),
	})

	c.Conn.Close()
}

func (c *clientConnection) sendAnswer(ans answer) {
	result, _ := json.Marshal(ans)
	c.Conn.WriteMessage(1, result)
}

func log() *logrus.Entry {
	return logrus.WithField("module", "api.ws")
}
