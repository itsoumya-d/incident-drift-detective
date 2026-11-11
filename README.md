# incident-drift-detective

Go + TypeScript incident triage library that correlates production alerts against recent deployments and infrastructure drift diffs, emitting deterministic postmortem timelines and customer status-page updates.

## Modules

- **`internal/detective/correlator.go`** — Ranks candidate deploys within a lookback window by temporal proximity and changed configuration/code blast radius, producing deterministic `IncidentBrief` timelines.
- **`internal/detective/drift.go`** — Summarizes Terraform/runtime attribute drift diffs, elevating deletes and high-impact mutations into markdown digests for incident commanders.
- **`ts/src/statusUpdate.ts`** — Strict TypeScript formatter for external customer status updates.

## Verification

```bash
go test -v ./...
cd ts && npm install --no-package-lock && npm run build && npm test
```
