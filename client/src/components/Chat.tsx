import { useEffect, useRef, useState } from 'react';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { MESSAGE, Message, NEW_CHAT } from '../types';
import { ChatsList } from './ChatsList';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

export function Chat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const ws = useRef<WebSocket>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate()
  const [currentChatId, setCurrentChatId] = useState<number>(1);
  const [currentUser, setCurrentUser] = useState<number>(0)
  const [isCreatingNewChat, setIsCreatingNewChat] = useState<boolean>(false);

  const fetchMEssagesHistory = async () => {
    await axios.get(`http://localhost:5050/api/chat/history/${currentChatId}`, {
      method: "GET"
    }).then((res)=> {
      if(res.status === 200) {
        res.data.result.forEach((element: any) => {
          const message = {
            id: element.id,
            text: element.message,
            sender: element.IsThisUserSender ? 'user' : 'other',
            timestamp: element.timestamp,
          } as Message
          // console.log(element)
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
        const text = event.data;
        const newMessage: Message = {
          id: Date.now().toString(),
          text,
          sender: 'other',
          timestamp: new Date(),
        };
        setMessages((prev) => [...prev, newMessage]);
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
          timestamp: new Date(),
      };
      setMessages((prev) => [...prev, newMessage]);
    }
  };

  return (
    <div className='chat-body'>
      <ChatsList 
        setCurrentChatId ={setCurrentChatId} 
        setMessages={setMessages} 
        setIsCreatingNewChat={setIsCreatingNewChat} 
        setCurrentUser={setCurrentUser}
      />
      <div className="chat-container">
        <h1>Минималистичный Чат</h1>
        <MessageList ref={messagesEndRef} messages={messages} />
        <MessageInput onSend={sendMessage} />
      </div>
    </div>
  );
}

export default Chat;