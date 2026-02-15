# OBS 集成

[English](./README.md)

cogmoteGO 与 OBS Studio 的集成模块，用于在直播画面上叠加实验数据文本。

## 前置要求

1. 安装 OBS Studio (28.0+)
2. 启用 OBS WebSocket 服务器：
   - `工具` → `WebSocket 服务器设置`
   - 端口：`4455`
   - 启用认证并设置密码

## 配置

### CLI 配置

```bash
# 交互式配置 scene、source 和密码
cogmoteGO obs set

# 查看当前配置状态
cogmoteGO obs show

# 删除保存的密码
cogmoteGO obs delete-password
```

### Docker 环境

使用环境变量传递密码：

```bash
docker run -e OBS_PASSWORD=yourpassword ...
```

## API 端点

### 状态查询

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/obs` | 获取 OBS 状态 (版本、是否直播中) |

### 初始化

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/obs/init` | 初始化 OBS 客户端连接 |

### 直播控制

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/obs/start` | 开始直播 |
| POST | `/api/obs/stop` | 停止直播 |

### 数据叠加

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/obs/data` | 更新叠加文本数据 |

## 使用流程

```
1. 启动 OBS 并确保 WebSocket 服务器运行
2. 调用 POST /api/obs/init 初始化连接
3. 调用 POST /api/obs/start 开始直播
4. 调用 POST /api/obs/data 更新实验数据
5. 调用 POST /api/obs/stop 停止直播
```

## Init 响应示例

```json
{
  "scene_name": "Scene",
  "source_name": "overlay_text",
  "scene_fallback": false,
  "source_created": true
}
```

- `scene_fallback`: 配置的 scene 不存在时使用第一个 scene
- `source_created`: 文本源是新创建的

## 数据格式

POST `/api/obs/data` 请求体：

```json
{
  "monkey_name": "monkey_001",
  "trial_id": 42,
  "start_time": "2024-01-15 10:30:00",
  "correct_rate": 0.85
}
```

叠加文本格式：`{hostname} {monkey_name} {trial_id} {correct_rate}% {start_time}`

## 注意事项

- Scene 和 Source 不能使用相同名称
- 文本源会自动定位在画面底部 (1920x1080 假设)
- 仅支持连接本地 OBS (`localhost:4455`)
