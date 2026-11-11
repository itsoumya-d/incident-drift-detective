package detective

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

type DeployEvent struct {
	DeployID    string
	Service     string
	Version     string
	Author      string
	Timestamp   time.Time
	ChangedKeys []string
}

type AlertEvent struct {
	AlertID   string
	Service   string
	Severity  string
	Message   string
	Timestamp time.Time
}

type SuspectChange struct {
	DeployID   string
	Service    string
	Version    string
	Author     string
	Score      int
	DeltaSecs  float64
	Reason     string
}

type IncidentBrief struct {
	IncidentKey   string
	Service       string
	Severity      string
	Suspects      []SuspectChange
	TimelineLines []string
}

// CorrelateDeploys ranks candidate deployments preceding an alert within lookbackWindow.
func CorrelateDeploys(alert AlertEvent, deploys []DeployEvent, lookbackWindow time.Duration) IncidentBrief {
	var suspects []SuspectChange
	for _, dep := range deploys {
		if dep.Service != alert.Service {
			continue
		}
		delta := alert.Timestamp.Sub(dep.Timestamp)
		if delta < 0 || delta > lookbackWindow {
			continue
		}
		proximityWeight := int((lookbackWindow - delta).Minutes()) + 10
		blastWeight := len(dep.ChangedKeys) * 15
		score := proximityWeight + blastWeight
		reason := fmt.Sprintf("deployed %s %.0fs prior (%d config/code keys)", dep.Version, delta.Seconds(), len(dep.ChangedKeys))
		suspects = append(suspects, SuspectChange{
			DeployID:  dep.DeployID,
			Service:   dep.Service,
			Version:   dep.Version,
			Author:    dep.Author,
			Score:     score,
			DeltaSecs: delta.Seconds(),
			Reason:    reason,
		})
	}

	sort.Slice(suspects, func(i, j int) bool {
		if suspects[i].Score != suspects[j].Score {
			return suspects[i].Score > suspects[j].Score
		}
		return suspects[i].DeployID < suspects[j].DeployID
	})

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", alert.Service, alert.Severity, alert.AlertID)))
	timeline := []string{
		fmt.Sprintf("%s [ALERT:%s] %s — %s", alert.Timestamp.UTC().Format(time.RFC3339), alert.Severity, alert.Service, alert.Message),
	}
	for _, s := range suspects {
		timeline = append(timeline, fmt.Sprintf("[SUSPECT score=%d] %s by @%s — %s", s.Score, s.DeployID, s.Author, s.Reason))
	}

	return IncidentBrief{
		IncidentKey:   "inc_" + hex.EncodeToString(sum[:6]),
		Service:       alert.Service,
		Severity:      alert.Severity,
		Suspects:      suspects,
		TimelineLines: timeline,
	}
}
