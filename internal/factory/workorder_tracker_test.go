package factory

import (
	"testing"
	"time"
)

func TestEvaluateWorkOrderNominal(t *testing.T) {
	now := time.Now()
	wo := WorkOrder{
		ID:         "WO-SG-2026-001",
		ProductSKU: "SST-MV-100KW",
		Stage:      StageCalibration,
		CreatedAt:  now.Add(-2 * time.Hour),
		UpdatedAt:  now,
		Samples: []TelemetrySample{
			{
				Timestamp:   now.Add(-30 * time.Minute),
				MetricKey:   "thermal_junction_degC",
				ObservedVal: 54.2,
				TargetVal:   55.0,
				Tolerance:   2.0,
			},
			{
				Timestamp:   now.Add(-15 * time.Minute),
				MetricKey:   "voltage_output_rms_v",
				ObservedVal: 400.1,
				TargetVal:   400.0,
				Tolerance:   1.0,
			},
		},
	}

	triage := EvaluateWorkOrder(wo)

	if triage.Recommended != "PROCEED_TO_NEXT_STAGE" {
		t.Fatalf("expected PROCEED_TO_NEXT_STAGE, got %s", triage.Recommended)
	}
	if len(triage.Anomalies) != 0 {
		t.Fatalf("expected 0 anomalies, got %d", len(triage.Anomalies))
	}
	if triage.CurrentStage != StageCalibration {
		t.Fatalf("expected stage StageCalibration, got %s", triage.CurrentStage)
	}
	if len(triage.Digest) != 16 {
		t.Fatalf("expected 16-hex digest, got %s", triage.Digest)
	}
}

func TestEvaluateWorkOrderCriticalToleranceBreach(t *testing.T) {
	now := time.Now()
	wo := WorkOrder{
		ID:         "WO-SG-2026-002",
		ProductSKU: "SST-MV-100KW",
		Stage:      StageCalibration,
		CreatedAt:  now.Add(-1 * time.Hour),
		UpdatedAt:  now,
		Samples: []TelemetrySample{
			{
				Timestamp:   now.Add(-10 * time.Minute),
				MetricKey:   "leakage_inductance_uh",
				ObservedVal: 15.8,
				TargetVal:   10.0,
				Tolerance:   1.5, // deviation 5.8 > 2*1.5=3.0 -> CRITICAL
			},
			{
				Timestamp:   now.Add(-5 * time.Minute),
				MetricKey:   "carrier_frequency_khz",
				ObservedVal: 102.5,
				TargetVal:   100.0,
				Tolerance:   2.0, // deviation 0.5 > tol -> WARNING
			},
		},
	}

	triage := EvaluateWorkOrder(wo)

	if triage.Recommended != "QUARANTINE_FOR_RECALIBRATION" {
		t.Fatalf("expected QUARANTINE_FOR_RECALIBRATION, got %s", triage.Recommended)
	}
	if triage.CurrentStage != StageQuarantined {
		t.Fatalf("expected stage StageQuarantined, got %s", triage.CurrentStage)
	}
	if len(triage.Anomalies) != 2 {
		t.Fatalf("expected 2 anomalies, got %d", len(triage.Anomalies))
	}
	if triage.Anomalies[0].Severity != "CRITICAL" {
		t.Fatalf("expected first anomaly to be CRITICAL, got %s", triage.Anomalies[0].Severity)
	}
}
