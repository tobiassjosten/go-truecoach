package truecoach

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Metrics map[int]*Metric

type Metric struct {
	ID          int
	Name        string
	Assessments map[int]*Assessment
}

type Assessment struct {
	ID      int
	Name    string
	Unit    string
	Samples []Sample
}

type Sample struct {
	ID    int
	Value float64
	Date  time.Time
}

func (ms *Metrics) Metric(name string) (*Metric, bool) {
	pattern := regexp.MustCompile(fmt.Sprintf(
		`^%s(?: \((?:.*)\))?$`,
		strings.ToLower(name),
	))

	for _, metric := range *ms {
		if pattern.MatchString(strings.ToLower(metric.Name)) {
			return metric, true
		}
	}

	return &Metric{}, false
}

func (m *Metric) Assessment(name string) (*Assessment, bool) {
	pattern := regexp.MustCompile(fmt.Sprintf(
		`^%s(?: \((?:.*)\))?$`,
		strings.ToLower(name),
	))

	for _, assessment := range m.Assessments {
		if pattern.MatchString(strings.ToLower(assessment.Name)) {
			return assessment, true
		}
	}

	return &Assessment{}, false
}

func (a *Assessment) ByWeek() map[string][]Sample {
	return map[string][]Sample{
		"202105": {},
	}
}

type MetricsResponse struct {
	Assessments []struct {
		ID        int    `json:"id"`
		GroupID   int    `json:"assessment_group_id"`
		Name      string `json:"name"`
		Unit      string `json:"units"`
		CreatedBy string `creator:"created_by"` // "Trainer"
	} `json:"assessments"`
	AssessmentItems []struct {
		ID           int          `json:"id"`
		AssessmentID int          `json:"assessment_id"`
		Value        paddedString `json:"value"`
		Date         time.Time    `json:"date"`
	} `json:"assessment_items"`
	AssessmentGroups []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"assessment_groups"`
}

type paddedString string

func (s *paddedString) UnmarshalJSON(data []byte) error {
	*s = paddedString(bytes.TrimSpace(bytes.Trim(data, `"`)))
	return nil
}

func (resp *MetricsResponse) Metrics() (Metrics, error) {
	metrics := Metrics{}

	for _, group := range resp.AssessmentGroups {
		metrics[group.ID] = &Metric{
			ID:          group.ID,
			Name:        group.Name,
			Assessments: map[int]*Assessment{},
		}
	}

	assessmentGroups := map[int]int{}

	for _, assessment := range resp.Assessments {
		if _, ok := metrics[assessment.GroupID]; !ok {
			return metrics, fmt.Errorf("missing assessment group '%d'", assessment.GroupID)
		}

		metrics[assessment.GroupID].Assessments[assessment.ID] = &Assessment{
			ID:      assessment.ID,
			Name:    assessment.Name,
			Unit:    assessment.Unit,
			Samples: []Sample{},
		}

		assessmentGroups[assessment.ID] = assessment.GroupID
	}

	for _, item := range resp.AssessmentItems {
		if _, ok := assessmentGroups[item.AssessmentID]; !ok {
			return metrics, fmt.Errorf("missing assessment group map '%d'", item.AssessmentID)
		}

		groupID := assessmentGroups[item.AssessmentID]
		if _, ok := metrics[groupID]; !ok {
			return metrics, fmt.Errorf("missing item group '%d'", groupID)
		}

		group := metrics[groupID]
		if _, ok := group.Assessments[item.AssessmentID]; !ok {
			return metrics, fmt.Errorf("missing item group '%d'", groupID)
		}

		assessment := group.Assessments[item.AssessmentID]

		var value float64

		switch assessment.Unit {
		case "calories":
			fallthrough
		case "kilograms":
			fallthrough
		case "pounds":
			v, err := strconv.ParseFloat(string(item.Value), 64)
			if err != nil {
				return metrics, fmt.Errorf("failed parsing value '%s': %w", item.Value, err)
			}

			value = v

		case "yes/no":
			if item.Value == "yes" {
				value = 1
			}

		case "beats per minute":
			// skip
		case "centimeters":
			// skip
		case "inches":
			// skip
		case "percent":
			// skip
		case "reps":
			// skip
		case "time":
			// skip
		case "other":
			// skip

		default:
			fmt.Printf("unsupported unit '%s'\n", assessment.Unit)
		}

		metrics[groupID].Assessments[item.AssessmentID].Samples = append(
			metrics[groupID].Assessments[item.AssessmentID].Samples,
			Sample{
				ID:    item.ID,
				Value: value,
				Date:  item.Date,
			},
		)
	}

	return metrics, nil
}

func (tc *Service) Metrics(clientID int) (Metrics, error) {
	metrics := Metrics{}

	path := fmt.Sprintf("/clients/%d/assessment_groups", clientID)

	resp, err := tc.get(path)
	if err != nil {
		return metrics, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return metrics, fmt.Errorf("couldn't read response body: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return metrics, fmt.Errorf("response %d: %s", resp.StatusCode, body)
	}

	var metricsResp MetricsResponse
	err = json.Unmarshal(body, &metricsResp)
	if err != nil {
		return metrics, fmt.Errorf("could not unmarshal assessment groups: %w (%+v)", err, string(body))
	}

	metrics, err = metricsResp.Metrics()
	if err != nil {
		return metrics, fmt.Errorf("could not hydrate metrics: %w", err)
	}

	return metrics, nil
}
