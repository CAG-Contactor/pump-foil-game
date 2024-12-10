package main

import (
	docs "backend/docs"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ContestantDTO struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GameResultDTO struct {
	SplitTime float64 `json:"splitTime"`
	EndTime   float64 `json:"endTime"`
}

type LeaderboardEntryDTO struct {
	Contestant ContestantDTO `json:"contestant"`
	GameResult GameResultDTO `json:"result"`
}

type QueueItemDTO struct {
	Timestamp  int64         `json:"timestamp"`
	Contestant ContestantDTO `json:"contestant"`
}

var allContestants = []ContestantDTO{
	{Email: "RJj8W@example.com", Name: "John Doe"},
	{Email: "J2X9j@example.com", Name: "Jane Doe"},
}
var allResults []LeaderboardEntryDTO
var gameState *QueueItemDTO = nil
var queue = []QueueItemDTO{}

// getContestants returns a list of ContestantDTOs filtered by the given string
//
// The filter string can have the following values:
//
// - "NOT_ENQUEUED": returns all contestants that are not in the queue
// - "ENQUEUED": returns all contestants that are in the queue
// - "ALL": returns all contestants
//
// If the filter string is not recognized, the function returns nil.
func getContestants(filter string) []ContestantDTO {

	switch filter {
	case "NOT_ENQUEUED":
		contenstantList := []ContestantDTO{}
		for _, contestant := range allContestants {
			found := false
			for _, queueItem := range queue {
				if queueItem.Contestant.Email == contestant.Email {
					found = true
					break
				}
			}
			if !found {
				contenstantList = append(contenstantList, contestant)
			}
		}
		return contenstantList
	case "ENQUEUED":
		contenstantList := []ContestantDTO{}
		for _, contestant := range allContestants {
			found := false
			for _, queueItem := range queue {
				if queueItem.Contestant.Email == contestant.Email {
					found = true
					break
				}
			}
			if found {
				contenstantList = append(contenstantList, contestant)
			}
		}
		return contenstantList
	case "ALL":
		return allContestants
	default:
		return nil
	}
}

// addContestant adds a new contestant to the list of all contestants and enqueues them.
// If the contestant already exists, they are only enqueued.
// Returns the QueueItemDTO for the enqueued contestant or an error.
func addContestant(contestant ContestantDTO) (QueueItemDTO, error) {
	for _, c := range allContestants {
		if c.Email == contestant.Email {
			newQueueItem, _ := enqueueContestant(contestant)
			return newQueueItem, nil
		}
	}
	allContestants = append(allContestants, contestant)
	newQueueItem, _ := enqueueContestant(contestant)
	return newQueueItem, nil
}

// enqueueContestant adds a contestant to the queue with the current timestamp.
// Returns the QueueItemDTO for the enqueued contestant and an error if any occurs.
func enqueueContestant(contestant ContestantDTO) (QueueItemDTO, error) {
	newQueueItem := QueueItemDTO{time.Now().UnixNano(), contestant}
	queue = append(queue, newQueueItem)
	return newQueueItem, nil
}

// startGame starts a game with an optional timestamp to start a game for any entry in the queue.
// If no timestamp is given, the function starts a game for the first item in the queue.
// If the given timestamp is not found in the queue, the function returns an error.
// If the game is already in progress, the function returns an error.
// If the queue is empty, the function returns an error.
func startGame(timestamp *int64) (*QueueItemDTO, error) {
	// Check if the queue is not empty
	if len(queue) > 0 {
		// If the game is not in progress, start it with the first item in the queue
		if gameState == nil {
			if timestamp != nil {
				for i, entry := range queue {
					if entry.Timestamp == *timestamp {
						gameState = &entry
						queue = append(queue[:i], queue[i+1:]...)
						return &entry, nil
					}
				}
				return nil, errors.New("timestamp not found in queue")
			} else {
				queueItem := queue[0]
				queue = queue[1:]
				gameState = &queueItem
				return &queueItem, nil
			}
		} else {
			return gameState, errors.New("game is already in progress")
		}

	} else {
		return gameState, errors.New("queue is empty")
	}
}

// finishGame finishes the current game with the given result and updates the leaderboard.
// It returns the updated leaderboard or an error if the game is not in progress.
// If the contestant already has a result in the leaderboard, the new result is only added if it is better than the previous result.
func finishGame(result GameResultDTO) ([]LeaderboardEntryDTO, error) {
	// Check if the game is in progress
	if gameState != nil {
		leaderboardEntry := LeaderboardEntryDTO{gameState.Contestant, result}

		// Check if the contestant already has a result
		for i, entry := range allResults {
			if entry.Contestant.Email == leaderboardEntry.Contestant.Email {
				// Update the result if the new result is better
				if leaderboardEntry.GameResult.EndTime < entry.GameResult.EndTime {
					allResults[i] = leaderboardEntry
					gameState = nil
					return allResults, nil
				}
			}
		}
		// Otherwise, add the new result
		allResults = append(allResults, leaderboardEntry)

		gameState = nil
		return allResults, nil
	} else {
		return nil, errors.New("game is not in progress")
	}
}

// abortGame aborts the current game by resetting the game state.
// It returns nil if successful, indicating the game state has been reset.
func abortGame() error {
	// Reset the game state
	gameState = nil
	return nil
}

// deleteQueueItem deletes the queue item with the given timestamp.
// It returns nil if the item was found and deleted, or an error if the item was not found.
func deleteQueueItem(timestamp int64) error {

	for i, entry := range queue {
		if entry.Timestamp == timestamp {
			queue = append(queue[:i], queue[i+1:]...)
			return nil
		}
	}

	return errors.New("queue item not found")
}

// @BasePath /api/v1

// getContestantsHandler godoc
// @Summary Get contestants
// @Schemes http https
// @Description Get contestants based on filter one of ALL, NOT_ENQUEUED or ENQUEUED
// @Tags example
// @Accept json
// @Produce json
// @Param filter query string  false  "One of ALL, NOT_ENQUEUED or ENQUEUED. If omitted ALL is used."
// @Success 200 {array} ContestantDTO
// @Router /contestants [get]
func getContestantsHandler(g *gin.Context) {
	// GET /contestants?filter={ALL,NOT_ENQUEUED,ENQUEUED}
	filter := strings.TrimSpace(g.DefaultQuery("filter", "ALL"))

	contestants := getContestants(filter)
	if contestants == nil {
		g.AbortWithStatus(http.StatusBadRequest)
		return
	}

	g.JSON(http.StatusOK, contestants)
}

// addContestantHandler godoc
// @Summary Add a contestant
// @Schemes
// @Description Add a contestant to the database, the contestant will also be added to the queue
// @Tags example
// @Accept json
// @Produce json
// @Param contestant body ContestantDTO true "Contestant to add"
// @Success 200 {object} QueueItemDTO
// @Router /contestants [post]
func addContestantHandler(g *gin.Context) {
	contestant := ContestantDTO{}
	if err := g.BindJSON(&contestant); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newQueueItem, err := addContestant(contestant)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, newQueueItem)
}

