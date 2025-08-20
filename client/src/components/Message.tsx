import { Message } from '../types';

interface MessageProps {
  message: Message;
}

export default function MessageComponent({ message }: MessageProps) {
  return (
    <div className={`message ${message.sender}`}>
      <div className="message-text">{message.text}</div>
      <div className="message-time">
        {new Date(message.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
      </div>
    </div>
  );
}