import { test } from "node:test";
import assert from "node:assert/strict";
import { formatFactoryTriageCard } from "../src/factoryStatus.ts";

test("formats nominal factory triage card", () => {
  const card = formatFactoryTriageCard({
    workOrderId: "WO-SG-2026-001",
    sku: "SST-MV-100KW",
    stage: "CALIBRATION",
    recommended: "PROCEED_TO_NEXT_STAGE",
    anomalies: [],
  });

  assert.match(card, /\[FACTORY TRIAGE\] WO-SG-2026-001/);
  assert.match(card, /Action: PROCEED_TO_NEXT_STAGE/);
  assert.match(card, /All telemetry metrics nominal/);
});

test("formats critical anomaly factory triage card", () => {
  const card = formatFactoryTriageCard({
    workOrderId: "WO-SG-2026-002",
    sku: "SST-MV-100KW",
    stage: "QUARANTINED",
    recommended: "QUARANTINE_FOR_RECALIBRATION",
    anomalies: [
      {
        metricKey: "leakage_inductance_uh",
        severity: "CRITICAL",
        deviation: 4.3,
        message: "exceeded tolerance by 4.300",
      },
    ],
  });

  assert.match(card, /\[CRITICAL\] leakage_inductance_uh: dev=\+4\.300/);
  assert.match(card, /Action: QUARANTINE_FOR_RECALIBRATION/);
});
