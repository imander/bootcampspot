package bcs

import (
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Grades []struct {
	AssignmentTitle string `json:"assignmentTitle"`
	StudentName     string `json:"studentName"`
	Submitted       bool   `json:"submitted"`
	Grade           string `json:"grade"`
}

func GetGrades() (Grades, error) {
	data := CourseBody{ID: CourseID}
	req := RestRequest{
		Method: http.MethodPost,
		Path:   gradeResource,
		Data:   data,
	}

	body := Grades{}
	err := req.Send(&body)
	return body, err
}

type GradesMetric struct {
	Submitted    int
	NotSubmitted int
	Incomplete   int
}

type GradesMetrics map[string]GradesMetric

func (grades Grades) Metrics() (GradesMetrics, error) {
	metrics := GradesMetrics{}
	assignments, err := GetAssignments()
	if err != nil {
		return metrics, err
	}
	am := assignments.Academic().nameMap()

	for _, g := range grades {
		if !am[g.AssignmentTitle] {
			continue
		}
		name := strings.Title(strings.ToLower(g.StudentName))
		m := metrics[name]
		if g.Submitted {
			m.Submitted += 1
		} else {
			m.NotSubmitted += 1
		}
		if g.Grade == "Incomplete" {
			m.Incomplete += 1
		}
		metrics[name] = m
	}

	return metrics.sort(), nil
}

func (gm GradesMetrics) sort() GradesMetrics {
	keys := make([]string, 0, len(gm))
	for k := range gm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	g := GradesMetrics{}
	for _, k := range keys {
		g[k] = gm[k]
	}

	return g
}

func (gm GradesMetrics) Print() {
	header := []string{"Student", "Submitted", "Not Submitted", "Incomplete"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	keys := make([]string, 0, len(gm))
	for k := range gm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, student := range keys {
		m := gm[student]
		if m.Submitted == 0 {
			continue
		}
		row := []string{student, strconv.Itoa(m.Submitted), strconv.Itoa(m.NotSubmitted), strconv.Itoa(m.Incomplete)}
		table.Append(row)
	}

	table.Render()
}
