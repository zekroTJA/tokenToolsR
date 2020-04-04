/** @format */

import React, { Component } from 'react';

import { ReactComponent as Logo } from '../../img/logo-wide.svg';
import './About.scss';

export default class AboutRoute extends Component {
  public render() {
    return (
      <div className="about-wrapper">
        <div className="about-header">
          <Logo className="mh-auto" height={40 * 1.5} width={283 * 1.5} />
        </div>
        TokenTools R&nbsp; (the <span className="embed">R</span> is, on the one
        side, for&nbsp;
        <span className="embed">Rewrite</span> and on the other for&nbsp;
        <span className="embed">React</span>) &nbsp; is a tool Discord Bot
        developers can use to validate the state and get informations about
        their Bot API tokens.
        <br />
        <br />
        <strong>
          This tool is not ment to be used for validating tokens found in public
          domains to abuse them! Please report found and valid tokens to the
          authors of the project!
        </strong>
        <h2>Credits</h2>
        Thanks to Subby (<span className="embed">Subby#8883</span> on Discord)
        for creating the Logo for this application.
        <br />
        <br />
        All other used icons are created by Ringo Hoffmann (zekro).
        <br />
        <br />
        Used fonts are&nbsp;
        <a href="https://www.dafont.com/larabie-font.font" target="_blank">
          Larabie Font
        </a>
        &nbsp;(Logo) and&nbsp;
        <a href="https://fonts.google.com/specimen/Barlow" target="_blank">
          Barlow
        </a>
        &nbsp;(UI).
        <h2>Techniques</h2>
        The back end of TokenTools R is created with&nbsp;
        <a href="https://golang.org" target="_blank">
          Go
        </a>
        ,&nbsp;
        <a href="https://github.com/gorilla/mux" target="_blank">
          gorilla/mux
        </a>
        &nbsp;as router and&nbsp;
        <a href="https://github.com/gorilla/websocket" target="_blank">
          gorilla/websocket
        </a>
        &nbsp;as web socket server wrapper.
        <br />
        <br />
        <a href="https://reactjs.org" target="_blank">
          React
        </a>
        &nbsp;is used to create the front end in combination with&nbsp;
        <a href="https://reacttraining.com/react-router/" target="_blank">
          React Router
        </a>
        &nbsp;for SPA routing.
        <h2>Copyright &amp; Licence</h2>
        Created and maintained by Ringo Hoffmann (zekro Development). Covered by
        the&nbsp;
        <a
          href="https://github.com/zekroTJA/tokenToolsR/blob/master/LICENCE"
          target="_blank"
        >
          MIT Licence
        </a>
        .
        <br />
        <br />© 2020 Ringo Hoffmann (zekro Development).
      </div>
    );
  }
}
