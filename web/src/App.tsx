/** @format */

import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import WebSocketAPI from './api/ws';
import MainRoute from './routes/main/Main';
import Header from './components/header/Header';

import './App.scss';

export default class App extends Component<{ wsapi: WebSocketAPI }> {
  public render() {
    return (
      <Router>
        <Header />
        <Route exact path="/test" render={() => <p>kekw</p>} />
        <Route path="/" render={() => <MainRoute wsapi={this.props.wsapi} />} />
      </Router>
    );
  }
}
