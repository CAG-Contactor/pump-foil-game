package main

import (
	docs "backend/docs"
	"errors"
	"net/http"
	"strconv"
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

func getContestants(filter string) []ContestantDTO {

	switch filter {
	case "NOT_ENQUEUED":
		return allContestants
	case "ENQUEUED":
		return allContestants
	case "ALL":
		return allContestants
	default:
		return nil
	}
}

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

func enqueueContestant(contestant ContestantDTO) (QueueItemDTO, error) {
	newQueueItem := QueueItemDTO{time.Now().UnixNano(), contestant}
	queue = append(queue, newQueueItem)
	return newQueueItem, nil
}

func startGame() error {
	if len(queue) > 0 {
		queueItem := queue[0]

		if gameState == nil {
			gameState = &queueItem
			return nil
		} else {
			return errors.New("game is already in progress")
		}

	} else {
		return errors.New("queue is empty")
	}
}

func finishGame(result GameResultDTO) []LeaderboardEntryDTO {

	leaderboardEntry := LeaderboardEntryDTO{gameState.Contestant, result}

	allResults = append(allResults, leaderboardEntry)
	return allResults
}

func abortGame() error {
	return nil
}

func deleteQueueItem(timestamp int64) error {

	for _, entry := range queue {
		if entry.Timestamp == timestamp {

		}
	}

	return nil
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
	filter := g.DefaultQuery("filter", "ALL")

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
// @Description Add a contestant to the database
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
// @Description Start a game for a contestant
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} ContestantDTO
// @Router /game-start [post]
func gameStartHandler(g *gin.Context) {
	contestant := ContestantDTO{}
	if err := g.BindJSON(&contestant); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := startGame(); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, contestant)
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

	leaderboard := finishGame(gameResult)

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

// deleteQueueItemHandler godoc
// @Summary Delete a contestant from the queue
// @Schemes
// @Description Delete a contestant from the queue
// @Tags example
// @Accept json
// @Produce json
// @Param timestamp path int64 true "timestamp of the contestant to delete"
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

	deleteQueueItem(timestamp)

	g.Status(http.StatusOK)
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
		v1.DELETE("/queue-item/:timestamp", deleteQueueItemHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
