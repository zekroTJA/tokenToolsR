/** @format */

import React, { Component } from 'react';
import { WSGuild } from '../../api/model';
import DefAvatar from '../def-avatar/DefAvatar';

import './GuildTile.scss';

export default class GuildTile extends Component<{
  guild: WSGuild;
  onClick: () => void;
}> {
  public render() {
    const guild = this.props.guild;
    return (
      <a className="guild" href="#" onClick={this.props.onClick.bind(this)}>
        <DefAvatar
          src={this.guildIcon}
          width={60}
          height={60}
          alt={guild.name}
        />
        <h3>{guild.name}</h3>
        <p className="id">{guild.id}</p>
      </a>
    );
  }

  private get guildIcon(): string {
    const guild = this.props.guild;
    return `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`;
  }
}
