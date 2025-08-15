import { useEffect, useRef, useState } from 'react';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { Message } from '../types';

export function Chat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const ws = useRef<WebSocket>(null);

  useEffect(() => {
      if (!ws.current) {
        ws.current = new WebSocket('ws://localhost:5050/ws?user_id=fslkfjslkfjslfs&chat_id=1');

      }
      console.log('useeffect')
      ws.current.onopen = () => {
        console.log('WebSocket connected');
      };

      ws.current.onmessage = (event) => {
        console.log('Raw message:', event.data);
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
          }
      };
  }, []);

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
    <div className="chat-container">
      <h1>Минималистичный Чат</h1>
      <MessageList messages={messages} />
      <MessageInput onSend={sendMessage} />
    </div>
  );
}

export default Chat;