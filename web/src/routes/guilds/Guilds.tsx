/** @format */

import React, { Component } from 'react';
import WebSocketAPI, { EventHandlerRemover } from '../../api/ws';
import { WSGuilds, WSGuild } from '../../api/model';
import GuildTile from '../../components/guild-tile/GuildTile';

import './Guilds.scss';
import GuildInfoModal from '../../components/guild-info-modal/GuildInfoModal';

export default class GuildsRoute extends Component<{
  wsapi: WebSocketAPI;
  token: string;
}> {
  public state = {
    guilds: [] as WSGuilds,
    showGuild: (null as any) as WSGuild,
  };

  private unmounts: EventHandlerRemover[] = [];

  public componentDidMount() {
    this.unmounts.push(
      this.props.wsapi.onopen(() => {
        this.props.wsapi.send('init', this.props.token);
        setTimeout(() => this.props.wsapi.send('guildinfo'), 100);
      })
    );

    this.unmounts.push(
      this.props.wsapi.on('guildinfo', (data: any) => {
        this.state.guilds.push(data.guild);
        this.setState({});
      })
    );
  }

  public componentWillUnmount() {
    this.unmounts.forEach((u) => u());
  }

  public render() {
    const guilds = this.state.guilds.map((g) => (
      <GuildTile
        guild={g}
        onClick={() => this.setState({ showGuild: g })}
        key={'g_' + g.id}
      />
    ));
    return (
      <div>
        {this.state.showGuild && (
          <GuildInfoModal
            onClose={() => this.setState({ showGuild: null })}
            wsapi={this.props.wsapi}
            guild={this.state.showGuild}
          />
        )}
        <div className="guilds-container">{guilds}</div>
      </div>
    );
  }
}
