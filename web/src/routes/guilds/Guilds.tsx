/** @format */

import React, { Component } from 'react';

import './Guilds.scss';
import WebSocketAPI from '../../api/ws';

export default class GuildsRoute extends Component<{
  wsapi: WebSocketAPI;
  token: string;
}> {
  public render() {
    return <p>kekw</p>;
  }
}
