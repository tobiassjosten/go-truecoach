package truecoach

import (
	"encoding/json"
	"fmt"
	"io"
)

type NutritionPlan struct {
	Threshold   int
	Description string

	MonCarbohydrate int `json:"mon_carbs"`
	MonFat          int `json:"mon_fat"`
	MonProtein      int `json:"mon_protein"`
	MonFiber        int `json:"mon_fiber"`
	MonCalories     int `json:"mon_calories"`

	TueCarbohydrate int `json:"tue_carbs"`
	TueFat          int `json:"tue_fat"`
	TueProtein      int `json:"tue_protein"`
	TueFiber        int `json:"tue_fiber"`
	TueCalories     int `json:"tue_calories"`

	WedCarbohydrate int `json:"wed_carbs"`
	WedFat          int `json:"wed_fat"`
	WedProtein      int `json:"wed_protein"`
	WedFiber        int `json:"wed_fiber"`
	WedCalories     int `json:"wed_calories"`

	ThuCarbohydrate int `json:"thu_carbs"`
	ThuFat          int `json:"thu_fat"`
	ThuProtein      int `json:"thu_protein"`
	ThuFiber        int `json:"thu_fiber"`
	ThuCalories     int `json:"thu_calories"`

	FriCarbohydrate int `json:"fri_carbs"`
	FriFat          int `json:"fri_fat"`
	FriProtein      int `json:"fri_protein"`
	FriFiber        int `json:"fri_fiber"`
	FriCalories     int `json:"fri_calories"`

	SatCarbohydrate int `json:"sat_carbs"`
	SatFat          int `json:"sat_fat"`
	SatProtein      int `json:"sat_protein"`
	SatFiber        int `json:"sat_fiber"`
	SatCalories     int `json:"sat_calories"`

	SunCarbohydrate int `json:"sun_carbs"`
	SunFat          int `json:"sun_fat"`
	SunProtein      int `json:"sun_protein"`
	SunFiber        int `json:"sun_fiber"`
	SunCalories     int `json:"sun_calories"`
}

func (tc *Service) ClientNutritionPlan(clientID int) (NutritionPlan, error) {
	path := fmt.Sprintf("/clients/%d/nutrition_plan", clientID)

	resp, err := tc.get(path)
	if err != nil {
		return NutritionPlan{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NutritionPlan{}, fmt.Errorf("couldn't read response body: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return NutritionPlan{}, fmt.Errorf("response %d: %s", resp.StatusCode, body)
	}

	var nutritionPlanResponse struct {
		NutritionPlan NutritionPlan `json:"nutrition_plan"`
	}

	err = json.Unmarshal(body, &nutritionPlanResponse)
	if err != nil {
		return NutritionPlan{}, fmt.Errorf("could not unmarshal nutrition plan: %w (%+v)", err, string(body))
	}

	return nutritionPlanResponse.NutritionPlan, nil
}

type DailyNutritionLog struct {
	ID             int           `json:"id"`
	ClientID       int           `json:"client_id"`
	Due            JSONDate      `json:"due"`
	GoalCalories   int           `json:"goal_calories"`
	ActualCalories int           `json:"actual_calories"`
	GoalCarbs      int           `json:"goal_carbs"`
	ActualCarbs    int           `json:"actual_carbs"`
	GoalProtein    int           `json:"goal_protein"`
	ActualProtein  int           `json:"actual_protein"`
	GoalFat        int           `json:"goal_fat"`
	ActualFat      int           `json:"actual_fat"`
	GoalFiber      int           `json:"goal_fiber"`
	ActualFiber    int           `json:"actual_fiber"`
	Editable       bool          `json:"is_editable"`
	Notes          string        `json:"notes"`
	Attachments    []interface{} `json:"attachments"`
}

func (tc *Service) ClientDailyNutritionLogs(clientID int) ([]DailyNutritionLog, error) {
	path := fmt.Sprintf("/clients/%d/daily_nutrition_logs", clientID)

	resp, err := tc.get(path)
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

	var dailyNutritionLogsResponse struct {
		DailyNutritionLogs []DailyNutritionLog `json:"daily_nutrition_logs"`
	}

	err = json.Unmarshal(body, &dailyNutritionLogsResponse)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal daily nutrition logs: %w (%+v)", err, string(body))
	}

	return dailyNutritionLogsResponse.DailyNutritionLogs, nil
}
