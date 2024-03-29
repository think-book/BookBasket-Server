package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (

	// 本情報用構造体（POST、データ保存用）
	BookInfo struct {
		ISBN        uint64 `json:"ISBN" db:"ISBN"`
		Title       string `json:"title" db:"title"`
		Description string `json:"description" db:"description"`
	}

	// 本とユーザの関係用構造体
	UserBookRelation struct {
		UserID int `db:"userID"`
		ISBN   int `db:"ISBN"`
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

// GetBookMetaInfoForUser ユーザの本情報全取得
func GetBookMetaInfoForUser(c echo.Context) error {

	// message（bookMetaInfo配列） にメタ情報を格納
	message := []BookMetaInfo{}

	//sessionを見る
	sess, _ := session.Get("session", c)
	var userID int
	var err error

	//ログインしているか
	if b := sess.Values["auth"]; b != true {
		return c.String(http.StatusUnauthorized, "Not Logined")
	} else {
		userID, _ = sess.Values["userID"].(int)
	}

	var user UserInfoForReturn
	// userIDがデータベースにあるか確認
	err = db.Get(&user, "SELECT id, userName FROM userInfo WHERE id=?", userID)
	// ユーザが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// ユーザと本の関係を入れる用
	relation := []UserBookRelation{}

	//userIDのユーザが登録している本のISBN全件取得クエリ relationに結果をバインド
	err = db.Select(&relation, "SELECT * FROM userBookRelation WHERE userID=?", userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// userIDのユーザが本を一冊も登録していなかったとき（[]を返す）
	if len(relation) == 0 {
		return c.JSON(http.StatusOK, message)
	}

	// relationからISBNだけ抜き取る
	ISBNs := []int{}
	for _, r := range relation {
		ISBNs = append(ISBNs, r.ISBN)
	}

	//本取得クエリを生成するための処理
	//query: where inを含んだ新しいクエリ
	//args : 引数
	query, args, err := sqlx.In("SELECT ISBN, title FROM bookInfo WHERE ISBN IN (?)", ISBNs)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	// mysql用にクエリをリバインド？（っぽい）
	query = db.Rebind(query)

	//全件取得クエリ messageに結果をバインド
	err = db.Select(&message, query, args...)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, message)
}

// GetBookMetaInfoForOtherUser 指定した他のユーザの本取得
func GetBookMetaInfoForOtherUser(c echo.Context) error {
	// message（bookMetaInfo配列） にメタ情報を格納
	message := []BookMetaInfo{}

	// urlのユーザid取得
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// idがintでなければBadRequestを返す
		return c.String(http.StatusBadRequest, "User id must be an integer")
	}

	var user UserInfoForReturn
	// userIDがデータベースにあるか確認
	err = db.Get(&user, "SELECT id, userName FROM userInfo WHERE id=?", userID)
	// ユーザが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// ユーザと本の関係を入れる用
	relation := []UserBookRelation{}

	//userIDのユーザが登録している本のISBN全件取得クエリ relationに結果をバインド
	err = db.Select(&relation, "SELECT * FROM userBookRelation WHERE userID=?", userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// userIDのユーザが本を一冊も登録していなかったとき（[]を返す）
	if len(relation) == 0 {
		return c.JSON(http.StatusOK, message)
	}

	// relationからISBNだけ抜き取る
	ISBNs := []int{}
	for _, r := range relation {
		ISBNs = append(ISBNs, r.ISBN)
	}

	//本取得クエリを生成するための処理
	//query: where inを含んだ新しいクエリ
	//args : 引数
	query, args, err := sqlx.In("SELECT ISBN, title FROM bookInfo WHERE ISBN IN (?)", ISBNs)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	// mysql用にクエリをリバインド？（っぽい）
	query = db.Rebind(query)

	//全件取得クエリ messageに結果をバインド
	err = db.Select(&message, query, args...)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, message)
}

// GetUserLists ユーザリスト登録
func GetUserLists(c echo.Context) error {

	//sessionを見る
	sess, _ := session.Get("session", c)
	var userID int
	var err error

	//ログインしているか
	if b := sess.Values["auth"]; b != true {
		return c.String(http.StatusUnauthorized, "Not Logined")
	} else {
		userID, _ = sess.Values["userID"].(int)
	}

	var user UserInfoForReturn
	// userIDがデータベースにあるか確認
	err = db.Get(&user, "SELECT id, userName FROM userInfo WHERE id=?", userID)
	// ユーザが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// userlists（UserInfoForReturn配列） にメタ情報を格納
	userlists := []UserInfoForReturn{}

	//ユーザ情報獲得クエリ usersに結果をバインド
	err = db.Select(&userlists, "SELECT id, userName FROM userInfo")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, userlists)

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

	//sessionを見る
	sess, _ := session.Get("session", c)
	var userID int
	var err error

	//ログインしているか
	if b := sess.Values["auth"]; b != true {
		return c.String(http.StatusUnauthorized, "Not Logined")
	} else {
		userID, _ = sess.Values["userID"].(int)
	}

	// request bodyをBookInfo構造体にバインド
	if err := c.Bind(&info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if info.Title == "" || info.ISBN == 0 || info.Description == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// 一件挿入用クエリ（ユーザと本の関係）
	_, err = db.Exec("INSERT INTO userBookRelation (userID, ISBN) VALUES(?,?)", userID, info.ISBN)
	// PRIMARY KEY(userID, ISBN)がすでに存在した時（を想定）
	if err != nil {
		return c.String(http.StatusBadRequest, "Book has already been registerd")
	}

	// 一件挿入用クエリ（グローバルな本棚）
	_, err = db.Exec("INSERT INTO bookInfo (ISBN, title, description) VALUES(?,?,?)", info.ISBN, info.Title, info.Description)
	// PRIMARY KEY(ISBN)がすでに存在した時（を想定）。ユーザにとっては初めての登録のため、StatusOK
	if err != nil {
		return c.JSON(http.StatusOK, info)
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

	//sessionを見る
	sess, _ := session.Get("session", c)
	var userID int

	//ログインしているか
	if b := sess.Values["auth"]; b != true {
		return c.String(http.StatusUnauthorized, "Not Logined")
	} else {
		userID, _ = sess.Values["userID"].(int)
	}

	// request bodyをThreadMetaInfo構造体にバインド
	if err := c.Bind(info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if info.Title == "" {
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
	// userIDがデータベースにあるか確認
	err = db.Get(&user, "SELECT id, userName FROM userInfo WHERE id=?", userID)
	// ユーザが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// スレッドのISBN設定
	info.ISBN = uint64(isbn)

	// 一件挿入用クエリ
	_, err = db.Exec("INSERT INTO threadMetaInfo (userName, title, ISBN) VALUES(?,?,?)", user.UserName, info.Title, info.ISBN)
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
	info.UserName = user.UserName

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

	//sessionを見る
	sess, _ := session.Get("session", c)
	var userID int

	//ログインしているか
	if b := sess.Values["auth"]; b != true {
		return c.String(http.StatusUnauthorized, "Not Logined")
	} else {
		userID, _ = sess.Values["userID"].(int)
	}

	// request bodyをThreadMessage構造体にバインド
	if err := c.Bind(info); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	// ポストメッセージのフォーマットが不正
	if info.Message == "" {
		return c.String(http.StatusBadRequest, "Invalid Post Format")
	}

	var threadMeta ThreadMetaInfo
	// threadIDがデータベースにあるか確認
	err = db.Get(&threadMeta, "SELECT userName, title, ISBN FROM threadMetaInfo WHERE id=?", threadID)
	// threadIDが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "Thread doesn't exist")
	}

	var user UserInfoForReturn
	// userIDがデータベースにあるか確認
	err = db.Get(&user, "SELECT id, userName FROM userInfo WHERE id=?", userID)
	// ユーザが存在しなければBad request
	if err != nil {
		return c.String(http.StatusBadRequest, "User doesn't exist")
	}

	// メッセージのthreadID設定
	info.ThreadID = threadID
	info.UserName = user.UserName

	// 一件挿入用クエリ
	_, err = db.Exec("INSERT INTO threadMessage (userName, message, threadID) VALUES(?,?,?)", user.UserName, info.Message, info.ThreadID)
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
	userRe, err := regexp.Compile(`^[a-zA-Z0-9_\-\.]{3,15}$`)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// 正規表現でパスワードチェック
	passRe, err := regexp.Compile(`^[a-zA-Z0-9_\-\.!]{8,15}$`)
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
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	var registeredUser UserInfoForReturn
	// 登録されたユーザ取得
	err = db.Get(&registeredUser, "SELECT id, userName FROM userInfo WHERE userName=?", info.UserName)
	// エラーが起きたとき
	if err != nil {

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["auth"] = true
	sess.Values["userID"] = registeredUser.ID
	sess.Values["userName"] = registeredUser.UserName
	err = sess.Save(c.Request(), c.Response())
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
	userRe, err := regexp.Compile(`^[a-zA-Z0-9_\-\.]{3,15}$`)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// 正規表現でパスワードチェック
	passRe, err := regexp.Compile(`^[a-zA-Z0-9_\-\.!]{8,15}$`)
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

	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["auth"] = true
	sess.Values["userID"] = loginedUser.ID
	sess.Values["userName"] = loginedUser.UserName
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, loginedUser)
}
