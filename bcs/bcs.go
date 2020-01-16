package bcs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	baseURL    = &url.URL{}
	httpClient = &http.Client{}
	layoutISO  = "2006-01-02"

	// Course and Enrollment IDs are used for rest requests
	CourseID     int
	EnrollmentID int
)

func bcsURL() url.URL {
	var err error
	if baseURL != nil {
		return *baseURL
	}

	baseURL, err = url.Parse(config.BCS.URL)
	if err != nil {
		fmt.Printf("error: %s\n url: %s\n", err.Error(), config.BCS.URL)
		os.Exit(1)
	}

	return *baseURL
}

func getToken() {
	baseURL, _ = url.Parse(config.BCS.URL)

	ab := AuthBody{
		Email:    config.BCS.User,
		Password: config.BCS.Password,
	}
	ab.GetToken()
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
		getToken()
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
		return restError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(iff)
	if err == nil || err.Error() == "EOF" {
		return nil
	}

	return err
}

func (req *RestRequest) build() (r *http.Request, err error) {
	u := bcsURL()
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

func restError(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return errors.New(string(b))
}
