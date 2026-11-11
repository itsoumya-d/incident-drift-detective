package detective

import (
	"fmt"
	"sort"
	"strings"
)

type ResourceChange struct {
	Address     string
	ChangeType  string // "create" | "update" | "delete"
	Attribute   string
	LiveValue   string
	PlannedVal  string
	HighImpact  bool
}

type DriftSummary struct {
	TotalDrifted   int
	HighImpactKeys []string
	DigestMarkdown string
}

// SummarizeDrift builds a concise customer/responder-safe summary from Terraform/runtime drift diffs.
func SummarizeDrift(diffs []ResourceChange) DriftSummary {
	var highImpact []string
	var lines []string

	sorted := append([]ResourceChange(nil), diffs...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Address < sorted[j].Address
	})

	for _, d := range sorted {
		tag := "LOW"
		if d.HighImpact || strings.EqualFold(d.ChangeType, "delete") {
			tag = "HIGH"
			highImpact = append(highImpact, fmt.Sprintf("%s.%s", d.Address, d.Attribute))
		}
		lines = append(lines, fmt.Sprintf("- [%s] `%s.%s` (%s): `%s` -> `%s`", tag, d.Address, d.Attribute, d.ChangeType, d.LiveValue, d.PlannedVal))
	}

	return DriftSummary{
		TotalDrifted:   len(sorted),
		HighImpactKeys: highImpact,
		DigestMarkdown: strings.Join(lines, "\n"),
	}
}
