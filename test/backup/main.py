# /// script
# requires-python = ">=3.13"
# dependencies = []
# ///

import csv
import json
import os
import random
import sys
import time
import urllib.error
import urllib.request
import uuid
from datetime import UTC, datetime
from pathlib import Path


API_URL = "http://localhost:9012"
SOURCE_ROOT_ID = "backup-test-source"
SAMBA_ROOT_ID = "backup-test-samba"
DESTINATION_PATH = "manual-tests"
POLL_INTERVAL_SECONDS = 0.5
LARGE_FILE_COUNT = 5
LARGE_FILE_SIZE = 1_000_000_000
WRITE_CHUNK_SIZE = 8 * 1024 * 1024

SCRIPT_DIRECTORY = Path(__file__).parent
GENERATED_DIRECTORY = SCRIPT_DIRECTORY / "generated"


def write_jsonl(path: Path) -> None:
    with path.open("w", encoding="utf-8") as file:
        for trial in range(100):
            record = {
                "trial": trial,
                "timestamp": datetime.now(UTC).isoformat(),
                "response_time_ms": round(random.uniform(100, 900), 2),
                "correct": random.choice([True, False]),
            }
            file.write(json.dumps(record) + "\n")


def write_csv(path: Path) -> None:
    with path.open("w", newline="", encoding="utf-8") as file:
        writer = csv.writer(file)
        writer.writerow(["sample", "signal_a", "signal_b"])
        for sample in range(200):
            writer.writerow([sample, random.random(), random.random()])


def write_large_files(directory: Path) -> None:
    for index in range(1, LARGE_FILE_COUNT + 1):
        path = directory / f"payload-{index:02d}.bin"
        print(f"creating {path.name} ({LARGE_FILE_SIZE:,} bytes)")
        remaining = LARGE_FILE_SIZE
        with path.open("wb") as file:
            while remaining > 0:
                chunk_size = min(WRITE_CHUNK_SIZE, remaining)
                file.write(os.urandom(chunk_size))
                remaining -= chunk_size


def create_data() -> Path:
    run_name = datetime.now(UTC).strftime("%Y%m%d-%H%M%S") + "-" + uuid.uuid4().hex[:8]
    directory = GENERATED_DIRECTORY / run_name
    realdata = directory / "realdata"
    empty_directory = directory / "empty"

    realdata.mkdir(parents=True)
    empty_directory.mkdir()
    write_jsonl(realdata / "trials.jsonl")
    write_csv(realdata / "measurements.csv")
    (directory / "metadata.json").write_text(
        json.dumps(
            {
                "run_id": run_name,
                "created_at": datetime.now(UTC).isoformat(),
                "generator": "test/backup/main.py",
            },
            indent=2,
        )
        + "\n",
        encoding="utf-8",
    )
    (directory / "notes.txt").write_text("Manual backup test data.\n", encoding="utf-8")
    write_large_files(realdata)
    return directory


def request(method: str, path: str, payload: dict | None = None) -> dict:
    data = None if payload is None else json.dumps(payload).encode("utf-8")
    headers = {} if data is None else {"Content-Type": "application/json"}
    http_request = urllib.request.Request(f"{API_URL}{path}", data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(http_request, timeout=10) as response:
            return json.load(response)
    except urllib.error.HTTPError as error:
        detail = error.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"{method} {path} failed with HTTP {error.code}: {detail}") from error
    except urllib.error.URLError as error:
        raise RuntimeError(f"cannot reach cogmoteGO at {API_URL}: {error.reason}") from error


def print_progress(task: dict) -> None:
    phase = task["phase"]
    total = task["bytes_total"]
    transferred = task["bytes_transferred"]
    progress = 0 if total == 0 else transferred / total * 100
    current_path = task.get("current_path", "")
    print(
        f"{phase:<10} {progress:6.2f}% "
        f"files={task['files_completed']}/{task['files_total']} "
        f"{current_path}"
    )


def upload(relative_path: str) -> int:
    task = request(
        "POST",
        "/api/backups",
        {
            "source": {
                "root_id": SOURCE_ROOT_ID,
                "entries": [relative_path],
            },
            "destination": {
                "type": "samba",
                "root_id": SAMBA_ROOT_ID,
                "path": DESTINATION_PATH,
            },
        },
    )
    print(f"created backup task {task['id']}")

    while task["status"] == "running":
        print_progress(task)
        time.sleep(POLL_INTERVAL_SECONDS)
        task = request("GET", "/api/backups")

    print_progress(task)
    if task["status"] == "succeeded":
        print("backup succeeded")
        return 0

    print(f"backup {task['status']}: {task.get('error', '')}", file=sys.stderr)
    return 1


def main() -> int:
    directory = create_data()
    relative_path = directory.relative_to(SCRIPT_DIRECTORY).as_posix()
    print(f"generated {directory}")
    return upload(relative_path)


if __name__ == "__main__":
    raise SystemExit(main())
