package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/think-book/BookBasket-Server/handler"
)

func main() {
	// mysqlに接続
	db, err := sqlx.Open("mysql", "root:password@tcp(172.19.0.2:3306)/bookbasket")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// handlerにデータベースの参照を渡す。
	handler.SetDB(db)

	// Echoのインスタンス作る
	e := echo.New()

	// セッション使うため
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// 静的ファイルのルーティング
	e.Static("/", "web")

	// ルーティング
	e.GET("/books", handler.GetBookMetaInfoForUser)
	e.GET("/books/all", handler.GetBookMetaInfoAll)
	e.GET("/books/:ISBN", handler.GetBookProfile)
	e.GET("/books/:ISBN/threads", handler.GetThreadTitles)
	e.GET("/threads/:threadID", handler.GetThreadMessages)
	//e.GET("/users/:userID/books", handler.GetBookMetaInfoForUser)
	e.POST("/books", handler.PostBookInfo)
	e.POST("/books/:ISBN/threads", handler.PostThreadTitle)
	e.POST("/threads/:threadID", handler.PostThreadMessage)
	e.POST("/users/registration", handler.RegisterUser)
	e.POST("/users/login", handler.AuthenticateUser)

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
