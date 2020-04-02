/** @format */

import React, { Component } from 'react';
import { WSTokenValid } from '../../api/model';

import './Info.scss';

export default class Info extends Component<{ data: WSTokenValid }> {
  public render() {
    return (
      <div className="tile">
        <img src={this.props.data.avatar} alt="account avatar" />
        <h3>{this.nameTag}</h3>
        <p className="id">{this.props.data.id}</p>
        <span className="embed">{this.props.data.guilds} GUILDS</span>
      </div>
    );
  }

  private get nameTag(): string {
    return `${this.props.data.username}#${this.props.data.discriminator}`;
  }
}
