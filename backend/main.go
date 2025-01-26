package main

import (
	"backend/db"
	docs "backend/docs"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

var gameState *db.QueueItemDTO = nil
var dbClient *mongo.Client

// getContestants returns a list of ContestantDTOs filtered by the given string
//
// The filter string can have the following values:
//
// - "NOT_ENQUEUED": returns all contestants that are not in the queue
// - "ENQUEUED": returns all contestants that are in the queue
// - "ALL": returns all contestants
//
// If the filter string is not recognized, the function returns nil.
func getContestants(filter string) []db.ContestantDTO {

	switch filter {
	case "NOT_ENQUEUED":
		contestants, err := db.GetContestants(dbClient, db.GetContestantsNotEnqueued)
		if err != nil {
			return nil
		}
		return contestants
	case "ENQUEUED":
		contestants, err := db.GetContestants(dbClient, db.GetContestantsEnqueued)
		if err != nil {
			return nil
		}
		return contestants
	case "ALL":
		contestants, err := db.GetContestants(dbClient, db.GetContestantsAll)
		if err != nil {
			return nil
		}
		return contestants
	default:
		return nil
	}
}

// addContestant adds a new contestant to the list of all contestants and enqueues them.
// If the contestant already exists, they are only enqueued.
// Returns the QueueItemDTO for the enqueued contestant or an error.
func addContestant(contestant db.ContestantDTO) (*db.QueueItemDTO, error) {
	dbContestant, err := db.AddContestant(dbClient, contestant.Email, contestant.Name)
	if err != nil {
		return nil, err
	}

	queueItem, err := db.EnqueueContestant(dbClient, dbContestant.Id)
	if err != nil {
		return nil, err
	}

	return &db.QueueItemDTO{Timestamp: queueItem.Timestamp, Contestant: contestant}, nil
}

// startGame starts a game with an optional timestamp to start a game for any entry in the queue.
// If no timestamp is given, the function starts a game for the first item in the queue.
// If the given timestamp is not found in the queue, the function returns an error.
// If the game is already in progress, the function returns an error.
// If the queue is empty, the function returns an error.
func startGame(timestamp *int64) (*db.QueueItemDTO, error) {
	queue, err := db.GetQueue(dbClient)
	if err != nil {
		return nil, err
	}
	// Check if the queue is not empty
	if len(queue) > 0 {
		// If the game is not in progress, start it with the first item in the queue
		if gameState == nil {
			if timestamp != nil {
				for _, entry := range queue {
					if entry.Timestamp == *timestamp {
						gameState = &entry
						err = db.DeleteQueueItem(dbClient, entry.Timestamp)
						if err != nil {
							return nil, err
						}
						return &entry, nil
					}
				}
				return nil, errors.New("timestamp not found in queue")
			} else {
				queueItem := queue[0]
				err = db.DeleteQueueItem(dbClient, queueItem.Timestamp)
				if err != nil {
					return nil, err
				}
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
func finishGame(result db.GameResultDTO) ([]db.LeaderboardEntryDTO, error) {
	// Check if the game is in progress
	if gameState != nil {
		//leaderboardEntry := db.LeaderboardEntryDTO{Contestant: gameState.Contestant, GameResult: result}

		_, err := db.CreateOrUpdateLeaderboardEntry(dbClient, *gameState, result)
		if err != nil {
			return nil, err
		}

		allResults, err := db.GetLeaderboard(dbClient)
		if err != nil {
			return nil, err
		}

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

// @BasePath /api/v1

// getContestantsHandler godoc
// @Summary Get contestants
// @Schemes http https
// @Description Get contestants based on filter one of ALL, NOT_ENQUEUED or ENQUEUED
// @Tags example
// @Accept json
// @Produce json
// @Param filter query string  false  "One of ALL, NOT_ENQUEUED or ENQUEUED. If omitted ALL is used."
// @Success 200 {array} db.ContestantDTO
// @Failure 400 {object} error
// @Router /contestants [get]
func getContestantsHandler(g *gin.Context) {
	// GET /contestants?filter={ALL,NOT_ENQUEUED,ENQUEUED}
	filter := strings.TrimSpace(g.DefaultQuery("filter", "ALL"))

	contestants := getContestants(filter)
	if contestants == nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid filter"})
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
// @Param contestant body db.ContestantDTO true "Contestant to add"
// @Success 200 {object} db.QueueItemDTO
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /contestants [post]
func addContestantHandler(g *gin.Context) {
	contestant := db.ContestantDTO{}
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
// @Success 200 {object} db.QueueItemDTO
// @Failure 500 {object} error
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
		notifyInitGame(*queueItem)
		g.JSON(http.StatusOK, queueItem)
	} else {
		queueItem, err := startGame(nil)
		if err != nil {
			g.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		notifyInitGame(*queueItem)
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
// @Param result body db.GameResultDTO true "Result of the game"
// @Success 200 {array} db.LeaderboardEntryDTO
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /game-finish [post]
func gameFinishHandler(g *gin.Context) {
	gameResult := db.GameResultDTO{}
	if err := g.BindJSON(&gameResult); err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	leaderboard, err := finishGame(gameResult)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	notifyEndGame(gameResult)
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

	notifyAbortGame()
	g.JSON(http.StatusOK, "aborted")
}

// getQueue godoc
// @Summary Get the queue
// @Schemes
// @Description Get the current queue
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {array} db.QueueItemDTO
// @Router /queue [get]
func getQueueHandler(g *gin.Context) {
	queue, err := db.GetQueue(dbClient)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if queue == nil {
		queue = []db.QueueItemDTO{}
	}
	g.JSON(http.StatusOK, queue)
}

// deleteContestantHandler godoc
// @Summary Delete a contestant
// @Schemes
// @Description Delete a contestant from the database based on the email parameter
// @Tags example
// @Accept json
// @Produce json
// @Param email path string true "Email of the contestant to delete"
// @Success 200
// @Failure 500 {object} error
// @Router /contestants/{email} [delete]
func deleteContestantHandler(g *gin.Context) {
	email := g.Param("email")
	err := db.DeleteContestant(dbClient, email)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.Status(http.StatusOK)
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
// @Failure 404 {object} error
// @Router /queue/{timestamp} [delete]
func deleteQueueItemHandler(g *gin.Context) {
	timestampString := g.Param("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = db.DeleteQueueItem(dbClient, timestamp)
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
// @Success 200 {array} db.LeaderboardEntryDTO
// @Router /leaderboard [get]
func getLeaderboardHandler(g *gin.Context) {
	allResults, err := db.GetLeaderboard(dbClient)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	g.JSON(http.StatusOK, allResults)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var websockets []*websocket.Conn

// websocketHandler godoc
// @Summary Handle websocket connections
// @Description Handle websocket connections. Websocket clients can connectusing the url ws://localhost:8080/api/v1/ws
// @Router /ws [get]
func websocketHandler(g *gin.Context) {
	conn, err := upgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		return
	}
	//defer conn.Close()
	websockets = append(websockets, conn)
}

func sendNotification(notification string) {
	for _, ws := range websockets {
		ws.WriteMessage(websocket.TextMessage, []byte(notification))
	}
}

type GameNotification struct {
	Type       string  `json:"type"`
	Userid     string  `json:"userid,omitempty"`
	Name       string  `json:"name,omitempty"`
	SplitTime  float64 `json:"splitTime,omitempty"`
	FinishTime float64 `json:"finishTime,omitempty"`
}

func notifyInitGame(queueItem db.QueueItemDTO) {
	notification := GameNotification{
		Type:   "InitGame",
		Userid: queueItem.Contestant.Email,
		Name:   queueItem.Contestant.Name,
	}
	notificationJson, _ := json.Marshal(notification)
	sendNotification(string(notificationJson))
}

func notifyAbortGame() {
	notification := GameNotification{
		Type: "AbortGame",
	}
	notificationJson, _ := json.Marshal(notification)
	sendNotification(string(notificationJson))
}

func notifyEndGame(result db.GameResultDTO) {
	notification := GameNotification{
		Type:       "EndGame",
		SplitTime:  result.SplitTime,
		FinishTime: result.EndTime,
	}
	notificationJson, _ := json.Marshal(notification)
	sendNotification(string(notificationJson))
}

// main starts a gin server and maps all the available endpoints
func main() {
	var err error

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	dbUrl := viper.GetString("mongodb.url")
	port := viper.GetString("port")

	dbClient, err = db.ConnectDB(dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.CloseDB(dbClient)

	log.Println("Connected to database using url:", dbUrl)

	r := gin.Default()
	r.Use(cors.Default())
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/contestants", getContestantsHandler)
		v1.POST("/contestants", addContestantHandler)
		v1.DELETE("/contestants/:email", deleteContestantHandler)
		v1.POST("/game-start", gameStartHandler)
		v1.POST("/game-finish", gameFinishHandler)
		v1.POST("/game-abort", gameAbortHandler)
		v1.GET("/queue", getQueueHandler)
		v1.DELETE("/queue/:timestamp", deleteQueueItemHandler)
		v1.GET("/leaderboard", getLeaderboardHandler)
		v1.GET("/ws", websocketHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":" + port)
}
