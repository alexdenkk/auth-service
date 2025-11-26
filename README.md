# Auth service
An authentication service that follows good coding practices

<br>

## Stack
- Go
- Echo
- DDD
- Zap
- uberfx
- JWT
- Docker(healthcheck)
- Caddy
- PostgreSQL
- Postman specification

<br>

## Architecture
```sh
├── cmd # program entrance
├── deployments # deploy configs
│   ├── docker
│   └── docker-compose
├── internal
│   ├── domain # business logic isolated layer
│   ├── handler # http layer
│   └── repository # database layer
├── migration # migrations for db
└── pkg
    └── config # service configuration
```

<br>

## Logs
 Logs stored in deployments/docker-compose/log/log.json file by default

<br>

## Run
#### 1. Clone repository
```sh
git clone https://github.com/alexdenkk/auth-service.git
```
#### 2. Configure .env and Caddyfile if needed
```sh
cd deployments/docker-compose
vi .env
vi Caddyfile
```
#### 3. Launch postgresql db
```sh
docker-compose up -d postgres
```
#### 4. Execute migrations
```sh
psql -h <host> -p <port> -U <user> -d users -f 000_create_migrations_table.sql
psql -h <host> -p <port> -U <user> -d users -f 001_create_users_table.sql
```
#### 5. Launch service
```sh
docker-compose up -d auth-service
```
