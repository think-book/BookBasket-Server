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
	metaInfoTestData = `{"id":1,"title":"cool book","ISBN":100},
{"id":2,"title":"awesome book","ISBN":200}`

	// GETProfile用確認データ
	bookProfileTestData = `{"ISBN":100,"title":"cool book","story":"A super hero beats monsters."}`

	// POST送信用データ
	postData1 = `{"title":"epic book","ISBN":300}`

	// POST送信用データ
	postData2 = `{"title":"boring book","ISBN":400}`

	// POST送信完了確認データ
	postReturnData = `{"id":3,"title":"epic book","ISBN":300}`

	// POST送信完了確認データ(メタ情報)
	metaDataAfterPost = `{"id":1,"title":"cool book","ISBN":100},
{"id":2,"title":"awesome book","ISBN":200},
{"id":3,"title":"epic book","ISBN":300},
{"id":4,"title":"boring book","ISBN":400}`

	// ダメなPOST
	invalidPostData = `{"foo":"bar"}`

	// エラーメッセージ
	notFound = `Not Found`

	// エラーメッセージ
	badRequest = `Invalid Post Format`
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
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, bookProfileTestData, rec.Body.String())
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
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFound, rec.Body.String())
	}
}

func TestPost(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postData1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostMetaInfo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postReturnData, rec.Body.String())
	}
}

func TestAfterPost(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(postData2))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")
	PostMetaInfo(c)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, GetBookMetaInfoAll(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, metaDataAfterPost, rec.Body.String())
	}
}

func TestPostWithInvalidArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(badRequest))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostMetaInfo(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, badRequest, rec.Body.String())
	}
}
