import { useEffect, useState } from "react";
import ChatCard from "./ChatCard";
import axios from "axios";
import { Chat, User } from "../types";
import { ParentForm } from "./ParentForm";
import UserCard from "./UserCard";
import { useNavigate } from "react-router-dom";

export function ChatsList({
        setCurrentChatId, 
        currentChatId, 
        setMessages, 
        setCurrentUser, 
        setIsCreatingNewChat, 
        setChatHeader, 
        fetchChats, 
        chats, 
        setChats,
        sendMultipleChatCreationNotify,
        usersOnline 
    }:any) {
   
    const [users, setUsers] = useState<User[]>([])
    const [selectedUsers, setSelectedUsers] = useState<number[]>([]);
    const navigate = useNavigate()
    
    const [groupName, setGroupName] = useState<string>('');
    const [isChatCreationOpen, setIsChatCreatonOpen] = useState<boolean>(false);
    const [isMultiChatCreationOpen, setIsMultiChatCreationOpen] = useState<boolean>(false);

    const fetchCurrentUsers = async () => {
        await axios.get('http://localhost:5000/api/user/get/except',{
          withCredentials: true,
        })
        .then((res) => {
            if(res.data.result){
                const friendList = res.data.result.map((user: any) => {
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

    const handleCheckboxChange = (value: number) => {
        setSelectedUsers(prev => {
            if (prev.includes(value)) {
                return prev.filter(item => item !== value);
            } else {
                return [...prev, value];
            }
        });
    };

    const createChat = async (user_id: number) => {
        setMessages([])
        setIsChatCreatonOpen(false);
        setIsCreatingNewChat(true);
        setCurrentUser(user_id)
    }

    const createMultipleChat = async() => {
        setMessages([])
        setIsChatCreatonOpen(false);
        setIsCreatingNewChat(true);

        await axios.post('http://localhost:5050/api/chat/chats',{
            IDs: selectedUsers, 
            GroupName: groupName
        },
        {
          withCredentials: true,
        })
        .then(()=> {
            sendMultipleChatCreationNotify(selectedUsers)
        })
        .catch((e)=> {
            if(e.status === 401){
                navigate('/auth', { replace: true});
            }
        });
        fetchChats();
        setGroupName('');
    }

    useEffect(()=> {
        fetchChats()
    },[])

    useEffect(()=> {
        if(isChatCreationOpen || isMultiChatCreationOpen) fetchCurrentUsers()
    },[isChatCreationOpen, isMultiChatCreationOpen])

    return (
        <div className="list">
            <ParentForm
                isDialog ={true}
                isOpen={isChatCreationOpen}
                setIsOpen={setIsChatCreatonOpen}
            >
                <div className="user-list">
                    <div className="d-flex flex-column align-items-stretch bg-body-tertiary w-100 h-100">
                        {users.map((val)=>{
                            return <UserCard 
                                user={val} 
                                createChatHandler = {createChat} 
                                handleCheckboxChange={handleCheckboxChange} 
                                checkboxAvaible = {false}
                            ></UserCard>
                        })}
                    </div>
                </div>
            </ParentForm>
            <ParentForm
                isDialog ={true}
                isOpen={isMultiChatCreationOpen}
                setIsOpen={setIsMultiChatCreationOpen}
            >
                <div className="user-list">
                    <button
                        className="btn btn-outline-secondary m-2" 
                        onClick={() => {
                            setIsMultiChatCreationOpen(false);
                            createMultipleChat()
                        }}
                    >
                        Создать чат с выбранными пользователями
                    </button>
                    
                    <div className="form-floating"> 
                        <input type="text" className="form-control mb-2" id="floatingGroupName"  
                            onChange={(event)=> {
                                setGroupName(event.target.value);
                            }}/> 
                        <label htmlFor="floatingGroupName">Имя группы</label> 
                    </div>
                   
                    <div className="d-flex flex-column align-items-stretch bg-body-tertiary w-100 h-100">
                        {users.map((val)=>{
                            return <UserCard 
                                user={val} 
                                createChatHandler = {createChat} 
                                handleCheckboxChange={handleCheckboxChange} 
                                checkboxAvaible = {true}
                            ></UserCard>
                        })}
                    </div>
                </div>
            </ParentForm>
            <div className="d-grid vh-100  w-100" style={{gridTemplateRows: 'auto 1fr'}}>
                <div className="buttons-group d-grid gap-2" style={{gridTemplateColumns: '1fr 1fr'}}>
                    <button 
                        className="btn btn-outline-secondary m-2" 
                        type="button" 
                        onClick={() => setIsChatCreatonOpen(true)}
                    >
                        Написать пользователю
                    </button>
                    <button 
                        className="btn btn-outline-secondary m-2" 
                        type="button" 
                        onClick={() => setIsMultiChatCreationOpen(true)}
                    >
                        Создать чат
                    </button>
                </div>
                <div className="chat-list">
                    <div className="d-flex flex-column align-items-stretch bg-body-tertiary w-100 h-100">
                        <a href="/" className="d-flex align-items-center flex-shrink-0 p-3 link-body-emphasis text-decoration-none border-bottom">
                            <svg className="bi pe-none me-2" width="30" height="24" aria-hidden="true">
                                <use href="#bootstrap"></use>
                            </svg>
                            <span className="fs-5 fw-semibold">Список чатов</span> 
                        </a>
                        {chats.map((Chat:Chat) => (
                            <ChatCard 
                                key={Chat.id} 
                                chat={Chat} 
                                setCurrentChatId={setCurrentChatId} 
                                currentChatId={currentChatId} 
                                setChatHeader={setChatHeader}
                                usersOnline={usersOnline}
                            />
                        ))}
                    </div>
                </div>
            </div>
        </div>

    )
}

export default ChatsList;