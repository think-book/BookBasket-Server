package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (

	// GET用

	// GetBookMetaInfoAll用確認データ
	metaInfoTestData = `[{"ISBN":100,"title":"cool book"},{"ISBN":200,"title":"awesome book"}]
`

	// GetBookProfile用確認データ
	bookProfileTestData = `{"ISBN":100,"title":"cool book","description":"A super hero beats monsters."}
`

	// GETThreadTitles用確認データ
	threadTitlesTestData = `[{"id":1,"userName":"Alice","title":"I don't understand p.32 at all.","ISBN":100},{"id":2,"userName":"Bob","title":"there is an awful typo on p.55","ISBN":100}]
`
	// 空配列確認データ
	emptyData = `[]
`
	// GETThreadMessages用確認データ
	threadMessagesTestData = `[{"userName":"Carol","message":"Me neither.","threadID":1},{"userName":"Charlie","message":"I think the author tries to say ...","threadID":1}]
`

	// POST用

	// POST送信用本情報
	bookInfoForPost = `{"title":"epic book","ISBN":300,"description":"funny"}`

	// POST送信用スレッドタイトル
	threadTitleForPost = `{"userName":"Alice","title":"I don't understand ..."}`

	// POST送信用スレッドメッセージ
	threadMessageForPost = `{"userName":"Alice","message":"Maybe it's because ..."}`

	// 本情報POST送信完了確認データ
	postReturnBookInfo = `{"ISBN":300,"title":"epic book","description":"funny"}
`
	// スレッドタイトルPOST送信完了確認データ
	postReturnThreadTitle = `{"id":3,"userName":"Alice","title":"I don't understand ...","ISBN":100}
`

	// スレッドメッセージPOST送信完了確認データ
	postReturnThreadMessage = `{"userName":"Alice","message":"Maybe it's because ...","threadID":1}
`

	// POSTした後のGET確認データ(メタ情報)
	metaDataAfterPost = `[{"ISBN":100,"title":"cool book"},{"ISBN":200,"title":"awesome book"},{"ISBN":300,"title":"epic book"}]
`

	// POSTした後のGET確認データ(詳細情報)
	profileDataAfterPost = `{"ISBN":300,"title":"epic book","description":"funny"}
`

	// POSTした後のGET確認データ（スレッドタイトル）
	threadTitlesAfterPost = `[{"id":1,"userName":"Alice","title":"I don't understand p.32 at all.","ISBN":100},{"id":2,"userName":"Bob","title":"there is an awful typo on p.55","ISBN":100},{"id":3,"userName":"Alice","title":"I don't understand ...","ISBN":100}]
`

	// POSTした後のGET確認データ（スレッドメッセージ）
	threadMessagesAfterPost = `[{"userName":"Carol","message":"Me neither.","threadID":1},{"userName":"Charlie","message":"I think the author tries to say ...","threadID":1},{"userName":"Alice","message":"Maybe it's because ...","threadID":1}]
`

	// ユーザ登録用
	testUser1 = `{"userName":"NewUser","password":"easypass"}`
	// ユーザ登録用
	testUser2 = `{"userName":"NewUser11111111","password":"easypass"}`
	// ユーザ登録用
	testUser3 = `{"userName":"New_00.-","password":"easypass"}`
	// ユーザ登録用
	testUser4 = `{"userName":"New","password":"easypass"}`
	// ユーザ登録用
	testUser5 = `{"userName":"NewUser1","password":"easypass1111111"}`
	// ユーザ登録用
	testUser6 = `{"userName":"NewUser2","password":"e_a.s!y-pass"}`
	// ユーザ登録用
	testUser7 = `{"userName":"NewUser3","password":"easypass111"}`

	// ユーザ登録/認証成功用
	testUserReturned1 = `{"id":5,"userName":"NewUser"}
`
	// ユーザ登録/認証成功用
	testUserReturned2 = `{"id":6,"userName":"NewUser11111111"}
`
	// ユーザ登録/認証成功用
	testUserReturned3 = `{"id":7,"userName":"New_00.-"}
`
	// ユーザ登録/認証成功用
	testUserReturned4 = `{"id":8,"userName":"New"}
`
	// ユーザ登録/認証成功用
	testUserReturned5 = `{"id":9,"userName":"NewUser1"}
`
	// ユーザ登録/認証成功用
	testUserReturned6 = `{"id":10,"userName":"NewUser2"}
`
	// ユーザ登録/認証成功用
	testUserReturned7 = `{"id":11,"userName":"NewUser3"}
`
	// ユーザ登録失敗用
	failUser1 = `{"userName":"!!Security!!","password":"easypass"}`

	// ユーザ登録失敗用
	failUser2 = `{"userName":"ab","password":"easypass"}`

	// ユーザ登録失敗用
	failUser3 = `{"userName":"ab11111111111111","password":"easypass"}`

	// ユーザ登録失敗用
	failUser4 = `{"userName":"test","password":"ab"}`

	// ユーザ登録失敗用
	failUser5 = `{"userName":"test","password":"ab"}`

	// ユーザ登録用(大文字小文字の区別なく重複))
	testSimilarUser = `{"userName":"Newuser","password":"Easypass"}`

	// ユーザ登録用(重複))
	testSameUser = `{"userName":"NewUser","password":"passpass"}`

	// ユーザ登録失敗用
	loginFailUser1 = `{"userName":"DifferentUser","password":"easypass"}`

	// ユーザ登録失敗用
	loginFailUser2 = `{"userName":"DifferentUser2","password":"fooooooo"}`

	// ダメなPOST
	invalidPostData = `{"foo":"bar"}`

	// 惜しいPOST（本情報）
	closePostBookData = `{"title":"epic book","ISBN":"300","description":"funny"}`

	// やる気のないPOST
	badPostData = `hello world`

	// ユーザがいないスレッドタイトル
	threadTitleMissingUser = `{"userName":"NoName","title":"foo"}`

	// ユーザがいないスレッドメッセージ
	threadMessageMissingUser = `{"userName":"NoName","message":"foo"}`

	// ヘッダ

	// JSON ヘッダ
	jsonHeader = `application/json; charset=UTF-8`

	// プレーンテキストヘッダ
	plainTextHeader = `text/plain; charset=UTF-8`

	// エラーメッセージ集

	// エラーメッセージ
	invalidThreadID = `ThreadID must be an integer`

	// エラーメッセージ
	notFound = `Not Found`

	// エラーメッセージ
	invalidFormat = `Invalid Post Format`

	// エラーメッセージ
	invalidISBN = `ISBN must be an integer`

	// エラーメッセージ
	bookInfoExists = `Book info already exists`

	// エラーメッセージ
	noUser = `User doesn't exist`

	// エラーメッセージ
	userExists = `User already exists`

	// エラーメッセージ
	noBook = `Book doesn't exist`

	// エラーメッセージ
	noThread = `Thread doesn't exist`

	// エラーメッセージ
	loginFailed = `Login Failed`
)

