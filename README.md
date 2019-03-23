# Location tracker server
Server for [location tracker](https://github.com/moritzschramm/location-tracker)

## Installation
```
sudo ./install.sh
```

## Project structure
- `/api`: router and controllers
- `/certs`: certificates used by the server
- `/config`: configuration and environment files for MQTT and API
- `/database`: database file (when using SQLite3) and initialization statement
- `/frontend`: vue.js components and build HTML files
- `/model`: database models
- `/mqtt`: client to connect to the MQTT broker