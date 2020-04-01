/** @format */

import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import WebSocketAPI from './api/ws';
import './index.scss';

const wsuri =
  process.env.NODE_ENV === 'development'
    ? 'ws://localhost:8081/ws'
    : window.location.href.replace(
        /((http)|(https)):\/\//gm,
        window.location.href.startsWith('http://') ? 'ws:/' : 'wss://'
      ) + 'ws';

const wsapi = new WebSocketAPI(wsuri);

ReactDOM.render(
  <React.StrictMode>
    <App wsapi={wsapi} />
  </React.StrictMode>,
  document.getElementById('root')
);
