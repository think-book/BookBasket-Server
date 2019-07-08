package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (

	// 本情報用構造体（POST、データ保存用）
	BookInfo struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ISBN        int    `json:"ISBN"`
	}

	// 本メタ情報用構造体（GET用）
	BookMetaInfo struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		ISBN  int    `json:"ISBN"`
	}

	// 本詳細情報用構造体（GET用）
	BookProfileInfo struct {
		ISBN        int    `json:"ISBN"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	// スレッドメタ情報
	ThreadMetaInfo struct {
		ID     int    `json:"id"`
		UserID int    `json:"userID"`
		Title  string `json:"title"`
		ISBN   int    `json:"ISBN"`
	}

	// スレッド発言情報
	ThreadMessage struct {
		ID       int    `json:"id"`
		UserID   int    `json:"userID"`
		Message  string `json:"message"`
		ThreadID int    `json:"threadID"`
	}

	// ユーザ情報
	UserInfo struct {
		ID       int    `json:"id"`
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
)

var (
	tmpData1 = BookInfo{
		ID:          1,
		Title:       "cool book",
		Description: "A super hero beats monsters.",
		ISBN:        100,
	}

	tmpData2 = BookInfo{
		ID:          2,
		Title:       "awesome book",
		Description: "A text book of go langage.",
		ISBN:        200,
	}

	// 本情報格納用配列　（そのうちデータベースに移行）
	bookDataBase = []BookInfo{
		tmpData1,
		tmpData2,
	}

	tmpThreadMeta1 = ThreadMetaInfo{
		ID:     1,
		UserID: 1,
		Title:  "I don't understand p.32 at all.",
		ISBN:   100,
	}

	tmpThreadMeta2 = ThreadMetaInfo{
		ID:     2,
		UserID: 2,
		Title:  "there is an awful typo on p.55",
		ISBN:   100,
	}

	// スレッドのメタ情報格納用配列　（そのうちデータベースに移行）
	threadMetaInfoDataBase = []ThreadMetaInfo{
		tmpThreadMeta1,
		tmpThreadMeta2,
	}

	tmpThreadMessage1 = ThreadMessage{
		ID:       1,
		UserID:   11,
		Message:  "Me neither.",
		ThreadID: 1,
	}

	tmpThreadMessage2 = ThreadMessage{
		ID:       2,
		UserID:   12,
		Message:  "I think the author tries to say ...",
		ThreadID: 1,
	}

	// スレッドのメッセージ情報格納用配列　（そのうちデータベースに移行）
	threadMessagesDataBase = []ThreadMessage{
		tmpThreadMessage1,
		tmpThreadMessage2,
	}

	tmpUser1 = UserInfo{
		ID:       1,
		UserName: "Alice",
		Password: "pass",
	}

	tmpUser2 = UserInfo{
		ID:       2,
		UserName: "Bob",
		Password: "word",
	}

	tmpUser3 = UserInfo{
		ID:       11,
		UserName: "Carol",
		Password: "qwer",
	}

	tmpUser4 = UserInfo{
		ID:       12,
		UserName: "Charlie",
		Password: "tyui",
	}

	// ユーザ情報格納用配列（そのうちデータベースに以降）
	userInfoDataBase = []UserInfo{
		tmpUser1,
		tmpUser2,
		tmpUser3,
		tmpUser4,
	}
)

// GetBookMetaInfoAll 本情報全取得
func GetBookMetaInfoAll(c echo.Context) error { //c をいじって Request, Responseを色々する

	// message（bookMetaInfo配列） にメタ情報を順次格納していく
	message := []BookMetaInfo{}

	for _, m := range bookDataBase {
		tmp := BookMetaInfo{
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
			tmp := BookProfileInfo{
				Title:       b.Title,
				ISBN:        b.ISBN,
				Description: b.Description,
			}
			return c.JSON(http.StatusOK, tmp)
		}
	}

	return c.String(http.StatusNotFound, "Not Found")

}

// PostBookInfo 本情報Post用メソッド
func PostBookInfo(c echo.Context) error {
	info := new(BookInfo)

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

// GetThreadTitles 本の詳細ページに表示するために使う、スレッドのタイトル取得用メソッド
func GetThreadTitles(c echo.Context) error {

	// urlのisbn取得
	isbn, err := strconv.Atoi(c.Param("ISBN"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ISBN must be an integer")
	}

	// 本データベースに該当のISBNの本が登録されているか確認
	for _, b := range bookDataBase {
		if isbn == b.ISBN {
			message := []ThreadMetaInfo{}

			// 該当のISBNに対応するスレッドタイトルを検索
			for _, f := range threadMetaInfoDataBase {
				if b.ISBN == f.ISBN {
					message = append(message, f)
				}
			}
			return c.JSON(http.StatusOK, message)
		}
	}
	return c.String(http.StatusNotFound, "Not Found")
}

// GetThreadMessages 各スレッドのメッセージ取得用メソッド
func GetThreadMessages(c echo.Context) error {

	// urlのisbn取得
	threadID, err := strconv.Atoi(c.Param("threadID"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ThreadID must be an integer")
	}

	// スレッドメタ情報データベースに該当のthreadIDをもつものが登録されているか確認
	for _, f := range threadMetaInfoDataBase {
		if threadID == f.ID {
			message := []ThreadMessage{}
			// 該当のthreadIDに対応するメッセージを検索
			for _, m := range threadMessagesDataBase {
				if threadID == m.ThreadID {
					message = append(message, m)
				}
			}

			return c.JSON(http.StatusOK, message)
		}
	}
	return c.String(http.StatusNotFound, "Not Found")
}

// PostThreadTitle スレッドタイトルPost用メソッド
func PostThreadTitle(c echo.Context) error {
	info := new(ThreadMetaInfo)

	// urlのisbn取得
	isbn, err := strconv.Atoi(c.Param("ISBN"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ISBN must be an integer")
	}

	// request bodyをThreadMetaInfo構造体にバインド
	if err := c.Bind(info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if info.Title == "" || info.UserID == 0 {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ISBNがデータベースにあるか確認
	isbnExists := false
	for _, b := range bookDataBase {
		if isbn == b.ISBN {
			isbnExists = true
			break
		}
	}
	// ISBNが存在しなければBad request
	if !isbnExists {
		return c.String(http.StatusBadRequest, "Book doesn't exist")
	}

	// userIDがデータベースにあるか確認
	userExists := false
	for _, u := range userInfoDataBase {
		if info.UserID == u.ID {
			userExists = true
			break
		}
	}
	// ユーザが存在しなければBad request
	if !userExists {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// スレッドのISBN設定
	info.ISBN = isbn

	// スレッドタイトル情報が既に登録ずみならBad request
	for _, b := range threadMetaInfoDataBase {
		if info.ISBN == b.ISBN && info.Title == b.Title {
			return c.String(http.StatusBadRequest, "Thread title already exists")
		}
	}

	id := threadMetaInfoDataBase[len(threadMetaInfoDataBase)-1].ID + 1

	info.ID = id

	threadMetaInfoDataBase = append(threadMetaInfoDataBase, *info)

	return c.JSON(http.StatusOK, info)
}

// PostThreadMessage スレッドタイトルPost用メソッド
func PostThreadMessage(c echo.Context) error {
	info := new(ThreadMessage)

	// urlのthreadID取得
	threadID, err := strconv.Atoi(c.Param("threadID"))
	if err != nil {
		// threadIDがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ThreadID must be an integer")
	}

	// request bodyをThreadMessage構造体にバインド
	if err := c.Bind(info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if info.UserID == 0 || info.Message == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// threadIDがデータベースにあるか確認
	threadExists := false
	for _, t := range threadMetaInfoDataBase {
		if threadID == t.ID {
			threadExists = true
			break
		}
	}
	// threadIDが存在しなければBad request
	if !threadExists {
		return c.String(http.StatusBadRequest, "Thread doesn't exist")
	}

	// userIDがデータベースにあるか確認
	userExists := false
	for _, u := range userInfoDataBase {
		if info.UserID == u.ID {
			userExists = true
			break
		}
	}
	// ユーザが存在しなければBad request
	if !userExists {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// メッセージのthreadID設定
	info.ThreadID = threadID

	id := threadMessagesDataBase[len(threadMessagesDataBase)-1].ID + 1

	info.ID = id

	threadMessagesDataBase = append(threadMessagesDataBase, *info)

	return c.JSON(http.StatusOK, info)

}
