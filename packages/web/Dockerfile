FROM node:16.17.0-buster

RUN npm install -g pnpm
RUN npm install -g serve

WORKDIR /client

# Install dependencies
COPY package.json .
RUN pnpm install
COPY [".eslintrc.json", ".prettierrc", "tsconfig.json", "/client/" ]
COPY . .

# Build the app
RUN pnpm build
CMD ["./start.sh"]