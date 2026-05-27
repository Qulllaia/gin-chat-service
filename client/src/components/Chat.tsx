import { useEffect, useRef, useState } from 'react';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { Chat, ChatMember, MESSAGE, Message, NEW_CHAT, NEW_MULTIPLE_CHAT } from '../types';
import { ChatsList } from './ChatsList';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { ParentForm } from './ParentForm';
import ChunkedAudioPlayer from './ChunkedAudioPlayer';
import '../styles/Chat.css';

interface PreviewItem {
  id: string;
  url: string;
  name: string;
  file: File;
}

export function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([]);
  const ws = useRef<WebSocket>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate()
  const [currentChatId, setCurrentChatId] = useState<number>(0);
  const currentChatIdRef = useRef(currentChatId);
  const [currentUser, setCurrentUser] = useState<number>(0)
  const [chatHeader, setChatHeader] = useState<string>('Минималистичный Чат')
  const [isCreatingNewChat, setIsCreatingNewChat] = useState<boolean>(false);
  const [chats, setChats] = useState<Chat[]>([])
  const [isBackgroundUpdateOpen, setIsBackgroundUpdateOpen] = useState<boolean>(false);
  const [isMembersModalOpen, setIsMembersModalOpen] = useState<boolean>(false);
  const [chatMembers, setChatMembers] = useState<ChatMember[]>([]);
  const [membersLoading, setMembersLoading] = useState<boolean>(false);
  const [membersError, setMembersError] = useState<string>('');

  const [previews, setPreviews] = useState<PreviewItem[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [usersOnline, setUsersOnline] = useState<number[]>([])

  const resetChatBackgroundToDefault = () => {
    const backgroundDiv = document.getElementById('background-div') as HTMLElement | null;
    if (!backgroundDiv) return;

    backgroundDiv.style.backgroundImage = 'none';
    backgroundDiv.style.backgroundSize = '';
    backgroundDiv.style.backgroundPosition = '';
    backgroundDiv.style.backgroundRepeat = '';
  };

  const applyChatBackground = (backgroundUrl?: string | null) => {
    const backgroundDiv = document.getElementById('background-div') as HTMLElement | null;
    if (!backgroundDiv) return;

    if (backgroundUrl) {
      backgroundDiv.style.backgroundImage = `url(http://localhost:5050${backgroundUrl})`;
      backgroundDiv.style.backgroundSize = 'cover';
      backgroundDiv.style.backgroundPosition = 'center';
      backgroundDiv.style.backgroundRepeat = 'no-repeat';
    } else {
      resetChatBackgroundToDefault();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    if (!e.target.files) return;

    const files = Array.from(e.target.files);
    createPreviews(files);
  };

  const createPreviews = (files: File[]): void => {
    const newPreviews: PreviewItem[] = [];
    const imageFiles = files.filter(file => file.type.startsWith('image/'));

    if (imageFiles.length === 0) return;

    let processedCount = 0;

    imageFiles.forEach(file => {
      const reader = new FileReader();

      reader.onload = (e: ProgressEvent<FileReader>) => {
        if (!e.target?.result) return;

        newPreviews.push({
          id: Math.random().toString(36).substr(2, 9),
          url: e.target.result as string,
          name: file.name,
          file: file
        });

        processedCount++;

        if (processedCount === imageFiles.length) {
          setPreviews(newPreviews);
        }
      };

      reader.readAsDataURL(file);
    });
  };

  useEffect(() => {
    if (fileInputRef.current && fileInputRef.current.files && fileInputRef.current.files.length > 0) {
      createPreviews(Array.from(fileInputRef.current.files));
    }
  }, []);

  const removePreview = (id: string): void => {
    setPreviews(prev => prev.filter(item => item.id !== id));

    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const getFileSize = (bytes: number): string => {
    if (bytes < 1024) return bytes + ' bytes';
    else if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB';
    else return (bytes / 1048576).toFixed(1) + ' MB';
  };

  const fetchMEssagesHistory = async () => {
    await axios.get(`http://localhost:5050/api/chat/history/${currentChatId}`, {
      method: "GET",
      withCredentials: true,
    }).then((res) => {
      if (res.status === 200) {
        res.data.result.forEach((element: any) => {
          const message = {
            id: element.id,
            text: element.message,
            sender: element.IsThisUserSender ? 'user' : 'other',
            timestamp: element.timestamp,
          } as Message
          setMessages((messages) => [...messages, message])
        });
      } else {
        navigate('/auth', { replace: true });
      }
    }).then(() => scrollToBottom())
      .catch((e) => {
        if (e.status === 401) {
          navigate('/auth', { replace: true });
        }
      });
  }

  const fetchChats = async () => {
    await axios.get('http://localhost:5050/api/chat/chats',
      {
        withCredentials: true,
      }
    ).then((res) => {
      if (res.data.result) {
        const friendList = res.data.result.map((chat: any) => {
          return {
            id: chat.ID,
            name: chat.Name,
            backgroundUrl: chat.Chat_background,
            chatType: chat.Chat_type_id,
            userId: chat.User_id,
            lastMessage: chat.LastMessage,
          } as Chat;
        })
        setChats(friendList);
      }
    })
      .catch((e) => {
        if (e.status === 401) {
          navigate('/auth', { replace: true });
        }
      });
  }

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollTo({
        top: messagesEndRef.current.scrollHeight,
        behavior: 'auto'
      });

    }
  };

  useEffect(() => {
    currentChatIdRef.current = currentChatId;
  }, [currentChatId])

  useEffect(() => {
    if (!ws.current) {
      ws.current = new WebSocket('ws://localhost:5050/api/chat/ws?user_id=fslkfjslkfjslfs&chat_id=1');
    }
    // console.log('useeffect')
    ws.current.onopen = () => {
      fetchMEssagesHistory();
    };

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === "MESSAGE") {
        const newMessage: Message = {
          id: Date.now().toString(),
          text: data.message,
          sender: 'other',
          chat_id: data.chat_id,
          timestamp: new Date(),
        };
        if (currentChatIdRef.current === newMessage.chat_id)
          setMessages((prev) => [...prev, newMessage]);
      }
      else if (data.type === 'NEW_CHAT') {
        resetChatBackgroundToDefault();
        setCurrentChatId(data.chat_id as number);
        fetchChats();
      }
      else if (data.type === "NEW_MULTIPLE_CHAT") {
        resetChatBackgroundToDefault();
        fetchChats();
      }
      // TODO: Привести входные данные к общему формату
      else if (data.online) {
        setUsersOnline(prev => [data.online, ...prev])
      }
      else if (data.offline) {
        setUsersOnline(prev => prev.filter((userId) => userId !== data.offline))
      } else if (data.activeUsersIds) {
        setUsersOnline(data.activeUsersIds)
      }
    };

    ws.current.onclose = () => {
      console.log('WebSocket disconnected');
    };

    ws.current.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.close();
        navigate('/auth', { replace: true });
      }
    };
  }, []);


  useEffect(() => {
    scrollToBottom();
  }, [messages])

  useEffect(() => {
    if (isCreatingNewChat || (currentChatId === 0 && currentUser !== 0) || currentChatId === 0) {
      resetChatBackgroundToDefault();
      return;
    }

    const activeChat = chats.find((chat) => chat.id === currentChatId);
    applyChatBackground(activeChat?.backgroundUrl || null);
  }, [currentChatId, chats, isCreatingNewChat, currentUser]);

  useEffect(() => {
    if (currentChatId === 0 || isCreatingNewChat) {
      return;
    }

    setMessages([]);
    fetchMEssagesHistory();
    currentChatIdRef.current = currentChatId;
  }, [currentChatId])

  const sendMessage = (text: string) => {
    if (ws.current && text.trim()) {

      if (isCreatingNewChat) {
        ws.current.send(JSON.stringify({
          type: NEW_CHAT,
          user_id: currentUser.toString(),
          messages: text
        }));
        setIsCreatingNewChat(false);
      } else {
        ws.current.send(JSON.stringify({
          type: MESSAGE,
          chat_id: currentChatId.toString(),
          messages: text
        }));
      }

      const newMessage: Message = {
        id: Date.now().toString(),
        text,
        sender: 'user',
        chat_id: currentChatId,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, newMessage]);
    }
  };

  const sendMultipleChatCreationNotify = (userIds: number[]) => {
    if (ws.current) {
      ws.current.send(JSON.stringify({
        type: NEW_MULTIPLE_CHAT,
        user_ids: userIds
      }));
      setIsCreatingNewChat(false);
    }
  }

  const getCurrentChat = (): Chat | undefined =>
    chats.find((chat) => chat.id === currentChatId);

  const isGroupChat = (): boolean => getCurrentChat()?.chatType === 'GROUPCHAT';

  const openMembersModal = async () => {
    if (currentChatId === 0 || isCreatingNewChat || !isGroupChat()) {
      return;
    }

    setIsMembersModalOpen(true);
    setMembersLoading(true);
    setMembersError('');
    setChatMembers([]);

    try {
      const res = await axios.get(`http://localhost:5050/api/chat/chats/${currentChatId}/members`, {
        withCredentials: true,
      });

      if (res.data?.done && res.data.result) {
        setChatMembers(res.data.result);
      } else {
        setMembersError('Не удалось загрузить участников');
      }
    } catch (e: any) {
      if (e?.response?.status === 401) {
        navigate('/auth', { replace: true });
        return;
      }
      setMembersError('Не удалось загрузить участников');
    } finally {
      setMembersLoading(false);
    }
  };

  const isOnlineStatusVisible = (): boolean => {
    const data = chats.find((value: Chat, index: number, chats: Chat[]) => {
      if (value.id === currentChatId) {
        return value
      }
    })
    if (data) {
      if (data.chatType === "PRIVATECHAT") {
        return true
      }
    }
    return false
  }
  return (
    <div className="chat-page">
    <div className='chat-body'>
      <div className="chat-audio-hidden" aria-hidden="true">
        <ChunkedAudioPlayer ws={ws} />
      </div>
      <ChatsList
        setCurrentChatId={setCurrentChatId}
        currentChatId={currentChatId}
        setMessages={setMessages}
        setIsCreatingNewChat={setIsCreatingNewChat}
        setCurrentUser={setCurrentUser}
        setChatHeader={setChatHeader}
        fetchChats={fetchChats}
        chats={chats}
        setChats={setChats}
        sendMultipleChatCreationNotify={sendMultipleChatCreationNotify}
        usersOnline={usersOnline}
      />
      {currentChatId === 0 && currentUser === 0 ? (
        <div className="chat-empty-placeholder">
          Выберите чат из списка или создайте новый
        </div>
      ) : (
      <div className="chat-container" id='background-div'>
        <div className='chat-header-panel'>
          <div className="avatar-container">
            <img className="avatar" alt="" />
            {isOnlineStatusVisible() && (
              <div className="status-indicator status-online" />
            )}
          </div>
          {isGroupChat() ? (
            <button
              type="button"
              className="chat-header-title"
              onClick={() => void openMembersModal()}
              disabled={currentChatId === 0 || isCreatingNewChat}
              title="Участники беседы"
            >
              {chatHeader}
            </button>
          ) : (
            <h5 className="chat-header-title chat-header-title--static">{chatHeader}</h5>
          )}
          <button
            type="button"
            className="chat-btn"
            onClick={() => setIsBackgroundUpdateOpen(true)}
          >
            Изменить фон
          </button>
        </div>
        <ParentForm
          isDialog={true}
          setIsOpen={setIsMembersModalOpen}
          isOpen={isMembersModalOpen}
          backdropClassName="chat-modal-backdrop"
          contentClassName="chat-modal-content"
        >
          <div className="chat-members-modal">
            <h2 className="h3 mb-3 fw-normal">Участники беседы</h2>
            <p className="chat-modal-hint">{chatHeader}</p>
            {membersLoading && <div className="chat-members-loading">Загрузка...</div>}
            {membersError && <div className="chat-members-error">{membersError}</div>}
            {!membersLoading && !membersError && (
              <ul className="chat-members-list">
                {chatMembers.map((member) => (
                  <li key={member.id} className="chat-member-item">
                    <span className="chat-member-avatar" aria-hidden>
                      {member.name.charAt(0).toUpperCase()}
                    </span>
                    <span className="chat-member-name">{member.name}</span>
                    {usersOnline.includes(member.id) && (
                      <span className="chat-member-badge">онлайн</span>
                    )}
                  </li>
                ))}
              </ul>
            )}
          </div>
        </ParentForm>
        <ParentForm
          isDialog={true}
          setIsOpen={setIsBackgroundUpdateOpen}
          isOpen={isBackgroundUpdateOpen}
          backdropClassName="chat-modal-backdrop"
          contentClassName="chat-modal-content"
        >
          <div className='scroll-controller'>
            <h1 className="h3 mb-3 fw-normal">Вставьте картинку</h1>
            <div className="chat-file-upload">
              <input
                type="file"
                id="floatingImage"
                className="chat-file-upload-input"
                ref={fileInputRef}
                multiple
                accept="image/*"
                onChange={handleFileChange}
              />
              <label htmlFor="floatingImage" className="chat-file-upload-label">
                <span className="chat-file-upload-icon" aria-hidden>
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </span>
                <span className="chat-file-upload-title">
                  {previews.length > 0 ? 'Выбрать другие файлы' : 'Нажмите или перетащите изображение'}
                </span>
                <span className="chat-file-upload-hint">PNG, JPG, WEBP</span>
              </label>
            </div>
            <div className="image-upload">
              <div className="preview-grid">
                {previews.map(preview => (
                  <div key={preview.id} className="preview-item">
                    <img
                      src={preview.url}
                      alt={preview.name}
                      className="preview-image"
                    />
                    <div className="preview-info">
                      <div className="file-name">{preview.name}</div>
                      <div className="file-size">{getFileSize(preview.file.size)}</div>
                    </div>
                    <button
                      onClick={() => removePreview(preview.id)}
                      className="chat-btn"
                      type="button"
                    >
                      Удалить
                    </button>
                    <button
                      onClick={() => {
                        const fileInput = document.getElementById('floatingImage') as HTMLInputElement;
                        const file = fileInput.files![0];
                        const formData = new FormData();
                        formData.append('image', file);
                        formData.append('chat_id', String(currentChatId));

                        axios.post('http://localhost:5050/api/chat/background', formData,
                          {
                            withCredentials: true,
                            headers: {
                              'Content-Type': 'multipart/form-data',
                            },
                          })
                          .then((res) => {
                            setIsBackgroundUpdateOpen(false)
                            const fullImageUrl = res.data.result.full_url;

                            const backgroundDiv = document.getElementById('background-div') as HTMLElement;
                            backgroundDiv.style.backgroundImage = `url(http://${fullImageUrl})`;
                            backgroundDiv.style.backgroundSize = 'cover';
                            backgroundDiv.style.backgroundPosition = 'center';
                            backgroundDiv.style.backgroundRepeat = 'no-repeat';
                            let temp_chat_list = []
                            chats.forEach(chat => {
                              if (chat.id === currentChatId) {
                                chat.backgroundUrl = res.data.result.url;
                              }
                              temp_chat_list.push(chat)
                            })

                          })
                          .catch((e) => {
                            console.log(e);
                          });
                      }}
                      className="chat-btn"
                      type="button"
                    >
                      Изменить фон
                    </button>
                  </div>
                ))}
              </div>

              {previews.length === 0 && (
                <div className="empty-state">
                  Изображения не выбраны
                </div>
              )}
            </div>
          </div>
        </ParentForm>
        <MessageList ref={messagesEndRef} messages={messages} />
        <MessageInput onSend={sendMessage} />
      </div>
      )}
    </div>
    </div>
  );
}

export default ChatPage;
