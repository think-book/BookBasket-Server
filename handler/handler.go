package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (

	// 本のメタ情報用構造体
	metaInfo struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		ISBN  int    `json:"ISBN"`
	}

	// 本の詳細情報用構造体
	bookProfile struct {
		ISBN  int    `json:"ISBN"`
		Title string `json:"title"`
		Story string `json:"story"`
	}
)

var (
	tmpMeta1 = metaInfo{
		Id:    1,
		Title: "cool book",
		ISBN:  100,
	}

	tmpMeta2 = metaInfo{
		Id:    2,
		Title: "awesome book",
		ISBN:  200,
	}

	// 本のメタ情報格納用配列　（そのうちデータベースに移行）
	metaInfoDataBase = []metaInfo{
		tmpMeta1,
		tmpMeta2,
	}

	tmpProfile1 = bookProfile{
		ISBN:  100,
		Title: "cool book",
		Story: "A super hero beats monsters.",
	}

	tmpProfile2 = bookProfile{
		ISBN:  200,
		Title: "awesome book",
		Story: "A text book of go langage.",
	}

	// 本の詳細情報格納用配列　（そのうちデータベースに移行）
	bookProfileDataBase = []bookProfile{
		tmpProfile1,
		tmpProfile2,
	}

	foo = []int{
		1,
		2,
	}
)

// GetBookInfoAll 本情報全取得
func GetBookMetaInfoAll(c echo.Context) error { //c をいじって Request, Responseを色々する

	// message にinfoを順次ぶち込んでいく
	message := ""

	for i, m := range metaInfoDataBase {
		//構造体をjsonのバイナリに変換
		jsonBinary, _ := json.Marshal(m)

		message += string(jsonBinary)

		if i != len(metaInfoDataBase)-1 {
			message += ",\n"
		}
	}

	return c.String(http.StatusOK, message)
}

// GetBookInfo 本情報１件取得
func GetBookProfile(c echo.Context) error {
	// urlのid取得
	isbn, _ := strconv.Atoi(c.Param("ISBN"))

	for _, b := range bookProfileDataBase {
		if isbn == b.ISBN {
			//構造体をjsonのバイナリに変換
			jsonBinary, _ := json.Marshal(b)

			return c.String(http.StatusOK, string(jsonBinary))
		}
	}

	return c.String(http.StatusNotFound, "Not Found")

}

func PostMetaInfo(c echo.Context) error {
	meta := new(metaInfo)

	if err := c.Bind(meta); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	id := metaInfoDataBase[len(metaInfoDataBase)-1].Id + 1

	meta.Id = id

	//構造体をjsonのバイナリに変換
	jsonBinary, _ := json.Marshal(meta)

	message := string(jsonBinary)
	metaInfoDataBase = append(metaInfoDataBase, *meta)

	return c.String(http.StatusOK, message)
}
