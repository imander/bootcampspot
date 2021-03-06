package bcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/imander/bootcampspot/config"
)

var (
	assignmentResource = "/assignments"
	attendanceResource = "/attendance"
	authResource       = "/login"
	feedbackResource   = "/weeklyFeedback"
	gradeResource      = "/grades"
	sessionResource    = "/sessions"
	userResource       = "/me"

	authToken  string
	baseURL    *url.URL
	httpClient = &http.Client{}
	layoutISO  = "2006-01-02"

	// Course and Enrollment IDs are used for rest requests
	CourseID     int
	EnrollmentID int
)

func bcsURL() (url.URL, error) {
	var err error
	if baseURL != nil {
		return *baseURL, nil
	}

	baseURL, err = url.Parse(config.BCS.URL)
	return *baseURL, err
}

func getToken() error {
	ab := AuthBody{
		Email:    config.BCS.User,
		Password: config.BCS.Password,
	}
	return ab.GetToken()
}

// RestRequest stores all the needed data to send a REST call to appsec services
type RestRequest struct {
	Method string
	Path   string
	Params map[string]string
	Data   interface{}
}

type CourseBody struct {
	ID int `json:"courseId"`
}

type EnrollementBody struct {
	ID int `json:"enrollmentId"`
}

func (req *RestRequest) Send(iff interface{}) error {
	if authToken == "" && !strings.HasSuffix(req.Path, authResource) {
		err := getToken()
		if err != nil {
			return err
		}
	}

	r, err := req.build()
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("REST error: code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(iff)
	if err == nil || err.Error() == "EOF" {
		return nil
	}

	return err
}

func (req *RestRequest) build() (r *http.Request, err error) {
	u, err := bcsURL()
	if err != nil {
		return nil, err
	}
	u.Path = u.Path + req.Path

	if req.Data != nil {
		var body []byte
		body, err = json.Marshal(req.Data)
		if err != nil {
			return
		}

		r, err = http.NewRequest(req.Method, u.String(), bytes.NewBuffer(body))
		if err != nil {
			return
		}
		r.Header.Add("Content-Type", "application/json")
	} else {
		r, err = http.NewRequest(req.Method, u.String(), nil)
		if err != nil {
			return
		}
	}

	r.Header.Add("authToken", authToken)
	r.Header.Add("Accept", "application/json")
	return
}
