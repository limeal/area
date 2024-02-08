# Area - List of services

## Spotify

### Oauth

Method: Authorization Code Flow
Authorization Endpoint: https://api.spotify.com/authorize
=======

Request Body:

content-type: application/json
data:
```json
{
    "client_id": "b6dc6e8e6ab840fcaea14337b91265fa",
    "response_type": "token",
    "redirect_uri": "http://localhost:8081",
    "scope": [...],
    "show_dialog": false
}
```

Response Body:

content-type: application/json
data:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "token_type": "Bearer",
    "expires_in": 200
}
```

## Discord

Doc -> https://discord.com/developers/docs/topics/oauth2

### Oauth

Method: Authorization Code Flow
Authorization Endpoint: https://discord.com/oauth2/authorize
=======

Request Body:

content-type: application/x-www-form-urlencoded
data:
```json
{
    "client_id": "b6dc6e8e6ab840fcaea14337b91265fa",
    "response_type": "token",
    "scope": "identify"
}
```

Response Body:

content-type: application/json
data:
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "token_type": "Bearer",
    "expires_in": 200,
    "scope": null
}
```

## Facebook

WIP -> https://developers.facebook.com/docs/javascript/quickstart

## Reddit

WIP -> https://github.com/reddit-archive/reddit/wiki/OAuth2

## YouTube

WIP -> https://developers.google.com/youtube/v3/guides/auth/client-side-web-apps

## Github

WIP -> https://docs.github.com/fr/developers/apps/building-oauth-apps/authorizing-oauth-apps

## Clickup
## Instagram
## Twitch