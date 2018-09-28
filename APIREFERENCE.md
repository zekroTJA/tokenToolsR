# TOKEN TOOLS R - API REFERENCE

---

The general API root point is
```
https://tokentools.zekro.de/api
```

The API endpoints will always respond in following JSON format: 

> This is an example successful response
```json
{
  "code": 200,
  "message": "OK",
  "data": {}
}
```

Currently used reponse codes are:

| Code | Message | Description |
|------|---------|-------------|
| `200` | `OK` | Everything is fine, request got proceed and responded. |
| `400` | `BAD_REQUEST` | Somwthign is wrong with the format of your request. Further information may be included in `data` value. |
| `900` | `RATE_LIMITED` | You have send to many requests in to short time distances. Calm down and wait some time until you send another request. |

---

## Endpoints

### `GET /check/:token`

**Parameters**

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `token` | `string` | The discord bot token to check | Yes |

**Response:**

> If token is valid:
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "valid": true,
    "id": "272336949841362944",
    "id": "726485837482472384",
    "username": "imNotThatCreative",
    "discriminator": "1234",
    "avatar": "https://cdn.discordapp.com/avatars/272336949841362944/aas242asd13equdasd12987eq224iad.png",
    "guilds": 69
  }
}
```

> If token is invalid:
```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "valid": false,
    "id": "",
    "id": "",
    "username": "",
    "discriminator": "",
    "avatar": "",
    "guilds": 0
  }
}
```

### `GET /guilds/:token`

**Parameters**

> If you use this endpoint without using the `GET /check/:token` endpoint before, it could take up some time until the API will respond. This is because the Discord API needs to be reuqested for the IDs of the Guilds the token is connected with, then, the Discord API will be requested for the Information of the guilds after. To not run into a rate limit, this takes up some waiting time in between. So, if you use the check endpoint before, the IDs will be cached and re-used after calling the guilds endpoint.

| Parameter | Type | Description | Required |
|-----------|------|-------------|----------|
| `token` | `string` | The discord bot token to check | Yes |

**Response:**

> If the token is invalid, the API will respond with a `200` `BAD_REQUEST` error and `"token invalid"` as data value.

```json
{
  "code": 200,
  "message": "OK",
  "data": [
      {
      "id": "12312312323413413123",
      "name": "some",
      "owner": "2849237897234234234",
      "members": 31
    },
    {
      "id": "4568894738953875347",
      "name": "stupid",
      "owner": "2374892643892374989",
      "members": 65
    },
    {
      "id": "389583095868974",
      "name": "guilds",
      "owner": "237489637865892374",
      "members": 98
    }
  ]
}
```