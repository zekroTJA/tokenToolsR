/** @format */

import React, { Component } from 'react';
import WebSocketAPI, { EventHandlerRemover } from '../../api/ws';
import { WSGuild, WSUser } from '../../api/model';
import { ReactComponent as CloseIcon } from '../../img/close.svg';

import './GuildInfoModal.scss';
import DefAvatar from '../def-avatar/DefAvatar';

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
    const guild = this.props.guild;
    const owner = this.state.owner;

    return (
      <div id="modal-wrapper" onClick={this.onBgClick.bind(this)}>
        <div id="modal">
          <div className="header">
            <h2>{guild.name}</h2>
            <a
              href="#"
              className="close-btn"
              onClick={this.props.onClose.bind(this)}
            >
              <CloseIcon width="25" />
            </a>
          </div>
          <div className="info-head">
            GUILD-ID: <span className="embed">{guild.id}</span>
          </div>
          {owner && (
            <div className="owner-tile">
              <DefAvatar
                src={owner.avatar}
                width={100}
                height={100}
                alt={owner.username}
              />
              <h2 className="heading">
                {owner.username}#{owner.discriminator}
              </h2>
              <p className="id">({owner.id})</p>
            </div>
          )}
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
