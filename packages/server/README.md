# Server Folder

## 1. How to start

a. Launch docker-compose in root folder
```
docker-compose up --build
```

### MIGRATION

You can migrate databases (reset) with script `migrate.sh` in folder scripts/ locate at root of the repository

### 1. API (Routes) => /api

In this folder, you can access to all the different routes that are available to the application

For more information about routes, see the documentation about endpoints

### 2. Classes => /classes

In this folder, you can access to all the classes that we used for the management of each different
services and authenticator => /static
The trigger folder that implements class for the main controller of the application (link action and reaction) => /triggers
The shared folder that implements the logger class and the response class of each action & reaction => /shared

### 3. Config => /config

In this folder, you can access all the config of the server to change different state of the server

### 4. Database => /db

In this folder, you can access to all the mangement of database table that we used for the project