/** @format */

import React, { Component, FormEvent } from 'react';
import WebSocketAPI, { EventHandlerRemover } from '../../api/ws';

import './Main.scss';

enum VALIDITY {
  UNSET,
  VALID,
  INVALID,
}

export default class MainRoute extends Component<{ wsapi: WebSocketAPI }> {
  public state = {
    tokenInput: '',
    valid: VALIDITY.UNSET,
    loading: false,
  };

  private unmounts: EventHandlerRemover[] = [];

  constructor(props: Readonly<any>) {
    super(props);

    this.unmounts.push(
      this.props.wsapi.on('tokenInvalid', () => {
        this.setValidity(VALIDITY.INVALID);
      })
    );

    this.unmounts.push(
      this.props.wsapi.on('tokenValid', () => {
        this.setValidity(VALIDITY.VALID);
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
              className="token-input"
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
      </div>
    );
  }

  private onTokenChange(event: FormEvent<HTMLInputElement>) {
    const inpt = (event.target as HTMLInputElement).value.trim();
    this.setState({
      tokenInput: inpt,
      valid: inpt.length === 0 ? VALIDITY.UNSET : undefined,
    });
  }

  private onCheckSubmit() {
    const inpt = this.state.tokenInput.trim();

    if (inpt.length === 0) {
      return;
    }

    this.setValidity(VALIDITY.UNSET, true);

    if (inpt.length <= 40) {
      setTimeout(() => this.setValidity(VALIDITY.INVALID), 500);
    } else {
      this.props.wsapi.send('checkToken', inpt);
    }
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

  private setValidity(valid: VALIDITY, loading: boolean = false) {
    this.setState({ valid, loading });
  }
}
