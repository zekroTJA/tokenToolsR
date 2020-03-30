import { WSMessage } from "./model";

export type EventHandler = (...args: any) => void;
export type EventHandlerRemover = () => void;

export default class WebSocketAPI {
  private ws: WebSocket;
  private handlers: { [key: string]: EventHandler[] } = {};

  constructor(url: string) {
    this.ws = new WebSocket(url);

    this.ws.onmessage = response => {
      try {
        const data = JSON.parse(response.data) as WSMessage;
        if (data) {
          const handler = this.handlers[data.event];
          this.emit(data.event, data.data);
        }
      } catch (err) {
        this.emit("error", err);
      }
    };

    this.ws.onerror = error => {
      this.emit("error", error);
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
    return this.on("error", handler);
  }

  public send(event: string, data: any) {
    const rawData = JSON.stringify({
      event,
      data
    } as WSMessage);

    this.ws.send(rawData);
  }

  private emit(event: string, ...args: any) {
    if (this.handlers[event]) {
      this.handlers[event].forEach(h => {
        if (h) h(args);
      });
    }
  }
}
