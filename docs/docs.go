package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "Go Packing Service API",
        "description": "API for calculating optimized pack allocations and managing pack-size configuration.",
        "version": "1.0"
    },
    "basePath": "/",
    "schemes": ["http"],
    "paths": {
        "/api/v1/calculate": {
            "post": {
                "summary": "Calculate pack breakdown",
                "tags": ["Calculate"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/CalculateRequest"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/PackBreakdown"}
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {"$ref": "#/definitions/ErrorResponse"}
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {"$ref": "#/definitions/ErrorResponse"}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {"$ref": "#/definitions/ErrorResponse"}
                    }
                }
            }
        },
        "/api/v1/pack-sizes": {
            "get": {
                "summary": "Get current pack sizes",
                "tags": ["Pack Sizes"],
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/PackSizesResponse"}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {"$ref": "#/definitions/ErrorResponse"}
                    }
                }
            },
            "put": {
                "summary": "Replace pack sizes",
                "tags": ["Pack Sizes"],
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/PackSizesRequest"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/PackSizesResponse"}
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {"$ref": "#/definitions/ErrorResponse"}
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {"$ref": "#/definitions/ErrorResponse"}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {"$ref": "#/definitions/ErrorResponse"}
                    }
                }
            }
        }
    },
    "definitions": {
        "CalculateRequest": {
            "type": "object",
            "required": ["amount"],
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 251
                }
            }
        },
        "PackSizesRequest": {
            "type": "object",
            "required": ["pack_sizes"],
            "properties": {
                "pack_sizes": {
                    "type": "array",
                    "items": {"type": "integer"},
                    "example": [250, 500, 1000, 2000, 5000]
                }
            }
        },
        "PackSizesResponse": {
            "type": "object",
            "properties": {
                "pack_sizes": {
                    "type": "array",
                    "items": {"type": "integer"}
                }
            }
        },
        "PackBreakdown": {
            "type": "object",
            "properties": {
                "size": {"type": "integer"},
                "count": {"type": "integer"}
            }
        },
        "ErrorBody": {
            "type": "object",
            "properties": {
                "code": {"type": "string"},
                "message": {"type": "string"}
            }
        },
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {"$ref": "#/definitions/ErrorBody"}
            }
        }
    }
}`

type swaggerInfo struct{}

func (s *swaggerInfo) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register(swag.Name, &swaggerInfo{})
}
