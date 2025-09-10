package websockets

import (
	"encoding/json"
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
	conns_to_user_ids map[*websocket.Conn]int
	clientsMu sync.RWMutex
	WSQ *queries.WSQueries
}

func NewConnectorActor(wsq *queries.WSQueries) *ConnectorActor {
    actor := &ConnectorActor{
        mailbox:  make(chan mailboxObject, 100),
        stopChan: make(chan struct{}),
        user_id_conns_with_chat_ids:  make(map[int]UserWSData),
        conns_to_user_ids:  make(map[*websocket.Conn]int),
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
	switch(msg.Type) {
		case "MESSAGE": 
			a.broadcastMessage(msg, messageType, conn);
		case "NEW_CHAT": 
			a.createChatWithMessage(msg, conn, messageType);
			a.broadcastMessage(msg, messageType, conn);
	}
}

func (a *ConnectorActor) createChatWithMessage(msg MessageWS, conn *websocket.Conn, messageType int) {
	current_user_id := a.conns_to_user_ids[conn];
	to_user_id, _ := strconv.Atoi(msg.User_id);


	err, chat_id := a.WSQ.CreateChatWithMessage(current_user_id, to_user_id);
	if err != nil {
		println("Error while creating a chat", err.Error());
	}

	responseData := map[string]interface{} {
		"type": "NEW_CHAT",
		"chat_id": chat_id,
	}

	chatCreationResponse, err := json.Marshal(responseData);

	if err != nil {
		println("Ошибка формирования ответа при создании нового чата", err.Error())
	}

	a.WSQ.InsertMessageIntoChatHistory(chat_id, current_user_id, msg.Message);

    a.clientsMu.Lock()
    defer a.clientsMu.Unlock()

	var current_user User = User{ ID: current_user_id}

	if err := a.WSQ.GetUserChatList(current_user_id, &current_user); err != nil {
		println("GetUserChatListError", err.Error())
	}
	
    a.user_id_conns_with_chat_ids[current_user_id] = UserWSData{
		chat_ids: current_user.Chat_list,
		ws: conn,
	}

	conn.WriteMessage(messageType, []byte(chatCreationResponse))

	if toUserData, exists := a.user_id_conns_with_chat_ids[to_user_id]; exists  {
		toUserConn := toUserData.ws;

		var to_user User = User{ ID: to_user_id}

		if err := a.WSQ.GetUserChatList(to_user_id, &to_user); err != nil {
			println("GetUserChatListError", err.Error())
		}
		
		a.user_id_conns_with_chat_ids[to_user_id] = UserWSData{
			chat_ids: to_user.Chat_list,
			ws: toUserConn,
		}
		toUserConn.WriteMessage(messageType, []byte(chatCreationResponse))
	}
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

	var user User = User{ ID: user_id}

	if err := a.WSQ.GetUserChatList(user_id, &user); err != nil {
		println("GetUserChatListError", err.Error())
	}
	
    a.user_id_conns_with_chat_ids[user_id] = UserWSData{
		chat_ids: user.Chat_list,
		ws: conn,
	}

	a.conns_to_user_ids[conn] = user_id;
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
				responseData, err := json.Marshal(map[string]interface{} {
					"message": message.Message,
					"type": "MESSAGE",
				})
				if err = getter_conn.WriteMessage(messageType, []byte(responseData)); err != nil {
					println(err.Error());	
				}
			} else {
				println(chat_id);
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