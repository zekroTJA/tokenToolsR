import React, { Props, Component } from 'react';
import './App.css';
import WebSocketAPI, { EventHandlerRemover } from './api/ws';

export default class App extends Component<{ wsapi: WebSocketAPI }> {

  public state = { valid: '' };

  private unmounts: EventHandlerRemover[] = [];

  constructor(props: Readonly<any>) {
    super(props);

    this.unmounts.push(this.props.wsapi.on('tokenInvalid', () => {
      this.setState({ valid: 'invalid' });
    }));

    this.unmounts.push(this.props.wsapi.on('tokenValid', () => {
      this.setState({ valid: 'valid' });
    }));
  }

  public componentWillUnmount() {

  }
  
  public render() {
    return <div>
      <a onClick={ () => this.onSubmit() } href="#">SUBMIT</a>
      <p>{ this.state.valid }</p>
      </div>;
  }

  private onSubmit() {
    this.props.wsapi.send('checkToken', 'sdhjkfhasdfjklhasdflhjk');
  }
}
