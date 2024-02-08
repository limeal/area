# How to add a new service to the application

## 0. Remove an existing service

1. If you want to keep the source code because you don't know if later you will use it

Go to ServiceRegistry.go file and comment the line of the service you want to remove

2. If you want to delete completely the service

Delete the folder associated with the service
WARNING: Do not delete **common** folder

## I.  Service Descriptor 

1. Create a sub folder in this directory with the service name you want to use

```sh
mkdir -p (name)
```

2. Create a go file with the service name (e.g. spotify.go)
```sh
cd (name) && touch (name).go
```

3. Add a Descriptor for the service inside this file
```go

func Descriptor() static.Service {
	return static.Service{
		Name: (Name of the service), // Mandatory
		Description: (Description of the service), // Mandatory
		Authenticator: (Authenticator of the service, can be null if service doesn't need auth),
		RateLimit: (RateLimit of the service, correspond to how many time he will wait before recheck, eg. 3 -> 30/3 he will wait 10 seconds)
		Validators: (Validators for the service, used to verify parameters when creating a new applet -> #Validators)
		Endpoints: (List of endpoints used by the server for the service -> #Endpoint) // Mandatory
		Routes: (List of routes for the service api accessible at /services/{name}/api/
		Actions: (List of actions of the service) // Mandatory
		Reactions (List of reactions of the service) // Mandatory
	}
}
```

4. Go back to services folder and add the descriptor in the `ServiceRegistry.go` file

## 2. Service Endpoints

1. In the same folder of the descriptor create a file following this convention -> (name)_endpoints.go

2. 
a. Inside  this file create a function that will register the differents endpoints, func name must follow the same convention that below:

```go
func (NameOfService)Endpoints static.ServiceEndpoints {
	return static.ServiceEndpoints{
		"(NameOfEndpoint)Endpoint": {
			BaseURL: (Url on service that correspond), //Mandatory
			Params: (NameOfEndpoint)EndpointParams, //Can be nil, but not recommended for specific case
			ExpectedStatus: []int{(status)} // Mandatory
		},
		...
	}
}

example:

func SpotifyEndpoints static.ServiceEndpoints {
	return static.ServiceEndpoints{
		"GetUserProfileEndpoint": {
			BaseURL: "https://api.spotify.com/v1/me",
			Params: GetUserProfileEndpointParams,
			ExpectedStatus: []int{200},
		},
		"FindTrackByIDEndpoint": {
			BaseURL: "https://api.spotify.com/v1/tracks/${id}",
			Params: FindByIDEndpointParams,
			ExpectedStatus: []int{200},
		},
	}
}

```

b. Params parameter correspond to the params of the request and must follow this convention

```go
func (NameOfEndpoint)EndpointParams(
	params []interface{} // Correspond to your params that you will add when calling the request
) *utils.RequestParams {
	return &utils.RequestParams{
		Method: (HTTPRequestMode -> 'GET' | 'POST' | 'PUT' | 'DELETE' | ...)
		Headers: (Headers of the request in a map[string]interface{},
		UrlParams: (Url parameters of the request it will replace the variable inside the url like the `${id}` in the example above)
		QueryParams: (Query Parameters of the request in a map[string]string)
		Body: body of the request (string)
	}
}

example (with the endpoint of before):

func GetUserProfileEndpointParams(params []interface{}) *utils.RequestParams {
	return &utils.RequestParams{
		Method: 'GET',
		Headers: map[string]string{
		"Authorization": "Bearer "  + params[0](*models.Authorization).AccessToken,
		"Content-Type": "application/json",
		},
	}
}

func  FindByIDEndpointParams(params []interface{}) *utils.RequestParams {
	return  &utils.RequestParams{
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer "  + params[0](*models.Authorization).AccessToken,
			"Content-Type": "application/json",
		},
		UrlParams: map[string]string{
			"id": params[1].(string),
		},
	}
}
```

## 3. Service Validators

Used to check validity of parameters provided by the client when creating applet (prevent crash of the application before creation)

1. In the same folder of the descriptor create a file following this convention -> (name)_validators.go

2. 
a. Inside  this file create a function that will register the differents validators, func name must follow the same convention that below:

```go
func (NameOfService)Validators static.ServiceValidator {
	return static.ServiceValidator{
		<nameOfField>: <validator>
	}
}
```
nameOfField: Correspond to the name of an element from the RequestStore of an action / reaction
validator: Function prototyped has follow
```go
func  (Name)Validator(authorization *models.Authorization, service *static.Service, value interface{}) bool
```
If the bool returned is false then it is incorrect, otherwise is ok

## 4. Service Routes

Used by service for declaring is own route that will be added at path: `/services/{name}/api/`, it work in same way that normal route

# How to add a new action/reaction to a service

## I. Declaration of Action & Reaction 

1. Inside the service folder when you want to add an action create a directory named `actions` and a directory name `reactions` and move into one of them depend on if you want to create an action or a reactions