func TestMain(m *testing.M) {

	// mysqlに接続
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/bookbasket")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// handlerにデータベースの参照を渡す。
	SetDB(db)

	// 全テスト実行
	code := m.Run()

	os.Exit(code)
}

func TestGetAll(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, GetBookMetaInfoAll(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, metaInfoTestData, rec.Body.String())
	}
}

func TestGetProfile(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, GetBookProfile(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, bookProfileTestData, rec.Body.String())
	}
}

func TestGetProfileWithStringISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("foo")

	// Assertions
	if assert.NoError(t, GetBookProfile(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidISBN, rec.Body.String())
	}
}

func TestGetProfileWithInvalidISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("110")

	// Assertions
	if assert.NoError(t, GetBookProfile(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFound, rec.Body.String())
	}
}

func TestPostBookInfo(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bookInfoForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostBookInfo(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postReturnBookInfo, rec.Body.String())
	}
}

func TestAfterPostMetaData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, GetBookMetaInfoAll(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, metaDataAfterPost, rec.Body.String())
	}
}

func TestAfterPostProfileData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("300")

	// Assertions
	if assert.NoError(t, GetBookProfile(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, profileDataAfterPost, rec.Body.String())
	}
}

func TestPostBookInfoMultipleTimes(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bookInfoForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostBookInfo(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, bookInfoExists, rec.Body.String())
	}
}

func TestPostBookInfoWithInvalidArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostBookInfo(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostBookInfoWithInvalidButCloseArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(closePostBookData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostBookInfo(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostBookInfoWithBadArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(badPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostBookInfo(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestGetThreadTitles(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, GetThreadTitles(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadTitlesTestData, rec.Body.String())
	}
}

func TestGetEmptyThreadsTitles(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("200")

	// Assertions
	if assert.NoError(t, GetThreadTitles(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, emptyData, rec.Body.String())
	}
}

func TestGetThreadTitlesWithInvalidISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("foo")

	// Assertions
	if assert.NoError(t, GetThreadTitles(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidISBN, rec.Body.String())
	}
}

func TestGetThreadTitlesMissingBookData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("500")

	// Assertions
	if assert.NoError(t, GetThreadTitles(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFound, rec.Body.String())
	}
}

func TestGetThreadMessages(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, GetThreadMessages(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadMessagesTestData, rec.Body.String())
	}
}

func TestGetEmptyThreadMessages(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("2")

	// Assertions
	if assert.NoError(t, GetThreadMessages(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, emptyData, rec.Body.String())
	}
}

func TestGetThreadMessagesWithInvalidThreadID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/thread/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("foo")

	// Assertions
	if assert.NoError(t, GetThreadMessages(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidThreadID, rec.Body.String())
	}
}

func TestGetThreadMessagesMissingThreadTitle(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("5")

	// Assertions
	if assert.NoError(t, GetThreadMessages(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFound, rec.Body.String())
	}
}

func TestPostThreadTitle(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadTitleForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, PostThreadTitle(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postReturnThreadTitle, rec.Body.String())
	}
}

