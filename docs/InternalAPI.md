Internal API - AREA v1.1

# No-Authentified Routes

## Area Information - About.json

================================
GET - /about.json
================================

```json
Response Body:
{
  "client": {
      "host": "10.8.4.2",
  },
  "server": {
    "current_time": 123423,
    "services": [
      {
        "name": "example",
        "authenticator": { // Can be null
          "name": "spotify",
          "authorization_uri": "https://example.com/authorize",
          "enabled": false,
          "more": {
            "avatar_uri": "http://example.com/logo.png",
            "color": "#3453",
          }
        },
        "more": {
          "avatar_uri": "http://example.com/logo.png",
          "color": "#3453",
        }, // Optional (use the one in authenticator if not provided, but this one is prioritary)
        "actions": [
          {
            "name": "check season",
            "description": "check if we are in winter",
            "components": [...],
            "store": {
                "current_season": {
                    "description": "Param that describe the current season",
                    "required": true,
                    "type": "select_uri",
                    "type_ctor": "/lol/test",
                    "need_fields": [...],
                    "value": ""
                }
             }
          }
        ],
        "reactions": <-like actions
      }
    ],
    "authenticators": [
      {
        "name": "spotify",
        "authorization_uri": "https://example.com/authorize",
        "enabled": false,
        "more": {
          "avatar_uri": "http://example.com/logo.png",
          "color": "#3453",
        }
      }
    ]
  }
}
```

## Authentication (How to authenticate)

### External (Use a service, e.g Github, Google)

================================
POST - /auth/external
================================

```json
Request Body:
{
   "service": "github",
   "code": "ekleklekek,cjnhbehyf",
   "redirect_uri": "http://localhost:8081"
}
```

### Login

================================
POST - /auth/login
================================

```json
Request Body:
{
  "email": "example@test.com",
  "encoded_password": "example"
}
```

### Register

================================
POST - /auth/register
================================

```json
Request Body:
{
  "email": "example@test.com",
  "encoded_password": "example"
}
```

Response Body:

```json
{
  "code": 201,
  "data": {
    "message": "Account created !",
    "token": "exampletoken"
  }
}
```

## Store (Get all public applets)

================================
GET - /store
================================

Response Body:

```json
{
    "code": 200,
    "data": {
        "applets": <applets>, (Look applet section for details)
	},
}
```

# Authentified Routes ➝ Authorization: Bearer <token>

## Profile (Get all information from user connected and make some action with it)

### Get Account

================================
GET - /me (Retrieve account information)
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "user": {
      "email": "example@example.com",
      "username": "example"
    }
  }
}
```

### Modify Account

================================
PUT - /me (Only username at the moment)
================================

Request Body:
All fields are optional

```json
{
  "username": "hugolezozo"
}
```

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "User updated"
  }
}
```

### Logout Account (Not used if token mode)

================================
POST - /logout
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "User logged out"
  }
}
```

### Delete Account

================================
DELETE - /me
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "User deleted"
  }
}
```

### Get Avatar

================================
GET - /me/avatar
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "uri": "https://www.pngitem.com/pimgs/m/30-307416_profile-icon-png-image-free-download-searchpng-employee.png"
  }
}
```

### Modify Avatar

================================
PUT - /me/avatar
Type: multipart/form
================================

Request Body:

```multipart/form
avatar: image
```

## Authorization (Manage authorization needed by some services and check which service is already authorized)

### Get authorizations

================================
GET - /authorization
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "authorizations": [
      {
        "name": <name>,
        "type:": "oauth2",
        "permanent": <true | false>,
        "expired_at": <time>
      }
    ]
  }
}
```

### Create authorization

================================
POST - /authorization
================================

Request Body:

```json
{
  "authenticator": <service>,
  "code": <code>
  "redirect_uri": <redirect_uri>
}
```

```json
Response Body:
{
  "code": 201,
  "data": {
    "message": "Authorization created !",
  }
}
```

### Delete authorization

================================
DELETE - /authorization/:authorization_id
================================

```json
{
  "code": 200,
  "data": {
    "message": "Authorization deleted !"
  }
}
```

### List services authorized from User

================================
GET - /authorization/services
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "services": {
      "example": true,
      "something": false
    }
  }
}
```

## Applets

### Get applets

================================
GET - /applet
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "applets": [
      ...
     ]
  }
}
```

### Get New Applet (Step 1)

================================
GET - /applet/new
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    ?"reaction": <reaction_info>,
    ?"action": <action_info>
  }
}
```

### Add New State To Applet (Step 2 - 2x)

================================
PUT - /applet/new
================================

Request Body:

```json
{
  "service": "spotify",
  "area_type": "action",
  "area_item": "new_saved_track",
  "area_settings": {
    "playlist_id": "5ykcBjrOlzH9RP4F3crnB8",
    ...
  }
}
```

```json
Response Body:
{
  "code": 200,
  "data": {
    "message": "Applet updated !"
  }
}
```

### Submit New Applet (Step 3)

================================
POST - /applet/new
================================

Request Body:

```json
{
  "name": "test",
  "description": "this is a test applet",
  "public: "public" | "private",
}
```

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "Applet submitted !"
  }
}
```