```sh
mkdir -p actions && mkdir -p reactions && cd [folder]
```

2. Create a go file with the action name, convention: `Action(Name)` or `Reaction(Name)` depend on what you want to create

```sh
touch Action(Name).go | touch Reaction(Name).go
```

3. Add a Descriptor for the action inside this file

```go
func DescriptorFor(Service)(Type)(Name)() static.ServiceArea {
	return static.ServiceArea{
		Name: (NameOfTheAction),
		Description: (DescriptionOfTheAction)
		RequestStore: (Field that will be demand by client more in the #RequestStore section)
		Method: (Method that will be call that the trigger and will be the core of your action)
		Components: (Component exported that will be used by client to know which variable it can use for reaction)
	}
}

example:
func  DescriptorForSpotifyActionNewTrackAddedToPlaylist() static.ServiceArea {
	return static.ServiceArea{
		Name: "new_track_added_to_playlist",
		Description: "Triggered when a playlist has a new track added to it",
		RequestStore: map[string]static.StoreElement{
			"req:playlist:id": {
				Description: "ID of the playlist to watch for new tracks",
				Required: true,
				Type: "select_uri",
				Values: []string{"/playlists"},
			},
		},
		Method: hasNewTrackAddedToPlaylist,
		Components: []string{ //Only useful for action
			"spotify:playlist:id", // 37i9dQZF1DXcBWIGoYBM5M
			"spotify:playlist:name", // Rock Classics
			"spotify:playlist:snapshot:id", // 0QJ4q7q7X3ZJZy9Y8X5Z0w
			// Track
			"spotify:track:uri", // spotify:track:6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:name", // We Will Rock You
			"spotify:track:duration", // 354000
			"spotify:track:href", // https://api.spotify.com/v1/tracks/6rqhFgbbKwnb9MLmUQDhG6
			"spotify:track:preview:url", // https://p.scdn.co/mp3-preview/08b7e...
			// Album
			"spotify:album:id", // spotify:album:6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:name", // News Of The World
			"spotify:album:total_tracks", // 10
			"spotify:album:href", // https://api.spotify.com/v1/albums/6JWc4iAiJ9FjyK0B59ABb4
			"spotify:album:release:date", // 1977-11-10
			"spotify:album:type",
		},
	}
}
```

4. Add the action / reaction in the service Actions or Reactions Field

## II. Request Store

The request store correspond to the field that the client must provide, each field has a different type depending on his utility

In this example:
```go
	RequestStore: map[string]static.StoreElement{
			"req:playlist:id": {
				Description: "ID of the playlist to watch for new tracks",
				Required: true,
				Type: "select_uri",
				Values: []string{"/playlists"},
	},
},
```
The store contains 1 field named req:playlist:id, that have different properties:

a. Description: The description of the field (Can be used by client)
b. Required: Is this field mandatory or not
c. Type: One of (string | email | number | select_uri | select | date)

### Types

1. String -> Will take an normal string as value
2. Select -> Will take on of the following value contains in `Values` field
3. Select_uri -> Will take on of the following value of the endpoint (Service Route) describe in `Values[0]`

d. Values: Used by select & select uri to contains additional data
e. NeedField: Need to provide this if you want to wait that all of fields contains here is filled
f. AllowedComponents: The component that are allowed for this field [Not implemented yet]

### Warning: Naming convention for request Store

All variable must be prototyped has follow req:(name), name must not contains any `_`, replace it with `:`

## III. Method

It is the main entry of your action/reaction it is prototyped has follow:
(name) can be want you want but i recommend that it describe the action/ reaction you are creating

```go
func (name)(req static.AreaRequest) shared.AreaResponse {
}
```

### AreaRequest Structure

```go
type AreaRequest struct {
	AppletID uuid.UUID
	Self *ServiceArea
	Authorization *models.Authorization
	Service *Service
	Logger *shared.Logger
	Store *map[string]interface{}
	ExternalData map[string]interface{}
}
```

AppletID -> Correspond to the id of the applet that the action / reaction use
Self -> Pointer on self
Authorization -> The authorization of the service, that can be used by the action / reaction to make request
Service -> Pointer on the service of the action / reaction
Logger -> Pointer on the logger that can be used to log data for the applet
Store -> Modifiable store that combines value provided by the client and your context variables
ExternalData -> Only filled for Reaction, it is the variable exported by the action when triggered

### AreaResponse Structure

```go
type AreaResponse struct {
		Error error  `json:"error"`
		Success bool  `json:"sucess"`
		Data map[string]interface{} `json:"data"`
}
```
Error: If not nil, then the application will crash with the error
Success: Only used by action to say if the action has been trigered or not
Data: Only used by action to export variable that will be put in ExternalData for the reaction

# Naming Conventions For Actions & Reactions store

- req -> Use for the request store (Data provided by the client)
- ctx -> Use for the context store (Data being used by the program, like a cache store [Not repeating])
- <service> -> Use for the service store (Data being returned by the service)
