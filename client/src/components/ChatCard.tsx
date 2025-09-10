import { Chat } from "../types";
import '../styles/FriendCard.css'
interface ChatCardProps {
    friend: Chat
    setCurrentChatId:  React.Dispatch<React.SetStateAction<number>>
    setChatHeader:  React.Dispatch<React.SetStateAction<string>>
    currentChatId: number
}
export function ChatCard({friend, setCurrentChatId, currentChatId, setChatHeader}: ChatCardProps) {
    return (
        <div className={currentChatId === friend.id ? 
        "list-group list-group-flush border-bottom scrollarea" : 
        "list-group-item list-group-item-action py-3 lh-sm"} 
        key={friend.id} onClick={()=>{
                setCurrentChatId(friend.id)
                setChatHeader(friend.name)
            }}>
            <a href="#" className="list-group-item list-group-item-action active py-3 lh-sm" aria-current="true">
                <div className="d-flex w-100 align-items-center justify-content-between"> 
                    <strong className="mb-1">{friend.name}</strong> 
                    <small>{friend.id}</small> 
                    </div>
                <div className="col-10 mb-1 small">Some placeholder </div>
            </a>
        </div>
    )
}

export default ChatCard;