import { Message } from '../types';
import MessageComponent from './Message';

interface MessageListProps {
  messages: Message[];
}

export default function MessageList({ messages }: MessageListProps) {
  return (
    <div className="message-list">
      {messages.map((message) => (
        <MessageComponent key={message.id} message={message} />
      ))}
    </div>
  );
}