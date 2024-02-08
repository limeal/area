# Web (Area)

## 1. About the web application

How to start the web application:

```sh
pnpm install
pnpm start
```

How to build the web application:

```sh
pnpm install
pnpm build
```

No business logic is provided for this web application, all the different services and tools we used
is located in the server folder.

Packages dependencies with used:

-   styled-components (For styling all different components, don't use normal css)
-   eslint (For eslint support)
-   prettier (For prettier support)
-   react-toastify (For usage of toast)
-   react-icons (For usage of icons)
-   react-redux (For help with api interactions)
-   framer-motion (For animation support)
-   chart.js (For chart and graph visualization)
-   normalize.css (For normalization)
-   ncrypta (For password encryption)

## 2. Architecture

This is how we decide to develop the application:

-   Fonts Folder: It is the folder where is located all the fonts we used
-   Hooks Folder: It is the folder where is located all the custom hooks we used
-   Interface Folder: It is the folder where is located all the interfaces we used
-   Pages Folder: It is the folder where is located all the different pages of the application
-   Redux Folder: It is the folder where is located all the different Redux query and tools
-   Utilities Folder: It is the folder where is located all the different utilities functions and helpers

Entrypoint file: index.tsx

## 3. Lint

We used eslint and prettier to lint the code and optimize the clarity of the application code

Config files:

```sh
.prettierc
.eslintrc.json
```

Lint commands:

```
pnpm lint
pnpm lint:fix
pnpm format
```
