package factory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// WorkOrderStage represents the lifecycle progression of an advanced power hardware unit.
type WorkOrderStage string

const (
	StageBOMVerified  WorkOrderStage = "BOM_VERIFIED"
	StageSubAssembly  WorkOrderStage = "SUB_ASSEMBLY"
	StageCalibration  WorkOrderStage = "CALIBRATION"
	StageQualityAudit WorkOrderStage = "QUALITY_AUDIT"
	StageDispatched   WorkOrderStage = "DISPATCHED"
	StageQuarantined  WorkOrderStage = "QUARANTINED"
)

// TelemetrySample captures high-rate sensor readings from manufacturing/testing rigs.
type TelemetrySample struct {
	Timestamp   time.Time
	MetricKey   string
	ObservedVal float64
	TargetVal   float64
	Tolerance   float64
}

// WorkOrder captures an active manufacturing operational batch.
type WorkOrder struct {
	ID         string
	ProductSKU string
	Stage      WorkOrderStage
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Samples    []TelemetrySample
}

// AnomalyReport highlights tolerance breaches requiring agentic or operator triage.
type AnomalyReport struct {
	WorkOrderID string
	MetricKey   string
	Severity    string
	Deviation   float64
	Message     string
}

// WorkOrderTriage summarizes status, anomalies, and agentic operational guidance.
type WorkOrderTriage struct {
	Digest       string
	WorkOrderID  string
	CurrentStage WorkOrderStage
	Anomalies    []AnomalyReport
	Recommended  string
}

// EvaluateWorkOrder analyzes sensor telemetry against tolerances and returns an operational triage brief.
func EvaluateWorkOrder(wo WorkOrder) WorkOrderTriage {
	var anomalies []AnomalyReport

	for _, s := range wo.Samples {
		delta := s.ObservedVal - s.TargetVal
		absDelta := delta
		if absDelta < 0 {
			absDelta = -absDelta
		}

		if absDelta > s.Tolerance {
			severity := "WARNING"
			if absDelta > s.Tolerance*2.0 {
				severity = "CRITICAL"
			}

			anomalies = append(anomalies, AnomalyReport{
				WorkOrderID: wo.ID,
				MetricKey:   s.MetricKey,
				Severity:    severity,
				Deviation:   absDelta - s.Tolerance,
				Message: fmt.Sprintf("%s exceeded tolerance by %.3f (observed=%.3f, target=%.3f, tol=%.3f)",
					s.MetricKey, absDelta-s.Tolerance, s.ObservedVal, s.TargetVal, s.Tolerance),
			})
		}
	}

	sort.Slice(anomalies, func(i, j int) bool {
		if anomalies[i].Severity != anomalies[j].Severity {
			return anomalies[i].Severity == "CRITICAL"
		}
		return anomalies[i].Deviation > anomalies[j].Deviation
	})

	recommended := "PROCEED_TO_NEXT_STAGE"
	currentStage := wo.Stage

	if len(anomalies) > 0 {
		hasCritical := false
		for _, a := range anomalies {
			if a.Severity == "CRITICAL" {
				hasCritical = true
				break
			}
		}

		if hasCritical {
			recommended = "QUARANTINE_FOR_RECALIBRATION"
			currentStage = StageQuarantined
		} else {
			recommended = "OPERATOR_INSPECTION_REQUIRED"
		}
	}

	hashInput := fmt.Sprintf("%s:%s:%d:%s", wo.ID, currentStage, len(anomalies), recommended)
	hash := sha256.Sum256([]byte(hashInput))
	digest := hex.EncodeToString(hash[:8])

	return WorkOrderTriage{
		Digest:       digest,
		WorkOrderID:  wo.ID,
		CurrentStage: currentStage,
		Anomalies:    anomalies,
		Recommended:  recommended,
	}
}
