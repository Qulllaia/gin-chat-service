import { useEffect, useState } from "react";
import ChatCard from "./ChatCard";
import axios from "axios";
import { Chat, User } from "../types";
import { ParentForm } from "./ParentForm";
import UserCard from "./UserCard";
import { useNavigate } from "react-router-dom";
import '../styles/ChatList.css';

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
}: any) {

  const [users, setUsers] = useState<User[]>([])
  const [selectedUsers, setSelectedUsers] = useState<number[]>([]);
  const navigate = useNavigate()

  const [groupName, setGroupName] = useState<string>('');
  const [isChatCreationOpen, setIsChatCreatonOpen] = useState<boolean>(false);
  const [isMultiChatCreationOpen, setIsMultiChatCreationOpen] = useState<boolean>(false);

  const logoutFunction = async () => {
    await axios.get('http://localhost:5000/api/auth/logout')
      .then(() => navigate('/auth', { replace: true }))
      .catch(() => console.log('error while logout'));
  }

  const fetchCurrentUsers = async () => {
    await axios.get('http://localhost:5000/api/user/get/except', {
      withCredentials: true,
    })
      .then((res) => {
        if (res.data.result) {
          const friendList = res.data.result.map((user: any) => {
            return {
              id: user.id,
              name: user.name,
            } as User;
          })
          setUsers(friendList);
        }
      })
      .catch((e) => {
        if (e.status === 401) {
          navigate('/auth', { replace: true });
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
    const existingPrivateChat = chats.find(
      (c: Chat) => c.chatType === 'PRIVATECHAT' && c.userId === user_id
    );

    setIsChatCreatonOpen(false);

    if (existingPrivateChat) {
      setMessages([]);
      setIsCreatingNewChat(false);
      setCurrentUser(0);
      setCurrentChatId(existingPrivateChat.id);
      setChatHeader(existingPrivateChat.name);
      return;
    }

    setMessages([]);
    setCurrentChatId(0);
    setIsCreatingNewChat(true);
    setCurrentUser(user_id);
  }

  const deleteChat = async (chatId: number, e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!window.confirm('Удалить этот чат и все сообщения?')) {
      return;
    }

    try {
      const res = await axios.delete(`http://localhost:5050/api/chat/chats/${chatId}`, {
        withCredentials: true,
      });

      if (res.status !== 200 || !res.data?.done) {
        return;
      }

      if (currentChatId === chatId) {
        setCurrentChatId(0);
        setMessages([]);
        setChatHeader('Минималистичный Чат');
      }

      setChats((prev: Chat[]) => prev.filter((c) => c.id !== chatId));
      await fetchChats();
    } catch (err) {
      console.error('Failed to delete chat', err);
    }
  };

  const openGroupChatModal = () => {
    setSelectedUsers([]);
    setGroupName('');
    setIsMultiChatCreationOpen(true);
  };

  const createMultipleChat = async () => {
    if (selectedUsers.length === 0) {
      window.alert('Выберите хотя бы одного пользователя');
      return;
    }

    if (!groupName.trim()) {
      window.alert('Введите имя группы');
      return;
    }

    const idsToSend = [...selectedUsers];

    setMessages([]);
    setCurrentChatId(0);
    setIsMultiChatCreationOpen(false);
    setIsCreatingNewChat(true);

    try {
      await axios.post(
        'http://localhost:5050/api/chat/chats',
        { ids: idsToSend, GroupName: groupName.trim() },
        { withCredentials: true },
      );
      sendMultipleChatCreationNotify(idsToSend);
      setSelectedUsers([]);
      setGroupName('');
      await fetchChats();
    } catch (e: any) {
      if (e?.response?.status === 401) {
        navigate('/auth', { replace: true });
      }
    }
  };

  useEffect(() => {
    fetchChats()
  }, [])

  useEffect(() => {
    if (isChatCreationOpen || isMultiChatCreationOpen) fetchCurrentUsers()
  }, [isChatCreationOpen, isMultiChatCreationOpen])

  return (
    <div className="list">
      <ParentForm
        isDialog={true}
        isOpen={isChatCreationOpen}
        setIsOpen={setIsChatCreatonOpen}
        backdropClassName="chat-modal-backdrop"
        contentClassName="chat-modal-content"
      >
        <div className="user-list">
          <h2 className="h3 mb-3 fw-normal">Выберите пользователя</h2>
          <div className="d-flex flex-column align-items-stretch w-100 h-100">
            {users.map((val) => {
              return <UserCard
                key={val.id}
                user={val}
                createChatHandler={createChat}
                handleCheckboxChange={handleCheckboxChange}
                checkboxAvaible={false}
              ></UserCard>
            })}
          </div>
        </div>
      </ParentForm>
      <ParentForm
        isDialog={true}
        isOpen={isMultiChatCreationOpen}
        setIsOpen={setIsMultiChatCreationOpen}
        backdropClassName="chat-modal-backdrop"
        contentClassName="chat-modal-content"
      >
        <div className="user-list">
          <h2 className="h3 mb-3 fw-normal">Новая группа</h2>
          <div className="chat-modal-field">
            <label className="chat-modal-label" htmlFor="floatingGroupName">Имя группы</label>
            <input
              type="text"
              className="chat-modal-input"
              id="floatingGroupName"
              placeholder="Например: Команда"
              value={groupName}
              onChange={(event) => setGroupName(event.target.value)}
            />
          </div>
          <p className="chat-modal-hint">
            Выбрано: {selectedUsers.length}
          </p>
          <div className="chat-modal-user-picker">
            {users.map((val) => (
              <UserCard
                key={val.id}
                user={val}
                createChatHandler={createChat}
                handleCheckboxChange={handleCheckboxChange}
                checkboxAvaible={true}
                isSelected={selectedUsers.includes(val.id)}
              />
            ))}
          </div>
          <button
            type="button"
            className="chat-btn w-100 mt-3"
            onClick={() => void createMultipleChat()}
          >
            Создать чат с выбранными
          </button>
        </div>
      </ParentForm>
      <div className="d-grid vh-100  w-100" style={{ gridTemplateRows: 'auto 1fr' }}>
        <div className="chat-sidebar-toolbar buttons-group">
          <button
            type="button"
            className="chat-btn chat-btn-logout"
            onClick={() => logoutFunction()}
            title="Выйти из аккаунта"
            aria-label="Выйти из аккаунта"
          >
            <svg className="chat-btn-logout-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H9m4 4v1a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h6a2 2 0 012 2v1" />
            </svg>
            <span className="chat-btn-logout-label">Выход</span>
          </button>
          <button
            className="chat-btn"
            type="button"
            onClick={() => setIsChatCreatonOpen(true)}
          >
            Написать
          </button>
          <button
            className="chat-btn"
            type="button"
            onClick={openGroupChatModal}
          >
            Группа
          </button>
        </div>
        <div className="chat-list">
          <div className="d-flex flex-column align-items-stretch w-100 h-100">
            <div className="chat-sidebar-title">Список чатов</div>
            {chats.map((Chat: Chat) => (
              <ChatCard
                key={Chat.id}
                chat={Chat}
                setCurrentChatId={setCurrentChatId}
                currentChatId={currentChatId}
                setChatHeader={setChatHeader}
                usersOnline={usersOnline}
                onDeleteChat={deleteChat}
              />
            ))}
          </div>
        </div>
      </div>
    </div>

  )
}

export default ChatsList;
