import { useEffect, useState } from "react";
import ChatCard from "./ChatCard";
import axios from "axios";
import { Chat, User } from "../types";
import { ParentForm } from "./ParentForm";
import UserCard from "./UserCard";
import { useNavigate } from "react-router-dom";

export function ChatsList({setCurrentChatId, currentChatId, setMessages, setCurrentUser, setIsCreatingNewChat, setChatHeader, fetchChats, chats, setChats}:any) {
    const [users, setUsers] = useState<User[]>([])
    const navigate = useNavigate()

    const [isChatCreationOpen, setIsChatCreatonOpen] = useState<boolean>(false);

    const fetchCurrentUsers = async () => {
        await axios.get('http://localhost:5000/api/user/get/except',{
          withCredentials: true,
        })
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
            <div className="d-grid vh-100  w-100" style={{gridTemplateRows: 'auto 1fr'}}>
                <button 
                    className="btn btn-outline-secondary m-2" 
                    type="button" 
                    onClick={() => setIsChatCreatonOpen(true)}
                >
                    Создать чат
                </button>
                <div className="chat-list">
                    <div className="d-flex flex-column align-items-stretch bg-body-tertiary w-100 h-100">
                        <a href="/" className="d-flex align-items-center flex-shrink-0 p-3 link-body-emphasis text-decoration-none border-bottom">
                            <svg className="bi pe-none me-2" width="30" height="24" aria-hidden="true">
                                <use href="#bootstrap"></use>
                            </svg>
                            <span className="fs-5 fw-semibold">Список чатов</span> 
                        </a>
                        {chats.map((Chat:Chat) => (
                            <ChatCard key={Chat.id} friend={Chat} setCurrentChatId={setCurrentChatId} currentChatId={currentChatId} setChatHeader={setChatHeader}/>
                        ))}
                    </div>
                </div>
            </div>
        </div>

    )
}

export default ChatsList;