package websockets

import (
	. "main/database/models"
	"main/database/queries"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

type ConnectorActor struct {
	mailbox  chan mailboxObject
	stopChan chan struct{}
	wg       sync.WaitGroup
	user_id_conns_with_chat_ids  map[int]UserWSData
	clientsMu sync.RWMutex
	WSQ *queries.WSQueries
}

func NewConnectorActor(wsq *queries.WSQueries) *ConnectorActor {
    actor := &ConnectorActor{
        mailbox:  make(chan mailboxObject, 100),
        stopChan: make(chan struct{}),
        user_id_conns_with_chat_ids:  make(map[int]UserWSData),
		WSQ: wsq,
    }
    
    actor.wg.Add(1)
    go actor.processMessages()
    
    return actor
}

func (a *ConnectorActor) processMessages() {
    defer a.wg.Done()
    
    for {
        select {
			case msg := <-a.mailbox:
				a.handleMessage(msg.message, msg.MessageType, msg.conn)
			case <-a.stopChan:
				return
        }
    }
}

func (a *ConnectorActor) handleMessage(msg MessageWS, messageType int, conn *websocket.Conn) {
	a.broadcastMessage(msg, messageType, conn);
}

func (a *ConnectorActor) Send(msg MessageWS, messageType int, conn *websocket.Conn) {
    a.mailbox <- mailboxObject{message: msg, MessageType: messageType, conn: conn }
}

func (a *ConnectorActor) Stop() {
    close(a.stopChan)
    a.wg.Wait()
}

func (a *ConnectorActor) AddClient(conn *websocket.Conn, user_id int) {
    a.clientsMu.Lock()
    defer a.clientsMu.Unlock()

	var user User

	if err := a.WSQ.GetUserChatList(user_id, &user); err != nil {
		println("GetUserChatListError", err.Error())
	}
	
    a.user_id_conns_with_chat_ids[user_id] = UserWSData{
		chat_ids: user.Chat_list,
		ws: conn,
	}
}

func (a *ConnectorActor) RemoveClient(user_id int) {
    a.clientsMu.Lock()
    defer a.clientsMu.Unlock()
    delete(a.user_id_conns_with_chat_ids, user_id)
}

func (a *ConnectorActor) broadcastMessage(message MessageWS, messageType int, conn *websocket.Conn) {
	var getter_conn *websocket.Conn;
	for user_id, value := range a.user_id_conns_with_chat_ids {
		chat_id, _ := strconv.Atoi(message.Chat_id); 

		if contains(value.chat_ids, int64(chat_id)) {
			getter_conn = value.ws;
			if getter_conn != conn {
				if err := getter_conn.WriteMessage(messageType, []byte(message.Message)); err != nil {
					println(err.Error());	
				}
			} else {
		
				if err := a.WSQ.InsertMessageIntoChatHistory(chat_id, user_id, string(message.Message)); err != nil {
					println(err.Error());	
				}
			}
		}
	}
}

func contains(slice pq.Int64Array, item int64) bool {
    for _, element := range slice {
        if element == item {
            return true
        }
    }
    return false
}