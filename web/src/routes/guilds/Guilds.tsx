/** @format */

import React, { Component } from 'react';
import WebSocketAPI, { EventHandlerRemover } from '../../api/ws';

import './Guilds.scss';
import { WSGuilds } from '../../api/model';
import GuildTile from '../../components/guild-tile/GuildTile';

export default class GuildsRoute extends Component<{
  wsapi: WebSocketAPI;
  token: string;
}> {
  public state = {
    guilds: [] as WSGuilds,
  };

  private unmounts: EventHandlerRemover[] = [];

  public componentDidMount() {
    this.unmounts.push(
      this.props.wsapi.on('open', () => {
        this.props.wsapi.send('getGuildInfo', this.props.token);
      })
    );

    this.unmounts.push(
      this.props.wsapi.on('guildInfo', (data: WSGuilds) => {
        this.setState({ guilds: data });
      })
    );
  }

  public componentWillUnmount() {
    this.unmounts.forEach((u) => u());
  }

  public render() {
    const guilds = this.state.guilds.map((g) => (
      <GuildTile guild={g} key={'g_' + g.id} />
    ));
    return <div className="guilds-container">{guilds}</div>;
  }
}
