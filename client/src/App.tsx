import React from 'react';
import './App.css';
import Chat from './components/Chat';
import './styles.css';
import { FriendsList } from './components/FriendsList';

function App() {
  return (
    <div className="App">
      <FriendsList/>
      <Chat/>
    </div>
  );
}

export default App;
