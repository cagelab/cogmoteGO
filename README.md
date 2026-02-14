<div align=center>
<h1><code>cogmoteGO</code></h1>
<b>"air traffic control" for remote neuroexperiments</b></br/>
</div>
<br/>

## Introduction

`cogmoteGO` is the "air traffic control" for remote neuroexperiments: a lightweight Go system coordinating distributed data streams, commands, and full experiment lifecycle management - from deployment to data collection.

## Bindings

- [for matlab](https://github.com/cagelab/matmoteGO.git)

## Installation

### For Linux & macOS

#### By install script

```sh
curl -sS https://raw.githubusercontent.com/cagelab/cogmoteGO/main/install.sh | sh
```

### For Windows

#### By install script

```sh
irm -Uri 'https://raw.githubusercontent.com/cagelab/cogmoteGO/main/install.ps1' | iex
```

#### By winget
> The winget version is relatively outdated; currently, we recommend installing via a script.

```sh
winget install ccccraz.cogmoteGO
```

## Getting started

### Run as service

#### For Linux & macOS

restart the service as user
> We recommend that you register cogmoteGO as a user service on Linux

```sh
cogmoteGO service -u
```

start the service as user

```sh
cogmoteGO service start -u
```

register the service

```sh
sudo cogmoteGO service
```

start the service

```sh
sudo cogmoteGO service start
```

#### For Windows

register the service

> note: you need to run the command as administrator

```sh
cogmoteGO service
```

start the service

```sh
cogmoteGO service start
```

restart the service as user

> note: the password is required for running the service as user

```sh
cogmoteGO service -u -p <your_password>
```

start the service as user

```sh
cogmoteGO service start -u
```

#### For all platforms

for more info about the service, run

```sh
cogmoteGO service --help
```


#### Test

```sh
curl --location --request GET 'http://localhost:9012/api/device'
```

## Email Configuration

`cogmoteGO` supports sending email notifications. Credentials are securely stored in the system keyring.

### Set email credentials

Interactive setup (prompts for email and password):

```sh
cogmoteGO email set
```

With SMTP server options:

```sh
cogmoteGO email set --host smtp.example.com --port 587
```

### Show current configuration

```sh
cogmoteGO email show
```

Example output:

```
Email Configuration:
  Email address : your@email.com
  SMTP host     : smtp.example.com
  SMTP port     : 587
  Recipients    : recipient1@example.com, recipient2@example.com
  Credentials   : configured
```

### Delete email configuration

```sh
cogmoteGO email delete
```

### For more info

```sh
cogmoteGO email --help
```

## Configuration

Manage application settings via CLI.

### Show configuration

Show all configuration:

```sh
cogmoteGO config show
```

Show specific key:

```sh
cogmoteGO config show port
```

### Set configuration

```sh
cogmoteGO config set port 8080
cogmoteGO config set proxy.max_retries 5
```

### Reset configuration

Reset specific key to default:

```sh
cogmoteGO config reset port
```

Reset all configuration (email settings not affected):

```sh
cogmoteGO config reset
```

### Available settings

| Key | Default | Description |
|-----|---------|-------------|
| `port` | 9012 | Server listening port |
| `instance_id` | auto | Unique instance identifier |
| `proxy.handshake_timeout` | 5000 | WebSocket handshake timeout (ms) |
| `proxy.msg_timeout` | 5000 | Message timeout (ms) |
| `proxy.max_retries` | 3 | Maximum retry attempts |
| `proxy.retry_interval` | 200 | Retry interval (ms) |

### For more info

```sh
cogmoteGO config --help
```
