/** @format */

import React, { Component, FormEvent } from 'react';
import WebSocketAPI, { EventHandlerRemover } from '../../api/ws';
import { Redirect, Link } from 'react-router-dom';
import { WSTokenValid } from '../../api/model';
import Info from '../../components/info/Info';

import './Main.scss';

enum VALIDITY {
  UNSET,
  VALID,
  INVALID,
}

export default class MainRoute extends Component<{
  wsapi: WebSocketAPI;
  token: string;
}> {
  public state = {
    tokenInput: '',
    valid: VALIDITY.UNSET,
    loading: false,
    showInvalid: false,
    tokenData: (null as any) as WSTokenValid,
    redirect: false,
  };

  private unmounts: EventHandlerRemover[] = [];

  public componentDidMount() {
    this.unmounts.push(
      this.props.wsapi.on('open', () => {
        const token = this.props.token;
        if (token.length > 10) {
          this.setState({ tokenInput: token });
          this.onCheckSubmit();
        }
      })
    );

    this.unmounts.push(
      this.props.wsapi.on('tokenInvalid', () => {
        this.setValidity(VALIDITY.INVALID);
      })
    );

    this.unmounts.push(
      this.props.wsapi.on('tokenValid', (data: WSTokenValid) => {
        this.setValidity(VALIDITY.VALID);
        this.setState({ tokenData: data });
      })
    );

    this.unmounts.push(
      this.props.wsapi.onerror(() => {
        this.setValidity(VALIDITY.UNSET);
      })
    );
  }

  public componentWillUnmount() {
    this.unmounts.forEach((u) => u());
  }

  public render() {
    return (
      <div>
        <div className="inpt-body-container">
          <div className="token-input-wrapper">
            <input
              type="password"
              value={this.state.tokenInput}
              onChange={this.onTokenChange.bind(this)}
              className={
                'token-input' + (this.state.showInvalid ? ' invalid' : '')
              }
            ></input>
            <div
              className={
                'validity-ident' +
                (this.state.valid !== VALIDITY.UNSET ? ' display' : '')
              }
              style={{
                backgroundColor: this.stateColor,
              }}
            ></div>
            {this.state.tokenData && (
              <div className="info-tile">
                <Link to={'/guilds/' + this.state.tokenInput}>
                  <Info data={this.state.tokenData} />
                </Link>
              </div>
            )}
            <div className="flex">
              <button
                className={'check-btn' + (this.state.loading ? ' loading' : '')}
                onClick={this.onCheckSubmit.bind(this)}
              >
                CHECK
              </button>
            </div>
          </div>
        </div>

        {this.state.redirect && (
          <Redirect to={'/check/' + this.state.tokenInput} />
        )}
      </div>
    );
  }

  private onTokenChange(event: FormEvent<HTMLInputElement>) {
    const inpt = (event.target as HTMLInputElement).value.trim();

    if (inpt.length === 0) {
      this.setState({
        tokenInput: inpt,
        valid: VALIDITY.UNSET,
        tokenData: null,
      });
    } else {
      this.setState({
        tokenInput: inpt,
      });
    }
  }

  private onCheckSubmit() {
    const inpt = this.state.tokenInput.trim();

    if (inpt.length === 0) {
      this.setInputInvalid();
      return;
    }

    this.setValidity(VALIDITY.UNSET, true);

    if (inpt.length <= 40) {
      setTimeout(() => this.setValidity(VALIDITY.INVALID), 500);
    } else {
      this.props.wsapi.send('checkToken', inpt);
      this.setState({ redirect: true });
    }
  }

  private setValidity(valid: VALIDITY, loading: boolean = false) {
    this.setState({ valid, loading });
  }

  private setInputInvalid() {
    this.setState({ showInvalid: true });
    setTimeout(() => this.setState({ showInvalid: false }), 1000);
  }

  private get stateColor(): string {
    switch (this.state.valid) {
      case VALIDITY.VALID:
        return '#CDDC39';
      case VALIDITY.INVALID:
        return '#f44336';
    }

    return '#00000000';
  }
}
