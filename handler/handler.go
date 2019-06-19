package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

// 本情報格納用配列
var info = []string{
	"{\"id\": 1, \"title\": \"cool book\", \"memo\": \"foo\"}",
	"{\"id\": 2, \"title\": \"awesome book\", \"memo\": \"bar\"}",
}

// GetBookInfoAll 本情報全取得
func GetBookInfoAll() echo.HandlerFunc {
	return func(c echo.Context) error { //c をいじって Request, Responseを色々する

		// message にinfoを順次ぶち込んでいく
		message := ""
		for i, s := range info {
			message += s
			if i != len(info)-1 {
				message += ", \n"
			}
		}
		message += "\n"

		return c.String(http.StatusOK, message)
	}
}

// GetBookInfo 本情報１件取得
func GetBookInfo() echo.HandlerFunc {
	return func(c echo.Context) error {

		// urlのid取得
		id, _ := strconv.Atoi(c.Param("id"))

		if id > len(info) || id <= 0 {
			return c.String(http.StatusNotFound, "Not Found\n")
		}
		return c.String(http.StatusOK, info[id-1]+"\n")

	}
}
