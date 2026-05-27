import { Chat } from "../types";
import '../styles/FriendCard.css'

interface ChatCardProps {
  chat: Chat
  setCurrentChatId: React.Dispatch<React.SetStateAction<number>>
  setChatHeader: React.Dispatch<React.SetStateAction<string>>
  currentChatId: number,
  usersOnline: number[]
  onDeleteChat: (chatId: number, e: React.MouseEvent) => void
}

export function ChatCard({ chat, setCurrentChatId, currentChatId, setChatHeader, usersOnline, onDeleteChat }: ChatCardProps) {
  return (
    <div className={currentChatId === chat.id ?
      "list-group-item list-group-item-action py-3 lh-sm blue" :
      "list-group-item list-group-item-action py-3 lh-sm"}
      key={chat.id} onClick={() => {
        setCurrentChatId(chat.id)
        setChatHeader(chat.name)
      }}>
      <div className="user-card chat-card-item">
        <div className="d-flex gap-4 align-items-center flex-wrap">
          <div className="avatar-container">
            <img className="avatar" alt="" />
            {chat.chatType === "PRIVATECHAT" ?
              usersOnline.includes(chat.userId!) ?
                <div className="status-indicator status-online"></div>
                : <div className="status-indicator status-offline"></div>
              : null}
          </div>
        </div>
        <div className="chat-card-body list-group-item list-group-item-action active py-3 lh-sm" aria-current="true">
          <div className="d-flex w-100 align-items-center justify-content-between">
            <strong className="mb-1">{chat.name}</strong>
          </div>
          <div className="col-10 mb-1 small">{chat.lastMessage}</div>
        </div>
        <button
          type="button"
          className="chat-card-delete"
          title="Удалить чат"
          aria-label="Удалить чат"
          onClick={(e) => onDeleteChat(chat.id, e)}
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden>
            <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>
  )
}

export default ChatCard;
