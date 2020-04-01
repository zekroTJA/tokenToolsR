/** @format */

import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import WebSocketAPI from './api/ws';
import MainRoute from './routes/main/Main';
import './App.scss';

export default class App extends Component<{ wsapi: WebSocketAPI }> {
  public render() {
    return (
      <Router>
        <Route
          path="/"
          exact={true}
          render={() => <MainRoute wsapi={this.props.wsapi} />}
        />
      </Router>
    );
  }
}
