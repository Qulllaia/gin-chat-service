package websockets

import (
	. "main/database/models"
	"main/database/queries"
	"main/redis"
	. "main/types"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectorActor struct {
	mailbox  chan MailboxObject
	stopChan chan struct{}
	wg       sync.WaitGroup
	User_id_conns_with_chat_ids  map[int]UserWSData
	conns_to_user_ids map[*websocket.Conn]int
	clientsMu sync.RWMutex
	WSQ *queries.WSQueries
	RDB *redis.RedisConnector
	handlerManager *HandlerManager
}

func NewConnectorActor(wsq *queries.WSQueries, rdb *redis.RedisConnector) *ConnectorActor {
    actor := &ConnectorActor{
        mailbox:  make(chan MailboxObject, 100),
        stopChan: make(chan struct{}),
        User_id_conns_with_chat_ids:  make(map[int]UserWSData),
        conns_to_user_ids:  make(map[*websocket.Conn]int),
		WSQ: wsq,
		RDB: rdb,
    }

	handlerManager := NewHandlerManager(actor)

	actor.handlerManager = handlerManager;
    
    actor.wg.Add(1)
    go actor.processMessages()
    
    return actor
}

func (a *ConnectorActor) processMessages() {
    defer a.wg.Done()
    
    for {
        select {
			case msg := <-a.mailbox:
				a.handleMessage(msg.Message, msg.MessageType, msg.Conn)
			case <-a.stopChan:
				return
        }
    }
}

func (a *ConnectorActor) GetUserConnections() map[int]UserWSData {
	result := make(map[int]UserWSData)
	for k, v := range a.User_id_conns_with_chat_ids {
		result[k] = UserWSData{
			WS:      v.WS,
			Chat_ids: v.Chat_ids,
		}
	}
	return result
}

func (a *ConnectorActor) GetConnectoinsToUsers() map[*websocket.Conn]int {
	result := make(map[*websocket.Conn]int)
	for k, v := range a.conns_to_user_ids {
		result[k] = v
	}
	return result
}

func (a *ConnectorActor) GetWSQ() *queries.WSQueries {
	return a.WSQ
}

func (a *ConnectorActor) GetRDB() *redis.RedisConnector {
	return a.RDB
}

func (a *ConnectorActor) handleMessage(msg MessageWS, messageType int, conn *websocket.Conn) {
	switch(msg.Type) {
		case "MESSAGE": 
			a.handlerManager.handlers["MESSAGE"].Handle(msg, messageType, conn, a)
		case "NEW_CHAT": 	
			a.handlerManager.handlers["NEW_CHAT"].Handle(msg, messageType, conn, a)
		case "NEW_MULTIPLE_CHAT":
			a.handlerManager.handlers["NEW_MULTIPLE_CHAT"].Handle(msg, messageType, conn, a)	
		case "USER_STATUS":
			a.handlerManager.handlers["USER_STATUS"].Handle(msg, messageType, conn, a)	
		}
}

func (a *ConnectorActor) Send(msg MessageWS, messageType int, conn *websocket.Conn) {
    a.mailbox <- MailboxObject{Message: msg, MessageType: messageType, Conn: conn }
}

func (a *ConnectorActor) Stop() {
    close(a.stopChan)
    a.wg.Wait()
}

func (a *ConnectorActor) AddClient(conn *websocket.Conn, user_id int) {
    a.clientsMu.Lock()
    defer a.clientsMu.Unlock()

	var user User = User{ ID: user_id}

	if err := a.WSQ.GetUserChatList(user_id, &user); err != nil {
		println("GetUserChatListError", err.Error())
	}
	
    a.User_id_conns_with_chat_ids[user_id] = UserWSData{
		Chat_ids: user.Chat_list,
		WS: conn,
	}

	a.conns_to_user_ids[conn] = user_id;
}

func (a *ConnectorActor) RemoveClient(user_id int) {
    a.clientsMu.Lock()
    defer a.clientsMu.Unlock()
    delete(a.User_id_conns_with_chat_ids, user_id)
}