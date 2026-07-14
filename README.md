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
| `internal_port` | 9011 | Local-only internal API listening port |
| `instance_id` | auto | Unique instance identifier |
| `proxy.handshake_timeout` | 5000 | WebSocket handshake timeout (ms) |
| `proxy.msg_timeout` | 5000 | Message timeout (ms) |
| `proxy.max_retries` | 3 | Maximum retry attempts |
| `proxy.retry_interval` | 200 | Retry interval (ms) |

### Internal API Port

`cogmoteGO` uses a separate listener for APIs that must only be available to
programs running on the same host. The listener always binds to `127.0.0.1` and
uses `internal_port`, which defaults to `9011`. Changing the port does not expose
the internal API to the network.

The public `port` and `internal_port` must be different. Configure the internal
port and restart the service to apply the change:

```nu
cogmoteGO config set internal_port 9011
```

The backup endpoints are registered only on the internal listener. Local Python,
MATLAB, and C# clients should therefore use
`http://127.0.0.1:<internal_port>/api/backups`. Requests to `/api/backups` on the
public port return `404 Not Found`.

When `cogmoteGO` is behind Nginx or another reverse proxy, proxy only the public
port. Do not proxy the internal port.

## Samba Incremental Backup

The backup API copies explicitly selected files or directories from trusted local
roots to trusted mounted Samba roots. The API only accepts root IDs and relative
paths, so callers cannot select arbitrary host filesystem paths.

Configure trusted roots before creating tasks:

```nu
cogmoteGO backup roots add source projectx-data /path/to/projectx/data
cogmoteGO backup roots add samba lab-nas /mnt/samba
```

List configured roots with `cogmoteGO backup roots list`. Restart the service
after changing them.

The resulting configuration is:

```json
{
  "backup": {
    "source_roots": [{"id": "projectx-data", "path": "/path/to/projectx/data"}],
    "samba_roots": [{"id": "lab-nas", "path": "/mnt/samba"}]
  }
}
```

Create an asynchronous backup task with Nushell:

```nu
http post http://127.0.0.1:9011/api/backups {
  source: {
    root_id: "projectx-data"
    entries: ["20260713" "20260714/realdata/result.jsonl"]
  }
  destination: {
    type: "samba"
    root_id: "lab-nas"
    path: "experiments/projectx/data"
  }
}
```

Each selected entry is mapped beneath the destination path. Existing destination
entries cause the task to fail; entries are never overwritten or merged. Files
are copied with SHA-256 verification. Only one backup task may run at a time;
a concurrent request receives `409 Conflict`. Failed transfers automatically
remove `.partial-<task-id>` data. Tasks are kept only in memory, so a service
restart requires the caller to submit the upload again. The service retains only
the latest task; creating a new task replaces the previous completed task.

Poll `GET /api/backups` every 500 ms to 1 s. Before the first task is created,
this endpoint returns `404 Not Found`. The response includes `status`
(`running`, `succeeded`, `partially_succeeded`, or `failed`) and
`phase` (`scanning`, `uploading`, `verifying`, `publishing`, `completed`),
`files_total`, `files_completed`, `bytes_total`, `bytes_transferred`, and
`current_path`. Failed tasks retain the phase in which the failure occurred.
Use `bytes_transferred / bytes_total` for the upload progress bar once the
scanning phase completes.

### For more info

```sh
cogmoteGO backup --help
```
