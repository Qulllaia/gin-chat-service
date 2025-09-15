import { useEffect, useRef, useState } from 'react';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { Chat, MESSAGE, Message, NEW_CHAT } from '../types';
import { ChatsList } from './ChatsList';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

export function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([]);
  const ws = useRef<WebSocket>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate()
  const [currentChatId, setCurrentChatId] = useState<number>(0);
  const currentChatIdRef = useRef(currentChatId);
  const [currentUser, setCurrentUser] = useState<number>(0)
  const [chatHeader, setChatHeader] = useState<string>('Минималистичный Чат')
  const [isCreatingNewChat, setIsCreatingNewChat] = useState<boolean>(false);
  const [chats, setChats] = useState<Chat[]>([])
  
  const fetchMEssagesHistory = async () => {
    await axios.get(`http://localhost:5050/api/chat/history/${currentChatId}`, {
      method: "GET",
      withCredentials: true,
    }).then((res)=> {
      if(res.status === 200) {
        res.data.result.forEach((element: any) => {
          const message = {
            id: element.id,
            text: element.message,
            sender: element.IsThisUserSender ? 'user' : 'other',
            timestamp: element.timestamp,
          } as Message
          setMessages((messages) => [...messages, message])
        });
      } else {
        navigate('/auth', { replace: true});
      }
    }).then(()=>scrollToBottom())
    .catch((e)=> {
      if(e.status === 401){
        navigate('/auth', { replace: true});
      }
    });
  }

  const fetchChats = async () => {
      await axios.get('http://localhost:5050/api/chat/chats',
        {
          withCredentials: true,
        }
      ).then((res) => {
          if(res.data.result){
              const friendList = res.data.result.map((chat: any) => {
                  return {
                      id: chat.ID,
                      name: chat.Name,
                  } as Chat;
              })
              setChats(friendList);
          }
      })
      .catch((e)=> {
          if(e.status === 401){
              navigate('/auth', { replace: true});
          }
      });
  }

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
        messagesEndRef.current.scrollTo({
          top: messagesEndRef.current.scrollHeight,
          behavior: 'auto'
        });

    }
  };

  useEffect(() => { 
      if (!ws.current) {
        ws.current = new WebSocket('ws://localhost:5050/api/chat/ws?user_id=fslkfjslkfjslfs&chat_id=1');
      }
      // console.log('useeffect')
      ws.current.onopen = () => {
        fetchMEssagesHistory();
      };

      ws.current.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if(data.type === "MESSAGE") {
          const newMessage: Message = {
            id: Date.now().toString(),
            text: data.message,
            sender: 'other',
            chat_id: data.chat_id,
            timestamp: new Date(),
          };
          if(currentChatIdRef.current === newMessage.chat_id)
            setMessages((prev) => [...prev, newMessage]);
        }
        else if (data.type === 'NEW_CHAT') {
          setCurrentChatId(data.chat_id as number);
          fetchChats()
        }
      };

      ws.current.onclose = () => {
        console.log('WebSocket disconnected');
      };

      ws.current.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      return () => {
          if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.close();
            navigate('/auth', { replace: true});
          }
      };
  }, []);


  useEffect(()=> {
    scrollToBottom();
  }, [messages])

  useEffect(()=>{
    setMessages([]);
    fetchMEssagesHistory();
    currentChatIdRef.current = currentChatId;
  }, [currentChatId])

  const sendMessage = (text: string) => {
    if (ws.current && text.trim()) {

      if(isCreatingNewChat) {
        ws.current.send(JSON.stringify({
          type: NEW_CHAT,
          user_id: currentUser.toString(),
          messages: text
        }));
        setIsCreatingNewChat(false);
      } else {
        ws.current.send(JSON.stringify({
          type: MESSAGE,
          chat_id: currentChatId.toString(),
          messages: text
        }));
      }
        
      const newMessage: Message = {
          id: Date.now().toString(),
          text,
          sender: 'user',
          chat_id: currentChatId,
          timestamp: new Date(),
      };
      setMessages((prev) => [...prev, newMessage]);
    }
  };

  return (
    <div className='chat-body'>
      <ChatsList 
        setCurrentChatId ={setCurrentChatId} 
        currentChatId={currentChatId}
        setMessages={setMessages} 
        setIsCreatingNewChat={setIsCreatingNewChat} 
        setCurrentUser={setCurrentUser}
        setChatHeader={setChatHeader}
        fetchChats={fetchChats}
        chats={chats}
        setChats={setChats}
      />
      <div className={currentChatId === 0 && currentUser === 0 ? "chat-container-hide" : "chat-container"} >
        <h5>{chatHeader}</h5>
        <MessageList ref={messagesEndRef} messages={messages} />
        <MessageInput onSend={sendMessage} />
      </div>
    </div>
  );
}

export default ChatPage;