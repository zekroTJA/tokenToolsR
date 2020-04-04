/** @format */

import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import { ReactComponent as Logo } from '../../img/logo-wide.svg';
import './Header.scss';

export default class Header extends Component {
  public render() {
    return (
      <div className="header">
        <Link to="/" className="header-logo header-link">
          <Logo height="25" width="177" />
        </Link>
        <a
          href="https://github.com/zekroTJA/tokenToolsR"
          target="_blank"
          className="header-logo header-link"
        >
          GITHUB
        </a>
        <Link to="/about" className="header-link">
          ABOUT
        </Link>
      </div>
    );
  }
}