### Delete State in New Applet

================================
DELETE - /applet/new
================================

Query Params:
type= action | reaction

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "Area deleted !"
  }
}
```

### Get an applet

================================
GET - /applet/:applet_id
================================

Response Body:

```json
{
  "id": <applet_id>,
  "name": <applet_name>,
  "description": <applet_name>,
  "active": true,
}
```

### Modify an applet (Not working actually)

================================
PUT - /applet/:applet_id
================================

### Modify an applet activity

================================
PATCH - /applet/:applet_id
================================

Query Params:
active= true | false

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "Applet activity updated"
  }
}
```

### Delete an applet

================================
DELETE - /applet/:applet_id
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "Applet deleted"
  }
}
```

### Start an applet

================================
PUT - /applet/:applet_id/start
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "Applet started"
  }
}
```

### Stop an applet

================================
PUT - /applet/:applet_id/stop
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "message": "Applet stopped"
  }
}
```

### Get context information on an applet

================================
GET - /applet/:applet_id/context
================================

Response Body:

```json
{
  "code": 200,
  "data": {
    "context": {
      "applet_id": "exampleid",
      "times_triggered": 1
    }
  }
}
```

### Get context logs report on an applet (Read log file)

================================
GET - /applet/:applet_id/context/logs
================================

Response Body:

```json
{
    "code": 200,
	"data": {
		"logs": <file content in array>
	},
}
```

## Authenticators

### Get Authenticators

================================
GET - /authenticators
================================

Response Body:

```json
[
    {
        "name": "example",
        "enabled": true,
        "more": {
            "avatar_uri": "http://example.com",
            "color": "green",
        },
        "authorization_uri": "http://example.com",
    },
    ...
]
```

### Get an Authenticator

================================
GET - /authenticators/:authenticator
================================

Response Body:

```json
{
    "name": "example",
    "enabled": true,
    "more": {
        "avatar_uri": "http://example.com",
        "color": "green",
    },
    "authorization_uri": "http://example.com",
},
```

## Services

### Get Services

================================
GET - /services (Get list of all existing services)
================================

Response Body:

```json
[
    {
        "name": "example",
        "description": "example service",
        "authenticator": <like above (Authenticator)>, // Can be nil
        "more": { // Optional
            "avatar_uri": "http://example.com",
            "color": "green",
        },
        "actions": [<Action>],
        "reactions": [<Reaction>],
    }
]
```

### Get a service

================================
GET - /services/:service_name (Get a specific service)
================================

Response Body:

```json
{
    "name": "example",
    "description": "this is an example service",
    "authenticator": <like above (Authenticator)>, // Can be nil
    "more": { // Optional
        "avatar_uri": "http://example.com",
        "color": "green",
    },
    "actions": [<like below (Action)>],
    "reactions": [<like below (Reaction)>],
}
```

### List all actions in a service

================================
GET - /services/:service_name/actions
================================

Response Body:

```json
[
    {
        "name": "action example",
        "description": "this is an example action",
        "use_gateway": false,
        "components": ["..."], // Array of string
        "store": {
            "example": {
                "description": "this is an example description for a parameter",
                "required": true,
                "type": "string", // Can be string | select | select_uri | textarea
                "need_fields": ["..."],
                "allowed_components": ["..."],
                "values": ["..."],
            }
        }
    },
    ...
]
```

### Get an action from a service

================================
GET - /services/:service_name/actions/:action_name
================================

Response Body:

```json
{
  "name": "action example",
  "description": "this is an example action",
  "use_gateway": false,
  "components": ["..."], // Array of string
  "store": {
    "example": {
      "description": "this is an example description for a parameter",
      "required": true,
      "type": "string", // Can be string | select | select_uri | textarea
      "need_fields": ["..."],
      "allowed_components": ["..."],
      "values": ["..."]
    }
  }
}
```

### List all reactions in a service

================================
GET - /services/:service_name/reactions
================================

Response Body:

```json
[
    {
        "name": "reaction example",
        "description": "this is an example reaction",
        "use_gateway": false,
        "components": ["..."], // Array of string
        "store": {
            "example": {
                "description": "this is an example description for a parameter",
                "required": true,
                "type": "string", // Can be string | select | select_uri | textarea
                "need_fields": ["..."],
                "allowed_components": ["..."],
                "values": ["..."],
            }
        }
    },
    ...
]
```

### Get a reaction from a service

================================
GET - /services/:service_name/reactions/:reaction_name
================================

Response Body:

```json
{
  "name": "reaction example",
  "description": "this is an example reaction",
  "use_gateway": false,
  "components": ["..."], // Array of string
  "store": {
    "example": {
      "description": "this is an example description for a parameter",
      "required": true,
      "type": "string", // Can be string | select | select_uri | textarea
      "need_fields": ["..."],
      "allowed_components": ["..."],
      "values": ["..."]
    }
  }
}
```

### Get all API route of a service
================================
GET - /services/:service_name/api
================================

```json
[
    {
        "endpoint": "/t",
        "method": "GET",
        "need_auth": true
    }
]
```

# Other are dynamic route generate by service
