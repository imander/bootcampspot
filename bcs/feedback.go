package bcs

import (
	"fmt"
	"net/http"
	"strconv"
)

type FeedbackResponse struct {
	SurveyDefinition struct {
		SurveyDescription string `json:"surveyDescription"`
		Steps             []struct {
			StepNumber string `json:"stepNumber"`
			Text       string `json:"text"`
		} `json:"steps"`
	} `json:"surveyDefinition"`
	Submissions []struct {
		Username string `json:"username"`
		Date     string `json:"date"`
		Answers  []struct {
			StepNumber string `json:"stepNumber"`
			Answer     struct {
				Value string `json:"value"`
			} `json:"answer"`
		} `json:"answers"`
	} `json:"submissions"`
}

func GetFeedback() (FeedbackResponse, error) {
	data := CourseBody{ID: CourseID}
	req := RestRequest{
		Method: http.MethodPost,
		Path:   feedbackResource,
		Data:   data,
	}

	body := FeedbackResponse{}
	err := req.Send(&body)
	return body, err
}

func (fr FeedbackResponse) questionMap() map[int]string {
	m := map[int]string{}

	for _, step := range fr.SurveyDefinition.Steps {
		num, _ := strconv.Atoi(step.StepNumber)
		m[num] = step.Text
	}

	return m
}

func (fr FeedbackResponse) Print() {
	questions := fr.questionMap()
	out := fmt.Sprintf("%s\n", fr.SurveyDefinition.SurveyDescription)

	for _, s := range fr.Submissions {
		out += fmt.Sprintf("\nUser: %s\n", s.Username)
		out += fmt.Sprintf("Date: %s\n", s.Date)
		for _, a := range s.Answers {
			index, _ := strconv.Atoi(a.StepNumber)
			q := questions[index]
			if len(a.Answer.Value) == 1 {
				out += fmt.Sprintf("\t%s - %s\n", a.Answer.Value, q)
			} else {
				if a.Answer.Value != "" {
					out += fmt.Sprintf("\t%s\n\t- %s\n", q, a.Answer.Value)
				}
			}
		}
	}

	fmt.Println(out)
}
