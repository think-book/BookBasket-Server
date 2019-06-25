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
		ID    int    `json:"id"`
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
		ID:    1,
		Title: "cool book",
		ISBN:  100,
	}

	tmpMeta2 = metaInfo{
		ID:    2,
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
)

// GetBookMetaInfoAll 本情報全取得
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

//GetBookProfile 本情報１件取得
func GetBookProfile(c echo.Context) error {
	// urlのisbn取得
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

// PostMetaInfo メタ情報Post用メソッド
func PostMetaInfo(c echo.Context) error {
	meta := new(metaInfo)

	// request bodyをmetaInfo構造体にバインド
	if err := c.Bind(meta); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if meta.Title == "" || meta.ISBN == 0 {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// メタ情報が既に登録ずみならBad request
	for _, m := range metaInfoDataBase {
		if meta.ISBN == m.ISBN {
			return c.String(http.StatusBadRequest, "Meta info already exists")
		}
	}

	id := metaInfoDataBase[len(metaInfoDataBase)-1].ID + 1

	meta.ID = id

	//構造体をjsonのバイナリに変換
	jsonBinary, _ := json.Marshal(meta)

	message := string(jsonBinary)
	metaInfoDataBase = append(metaInfoDataBase, *meta)

	return c.String(http.StatusOK, message)
}

// PostBookProfile 詳細情報Post用メソッド
func PostBookProfile(c echo.Context) error {
	// urlのisbn取得
	isbn, err := strconv.Atoi(c.Param("ISBN"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ISBN must be an integer")
	}

	profile := new(bookProfile)

	// request bodyをbookProfile構造体にバインド
	if err := c.Bind(profile); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if profile.Title == "" || profile.Story == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// urlとpostデータのISBNが一致していることを確認
	if isbn != profile.ISBN {
		return c.String(http.StatusBadRequest, "ISBN is inconsistent")
	}

	// 詳細情報が既に登録ずみならBad request
	for _, p := range bookProfileDataBase {
		if profile.ISBN == p.ISBN {
			return c.String(http.StatusBadRequest, "Book profile already exists")
		}
	}

	// メタ情報が登録されていることを確認
	for _, b := range metaInfoDataBase {
		if isbn == b.ISBN {

			//構造体をjsonのバイナリに変換
			jsonBinary, _ := json.Marshal(profile)

			message := string(jsonBinary)
			bookProfileDataBase = append(bookProfileDataBase, *profile)

			return c.String(http.StatusOK, message)
		}
	}

	return c.String(http.StatusNotFound, "Book Meta Data Not Found")

}
