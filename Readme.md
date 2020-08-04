# Quests

## Pre-requisites

This bot requires the following tools:

- docker
- docker-compose
- mongodb

### Setting up mongodb

MongoDB can be deployed simply by running the following command:

```bash
docker-compose up
```

## Bot Configuration

The bot can be configured by creating a `config.yaml` file like so:

```yaml
discord:
  token: <token>
  channels:
    teamRocketChannelID: 12345
bot:
  refreshInterval: 10s
database:
  kind: mongodb
  address: localhost:27017
  username: root
  password: example
  name: questrocket

```

More configuration options will be added soon
