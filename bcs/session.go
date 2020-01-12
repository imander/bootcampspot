package bcs

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type SessionResponse struct {
	CurrentWeekSessions []CurrentWeekSessions `json:"currentWeekSessions"`
	CalendarSessions    []CalendarSessions    `json:"calendarSessions"`
}

type Session struct {
	ID               int       `json:"id"`
	CourseID         int       `json:"courseId"`
	ContextID        int       `json:"contextId"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"shortDescription"`
	LongDescription  string    `json:"longDescription"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	Chapter          string    `json:"chapter"`
}

type Context struct {
	ID          int    `json:"id"`
	ContextCode string `json:"contextCode"`
	Name        string `json:"name"`
}

type EventType struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CurrentWeekSessions struct {
	Session   Session   `json:"session"`
	Context   Context   `json:"context"`
	EventType EventType `json:"eventType"`
}

type CalendarSessions struct {
	Session   Session   `json:"session"`
	Context   Context   `json:"context"`
	EventType EventType `json:"eventType"`
}

type Sessions []Session

func GetSessions() SessionResponse {
	data := EnrollementBody{ID: EnrollmentID}
	req := RestRequest{
		Method: http.MethodPost,
		Path:   sessionResource,
		Data:   data,
	}

	body := SessionResponse{}
	err := req.Send(&body)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	return body
}

func (sr SessionResponse) Academic() Sessions {
	sess := Sessions{}
	t := time.Now()
	for _, cs := range sr.CalendarSessions {
		if cs.Session.CourseID == CourseID && cs.Session.ContextID == 1 {
			if cs.Session.StartTime.Before(t) {
				sess = append(sess, cs.Session)
			}
		}
	}

	return sess
}

func (sessions Sessions) nameMap() map[string]bool {
	m := map[string]bool{}
	for _, session := range sessions {
		m[session.Name] = true
	}

	return m
}
