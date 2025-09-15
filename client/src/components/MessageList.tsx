import { Message } from '../types';
import MessageComponent from './Message';

interface MessageListProps {
  messages: Message[];
  ref: React.RefObject<HTMLDivElement | null>
}

export default function MessageList({ messages, ref }: MessageListProps) {
  return (
    <div ref={ref} className="message-list">
      {messages.map((message) => (
        <MessageComponent key={message.id +  Math.random().toString(16).slice(2)} message={message} />
      ))}
    </div>
  );
}