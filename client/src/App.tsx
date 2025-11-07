import React from 'react';
import './App.css';
import { ChatPage } from './components/Chat';
import './styles.css';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import Auth from './components/Auth';
import 'bootstrap/dist/css/bootstrap.min.css';
import { VerifyPage } from './components/VerifyPage';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Navigate to="/auth" replace />} />
          <Route path="/chat" element={<ChatPage/>} />
          <Route path="/auth" element={<Auth/>} />
          <Route path="/verify/:token" element={<VerifyPage/>} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
