/** @format */

import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Redirect } from 'react-router-dom';
import WebSocketAPI from './api/ws';
import MainRoute from './routes/main/Main';
import GuildsRoute from './routes/guilds/Guilds';
import Header from './components/header/Header';

import './App.scss';
import AboutRoute from './routes/about/About';

export default class App extends Component<{ wsapi: WebSocketAPI }> {
  public render() {
    return (
      <Router>
        <Header />

        <Route
          exact
          path="/check/:token"
          render={({ match }) => (
            <MainRoute wsapi={this.props.wsapi} token={match.params.token} />
          )}
        />
        <Route
          exact
          path="/guilds/:token"
          render={({ match }) => (
            <GuildsRoute wsapi={this.props.wsapi} token={match.params.token} />
          )}
        />
        <Route exact path="/about" render={() => <AboutRoute />} />

        <Route exact path="/" render={() => <Redirect to="/check/_" />} />
        <Route exact path="/check" render={() => <Redirect to="/check/_" />} />
      </Router>
    );
  }
}
