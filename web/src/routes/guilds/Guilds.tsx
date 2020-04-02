/** @format */

import React, { Component } from 'react';
import WebSocketAPI, { EventHandlerRemover } from '../../api/ws';
import { WSGuilds } from '../../api/model';
import GuildTile from '../../components/guild-tile/GuildTile';

import './Guilds.scss';

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
      this.props.wsapi.onopen(() => {
        this.props.wsapi.send('getGuildInfo', this.props.token);
      })
    );

    this.unmounts.push(
      this.props.wsapi.on('guildInfo', (data: any) => {
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
      <GuildTile guild={g} key={'g_' + g.id} />
    ));
    return <div className="guilds-container">{guilds}</div>;
  }
}
