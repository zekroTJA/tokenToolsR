import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';
import WebSocketAPI from './api/ws';

const wsuri = process.env.NODE_ENV === 'development' ? 
  'ws://localhost:8081/ws' :
  window.location.href.replace(/((http)|(https)):\/\//gm, 
    window.location.href.startsWith('http://') ? 'ws:/' : 'wss://') + 'ws';

const wsapi = new WebSocketAPI(wsuri);

ReactDOM.render(
  <React.StrictMode>
    <App wsapi={ wsapi } />
  </React.StrictMode>,
  document.getElementById('root')
);
