import { useEffect, useRef, useState } from 'react';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { Chat, MESSAGE, Message, NEW_CHAT, NEW_MULTIPLE_CHAT } from '../types';
import { ChatsList } from './ChatsList';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { ParentForm } from './ParentForm';

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
  const [backgroundUrl, setBackgroundUrl] = useState<string>('');

  const [previews, setPreviews] = useState<PreviewItem[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);

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
    }).then((res)=> {
      if(res.status === 200) {
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
        navigate('/auth', { replace: true});
      }
    }).then(()=>scrollToBottom())
    .catch((e)=> {
      if(e.status === 401){
        navigate('/auth', { replace: true});
      }
    });
  }

  const fetchChats = async () => {
      await axios.get('http://localhost:5050/api/chat/chats',
        {
          withCredentials: true,
        }
      ).then((res) => {
          if(res.data.result){
              const friendList = res.data.result.map((chat: any) => {
                  return {
                      id: chat.ID,
                      name: chat.Name,
                      backgroundUrl: chat.Chat_background
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

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
        messagesEndRef.current.scrollTo({
          top: messagesEndRef.current.scrollHeight,
          behavior: 'auto'
        });

    }
  };

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
        if(data.type === "MESSAGE") {
          const newMessage: Message = {
            id: Date.now().toString(),
            text: data.message,
            sender: 'other',
            chat_id: data.chat_id,
            timestamp: new Date(),
          };
          if(currentChatIdRef.current === newMessage.chat_id)
            setMessages((prev) => [...prev, newMessage]);
        }
        else if (data.type === 'NEW_CHAT') {
          setCurrentChatId(data.chat_id as number);
          fetchChats()
        }
        else if (data.type === "NEW_MULTIPLE_CHAT") {
          fetchChats()
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
            navigate('/auth', { replace: true});
          }
      };
  }, []);


  useEffect(()=> {
    scrollToBottom();
  }, [messages])

  useEffect(()=>{
    setMessages([]);
    fetchMEssagesHistory();
    let backgroundUrl;
    chats.forEach(chat => {
      if(chat.id === currentChatId) {
        backgroundUrl = chat.backgroundUrl;
      }
    }) 
    
    const backgroundDiv = document.getElementById('background-div') as HTMLElement;
    if(backgroundUrl) {  
      backgroundDiv.style.backgroundImage = `url(http://localhost:5050${backgroundUrl})`;
      backgroundDiv.style.backgroundSize = 'cover';
      backgroundDiv.style.backgroundPosition = 'center';
      backgroundDiv.style.backgroundRepeat = 'no-repeat';
      currentChatIdRef.current = currentChatId;
    } else { 
      backgroundDiv.style.backgroundImage = ``;
    }
  }, [currentChatId])

  const sendMessage = (text: string) => {
    if (ws.current && text.trim()) {

      if(isCreatingNewChat) {
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
  return (
    <div className='chat-body'>
      <ChatsList 
        setCurrentChatId ={setCurrentChatId} 
        currentChatId={currentChatId}
        setMessages={setMessages} 
        setIsCreatingNewChat={setIsCreatingNewChat} 
        setCurrentUser={setCurrentUser}
        setChatHeader={setChatHeader}
        fetchChats={fetchChats}
        chats={chats}
        setChats={setChats}
        sendMultipleChatCreationNotify = {sendMultipleChatCreationNotify}
      />
      <div className={currentChatId === 0 && currentUser === 0 ? "chat-container-hide" : "chat-container"}  id = 'background-div'>
        <div className='chat-header-panel'>
          <h5>{chatHeader}</h5>
          <button className="buttons-group d-grid gap-2" onClick={()=>{
            setIsBackgroundUpdateOpen(true);
          }}>Изменить фон чата</button>
        </div>
        <ParentForm isDialog={true} setIsOpen={setIsBackgroundUpdateOpen} isOpen ={isBackgroundUpdateOpen}>
          <div className='scroll-controller'>
              <h1 className="h3 mb-3 fw-normal">Вставьте картинку</h1> 
              <div className="form-floating"> 
                  <input type="file" className="form-control mb-2" id="floatingImage" placeholder="Password"  
                    ref={fileInputRef}
                    multiple
                    accept="image/*"
                    onChange={handleFileChange}
                  /> 
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
                      className="buttons-group d-grid gap-2"
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
                        // console.log(currentChatId) 
                        // console.log(currentChatIdRef) 
                        formData.append('chat_id', String(currentChatIdRef.current)); 

                        axios.post('http://localhost:5050/api/chat/background', formData,
                        {
                          withCredentials: true,
                          headers: {
                              'Content-Type': 'multipart/form-data',
                          },
                        })
                        .then((res)=> {
                          setIsBackgroundUpdateOpen(false)
                          const fullImageUrl = res.data.full_url;
                          
                          const backgroundDiv = document.getElementById('background-div') as HTMLElement;
                          backgroundDiv.style.backgroundImage = `url(http://${fullImageUrl})`;
                          backgroundDiv.style.backgroundSize = 'cover';
                          backgroundDiv.style.backgroundPosition = 'center';
                          backgroundDiv.style.backgroundRepeat = 'no-repeat';
                          let temp_chat_list = [] 
                          chats.forEach(chat => {
                            if (chat.id === currentChatId) {
                              chat.backgroundUrl = res.data.url;
                            }
                            temp_chat_list.push(chat)
                          })                           

                          setBackgroundUrl(fullImageUrl); 
                        })
                        .catch((e)=> { 
                          console.log(e); 
                        });
                      }}
                      className="buttons-group d-grid gap-2"
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
    </div>
  );
}

export default ChatPage;