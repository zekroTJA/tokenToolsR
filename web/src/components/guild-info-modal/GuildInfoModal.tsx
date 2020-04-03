/** @format */

import React, { Component } from 'react';
import WebSocketAPI, { EventHandlerRemover } from '../../api/ws';
import { WSGuild, WSUser } from '../../api/model';
import { ReactComponent as CloseIcon } from '../../img/close.svg';

import './GuildInfoModal.scss';

export default class GuildInfoModal extends Component<{
  wsapi: WebSocketAPI;
  guild: WSGuild;
  onClose: () => void;
}> {
  public state = {
    owner: (null as any) as WSUser,
  };

  public static defaultProps = {
    onClose: () => {},
  };

  private unmounts: EventHandlerRemover[] = [];

  public componentDidMount() {
    this.unmounts.push(
      this.props.wsapi.on('userInfo', (data: WSUser) => {
        this.setState({ owner: data });
      })
    );

    this.unmounts.push(
      this.props.wsapi.onopen(() => {
        this.props.wsapi.send('getUserInfo', this.props.guild.owner);
      })
    );
  }

  public componentWillUnmount() {
    this.unmounts.forEach((u) => u());
  }

  public render() {
    return (
      <div id="modal-wrapper" onClick={this.onBgClick.bind(this)}>
        <div id="modal">
          <div className="header">
            <h2>{this.props.guild.name}</h2>
            <a
              href="#"
              className="close-btn"
              onClick={this.props.onClose.bind(this)}
            >
              <CloseIcon width="25" />
            </a>
          </div>
        </div>
      </div>
    );
  }

  private onBgClick(e: React.MouseEvent<HTMLDivElement, MouseEvent>) {
    e.preventDefault();
    e.stopPropagation();

    const target = e.target as HTMLDivElement;
    if (target.id === 'modal-wrapper') {
      this.props.onClose();
    }
  }
}
