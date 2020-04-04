# Web Socket API Reference

The official API endpoint URL is
```
wss://tokentools.zekro.de/ws
```

The web socket API is based on JSON encoded UTF-8 encoded string messages in the following format:

| Property | Type | Description |
|----------|------|-------------|
| `event` | `string` | The name of the event. Events are always sent and handled lower-cased. |
| `cid` | `int` | A Command ID, set on Request - Response to this Request will contain the same Command ID |
| `data` | `any` | The data as JSON object |

Example:
```json
{
  "event": "valid",
  "cid": 4,
  "data": {
    "id": "524847123875889153",
    "name": "shinpuru",
    "discriminator": "4878",
    "avatar": "https://cdn.discordapp.com/avatars/524847123875889153/9ccd16b7555487648d4e1c38f27aba91.png",
    "guilds": 134
  }
}
```

## Objects

### Error

| Property | Type | Description |
|----------|------|-------------|
| `code` | `int` | The Error Code |
| `message` | `string` | The Error Message |

### User

| Property | Type | Description |
|----------|------|-------------|
| `id` | `string` | The User ID |
| `username` | `string` | The Users Name |
| `discriminator` | `string` | The Users Discriminator |
| `avatar` | `string` | The Users Avatar URL |
| `guilds` | `int` | The Ammount of Guilds the User is Member of |

### Guild

| Property | Type | Description |
|----------|------|-------------|
| `id` | `string` | The Guild ID |
| `name` | `string` | The Guilds Name |
| `owner` | `string` | The User ID of the Guilds Owner |
| `member` | `int` | The Ammount of Members of the Guild |
| `icon` | `string` | The Hash of the Guilds Icon |

### GuildTile

| Property | Type | Description |
|----------|------|-------------|
| `guild` | `Guild` | The Guild Object |
| `n` | `int` | The Number of the Guild |
| `nmax` | `int` | The total Number of expectable Guild Obects |

## Command Events

| Event | Data | Description |
|-------|------|-------------|
| `init` | `string` | Sets the API token for the connection - no Response |
| `check` | - | Checks if the set Token is valid - responses with `valid` or `invalid` |
| `guildinfo` | - | Requests the Guild Info of all Guilds in seperated Events - responses with `guildinfo` if successful |
| `userinfo` | `string` | Requests the User Info of the passed User ID if accessable - Responses with `userinfo` if successful |

## Response Events

| Event | Data | Description |
|-------|------|-------------|
| `error` | `Error` | Occurs if something fails |
| `valid` | `User` | Response to `check` command if the token is valid - Contains the User Object of the Bot Account |
| `invalid` | `null` | Response to `check` command if the token is invalid |
| `guildinfo` | `GuildTile` | Response to `guildinfo` command - Guilds will be sent in multiple events |
| `userinfo` | `User` | Reponse to `userinfo` command - contains the User Object of the requested User ID if accessable |