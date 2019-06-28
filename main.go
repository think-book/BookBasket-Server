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
	e.GET("/books/:ISBN/forum", handler.GetForumTitles)
	e.GET("/forum/:forumID", handler.GetForumMessages)
	e.POST("/books", handler.PostBookInfo)

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
