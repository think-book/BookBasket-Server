package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/think-book/BookBasket-Server/handler"

	//Localhost　からアクセスする用
	"net/http"
	"github.com/labstack/echo/middleware"
)

func main() {
	// mysqlに接続
	db, err := sqlx.Open("mysql", "root:root@tcp(my_db:3306)/bookbasket")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// handlerにデータベースの参照を渡す。
	handler.SetDB(db)

	// Echoのインスタンス作る
	e := echo.New()

	// 異なるプラットフォームからのアクセスを許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))

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
