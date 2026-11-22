export interface FactoryAnomaly {
  metricKey: string;
  severity: "WARNING" | "CRITICAL";
  deviation: number;
  message: string;
}

export interface FactoryTriageInput {
  workOrderId: string;
  sku: string;
  stage: string;
  recommended: string;
  anomalies: FactoryAnomaly[];
}

export function formatFactoryTriageCard(input: FactoryTriageInput): string {
  const header = `[FACTORY TRIAGE] ${input.workOrderId} (${input.sku}) - Stage: ${input.stage}`;
  const action = `Action: ${input.recommended}`;

  if (input.anomalies.length === 0) {
    return `${header}\n${action}\nAll telemetry metrics nominal within tolerance bounds.`;
  }

  const lines = input.anomalies.map(
    (a) => ` - [${a.severity}] ${a.metricKey}: dev=+${a.deviation.toFixed(3)} (${a.message})`
  );

  return `${header}\n${action}\nIdentified Anomalies (${input.anomalies.length}):\n${lines.join("\n")}`;
}
