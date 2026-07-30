# cogmoteGO Bruno collection

This OpenCollection YAML collection documents all HTTP routes registered by the
current cogmoteGO source tree.

## Open in Bruno

Use Bruno 3.0 or newer and open this `bruno` directory as a collection. Select
the `local` environment before sending requests.

The default environment targets:

- Public API: `http://localhost:9012/api`
- Loopback-only backup API: `http://127.0.0.1:9011/api`

No HTTP authentication is configured by the current server.

## Tags

- `smoke`: read-only requests suitable for a quick service check.
- `regression`: bounded requests with basic response tests.
- `manual`: requires state, local software, a worker, a file, or user review.
- `sse`: long-running Server-Sent Events request; cancel it manually.
- `internal`: available only through the loopback internal listener.
- `destructive`: deletes records or state and must be run intentionally.

## Suggested workflows

Broadcast:

1. Create broadcast endpoint.
2. Publish named broadcast data.
3. Get latest named broadcast data or subscribe through SSE.
4. Delete named broadcast endpoint.

Experiment:

1. Register local experiment (stores `experimentId` as a runtime variable).
2. Get or update the experiment.
3. Start and stop it only after reviewing its `execs`.
4. Delete it when finished.

OBS:

1. Start OBS application.
2. Initialize OBS client.
3. Start streaming.
4. Update overlay data.
5. Stop streaming and then stop OBS.

Backup:

1. Configure source and Samba roots in cogmoteGO.
2. Fill the `backup*` environment variables.
3. Create backup task on `internalBaseUrl`.
4. Poll the current backup task.

## Known server behavior

`POST /exps/:id/start/:nickname` currently calls `StartExperimentHandler` rather
than `StartSpecificExperimentHandler`, so the nickname is ignored by the current
server wiring.

The repository's SSE end-to-end test still references legacy routes and is not
used as the source of truth for this collection.

`endpoints.json` is the source inventory used to deterministically regenerate
the 51 request files.
