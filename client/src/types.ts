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