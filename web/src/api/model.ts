/** @format */

export interface WSMessage {
  event: string;
  data: any;
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
