import { Message } from '../types';
import MessageComponent from './Message';

interface MessageListProps {
  messages: Message[];
  ref: React.RefObject<HTMLDivElement | null>
}

export default function MessageList({ messages, ref }: MessageListProps) {
  return (
    <div
      ref={ref}
      className={`message-list${messages.length === 0 ? ' message-list--empty' : ''}`}
    >
      {messages.length === 0 ? (
        <div className="chat-messages-empty">
          <div className="chat-messages-empty-icon" aria-hidden>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75">
              <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </div>
          <p className="chat-messages-empty-title">В этом чате сообщений нет!</p>
          <p className="chat-messages-empty-text">
            Поспешите поздороваться с вашим другом
          </p>
        </div>
      ) : (
        messages.map((message) => (
          <MessageComponent key={message.id + Math.random().toString(16).slice(2)} message={message} />
        ))
      )}
    </div>
  );
}