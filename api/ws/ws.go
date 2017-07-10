package ws

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dokkur/swanager/core/auth"
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/lib"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type authMessage struct {
	Token string `json:"token"`
}

type answer struct {
	AnswerType string      `json:"type"`
	Data       interface{} `json:"data"`
}

type clientConnection struct {
	ID        string
	State     string
	User      *entities.User
	Conn      *websocket.Conn
	AuthError error
	Incoming  chan interface{}
}

const (
	stateWorking         = "working"
	stateUnauthenticated = "unauthenticated"

	answerTypeData          = "data"
	answerTypeError         = "error"
	answerTypeAuthenticated = "authenticated"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// clients[UserID][ConnectionID]
var clients = make(map[string]map[string]*clientConnection)

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
		ID:    lib.GenerateUUID(),
		State: stateUnauthenticated,
		Conn:  conn,
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log().Debugf("ws Error: %s", err.Error())
			if context.State == stateWorking {
				removeClient(&context)
			}
			break
		}

		log().Debugf("[%s] WS message: [%d] %s", context.State, msgType, msg)

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
		case data := <-c.Incoming:
			c.sendAnswer(answer{
				AnswerType: answerTypeData,
				Data:       data,
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

	incoming := make(chan interface{}, 10)
	c.Incoming = incoming
	c.State = stateWorking

	addClient(c)

	c.sendAnswer(answer{
		AnswerType: answerTypeAuthenticated,
		Data:       "Ok",
	})

	go c.listen()
}

func (c *clientConnection) authError() {
	log().Debugf("Auth error: %s", c.AuthError.Error())

	c.sendAnswer(answer{
		AnswerType: answerTypeError,
		Data:       c.AuthError.Error(),
	})

	c.Conn.Close()
}

func (c *clientConnection) sendAnswer(ans answer) {
	result, _ := json.Marshal(ans)
	c.Conn.WriteMessage(1, result)
}

func addClient(c *clientConnection) {
	if _, ok := clients[c.User.ID]; !ok {
		clients[c.User.ID] = make(map[string]*clientConnection)
	}
	clients[c.User.ID][c.ID] = c
}

func removeClient(c *clientConnection) {
	delete(clients[c.User.ID], c.ID)
	if len(clients[c.User.ID]) == 0 {
		delete(clients, c.User.ID)
	}
}

// Send - sends data to connected client
func Send(userID string, data interface{}) {
	if clientConnections, ok := clients[userID]; ok {
		for _, connection := range clientConnections {
			connection.Incoming <- data
		}
	}
}

// IsUserConnected returns true is user with UserID is connected
func IsUserConnected(userID string) bool {
	_, connected := clients[userID]
	return connected
}

func log() *logrus.Entry {
	return logrus.WithField("module", "api.ws")
}
