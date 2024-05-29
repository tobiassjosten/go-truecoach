package truecoach

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type User struct {
	ID      int `json:"id"`
	ImageID int `json:"image_id"`

	Image *Image

	Birthday         time.Time `json:"birthday_timestamp"`
	Demo             bool      `json:"demo"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	Gender           string    `json:"gender"`
	Height           int       `json:"height"`
	InPerson         bool      `json:"in_person"`
	InviteAcceptedAt time.Time `json:"invite_accepted_at"`
	InviteToken      string    `json:"invite_token"`
	LastName         string    `json:"last_name"`
	Location         string    `json:"location"`
	MFPConnected     bool      `json:"is_mfp_connected"`
	Online           bool      `json:"is_online"`
	PendingInvite    bool      `json:"has_pending_invite"`
	PhoneNumber      string    `json:"phone_number"`
	Skype            string    `json:"skype"`
	Timezone         string    `json:"timezone"`
	TimezoneOffset   int       `json:"timezone_offset"`
	Trainer          bool      `json:"has_trainer"`
	Units            string    `json:"units"`
	Weight           float64   `json:"weight"`
}

type Client struct {
	ID             int `json:"id"`
	OrganizationID int `json:"organization_id"`
	// StripePaymentSourceID ??? `json:"stripe_payment_source_id"`
	TrainerID int `json:"trainer_id"`
	UserID    int `json:"user_id"`

	User *User

	CompletedWorkoutsCount int       `json:"completed_workouts_count"`
	ComplianceRateMonth    float64   `json:"compliance_rate_month"`
	ComplianceRateQuarter  float64   `json:"compliance_rate_quarter"`
	ComplianceRateWeek     float64   `json:"compliance_rate_week"`
	CreatedAt              time.Time `json:"created_at"`
	Delinquent             bool      `json:"is_delinquent"`
	Due                    JSONDate  `json:"due"`
	DueDateLocked          bool      `json:"due_date_locked"`
	// Equipment ??? `json:"equipment"`
	// EquipmentAttachments []??? `json:"equipment_attachments"`
	// Goals ??? `json:"goals"`
	HideFromFeed bool `json:"hide_from_feed"`
	// Limitations ??? `json:"limitations"`
	MissedSessionsCount int       `json:"missed_sessions_count"`
	Slug                string    `json:"slug"`
	State               string    `json:"state"`
	Transferring        bool      `json:"is_transferring"`
	Type                string    `json:"type"`
	UpdatedAt           time.Time `json:"updated_at"`

	Links struct {
		AssessmentGroups    string `json:"assessment_groups"`
		Assessments         string `json:"assessments"`
		Conversation        string `json:"conversation"`
		DailyNutritionLogs  string `json:"daily_nutrition_logs"`
		HealthTrackings     string `json:"health_trackings"`
		Notes               string `json:"notes"`
		NutritionPlan       string `json:"nutrition_plan"`
		PhotoSessions       string `json:"photo_sessions"`
		Skeletons           string `json:"skeletons"`
		StripeSubscriptions string `json:"stripe_subscriptions"`
		WeightTrackings     string `json:"weight_trackings"`
		Workouts            string `json:"workouts"`
	} `json:"links"`

	Settings struct {
		CurrentWeekOnly      bool `json:"currentWeekOnly"`
		DailyWorkoutEmails   bool `json:"dailyWorkoutEmails"`
		MissedWorkoutsEmails bool `json:"missedWorkoutsEmails"`
		NewCommentEmails     bool `json:"newCommentEmails"`
		NewMessageEmails     bool `json:"newMessageEmails"`
		OverrideDefaults     bool `json:"overrideDefaults"`
		WeeklyDigest         bool `json:"weeklyDigest"`
		WorkoutsThreshold    int  `json:"workoutsThreshold"`
	} `json:"settings"`
}

type Image struct {
	ID int `json:"id"`

	FileSize   int    `json:"file_size"`
	MimeType   string `json:"mime_type"`
	URL        string `json:"url"`
	UploaderID int    `json:"uploaded_by_id"`

	Parent struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	} `json:"parent"`
}

func (tc *Service) Clients() ([]Client, error) {
	resp, err := tc.get("/clients")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("response %d: %s", resp.StatusCode, body)
	}

	var data struct {
		*BaseResponse
		Images  []Image  `json:"images"`
		Users   []User   `json:"users"`
		Clients []Client `json:"clients"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal clients: %w", err)
	}

	for i, u := range data.Users {
		for _, ii := range data.Images {
			if ii.ID == u.ImageID {
				data.Users[i].Image = &ii
				break
			}
		}
	}

	for i, c := range data.Clients {
		for _, u := range data.Users {
			if u.ID == c.UserID {
				data.Clients[i].User = &u
				break
			}
		}
	}

	return data.Clients, nil
}
