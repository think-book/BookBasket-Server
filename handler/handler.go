package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (

	// 本情報用構造体（POST用）
	bookInfo struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ISBN        int    `json:"ISBN"`
	}

	// 本メタ情報用構造体（GET用）
	bookMetaInfo struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		ISBN  int    `json:"ISBN"`
	}

	// 本詳細情報用構造体（GET用）
	bookProfileInfo struct {
		ISBN        int    `json:"ISBN"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
)

var (
	tmpData1 = bookInfo{
		ID:          1,
		Title:       "cool book",
		Description: "A super hero beats monsters.",
		ISBN:        100,
	}

	tmpData2 = bookInfo{
		ID:          2,
		Title:       "awesome book",
		Description: "A text book of go langage.",
		ISBN:        200,
	}

	// 本情報格納用配列　（そのうちデータベースに移行）
	bookDataBase = []bookInfo{
		tmpData1,
		tmpData2,
	}
)

// GetBookMetaInfoAll 本情報全取得
func GetBookMetaInfoAll(c echo.Context) error { //c をいじって Request, Responseを色々する

	// message（bookMetaInfo配列） にメタ情報を順次格納していく
	message := []bookMetaInfo{}

	for _, m := range bookDataBase {
		tmp := bookMetaInfo{
			ID:    m.ID,
			Title: m.Title,
			ISBN:  m.ISBN,
		}

		message = append(message, tmp)
	}

	return c.JSON(http.StatusOK, message)
}

//GetBookProfile 本情報１件取得
func GetBookProfile(c echo.Context) error {
	// urlのisbn取得
	isbn, err := strconv.Atoi(c.Param("ISBN"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ISBN must be an integer")
	}

	for _, b := range bookDataBase {
		if isbn == b.ISBN {
			tmp := bookProfileInfo{
				Title:       b.Title,
				ISBN:        b.ISBN,
				Description: b.Description,
			}
			return c.JSON(http.StatusOK, tmp)
		}
	}

	return c.String(http.StatusNotFound, "Not Found")

}

// PostBookInfo メタ情報Post用メソッド
func PostBookInfo(c echo.Context) error {
	info := new(bookInfo)

	// request bodyをmetaInfo構造体にバインド
	if err := c.Bind(info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if info.Title == "" || info.ISBN == 0 || info.Description == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// メタ情報が既に登録ずみならBad request
	for _, b := range bookDataBase {
		if info.ISBN == b.ISBN {
			return c.String(http.StatusBadRequest, "Book info already exists")
		}
	}

	id := bookDataBase[len(bookDataBase)-1].ID + 1

	info.ID = id

	bookDataBase = append(bookDataBase, *info)

	return c.JSON(http.StatusOK, info)
}