// gameStartHandler godoc
// @Summary Start a game for a contestant
// @Schemes
// @Description Start a game for a contestant if the optional query parameter timestamp is provided the specific queueitem will be started otherwise the first item in the queue will be started
// @Tags example
// @Accept json
// @Produce json
// @Param timestamp query string false "imestamp of queueitem to start"
// @Success 200 {object} QueueItemDTO
// @Router /game-start [post]
func gameStartHandler(g *gin.Context) {
	timestampString := g.Query("timestamp")
	if timestampString != "" {
		timestamp, err := strconv.ParseInt(timestampString, 10, 64)
		if err != nil {
			g.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		queueItem, err := startGame(&timestamp)
		if err != nil {
			g.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		g.JSON(http.StatusOK, queueItem)
	} else {
		queueItem, err := startGame(nil)
		if err != nil {
			g.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		g.JSON(http.StatusOK, queueItem)
	}
}

// gameFinishHandler godoc
// @Summary Finish a game and save the result
// @Schemes
// @Description Finish a game and save the result in the database
// @Tags example
// @Accept json
// @Produce json
// @Param result body GameResultDTO true "Result of the game"
// @Success 200 {array} LeaderboardEntryDTO
// @Router /game-finish [post]
func gameFinishHandler(g *gin.Context) {
	gameResult := GameResultDTO{}
	if err := g.BindJSON(&gameResult); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	leaderboard, err := finishGame(gameResult)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, leaderboard)
}

// gameAbortHandler godoc
// @Summary Abort the current game
// @Schemes
// @Description Abort the ongoing game and respond with a status message
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} string "aborted"
// @Failure 500 {object} error
// @Router /game-abort [post]
func gameAbortHandler(g *gin.Context) {

	if err := abortGame(); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, "aborted")
}

// getQueue godoc
// @Summary Get the queue
// @Schemes
// @Description Get the current queue
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {array} QueueItemDTO
// @Router /queue [get]
func getQueueHandler(g *gin.Context) {
	g.JSON(http.StatusOK, queue)
}

// deleteQueueItemHandler godoc
// @Summary Delete a contestant from the queue
// @Schemes
// @Description Delete a queue item from the queue
// @Tags example
// @Accept json
// @Produce json
// @Param timestamp path int64 false "timestamp of the queue item to delete"
// @Success 200
// @Failure 500 {object} error
// @Router /queue/{timestamp} [delete]
func deleteQueueItemHandler(g *gin.Context) {
	timestampString := g.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = deleteQueueItem(timestamp)
	if err != nil {
		g.AbortWithError(http.StatusNotFound, err)
		return
	}

	g.Status(http.StatusOK)
}

// getLeaderboardHandler godoc
// @Summary Get leaderboard
// @Schemes
// @Description Retrieve the current leaderboard
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {array} LeaderboardEntryDTO
// @Router /leaderboard [get]
func getLeaderboardHandler(g *gin.Context) {
	g.JSON(http.StatusOK, allResults)
}

// main starts a gin server and maps all the available endpoints
func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/contestants", getContestantsHandler)
		v1.POST("/contestants", addContestantHandler)
		v1.POST("/game-start", gameStartHandler)
		v1.POST("/game-finish", gameFinishHandler)
		v1.POST("/game-abort", gameAbortHandler)
		v1.GET("/queue", getQueueHandler)
		v1.DELETE("/queue/:timestamp", deleteQueueItemHandler)
		v1.GET("/leaderboard", getLeaderboardHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
