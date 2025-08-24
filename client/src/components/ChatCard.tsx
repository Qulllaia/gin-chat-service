import { Chat } from "../types";
import '../styles/FriendCard.css'
interface ChatCardProps {
    friend: Chat
    setCurrentChatId:  React.Dispatch<React.SetStateAction<number>>
}

export function ChatCard({friend, setCurrentChatId}: ChatCardProps) {
    return (
        <div key={friend.id} className="friend-card"
            onClick={()=>setCurrentChatId(friend.id)}
        >
            <p className="friend-p">{friend.id}</p>
            <p className="friend-p">{friend.name}</p>
        </div>
    )
}

export default ChatCard;