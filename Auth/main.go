package main

import server "AuthServ/Server"

// @title           Сервис аутентификации
// @version         1.0
// @description     Тестовое задание от MEDODS

// @contact.name   Evans Trein
// @contact.email  evanstrein@icloud.com
// @contact.url  https://github.com/EvansTrein

// @host      localhost:4000
// @schemes   http

// @securityDefinitions.apikey accessToken
// @in header
// @name Authorization
// @description Type Bearer  accessToken_example

// @securityDefinitions.apikey refreshRefresh
// @in header
// @name RefreshToken
// @description Type Bearer  refreshToken_example

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func init() {
	server.InitServer()
}

func main() {
	server.StartServer()
}
