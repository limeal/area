# Tutorial (How to add a new authenticator):

## I. Select method for authenticating (Oauth 2.0, OID Authentication, etc..)
Actually only oauth2 is supported
=============

Depend on what you choose go to the directory of the method
```sh
cd oauth2
```

## II. Create your Authenticator file

1. Create a go file with the authenticator name
```sh
touch (name)Authenticator.go
```

2. Add a Descriptor for the authenticator inside this file
```go

func (name)Authenticator() static.OAuth2Authenticator {
	return static.OAuth2Authenticator{
        Name: (Name of the authenticator),
        Enabled: ((true | false) => Activate the service in the Authentication Modal),
        More: static.More{
			Avatar: true,
			Color:  (Color in hexadecimal used for color the authenticator box in client),
		},
        AuthorizationURI: (The authorization URI of the authenticator that will be used),
        AuthEndpoints: {
            AccessToken: (Descriptor of the request that will retreive the access token, See the documentation about RequestLibrary)
        }
	}
}