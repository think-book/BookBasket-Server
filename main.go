package main

import (
	"github.com/labstack/echo"
	"github.com/think-book/BookBasket-Server/handler"
)

func main() {
	// Echoのインスタンス作る
	e := echo.New()

	// ルーティング
	e.GET("/books", handler.GetBookMetaInfoAll)
	e.GET("/books/:ISBN", handler.GetBookProfile)
	e.POST("/books", handler.PostMetaInfo)
	e.POST("/books/:ISBN", handler.PostBookProfile)

	e.Start(":8080")

	// サーバー起動

}
