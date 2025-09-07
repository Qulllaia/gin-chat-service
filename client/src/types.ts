export interface Message {
  id: string;
  text: string;
  sender: 'user' | 'other';
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
export const MESSAGE = 'MESSAGE'
