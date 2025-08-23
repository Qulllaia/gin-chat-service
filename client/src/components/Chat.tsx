import { useEffect, useRef, useState } from 'react';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { Message } from '../types';
import FriendsList from './FriendsList';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

export function Chat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const ws = useRef<WebSocket>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate()

  const fetchMEssagesHistory = async () => {
    await axios.get("http://localhost:5050/api/chat/history", {
      method: "GET"
    }).then((res)=> {
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
    }).then(()=>scrollToBottom());
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

  const sendMessage = (text: string) => {
    if (ws.current && text.trim()) {
      ws.current.send(text);
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
      <FriendsList/>
      <div className="chat-container">
        <h1>Минималистичный Чат</h1>
        <MessageList ref={messagesEndRef} messages={messages} />
        <MessageInput onSend={sendMessage} />
      </div>
    </div>
  );
}

export default Chat;