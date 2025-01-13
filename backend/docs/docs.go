// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/contestants": {
            "get": {
                "description": "Get contestants based on filter one of ALL, NOT_ENQUEUED or ENQUEUED",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Get contestants",
                "parameters": [
                    {
                        "type": "string",
                        "description": "One of ALL, NOT_ENQUEUED or ENQUEUED. If omitted ALL is used.",
                        "name": "filter",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.ContestantDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            },
            "post": {
                "description": "Add a contestant to the database, the contestant will also be added to the queue",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Add a contestant",
                "parameters": [
                    {
                        "description": "Contestant to add",
                        "name": "contestant",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.ContestantDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.QueueItemDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/game-abort": {
            "post": {
                "description": "Abort the ongoing game and respond with a status message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Abort the current game",
                "responses": {
                    "200": {
                        "description": "aborted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/game-finish": {
            "post": {
                "description": "Finish a game and save the result in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Finish a game and save the result",
                "parameters": [
                    {
                        "description": "Result of the game",
                        "name": "result",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.GameResultDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.LeaderboardEntryDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/game-start": {
            "post": {
                "description": "Start a game for a contestant if the optional query parameter timestamp is provided the specific queueitem will be started otherwise the first item in the queue will be started",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Start a game for a contestant",
                "parameters": [
                    {
                        "type": "string",
                        "description": "imestamp of queueitem to start",
                        "name": "timestamp",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.QueueItemDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/leaderboard": {
            "get": {
                "description": "Retrieve the current leaderboard",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Get leaderboard",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.LeaderboardEntryDTO"
                            }
                        }
                    }
                }
            }
        },
        "/queue": {
            "get": {
                "description": "Get the current queue",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Get the queue",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.QueueItemDTO"
                            }
                        }
                    }
                }
            }
        },
        "/queue/{timestamp}": {
            "delete": {
                "description": "Delete a queue item from the queue",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "Delete a contestant from the queue",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "timestamp of the queue item to delete",
                        "name": "timestamp",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/ws": {
            "get": {
                "description": "Handle websocket connections. Websocket clients can connectusing the url ws://localhost:8080/api/v1/ws",
                "summary": "Handle websocket connections",
                "responses": {}
            }
        }
    },
    "definitions": {
        "db.ContestantDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "db.QueueItemDTO": {
            "type": "object",
            "properties": {
                "contestant": {
                    "$ref": "#/definitions/db.ContestantDTO"
                },
                "timestamp": {
                    "type": "integer"
                }
            }
        },
        "main.ContestantDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.GameResultDTO": {
            "type": "object",
            "properties": {
                "endTime": {
                    "type": "number"
                },
                "splitTime": {
                    "type": "number"
                }
            }
        },
        "main.LeaderboardEntryDTO": {
            "type": "object",
            "properties": {
                "contestant": {
                    "$ref": "#/definitions/main.ContestantDTO"
                },
                "result": {
                    "$ref": "#/definitions/main.GameResultDTO"
                }
            }
        },
        "main.QueueItemDTO": {
            "type": "object",
            "properties": {
                "contestant": {
                    "$ref": "#/definitions/main.ContestantDTO"
                },
                "timestamp": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
