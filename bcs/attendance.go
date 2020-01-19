package bcs

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Attendance struct {
	SessionName string `json:"sessionName"`
	StudentName string `json:"studentName"`
	Pending     bool   `json:"pending"`
	Present     bool   `json:"present"`
	Remote      bool   `json:"remote"`
	Excused     bool   `json:"excused"`
}

type Attendances []Attendance

type AttendanceMetric struct {
	Present int
	Absent  int
	Remote  int
}

type AttendanceMetrics map[string]AttendanceMetric

func GetAttendance() (Attendances, error) {
	data := CourseBody{ID: CourseID}
	req := RestRequest{
		Method: http.MethodPost,
		Path:   attendanceResource,
		Data:   data,
	}

	body := Attendances{}
	err := req.Send(&body)
	return body, err
}

func (att Attendances) Metrics() (AttendanceMetrics, error) {
	metrics := AttendanceMetrics{}
	sessions, err := GetSessions()
	if err != nil {
		return metrics, err
	}
	sm := sessions.Academic().nameMap()

	for _, a := range att {
		if !sm[a.SessionName] {
			continue
		}
		name := strings.Title(strings.ToLower(a.StudentName))
		m := metrics[name]
		if a.Present {
			m.Present += 1
		} else {
			m.Absent += 1
		}
		if a.Remote {
			m.Remote += 1
		}
		metrics[name] = m
	}

	return metrics, nil
}

func (am AttendanceMetrics) sort() {
	keys := make([]string, 0, len(am))
	for k := range am {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, am[k])
	}
}

func (am AttendanceMetrics) Print() {
	header := []string{"Student", "Present", "Absent", "Remote", "Total Missed"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	keys := make([]string, 0, len(am))
	for k := range am {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, student := range keys {
		m := am[student]
		if m.Present < 10 {
			continue
		}
		total := m.Absent + m.Remote
		row := []string{
			student,
			strconv.Itoa(m.Present),
			strconv.Itoa(m.Absent),
			strconv.Itoa(m.Remote),
			strconv.Itoa(total),
		}
		table.Append(row)
	}

	table.Render()
}
