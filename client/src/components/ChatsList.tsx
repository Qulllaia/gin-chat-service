import { useEffect, useState } from "react";
import ChatCard from "./ChatCard";
import axios from "axios";
import { Chat } from "../types";

export function ChatsList({setCurrentChatId}:any) {

    const [chats, setChats] = useState<Chat[]>([])

    const fetchFriends = async () => {
        await axios.get('http://localhost:5050/api/chat/chats').then((res) => {
            if(res.data.result){
                const friendList = res.data.result.map((user: any) => {
                    return {
                        id: user.ID,
                        name: user.Name,
                    } as Chat;
                })
                setChats(friendList);
            }
        })
    }

    useEffect(()=> {
        fetchFriends()
    },[])

    return (
        <div className="list">
            {chats.map((Chat) => (
                <ChatCard key={Chat.id} friend ={Chat} setCurrentChatId={setCurrentChatId}/>
            ))}
        </div>
    )
}

export default ChatsList;