package detective

import (
	"testing"
	"time"
)

func TestCorrelateDeploysRanksClosestHighBlastDeploy(t *testing.T) {
	now := time.Date(2026, 11, 11, 10, 0, 0, 0, time.UTC)
	alert := AlertEvent{
		AlertID:   "al_503",
		Service:   "checkout-api",
		Severity:  "SEV1",
		Message:   "5xx rate > 4.2%",
		Timestamp: now,
	}
	deploys := []DeployEvent{
		{
			DeployID:    "dep_old",
			Service:     "checkout-api",
			Version:     "v1.4.0",
			Author:      "alex",
			Timestamp:   now.Add(-45 * time.Minute),
			ChangedKeys: []string{"ui.banner"},
		},
		{
			DeployID:    "dep_root",
			Service:     "checkout-api",
			Version:     "v1.4.1",
			Author:      "soumya",
			Timestamp:   now.Add(-4 * time.Minute),
			ChangedKeys: []string{"db.pool_max", "stripe.timeout_ms", "circuit_breaker.threshold"},
		},
		{
			DeployID:    "dep_other_svc",
			Service:     "billing-worker",
			Version:     "v2.0.0",
			Author:      "sam",
			Timestamp:   now.Add(-2 * time.Minute),
			ChangedKeys: []string{"queue.concurrency"},
		},
	}

	brief := CorrelateDeploys(alert, deploys, 60*time.Minute)
	if len(brief.Suspects) != 2 {
		t.Fatalf("expected 2 checkout-api suspects, got %d", len(brief.Suspects))
	}
	if brief.Suspects[0].DeployID != "dep_root" {
		t.Fatalf("expected dep_root first, got %s", brief.Suspects[0].DeployID)
	}
	if brief.IncidentKey == "" || len(brief.TimelineLines) != 3 {
		t.Fatalf("unexpected brief output: %+v", brief)
	}
}

func TestSummarizeDriftFlagsDeletesAndHighImpact(t *testing.T) {
	res := SummarizeDrift([]ResourceChange{
		{Address: "aws_iam_role.worker", ChangeType: "update", Attribute: "max_session_duration", LiveValue: "3600", PlannedVal: "7200", HighImpact: false},
		{Address: "aws_rds_cluster.prod", ChangeType: "delete", Attribute: "deletion_protection", LiveValue: "true", PlannedVal: "false", HighImpact: true},
	})
	if res.TotalDrifted != 2 || len(res.HighImpactKeys) != 1 || res.HighImpactKeys[0] != "aws_rds_cluster.prod.deletion_protection" {
		t.Fatalf("unexpected drift summary: %+v", res)
	}
}