func TestAfterPostThreadTitle(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, GetThreadTitles(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadTitlesAfterPost, rec.Body.String())
	}
}

func TestPostThreadTitleWithInvalidArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, PostThreadTitle(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostThreadTitleWithBadArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(badPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, PostThreadTitle(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostThreadTitleWithInvalidISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadTitleForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("foo")

	// Assertions
	if assert.NoError(t, PostThreadTitle(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidISBN, rec.Body.String())
	}
}

func TestPostThreadTitleWithMissingBook(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadTitleForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("120")

	// Assertions
	if assert.NoError(t, PostThreadTitle(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, noBook, rec.Body.String())
	}
}

func TestPostThreadTitleWithMissingUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadTitleMissingUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/threads")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, PostThreadTitle(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, noUser, rec.Body.String())
	}
}

func TestPostThreadMessage(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadMessageForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, PostThreadMessage(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postReturnThreadMessage, rec.Body.String())
	}
}

func TestAfterPostThreadMessage(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, GetThreadMessages(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, threadMessagesAfterPost, rec.Body.String())
	}
}

func TestPostThreadMessageWithInvalidArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, PostThreadMessage(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostThreadMessageWithBadArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(badPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, PostThreadMessage(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostThreadTitleWithInvalidThreadID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadMessageForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("foo")

	// Assertions
	if assert.NoError(t, PostThreadMessage(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidThreadID, rec.Body.String())
	}
}

func TestPostThreadMessageWithMissingThread(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadMessageForPost))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("5")

	// Assertions
	if assert.NoError(t, PostThreadMessage(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, noThread, rec.Body.String())
	}
}

func TestPostThreadMessageWithMissingUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(threadMessageMissingUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/threads/:threadID")
	c.SetParamNames("threadID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, PostThreadMessage(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, noUser, rec.Body.String())
	}
}

func TestUserRegistration(t *testing.T) {
	RegistrationHelper(t, testUser1, testUserReturned1)
	RegistrationHelper(t, testUser2, testUserReturned2)
	RegistrationHelper(t, testUser3, testUserReturned3)
	RegistrationHelper(t, testUser4, testUserReturned4)
	RegistrationHelper(t, testUser5, testUserReturned5)
	RegistrationHelper(t, testUser6, testUserReturned6)
	RegistrationHelper(t, testUser7, testUserReturned7)
}

func RegistrationHelper(t *testing.T, input, expectation string) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(input))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/register")

	// Assertions
	if assert.NoError(t, RegisterUser(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectation, rec.Body.String())
	}
}

func TestInvalidUserRegistration(t *testing.T) {
	RegistrationFailHelper(t, testSimilarUser, userExists)
	RegistrationFailHelper(t, testSameUser, userExists)
	RegistrationFailHelper(t, badPostData, invalidFormat)
	RegistrationFailHelper(t, invalidPostData, invalidFormat)
	RegistrationFailHelper(t, failUser1, invalidFormat)
	RegistrationFailHelper(t, failUser2, invalidFormat)
	RegistrationFailHelper(t, failUser3, invalidFormat)
	RegistrationFailHelper(t, failUser4, invalidFormat)
	RegistrationFailHelper(t, failUser5, invalidFormat)
}

func RegistrationFailHelper(t *testing.T, input, expectation string) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(input))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/register")

	// Assertions
	if assert.NoError(t, RegisterUser(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectation, rec.Body.String())
	}
}

func TestUserLogin(t *testing.T) {
	LoginHelper(t, testUser1, testUserReturned1)
	LoginHelper(t, testUser2, testUserReturned2)
	LoginHelper(t, testUser3, testUserReturned3)
	LoginHelper(t, testUser4, testUserReturned4)
	LoginHelper(t, testUser5, testUserReturned5)
	LoginHelper(t, testUser6, testUserReturned6)
	LoginHelper(t, testUser7, testUserReturned7)
}

func LoginHelper(t *testing.T, input, expectation string) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(input))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/login")

	// Assertions
	if assert.NoError(t, AuthenticateUser(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectation, rec.Body.String())
	}
}

func TestUserLoginFail(t *testing.T) {
	LoginFailHelper(t, testSimilarUser, loginFailed)
	LoginFailHelper(t, testSameUser, loginFailed)
	LoginFailHelper(t, invalidPostData, invalidFormat)
	LoginFailHelper(t, badPostData, invalidFormat)
	LoginFailHelper(t, failUser1, invalidFormat)
	LoginFailHelper(t, failUser2, invalidFormat)
	LoginFailHelper(t, failUser3, invalidFormat)
	LoginFailHelper(t, failUser4, invalidFormat)
	LoginFailHelper(t, failUser5, invalidFormat)
	LoginFailHelper(t, loginFailUser1, loginFailed)
	LoginFailHelper(t, loginFailUser2, loginFailed)
}

func LoginFailHelper(t *testing.T, input, expectation string) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(input))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/login")

	// Assertions
	if assert.NoError(t, AuthenticateUser(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, expectation, rec.Body.String())
	}
}
