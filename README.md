
# AREA - Team Project

Time: 2month
Members:
- @myselfandme => FrontEnd Developer (Mobile)
- @CorentinLan => FrontEnd Developer (Mobile)
- @limeal => Backend Developer & FrontEnd Contributor (Web, Mobile)
- @p0lar1s => FrontEnd Developer (Web) && Backend Content Contributor (Actions & Reactions)
- @calmifaire => Backend Content Contributor (Actions & Reactions)

## Organizations

[Clickup](https://app.clickup.com/9003027499/v/dc/8c9yk1b-142)

## Project Architecture

Project is nested in multiple subdirectories for each main functionality.

### Web Folder (/web)

It contains all the needed files to startup a react web application

- Language: Javascript => **Typescript**
- Main Framework: **React**
- Main Developer: **@Limeal**
- Work without docker: **yes**, but with fake datas (Faker.js)

Run command:

```sh
npm install
npm start
or
yarn install
yarn start
or
pnpm install
pnpm start
```

### Server Folder (/server)

It contains all the needed files to start a fiber web server
 
- Language: **Golang**
- Main Framework: **Fiber**
- Main Developer: **@Limeal**
- Other Member: **@calmifaire**, **@myselfandme**
- Work without docker: **no**, because of the link to databases (postgres and redis)

Run command:

```
At root: docker-compose up --build postgres redis server
In folder: Nope
```

### Mobile Folder (/mobile)

It contains all the needed files to start a flutter project

- Langauge: Dart
- Main Framework: Flutter
- Main Developer: **@limeal**
- Other Members: **@CorentinLan**, **@myselfandme**
- Work without docker: yes, but with fake datas

Run command:

```sh
flutter emulator --launch <emulator_name>
flutter run
```

## Test project

A Tutorial by @limeal

1. Copy .env.example file and fill it with the good stuff
2. Create Oauth2 Apps for each of this applications:

```
Discord
Dropbox
Github
Google (Gmail, Youtube)
OpenWeather
Reddit
Spotify
Time
Twitch
```

To add/remove services check the file Tutorial.md in server/services folder

3. Copy .services.example file and fill it with the good stuff
4. Build `docker-compose.yml` file
```sh
docker-compose build [--no-cache, -d]
```
5. Run `docker-compose.yml` file
```sh
docker-compose up [--force-recreate, -d]
```

## Bonus

- Websocket discord
- Applet:
    - Logs management (Websocket)
    - Start & Stop & Delete & OnOff
    - Multiple reactions
- Oauth2:
    - Can remove authorizations (+ Prevent used authorization to be deleted)
- Areas:
    - Client can select in a list a value
    - Validators of parameters (Prevent Crash of applet before start - ~90%)
    - Webhook management
    - Custom internal route for each service (Interfacing)
- Profile:
    - Avatar
