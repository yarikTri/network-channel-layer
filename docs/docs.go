// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Yaroslav Kuzmin",
            "email": "yarik1448kuzmin@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/code": {
            "post": {
                "description": "Осуществляет кодировку сообщения в код Хэмминга [15, 11], внесение ошибки в каждый закодированный 15-битовый кадр с вероятностью 7%, исправление внесённых ошибок, раскодировку кадров в изначальное сообщение. Затем отправляет результат в Procuder-сервис транспортного уровня. Сообщение может быть потеряно с вероятностью 1%.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Code"
                ],
                "summary": "Code network flow",
                "parameters": [
                    {
                        "description": "Сообщение",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Кодировка и отправка запущены"
                    },
                    "400": {
                        "description": "Ошибка при чтении сообщения"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{"https", "http"},
	Title:            "КР СТ АСОИУ Сервис Канального Уровня",
	Description:      "Сервис имитирует передачу данных через ненадёжную сеть с защитой данных с помощью кодировки Хэмминга [15, 11].",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}