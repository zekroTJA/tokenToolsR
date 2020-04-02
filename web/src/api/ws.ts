/** @format */

import { WSMessage } from './model';

export type EventHandler = (args: any) => void;
export type EventHandlerRemover = () => void;

export default class WebSocketAPI {
  private ws: WebSocket;
  private handlers: { [key: string]: EventHandler[] } = {};
  private open = false;

  constructor(url: string) {
    this.ws = new WebSocket(url);

    this.ws.onmessage = (response) => {
      try {
        const data = JSON.parse(response.data) as WSMessage;
        if (data) {
          this.emit(data.event, data.data);
        }
      } catch (err) {
        this.emit('error', err);
      }
    };

    this.ws.onerror = (error) => {
      this.emit('error', error);
    };

    this.ws.onopen = (event) => {
      this.open = true;
      this.emit('open', event);
    };
  }

  public on(event: string, handler: EventHandler): EventHandlerRemover {
    if (!this.handlers[event]) {
      this.handlers[event] = [];
    }

    this.handlers[event].push(handler);

    return () => {
      const i = this.handlers[event].indexOf(handler);
      this.handlers[event].splice(i, 1);
    };
  }

  public onerror(handler: EventHandler): EventHandlerRemover {
    return this.on('error', handler);
  }

  public onopen(handler: EventHandler): EventHandlerRemover {
    if (this.open) {
      handler(null);
      return () => {};
    }

    return this.on('open', handler);
  }

  public send(event: string, data: any) {
    const rawData = JSON.stringify({
      event,
      data,
    } as WSMessage);

    this.ws.send(rawData);
  }

  public get isOpen(): boolean {
    return this.open;
  }

  private emit(event: string, data: any) {
    if (this.handlers[event]) {
      this.handlers[event].forEach((h) => {
        if (h) h(data);
      });
    }
  }
}
