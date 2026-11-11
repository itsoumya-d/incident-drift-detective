export interface IncidentStateInput {
  incidentKey: string;
  service: string;
  severity: "SEV1" | "SEV2" | "SEV3";
  status: "investigating" | "identified" | "mitigated" | "resolved";
  topSuspectDeploy?: string;
}

export function renderCustomerStatusUpdate(input: IncidentStateInput): string {
  const lead =
    input.severity === "SEV1" ? "Critical service degradation" : "Service disruption";
  const action =
    input.status === "resolved"
      ? "Service health has returned to normal baseline."
      : input.topSuspectDeploy
        ? `Engineering is rolling back / isolating change ${input.topSuspectDeploy}.`
        : "Engineering is actively investigating root cause.";
  return `[${input.incidentKey}] ${lead} on ${input.service} (${input.status}): ${action}`;
}
