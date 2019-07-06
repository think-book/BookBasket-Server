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
	e.GET("/books/:ISBN/threads", handler.GetThreadTitles)
	e.GET("/threads/:threadID", handler.GetThreadMessages)
	e.POST("/books", handler.PostBookInfo)
	e.POST("/books/:ISBN/threads", handler.PostThreadTitle)
	e.POST("/threads/:threadID", handler.PostThreadMessage)

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
