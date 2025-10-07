export interface Message {
  id: string;
  text: string;
  sender: 'user' | 'other';
  chat_id: number;
  timestamp: Date;
}

export interface Chat {
    id: number, 
    name: string,
}

export interface User {
    id: number, 
    name: string,
}

export const NEW_CHAT = 'NEW_CHAT'
export const NEW_MULTIPLE_CHAT = 'NEW_MULTIPLE_CHAT'
export const MESSAGE = 'MESSAGE'

