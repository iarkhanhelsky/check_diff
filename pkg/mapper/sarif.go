package mapper

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/owenrumney/go-sarif/sarif"
)

func SarifBytesToIssues(bytes []byte) ([]core.Issue, error) {
	return parseSarif(bytes)
}

func parseSarif(reportBytes []byte) ([]core.Issue, error) {
	report, err := sarif.FromBytes(reportBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sarif report: %v: %s", err, string(reportBytes))
	}

	var issues []core.Issue
	for _, run := range report.Runs {
		for _, result := range run.Results {
			file, line, column := extractLocation(result.Locations)

			issues = append(issues, core.Issue{
				Tag:      *result.RuleID,
				Message:  *result.Message.Text,
				File:     file,
				Line:     line,
				Column:   column,
				Severity: "error",
			})
		}
	}

	return issues, nil
}

func extractLocation(locations []*sarif.Location) (string, int, int) {
	if len(locations) == 0 {
		return "", 0, 0
	}

	location := locations[0].PhysicalLocation

	var file string
	var line int
	var column int

	if location.ArtifactLocation.URI != nil {
		file = *location.ArtifactLocation.URI
	}

	if location.Region.StartLine != nil {
		line = *location.Region.StartLine
	}

	if location.Region.EndLine != nil {
		column = *location.Region.StartColumn
	}

	return file, line, column
}
