# OBS Integration

[中文文档](./README.zh_cn.md)

Integration module between cogmoteGO and OBS Studio for overlaying experiment data text on live streams.

## Prerequisites

1. Install OBS Studio (28.0+)
2. Enable OBS WebSocket server:
   - `Tools` → `WebSocket Server Settings`
   - Port: `4455`
   - Enable authentication and set a password

## Configuration

### CLI Configuration

```bash
# Interactive configuration for scene, source, and password
cogmoteGO obs set

# View current configuration status
cogmoteGO obs show

# Delete saved password
cogmoteGO obs delete-password
```

### Docker Environment

Pass password via environment variable:

```bash
docker run -e OBS_PASSWORD=yourpassword ...
```

## API Endpoints

### Status

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/obs` | Get OBS status (version, streaming state) |

### Initialization

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/obs/init` | Initialize OBS client connection |

### Streaming Control

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/obs/start` | Start streaming |
| POST | `/api/obs/stop` | Stop streaming |

### Data Overlay

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/obs/data` | Update overlay text data |

## Usage Flow

```
1. Start OBS and ensure WebSocket server is running
2. Call POST /api/obs/init to initialize connection
3. Call POST /api/obs/start to start streaming
4. Call POST /api/obs/data to update experiment data
5. Call POST /api/obs/stop to stop streaming
```

## Init Response Example

```json
{
  "scene_name": "Scene",
  "source_name": "overlay_text",
  "scene_fallback": false,
  "source_created": true
}
```

- `scene_fallback`: First scene used when configured scene doesn't exist
- `source_created`: Text source was newly created

## Data Format

POST `/api/obs/data` request body:

```json
{
  "monkey_name": "monkey_001",
  "trial_id": 42,
  "start_time": "2024-01-15 10:30:00",
  "correct_rate": 0.85
}
```

Overlay text format: `{hostname} {monkey_name} {trial_id} {correct_rate}% {start_time}`

## Notes

- Scene and Source cannot share the same name
- Text source is automatically positioned at the bottom of the screen (1920x1080 assumed)
- Only supports local OBS connection (`localhost:4455`)
