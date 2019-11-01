package handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type (

	// 本情報用構造体（POST、データ保存用）
	BookInfo struct {
		ISBN        uint64 `json:"ISBN" db:"ISBN"`
		Title       string `json:"title" db:"title"`
		Description string `json:"description" db:"description"`
	}

	// 本メタ情報用構造体（GET用）
	BookMetaInfo struct {
		ISBN  uint64 `json:"ISBN" db:"ISBN"`
		Title string `json:"title" db:"title"`
	}

	// 本詳細情報用構造体（GET用）
	BookProfileInfo struct {
		ISBN        uint64 `json:"ISBN" db:"ISBN"`
		Title       string `json:"title" db:"title"`
		Description string `json:"description" db:"description"`
	}

	// スレッドメタ情報
	ThreadMetaInfo struct {
		ID       int    `json:"id" db:"id"`
		UserName string `json:"userName" db:"userName"`
		Title    string `json:"title" db:"title"`
		ISBN     uint64 `json:"ISBN" db:"ISBN"`
	}

	// スレッド発言情報
	ThreadMessage struct {
		UserName string `json:"userName" db:"userName"`
		Message  string `json:"message" db:"message"`
		ThreadID int    `json:"threadID" db:"threadID"`
	}

	// ユーザ情報（登録用）
	UserInfo struct {
		UserName string `json:"userName" db:"userName"`
		Password string `json:"password" db:"password"`
	}

	// ユーザ情報（返信用）
	UserInfoForReturn struct {
		ID       int    `json:"id" db:"id"`
		UserName string `json:"userName" db:"userName"`
	}
)

var (
	//　データベースへの参照
	db *sqlx.DB
)

// SetDB データベースへの参照をセット
func SetDB(d *sqlx.DB) {
	db = d
}

// GetBookMetaInfoAll 本情報全取得
func GetBookMetaInfoAll(c echo.Context) error { //c をいじって Request, Responseを色々する

	// message（bookMetaInfo配列） にメタ情報を格納
	message := []BookMetaInfo{}

	//全件取得クエリ messageに結果をバインド
	err := db.Select(&message, "SELECT ISBN, title FROM bookInfo")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
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

	var profile BookProfileInfo
	// 一件取得用クエリ　profileに結果をバインド
	err = db.Get(&profile, "SELECT * FROM bookInfo WHERE ISBN=?", isbn)
	if err != nil {
		return c.String(http.StatusNotFound, "Not Found")
	}

	return c.JSON(http.StatusOK, profile)

}

// PostBookInfo 本情報Post用メソッド
func PostBookInfo(c echo.Context) error {
	var info BookInfo

	// request bodyをBookInfo構造体にバインド
	if err := c.Bind(&info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if info.Title == "" || info.ISBN == 0 || info.Description == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// 一件挿入用クエリ
	_, err := db.Exec("INSERT INTO bookInfo (ISBN, title, description) VALUES(?,?,?)", info.ISBN, info.Title, info.Description)
	// PRIMARY KEY(ISBN)がすでに存在した時（を想定）
	if err != nil {
		return c.String(http.StatusBadRequest, "Book info already exists")
	}

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

	var profile BookProfileInfo
	// 本データベースに該当のISBNの本が登録されているか確認
	err = db.Get(&profile, "SELECT * FROM bookInfo WHERE ISBN=?", isbn)
	if err != nil {
		return c.String(http.StatusNotFound, "Not Found")
	}

	// message（ThreadMetaInfo配列） にスレッド情報を格納
	message := []ThreadMetaInfo{}

	//全件取得クエリ messageに結果をバインド
	err = db.Select(&message, "SELECT * FROM threadMetaInfo WHERE ISBN=?", isbn)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, message)

}

// GetThreadMessages 各スレッドのメッセージ取得用メソッド
func GetThreadMessages(c echo.Context) error {

	// urlのisbn取得
	threadID, err := strconv.Atoi(c.Param("threadID"))
	if err != nil {
		// ISBNがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "ThreadID must be an integer")
	}

	var threadMeta ThreadMetaInfo
	// スレッドメタ情報データベースに該当のthreadIDをもつものが登録されているか確認
	err = db.Get(&threadMeta, "SELECT userName, title, ISBN FROM threadMetaInfo WHERE id=?", threadID)
	if err != nil {
		return c.String(http.StatusNotFound, "Not Found")
	}

	// message（ThreadMetaInfo配列） にスレッド情報を格納
	message := []ThreadMessage{}

	//全件取得クエリ messageに結果をバインド
	err = db.Select(&message, "SELECT userName, message, threadID FROM threadMessage WHERE threadID=?", threadID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, message)

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
	if info.Title == "" || info.UserName == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	var book BookInfo
	// ISBNがデータベースにあるか確認
	err = db.Get(&book, "SELECT * FROM bookInfo WHERE ISBN=?", isbn)
	// ISBNが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "Book doesn't exist")
	}

	var user UserInfoForReturn
	// userNameがデータベースにあるか確認
	err = db.Get(&user, "SELECT id, userName FROM userInfo WHERE userName=?", info.UserName)
	// ユーザが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// スレッドのISBN設定
	info.ISBN = uint64(isbn)

	// 一件挿入用クエリ
	_, err = db.Exec("INSERT INTO threadMetaInfo (userName, title, ISBN) VALUES(?,?,?)", info.UserName, info.Title, info.ISBN)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// 挿入したレコードのid取得
	var id int
	err = db.Get(&id, "SELECT LAST_INSERT_ID()")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	info.ID = id

	return c.JSON(http.StatusOK, info)
}

