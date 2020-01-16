package bcs

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type AssignmentResponse struct {
	CalendarAssignments []Assignment `json:"calendarAssignments"`
}

type AssignmentHeader struct {
	ID           int    `json:"id"`
	AssignmentID int    `json:"assignmentId"`
	Header       string `json:"header"`
}

type AcademicAssignment struct {
	ID      int  `json:"id"`
	Prework bool `json:"prework"`
}

type Assignment struct {
	ID                    int       `json:"id"`
	CourseID              int       `json:"courseId"`
	ContextID             int       `json:"contextId"`
	AssignmentDate        time.Time `json:"assignmentDate"`
	DueDate               time.Time `json:"dueDate"`
	EffectiveDueDate      time.Time `json:"effectiveDueDate"`
	Title                 string    `json:"title"`
	Required              bool      `json:"required"`
	RequiredForGraduation bool      `json:"requiredForGraduation"`
	Context               Context   `json:"context"`
}

func GetAssignments() AssignmentResponse {
	data := EnrollementBody{ID: EnrollmentID}
	req := RestRequest{
		Method: http.MethodPost,
		Path:   assignmentResource,
		Data:   data,
	}

	body := AssignmentResponse{}
	err := req.Send(&body)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	return body
}

type Assignments []Assignment

func (ar AssignmentResponse) Academic() Assignments {
	assign := Assignments{}
	t := time.Now()
	for _, a := range ar.CalendarAssignments {
		if a.CourseID == CourseID && a.ContextID == 1 {
			if a.EffectiveDueDate.Before(t) {
				assign = append(assign, a)
			}
		}
	}

	return assign
}

func (assignments Assignments) nameMap() map[string]bool {
	m := map[string]bool{}
	for _, assign := range assignments {
		if assign.Required || assign.RequiredForGraduation {
			m[assign.Title] = true
		}
	}

	return m
}
