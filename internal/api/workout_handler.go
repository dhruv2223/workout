package api

import (
	"encoding/json"
	"fmt"
	"internal/runtime/exithook"
	"net/http"
	"strconv"
	"workout/store"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}
func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam(r, "id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
		return

	}
	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		http.Error(w, "failed to get workout", http.StatusInternalServerError)
	}
	if workout == nil {
		http.NotFound(w, r)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)

}
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout 1", http.StatusInternalServerError)
		return
	}
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)

}
func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	paramWorkoutId := chi.URLParam("id")
	if paramWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(paramWorkoutId, 10, 64)
	if err != nil {
		http.Error(w, "failed to parse workout id", http.StatusBadRequest)
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"Description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		fmt.Println("Failed to update workout:", err)
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}
