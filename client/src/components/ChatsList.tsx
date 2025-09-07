import { useEffect, useState } from "react";
import ChatCard from "./ChatCard";
import axios from "axios";
import { Chat, User } from "../types";
import { ParentForm } from "./ParentForm";
import UserCard from "./UserCard";
import { useNavigate } from "react-router-dom";

export function ChatsList({setCurrentChatId, setMessages, setCurrentUser, setIsCreatingNewChat}:any) {

    const [chats, setChats] = useState<Chat[]>([])
    const [users, setUsers] = useState<User[]>([])
    const navigate = useNavigate()



    const [isChatCreationOpen, setIsChatCreatonOpen] = useState<boolean>(false);

    const fetchChats = async () => {
        await axios.get('http://localhost:5050/api/chat/chats').then((res) => {
            if(res.data.result){
                const friendList = res.data.result.map((chat: any) => {
                    return {
                        id: chat.ID,
                        name: chat.Name,
                    } as Chat;
                })
                setChats(friendList);
            }
        })
        .catch((e)=> {
            if(e.status === 401){
                navigate('/auth', { replace: true});
            }
        });
    }

    const fetchCurrentUsers = async () => {
        await axios.get('http://localhost:5000/api/user/get/except')
        .then((res) => {
            if(res.data.result){
                const friendList = res.data.result.map((user: any) => {
                    console.log(user)
                    return {
                        id: user.id,
                        name: user.name,
                    } as User;
                })
                setUsers(friendList);
            }
        })
        .catch((e)=> {
            if(e.status === 401){
                navigate('/auth', { replace: true});
            }
        });
    }

    const createChat = async (user_id: number) => {
        setMessages([])
        setIsChatCreatonOpen(false);
        setIsCreatingNewChat(true);
        setCurrentUser(user_id)
    }

    useEffect(()=> {
        fetchChats()
    },[])

    useEffect(()=> {
        if(isChatCreationOpen) fetchCurrentUsers()
    },[isChatCreationOpen])

    return (
        <div className="list">
            <ParentForm
                isDialog ={true}
                isOpen={isChatCreationOpen}
                setIsOpen={setIsChatCreatonOpen}
            >
                {users.map((val)=>{
                    return <UserCard user={val} createChatHandler={createChat}></UserCard>
                })}
            </ParentForm>
            <button onClick={()=>setIsChatCreatonOpen(true)}>Создать чат</button>
            {chats.map((Chat) => (
                <ChatCard key={Chat.id} friend ={Chat} setCurrentChatId={setCurrentChatId}/>
            ))}
        </div>
    )
}

export default ChatsList;