// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Nikos Tsitas",
            "email": "nktsitas@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/authorize": {
            "post": {
                "description": "Creates a new authorization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Creates a new authorization",
                "parameters": [
                    {
                        "description": "Create authorization",
                        "name": "authorization",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/gateway.Authorization"
                        }
                    },
                    {
                        "type": "string",
                        "description": "generated.jwt.token",
                        "name": "Token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.authResponse"
                        }
                    }
                }
            }
        },
        "/capture": {
            "post": {
                "description": "Captures amount from authorization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Captures amount from authorization",
                "parameters": [
                    {
                        "description": "Capture Amount",
                        "name": "captureRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.requestParams"
                        }
                    },
                    {
                        "type": "string",
                        "description": "generated.jwt.token",
                        "name": "Token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.actionsResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Logins a user and provides an authentication token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Logins a user and provides an authentication token",
                "parameters": [
                    {
                        "description": "User Credentials",
                        "name": "Credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.loginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.tokenResponse"
                        }
                    }
                }
            }
        },
        "/refund": {
            "post": {
                "description": "Refunds a previously captured amount from authorization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Refunds a previously captured amount from authorization",
                "parameters": [
                    {
                        "description": "Refund Amount",
                        "name": "refundRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.requestParams"
                        }
                    },
                    {
                        "type": "string",
                        "description": "generated.jwt.token",
                        "name": "Token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.actionsResponse"
                        }
                    }
                }
            }
        },
        "/status/ping": {
            "get": {
                "description": "Get a server status update",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Get a server status update",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/void": {
            "post": {
                "description": "Voids a transaction without charging the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Voids a transaction without charging the user",
                "parameters": [
                    {
                        "description": "Refund Amount",
                        "name": "voidRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.voidRequestParams"
                        }
                    },
                    {
                        "type": "string",
                        "description": "generated.jwt.token",
                        "name": "Token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.actionsResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.loginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "auth.tokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "generated.jwt.token"
                }
            }
        },
        "bank.CreditCard": {
            "type": "object",
            "properties": {
                "cvv": {
                    "type": "string"
                },
                "expiry": {
                    "type": "string"
                },
                "number": {
                    "type": "string"
                }
            }
        },
        "gateway.Authorization": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 100
                },
                "credit_card": {
                    "type": "object",
                    "$ref": "#/definitions/bank.CreditCard"
                },
                "currency": {
                    "type": "string",
                    "example": "EUR"
                }
            }
        },
        "handlers.actionsResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 100
                },
                "currency": {
                    "type": "string",
                    "example": "EUR"
                }
            }
        },
        "handlers.authResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 100
                },
                "currency": {
                    "type": "string",
                    "example": "EUR"
                },
                "id": {
                    "type": "string",
                    "example": "unique_authorization_id"
                }
            }
        },
        "handlers.requestParams": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 100
                },
                "id": {
                    "type": "string",
                    "example": "unique_authorization_id"
                }
            }
        },
        "handlers.voidRequestParams": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "unique_authorization_id"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:2012",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Checkout.com API Challenge",
	Description: "This is a simple Gateway service for Payments",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
