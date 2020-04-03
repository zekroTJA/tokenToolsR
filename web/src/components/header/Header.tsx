/** @format */

import React, { Component } from 'react';
import RestAPI from '../../api/rest';
import { Version } from '../../api/model';
import { Link } from 'react-router-dom';

import { ReactComponent as Logo } from '../../img/logo-wide.svg';
import './Header.scss';

export default class Header extends Component {
  // public state = {
  //   version: {} as Version,
  // };

  // public async componentDidMount() {
  //   this.setState({
  //     version: (await RestAPI.getVersion())?.data,
  //   });
  // }

  public render() {
    return (
      <div className="header">
        <Link to="/" className="header-logo">
          <Logo height="25" width="177" />
        </Link>
        {/* {this.state.version && (
          <span className="version-info">{this.version}</span>
        )} */}
      </div>
    );
  }

  // private get version(): string {
  //   return `build ${this.state.version.version} (${this.state.version.commit}) [${this.state.version.date}]`;
  // }
}
