/** @format */

export type WSGuilds = WSGuild[];

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

export interface WSGuild {
  id: string;
  name: string;
  owner: string;
  members: number;
  icon: string;
}

export interface WSUser {
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
