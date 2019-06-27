package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	// GETAll用確認データ
	metaInfoTestData = `[{"id":1,"title":"cool book","ISBN":100},{"id":2,"title":"awesome book","ISBN":200}]
`

	// GETProfile用確認データ
	bookProfileTestData = `{"ISBN":100,"title":"cool book","description":"A super hero beats monsters."}
`

	// POST送信用データ
	postData = `{"title":"epic book","ISBN":300,"description":"funny"}
`

	// POST送信完了確認データ
	postReturnData = `{"id":3,"title":"epic book","description":"funny","ISBN":300}
`

	// POST送信完了確認データ(メタ情報)
	metaDataAfterPost = `[{"id":1,"title":"cool book","ISBN":100},{"id":2,"title":"awesome book","ISBN":200},{"id":3,"title":"epic book","ISBN":300}]
`

	// ダメなPOST
	invalidPostData = `{"foo":"bar"}`

	// やる気のないPOST
	badPostData = `hello world`

	// JSON ヘッダ
	jsonHeader = `application/json; charset=UTF-8`

	// プレーンテキストヘッダ
	plainTextHeader = `text/plain; charset=UTF-8`

	// エラーメッセージ
	notFound = `Not Found`

	// エラーメッセージ
	invalidFormat = `Invalid Post Format`

	// エラーメッセージ
	invalidISBN = `ISBN must be an integer`

	// エラーメッセージ
	dataExists = `Book info already exists`
)

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

func TestPostData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostBookInfo(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postReturnData, rec.Body.String())
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

func TestPostDataMultipleTimes(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostBookInfo(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, dataExists, rec.Body.String())
	}
}

func TestPostDataWithInvalidArgument(t *testing.T) {
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

func TestPostDataWithBadArgument(t *testing.T) {
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
