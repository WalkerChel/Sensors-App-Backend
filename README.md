# Sensors App Backend Side

## Start up requirements
* Docker & Docker Compose

### How to start the service
1. Move to the project directory.
2. Create ```.env``` file and copy ```.env.example``` content into it(or use command ```mv .env.example .env```). After that simply fill in all required variables.
3. Run ```docker compose up --build app``` command to start application.

### Manual Start

For this case you need to make sure ```Go 1.22+```, ```make``` and ```migrate``` to be installed locally.

1. Make yourself familiar with available makefile commands. Run ```make help``` to see commands.
2. Use ```make infrastructure-up``` to start all dependencies.
3. Use ```make migrate-up``` to apply migrations to the database.
4. Run ```go run cmd/sensors-app/main.go``` to start the application.
