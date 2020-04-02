/** @format */

import React, { Component } from 'react';
import { WSGuild } from '../../api/model';

import './GuildTile.scss';

export default class GuildTile extends Component<{ guild: WSGuild }> {
  public render() {
    const guild = this.props.guild;
    return (
      <div className="guild">
        <img src={this.guildIcon} alt="guild icon" />
        <h3>{guild.name}</h3>
        <p className="id">{guild.id}</p>
      </div>
    );
  }

  private get guildIcon(): string {
    const guild = this.props.guild;
    return `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`;
  }
}
