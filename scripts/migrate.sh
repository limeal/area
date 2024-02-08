#!/bin/sh

# Down all containers
docker-compose down 2>/dev/null
# Up postgres container
docker-compose up -d postgres 2>/dev/null
CONTAINER=$(docker-compose ps -q postgres)
GOV=$(go version | grep "1.19.2")

if [ -z "$GOV" ]; then
  echo "Go version 1.19.2 is required to run migrations. Please install it and try again."
  exit 1
fi

ENV_PATH=$(docker inspect $CONTAINER | grep "com.docker.compose.project.working_dir" | cut -d '"' -f 4 | tr ":" " " | cut -d " " -f 3)
READY=$(docker logs $CONTAINER 2>&1 | grep "database system is ready to accept connections")

if [ -z "$READY" ]; then
  echo "Waiting for database to be ready..."
  sleep 5
  exec $0
fi

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "[${GREEN}OK${NC}] Database is ready to accept connections..."
echo "----------- JOBS -----------"

# Export environment variables
ENVFILE=$ENV_PATH/.env
export $(grep -v '^#' $ENVFILE | xargs)
export POSTGRES_HOST=localhost

cd $ENV_PATH/packages/server/migration/
go run migrate.go
cd $ENV_PATH
echo "----------------------------"
echo "[${GREEN}OK${NC}] Migrations complete. Downing containers..."
docker-compose down 2>/dev/null