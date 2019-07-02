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

	// フォーラムメタ情報
	forumMetaInfo struct {
		ID     int    `json:"id"`
		UserID int    `json:"userID"`
		Title  string `json:"title"`
		ISBN   int    `json:"ISBN"`
	}

	// フォーラム発言情報
	forumMessages struct {
		ID      int    `json:"id"`
		UserID  int    `json:"userID"`
		Message string `json:"message"`
		ForumID int    `json:"forumID"`
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

	tmpForumMeta1 = forumMetaInfo{
		ID:     1,
		UserID: 1,
		Title:  "I don't understand p.32 at all.",
		ISBN:   100,
	}

	tmpForumMeta2 = forumMetaInfo{
		ID:     2,
		UserID: 2,
		Title:  "there is an awful typo on p.55",
		ISBN:   100,
	}

	// フォーラムのメタ情報格納用配列　（そのうちデータベースに移行）
	forumMetaInfoDataBase = []forumMetaInfo{
		tmpForumMeta1,
		tmpForumMeta2,
	}

	tmpforumMessage1 = forumMessages{
		ID:      1,
		UserID:  11,
		Message: "Me neither.",
		ForumID: 1,
	}

	tmpforumMessage2 = forumMessages{
		ID:      2,
		UserID:  12,
		Message: "I think the author tries to say ...",
		ForumID: 1,
	}

	// フォーラムのメッセージ情報格納用配列　（そのうちデータベースに移行）
	forumMessagesDataBase = []forumMessages{
		tmpforumMessage1,
		tmpforumMessage2,
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

// GetForumTitles 本の詳細ページに表示するために使う、フォーラムのタイトル取得用メソッド
func GetForumTitles(c echo.Context) error {

	// urlのisbn取得
	isbn, err := strconv.Atoi(c.Param("ISBN"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ISBN must be an integer")
	}

	// 本データベースに該当のISBNの本が登録されているか確認
	for _, b := range bookDataBase {
		if isbn == b.ISBN {
			message := []forumMetaInfo{}

			// 該当のISBNに対応するフォーラムタイトルを検索
			for _, f := range forumMetaInfoDataBase {
				if b.ISBN == f.ISBN {
					message = append(message, f)
				}
			}
			return c.JSON(http.StatusOK, message)
		}
	}
	return c.String(http.StatusNotFound, "Not Found")
}

// GetForumMessages 本の詳細ページに表示するために使う、フォーラムのタイトル取得用メソッド
func GetForumMessages(c echo.Context) error {

	// urlのisbn取得
	forumID, err := strconv.Atoi(c.Param("forumID"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "forumID must be an integer")
	}

	// フォーラムメタ情報データベースに該当のforumIDをもつものが登録されているか確認
	for _, f := range forumMetaInfoDataBase {
		if forumID == f.ID {
			message := []forumMessages{}
			// 該当のforumIDに対応するメッセージを検索
			for _, m := range forumMessagesDataBase {
				if forumID == m.ForumID {
					message = append(message, m)
				}
			}

			return c.JSON(http.StatusOK, message)
		}
	}
	return c.String(http.StatusNotFound, "Not Found")
}
