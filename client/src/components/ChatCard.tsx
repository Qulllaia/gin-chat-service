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
        "list-group-item list-group-item-action py-3 lh-sm blue" : 
        "list-group-item list-group-item-action py-3 lh-sm"} 
        key={friend.id} onClick={()=>{
                setCurrentChatId(friend.id)
                setChatHeader(friend.name)
            }}>
            <div className="user-card">    
                <div className="d-flex gap-4 align-items-center flex-wrap">
                    <div className="avatar-container">
                        <img src="https://via.placeholder.com/80" className="avatar"/>
                        <div className="status-indicator status-online"></div>
                    </div>    
                </div>
                <a href="#" className="list-group-item list-group-item-action active py-3 lh-sm" aria-current="true">
                    <div className="d-flex w-100 align-items-center justify-content-between"> 
                        <strong className="mb-1">{friend.name}</strong> 
                        </div>
                    <div className="col-10 mb-1 small">Some placeholder </div>
                </a>
            </div>
        </div>
    )
}

export default ChatCard;