# Portier

```
NAME:
   portier - Provides convenience functions for corporate Yandex.Taxi accounts

USAGE:
   portier [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   users, u  Enables or disables orders from application for user role
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --sessionid value, -s value  value of Yandex Session_id cookie [$SESSION_ID]
   --clientid value, -c value   Yandex client id [$CLIENT_ID]
   --help, -h                   show help
   --version, -v                print the version
```

## Users

```
NAME:
   portier users - Enables or disables orders from application for user role

USAGE:
   portier users [command options] [arguments...]

OPTIONS:
   --role value  operate on users in this role only (default: "Кирпичников")
```
