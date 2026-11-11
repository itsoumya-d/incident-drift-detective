import assert from "node:assert/strict";
import test from "node:test";
import { renderCustomerStatusUpdate } from "../src/statusUpdate.ts";

test("renders customer-safe status updates with suspect deploy attribution", () => {
  const msg = renderCustomerStatusUpdate({
    incidentKey: "inc_a1b2c3",
    service: "checkout-api",
    severity: "SEV1",
    status: "identified",
    topSuspectDeploy: "v1.4.1",
  });
  assert.match(msg, /checkout-api \(identified\)/);
  assert.match(msg, /v1\.4\.1/);
});
