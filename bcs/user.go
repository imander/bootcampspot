package bcs

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/imander/bootcampspot/util"
	"github.com/olekukonko/tablewriter"
)

type User struct {
	UserAccount UserAccount `json:"userAccount"`
	Enrollments Enrollments `json:"enrollments"`
}

type BcsExtends struct {
	ID             int  `json:"id"`
	IsAdmin        bool `json:"isAdmin"`
	JobTrackActive bool `json:"jobTrackActive"`
}

type UserAccount struct {
	ID                    int        `json:"id"`
	UserName              string     `json:"userName"`
	FirstName             string     `json:"firstName"`
	LastName              string     `json:"lastName"`
	Email                 string     `json:"email"`
	GithubUserName        string     `json:"githubUserName"`
	PasswordResetRequired bool       `json:"passwordResetRequired"`
	Active                bool       `json:"active"`
	BcsExtends            BcsExtends `json:"bcsExtends"`
}

type University struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logoUrl"`
}

type ProgramType struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type Program struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	University  University  `json:"university"`
	ProgramType ProgramType `json:"programType"`
}

type Cohort struct {
	ID        int       `json:"id"`
	ProgramID int       `json:"programId"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Program   Program   `json:"program"`
}

type GraduationRequirements struct {
	MaxAbsence                  int `json:"maxAbsence"`
	MaxRemoteAttendance         int `json:"maxRemoteAttendance"`
	MaxMissedGeneralAssignment  int `json:"maxMissedGeneralAssignment"`
	MaxMissedRequiredAssignment int `json:"maxMissedRequiredAssignment"`
}

type Course struct {
	ID                     int                    `json:"id"`
	CohortID               int                    `json:"cohortId"`
	Name                   string                 `json:"name"`
	Code                   string                 `json:"code"`
	StartDate              time.Time              `json:"startDate"`
	EndDate                time.Time              `json:"endDate"`
	Location               string                 `json:"location"`
	EndSurveyEnabled       bool                   `json:"endSurveyEnabled"`
	Cohort                 Cohort                 `json:"cohort"`
	GraduationRequirements GraduationRequirements `json:"graduationRequirements"`
}

type CourseRole struct {
	ID             int    `json:"id"`
	CourseRoleCode string `json:"courseRoleCode"`
	Name           string `json:"name"`
}

type MidCourseSurveySchedule struct {
	ID       int    `json:"id"`
	CourseID int    `json:"courseId"`
	Survey   string `json:"survey"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type EndCourseSurveySchedule struct {
	ID       int    `json:"id"`
	CourseID int    `json:"courseId"`
	Survey   string `json:"survey"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type Enrollment struct {
	ID                                  int                     `json:"id"`
	CourseID                            int                     `json:"courseId"`
	UserAccountID                       int                     `json:"userAccountId"`
	CourseRoleID                        int                     `json:"courseRoleId"`
	Active                              bool                    `json:"active"`
	Course                              Course                  `json:"course"`
	CourseRole                          CourseRole              `json:"courseRole"`
	MidCourseSurveySchedule             MidCourseSurveySchedule `json:"midCourseSurveySchedule"`
	EndCourseSurveySchedule             EndCourseSurveySchedule `json:"endCourseSurveySchedule"`
	RemoteAttendanceCount               int                     `json:"remoteAttendanceCount"`
	PendingRemoteAttendanceRequestCount int                     `json:"pendingRemoteAttendanceRequestCount"`
}

type Enrollments []Enrollment

func GetUser() User {
	req := RestRequest{
		Method: http.MethodGet,
		Path:   userResource,
	}

	body := User{}
	err := req.Send(&body)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	return body
}

func (u User) PrintEnrollments() {
	header := []string{"Enrollment ID", "Course ID", "Course Name", "Role", "Start Date", "End Date"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, enr := range u.Enrollments {
		enrID := strconv.Itoa(enr.ID)
		courseID := strconv.Itoa(enr.CourseID)
		start := enr.Course.StartDate.Format(layoutISO)
		end := enr.Course.EndDate.Format(layoutISO)
		row := []string{enrID, courseID, enr.Course.Name, enr.CourseRole.Name, start, end}
		table.Append(row)
	}

	table.Render()
}

func (u User) ChooseEnrollment() Enrollment {
	enrs := u.Enrollments.Filter()
	if len(enrs) == 0 {
		fmt.Println("no enrolments found")
		os.Exit(1)
	}

	if len(enrs) == 1 {
		e := u.Enrollments[0]
		CourseID = e.CourseID
		EnrollmentID = e.ID
		return e
	}

	fmt.Println()
	for i, enr := range u.Enrollments {
		fmt.Printf("[%d] %s\n", i, enr.Course.Name)
	}

	i := util.ReadInt("\nChoose course")

	e := u.Enrollments[i]
	CourseID = e.CourseID
	EnrollmentID = e.ID

	return e
}

func (e Enrollments) Filter() Enrollments {
	if CourseID == -1 && EnrollmentID == -1 {
		return e
	}

	enrollments := Enrollments{}
	for _, enr := range e {
		if enr.CourseID == CourseID || enr.ID == EnrollmentID {
			enrollments = append(enrollments, enr)
		}
	}

	if len(enrollments) == 0 {
		return e
	}

	return enrollments
}
