/** @format */

export interface WSMessage {
  event: string;
  data: any;
}

export interface WSTokenValid {
  avatar: string;
  discriminator: string;
  guilds: number;
  id: string;
  username: string;
}

export interface RESTResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface Version {
  version: string;
  commit: string;
  date: string;
}
