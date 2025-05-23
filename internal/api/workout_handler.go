package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"workout/internal/utils"
	"workout/store"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}
func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: ReadIdParam : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id "})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "internal server error"})
		return

	}
	if workout == nil {
		wh.logger.Printf("ERROR: getWorkoutById : workout is nil: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "no such workout exists"})
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR : decodingCreateWorkout : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request sent"})
		return
	}
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR : creatingWorkout : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": createdWorkout})

}
func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: ReadIdParam : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: GetWorkoutById : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	if existingWorkout == nil {
		wh.logger.Printf("ERROR: GetWorkoutById (No such workout): %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "no such workoutExists"})
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
		wh.logger.Printf("ERROR: UpdateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}
func (wh *WorkoutHandler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: ReadIdParam : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}
	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: DeleteWorkout : %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "no such workout exists"})
		return
	}
	if err != nil {
		wh.logger.Printf("ERROR: DeleteWorkout : %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
