/** @format */

import React, { Component } from 'react';
import RestAPI from '../../api/rest';
import { ReactComponent as Logo } from '../../img/header-logo.svg';
import { Version } from '../../api/model';
import { Link } from 'react-router-dom';

import './Header.scss';

export default class Header extends Component {
  public state = {
    version: {} as Version,
  };

  public async componentDidMount() {
    this.setState({
      version: (await RestAPI.getVersion())?.data,
    });
  }

  public render() {
    return (
      <div className="header">
        <Link to="/">
          <Logo />
        </Link>
        {this.state.version && (
          <span className="version-info">{this.version}</span>
        )}
      </div>
    );
  }

  private get version(): string {
    return `build ${this.state.version.version} (${this.state.version.commit}) [${this.state.version.date}]`;
  }
}
