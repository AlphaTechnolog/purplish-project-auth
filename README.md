# Purplish Project - Auth

Micro for the auth functionality

## Setting up the project

Follow the next instructions to get the micro up and running

### Dependencies

Make sure you have the next dependencies on the target system:

- sqlite3
- go

### Deploying

Run the next commands on your system:

```
./.bin/run-migrations.sh
go mod tidy
PORT=8007 go run .
```
