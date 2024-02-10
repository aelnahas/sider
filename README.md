# Sider 

A hobbyist project aimed at building redis from the ground up using golang. This project is for educational / entertainment purposes and is not meant to be used in production code.

## Intro

I am using [RESP](https://redis.io/docs/reference/protocol-spec/), which is the same wire protocol that redis uses. This means that for any implemented command you shouldbe able to use a redis client to communicate with this db server. Note that currently the code has been using RESP v2 and does not support v3. In addition to that, the code implementation predates the new client handling scheme that redis currently uses.

Currently the supported commands include:
- PING
- ECHO 
- GET
- SET
- DEL
- PUB/SUB

Currently the db stores the data strictly in memory, therefore the data is not durable.


## Getting Started

### Installation

```bash
go install github.com/aelnahas/sider@latest
```

### Starting the server

```bash
# start the server in background
sider start -d 
```

### Stopping the server
```bash
sider stop
```


## Testing

If you wish to test this with a real redis client you should:

- Install redis older than v7 
- run the server `sider start`
- run cli client `redis-cli`
- execute one of the supported commands

