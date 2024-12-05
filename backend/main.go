package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	//docs "myserver/docs"
	//swaggerfiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Contestant struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GameResult struct {
	SplitTime float64 `json:"splitTime"`
	EndTime   float64 `json:"endTime"`
}

type LeaderboardEntry struct {
	Contestant Contestant `json:"contestant"`
	GameResult GameResult `json:"result"`
}

type QueueItem struct {
	Timestamp  int64      `json:"timestamp"`
	Contestant Contestant `json:"contestant"`
}

var allContestants []Contestant
var allResults []LeaderboardEntry
var gameState *QueueItem = nil
var queue []QueueItem

func getContestants(filter string) []Contestant {

	switch filter {
	case "NOT_ENQUEUED":
		return allContestants
	case "ENQUEUED":
		return allContestants
	default:
		return allContestants
	}
}

func addContestant(contestant Contestant) error {
	allContestants = append(allContestants, contestant)
	return nil
}

func startGame(contestant Contestant) error {
	if gameState == nil {
		gameState = &QueueItem{time.Now().UnixNano(), contestant}
	}

	return nil
}

func finishGame(result GameResult) []LeaderboardEntry {

	leaderboardEntry := LeaderboardEntry{gameState.Contestant, result}

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
// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

func getContestantsHandler(g *gin.Context) {
	// GET /contestants?filter={ALL,NOT_ENQUEUED,ENQUEUED}
	filter := g.DefaultQuery("filter", "ALL")

	contestants := getContestants(filter)

	g.JSON(http.StatusOK, contestants)
}

func addContestantHandler(g *gin.Context) {
	contestant := Contestant{}
	if err := g.BindJSON(&contestant); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := addContestant(contestant); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, contestant)
}

func gameStartHandler(g *gin.Context) {
	contestant := Contestant{}
	if err := g.BindJSON(&contestant); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := startGame(contestant); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, contestant)
}

func gameFinishHandler(g *gin.Context) {
	gameResult := GameResult{}
	if err := g.BindJSON(&gameResult); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	leaderboard := finishGame(gameResult)

	g.JSON(http.StatusOK, leaderboard)
}

func gameAbortHandler(g *gin.Context) {

	if err := abortGame(); err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, "aborted")
}

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

func main() {
	r := gin.Default()
	//docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/contestants", getContestantsHandler)
		v1.POST("/contestants", addContestantHandler)
		v1.POST("/game-start", gameStartHandler)
		v1.POST("/game-finish", gameFinishHandler)
		v1.POST("/game-abort", gameAbortHandler)
		v1.DELETE("/queue-item/:timestamp", Helloworld)
	}
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