// PostThreadMessage スレッドメッセージPost用メソッド
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
	if info.UserName == "" || info.Message == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	var threadMeta ThreadMetaInfo
	// threadIDがデータベースにあるか確認
	err = db.Get(&threadMeta, "SELECT userName, title, ISBN FROM threadMetaInfo WHERE id=?", threadID)
	// threadIDが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "Thread doesn't exist")
	}

	var user UserInfo
	// userNameがデータベースにあるか確認
	err = db.Get(&user, "SELECT userName FROM userInfo WHERE userName=?", info.UserName)
	// ユーザが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// メッセージのthreadID設定
	info.ThreadID = threadID

	// 一件挿入用クエリ
	_, err = db.Exec("INSERT INTO threadMessage (userName, message, threadID) VALUES(?,?,?)", info.UserName, info.Message, info.ThreadID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, info)

}

// RegisterUser ユーザ登録用Post用メソッド
func RegisterUser(c echo.Context) error {
	info := new(UserInfo)

	// request bodyをUserInfo構造体にバインド
	if err := c.Bind(info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// 正規表現でユーザ名チェック
	userRe, err := regexp.Compile(`^[a-zA-Z0-9_\-.]{3,15}$`)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// 正規表現でパスワードチェック
	passRe, err := regexp.Compile(`^[a-zA-Z0-9_\-.!]{8,15}$`)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// ポストメッセージのフォーマットが不正
	if !userRe.Match([]byte(info.UserName)) || !passRe.Match([]byte(info.Password)) {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	var user UserInfo
	// userNameがデータベースにあるか確認
	err = db.Get(&user, "SELECT userName FROM userInfo WHERE userName=?", info.UserName)
	// ユーザが存在すればBadRequest
	if err == nil {
		return c.String(http.StatusBadRequest, "User already exists")
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	// 一件挿入用クエリ
	_, err = db.Exec("INSERT INTO userInfo (userName, password) VALUES(?,?)", info.UserName, string(newPassword))
	// userName衝突
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	var registeredUser UserInfoForReturn
	// 登録されたユーザ取得
	err = db.Get(&registeredUser, "SELECT id, userName FROM userInfo WHERE userName=?", info.UserName)
	// エラーが起きたとき
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, registeredUser)

}

// AuthenticateUser ユーザ認証用Post用メソッド
func AuthenticateUser(c echo.Context) error {
	info := new(UserInfo)

	// request bodyをUserInfo構造体にバインド
	if err := c.Bind(info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// 正規表現でユーザ名チェック
	userRe, err := regexp.Compile(`^[a-zA-Z0-9_\-.]{3,15}$`)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// 正規表現でパスワードチェック
	passRe, err := regexp.Compile(`^[a-zA-Z0-9_\-.!]{8,15}$`)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// ポストメッセージのフォーマットが不正
	if !userRe.Match([]byte(info.UserName)) || !passRe.Match([]byte(info.Password)) {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ユーザが存在しなければ、ダミー初期値を使用
	userName := "!!Security!!"
	password := "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	var user UserInfo
	// userNameがデータベースにあるか確認
	err = db.Get(&user, "SELECT userName, password FROM userInfo WHERE userName=?", info.UserName)
	// ユーザが存在すれば、データをとってくる
	if err == nil {
		userName = user.UserName
		password = user.Password
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(info.Password))
	// ログイン失敗
	if err != nil || userName == "!!Security!!" {
		return c.String(http.StatusBadRequest, "Login Failed")
	}

	var loginedUser UserInfoForReturn
	// 登録されたユーザ取得
	err = db.Get(&loginedUser, "SELECT id, userName FROM userInfo WHERE userName=?", info.UserName)
	// エラーが起きたとき
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, loginedUser)

}
