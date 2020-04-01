/** @format */

import { RESTResponse, Version } from './model';

const PREFIX =
  process.env.NODE_ENV === 'development' ? 'http://localhost:8081' : '';

export default class RestAPI {
  public static async getVersion(): Promise<RESTResponse<Version>> {
    return window.fetch(PREFIX + '/api/info').then((res) => res.json());
  }
}
