import React from 'react';
import './App.css';
import Chat from './components/Chat';
import './styles.css';
import { FriendsList } from './components/FriendsList';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import Auth from './components/Auth';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Navigate to="/auth" replace />} />
          <Route path="/chat" element={<Chat/>} />
          <Route path="/auth" element={<Auth/>} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
