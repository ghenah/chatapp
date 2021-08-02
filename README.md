# chatapp - v0.3.0

## **Installation**

All command line operations below are assumed to be executed from the project's root folder. Adjust accordingly.

### 1. **Environment setup**

The web client will need three environment variables to be able to talk to the server. Copy the `.env_template` file to `.env` inside the client folder (the web client has a separate environment file).

```
cp ./client/.env_template ./client/.env
```

Make sure the port matches the one in the main environment file (`./.env`). Example of a filled in `./client/.env`:

```
VUE_APP_SCHEMA=http://
VUE_APP_DOMAIN=localhost
VUE_APP_PORT=8010
```

The main `.env` file currently contains the settings for the docker-compose file, MySQL, and the app's http server. You will need to copy the `.env_template` file to `.env`

```
cp ./.env_template ./.env
```

The `MYSQL_DB_HOST`, `MYSQL_DB_DATABASE_SETTINGS`, and `APP_ADDRESS_PORT` variables are predefined. An example of the file is provided below (`APP_WS_ORIGIN_*` variables values should match the ones in the web client's `.env` file; `<fill in yourself>` values should be substituted for actual values - secure passwords and secret keys):

```
MYSQL_DB_ROOT_PASSWORD=<fill in yourself>
MYSQL_DB_USER=chatappadmin
MYSQL_DB_PASSWORD=<fill in yourself>
MYSQL_DB_HOST=mysql
MYSQL_DB_DATABASE=chatapp
MYSQL_DB_DATABASE_SETTINGS=?parseTime=True

APP_ADDRESS_HOSTNAME=
APP_ADDRESS_PORT=8081
APP_JWT_SECRET_KEY=<fill in yourself>
APP_JWT_REFRESH_TOKEN_SECRET_KEY=<fill in yourself>
APP_JWT_WEB_SOCKET_SECRET_KEY=<fill in yourself>
APP_WS_ORIGIN_SCHEMA=http://
APP_WS_ORIGIN_DOMAIN=localhost
APP_WS_ORIGIN_PORT=8010
```
