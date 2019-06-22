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

	// サーバー起動
	go func() {
		if err := e.Start(":8080"); err != nil { //ポート番号指定してね
			e.Logger.Info("Shutting down the server")
		}
	}()

}
