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
                    }
                }
            },
            "post": {
                "description": "Add a contestant to the database",
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
                    }
                }
            }
        },
        "/game-start": {
            "post": {
                "description": "Start a game for a contestant",
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
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.QueueItemDTO"
                        }
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
                                "$ref": "#/definitions/main.QueueItemDTO"
                            }
                        }
                    }
                }
            }
        },
        "/queue/{timestamp}": {
            "delete": {
                "description": "Delete a contestant from the queue",
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
                        "description": "timestamp of the contestant to delete",
                        "name": "timestamp",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
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
