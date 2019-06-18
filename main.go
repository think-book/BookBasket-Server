package main

import (
	"github.com/labstack/echo"
	"github.com/think-book/BookBasket-Server/handler"
)

func main() {
	// Echoのインスタンス作る
	e := echo.New()

	// ルーティング
	e.GET("/api/v1/event", handler.GetBookInfoAll())
	e.GET("/api/v1/event/:id", handler.GetBookInfo())

	// サーバー起動
	e.Start(":8080") //ポート番号指定してね
}
