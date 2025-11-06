import { Chat } from "../types";
import '../styles/FriendCard.css'
interface ChatCardProps {
    chat: Chat
    setCurrentChatId:  React.Dispatch<React.SetStateAction<number>>
    setChatHeader:  React.Dispatch<React.SetStateAction<string>>
    currentChatId: number,
    usersOnline: number[]
}
export function ChatCard({chat, setCurrentChatId, currentChatId, setChatHeader, usersOnline}: ChatCardProps) {
    return (
        <div className={currentChatId === chat.id ? 
        "list-group-item list-group-item-action py-3 lh-sm blue" : 
        "list-group-item list-group-item-action py-3 lh-sm"} 
        key={chat.id} onClick={()=>{
                console.log(usersOnline)
                console.log(chat.userId)
                setCurrentChatId(chat.id)
                setChatHeader(chat.name)
            }}>
            <div className="user-card">    
                <div className="d-flex gap-4 align-items-center flex-wrap">
                    <div className="avatar-container">
                        <img className="avatar"/>
                        {chat.chatType === "PRIVATECHAT" ?
                            usersOnline.includes(chat.userId!) ? 
                            <div className="status-indicator status-online"></div> 
                            : <div className="status-indicator status-offline"></div> 
                         : null}
                    </div>    
                </div>
                <a href="#" className="list-group-item list-group-item-action active py-3 lh-sm" aria-current="true">
                    <div className="d-flex w-100 align-items-center justify-content-between"> 
                        <strong className="mb-1">{chat.name}</strong> 
                        </div>
                    <div className="col-10 mb-1 small">Some placeholder </div>
                </a>
            </div>
        </div>
    )
}

export default ChatCard;