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
	metaInfoTestData = `[{"id":1,"title":"cool book","ISBN":100},{"id":2,"title":"awesome book","ISBN":200}]`

	// GETProfile用確認データ
	bookProfileTestData = `{"ISBN":100,"title":"cool book","story":"A super hero beats monsters."}`

	// GETForumTitles用確認データ
	forumTitlesTestData = `[{"id":1,"user":"user_X","title":"I don't understand p.32 at all.","ISBN":100},{"id":2,"user":"user_Y","title":"there is an awful typo on p.55","ISBN":100}]
`
	// 空配列確認データ
	emptyData = `[]
`
	// GETForumMessages用確認データ
	forumMessagesTestData = `[{"id":1,"user":"user_A","message":"Me neither.","forumID":1},{"id":2,"user":"user_B","message":"I think the author tries to say ...","forumID":1}]
`

	// POST送信用メタデータ
	postMetaData = `{"title":"epic book","ISBN":300}`

	// POST送信用詳細データ
	postProfileData = `{"ISBN":300,"title":"epic book","story":"funny"}`

	// POST送信用詳細データ(メタデータなし)
	postProfileDataMissingMeta = `{"ISBN":500, "title":"normal book", "story":"sad"}`

	// POST送信完了確認データ
	postReturnMetaData = `{"id":3,"title":"epic book","ISBN":300}`

	// POST送信完了確認データ(メタ情報)
	metaDataAfterPost = `[{"id":1,"title":"cool book","ISBN":100},{"id":2,"title":"awesome book","ISBN":200},{"id":3,"title":"epic book","ISBN":300}]`

	// ダメなPOST
	invalidPostData = `{"foo":"bar"}`

	// エラーメッセージ
	notFound = `Not Found`

	// エラーメッセージ
	invalidFormat = `Invalid Post Format`

	// エラーメッセージ
	invalidISBN = `ISBN must be an integer`

	// エラーメッセージ
	invalidforumID = `forumID must be an integer`

	// エラーメッセージ
	metaExists = `Meta info already exists`

	// エラーメッセージ
	profileExists = `Book profile already exists`

	// エラーメッセージ
	metaNotFound = `Book Meta Data Not Found`

	// エラーメッセージ
	inconsistentISBN = `ISBN is inconsistent`

	// JSON ヘッダ
	jsonHeader = `application/json; charset=UTF-8`

	// プレーンテキストヘッダ
	plainTextHeader = `text/plain; charset=UTF-8`
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
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFound, rec.Body.String())
	}
}

func TestPostMetaData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postMetaData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostMetaInfo(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postReturnMetaData, rec.Body.String())
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
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, metaDataAfterPost, rec.Body.String())
	}
}

func TestPostMetaDataMultipleTimes(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postMetaData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostMetaInfo(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, metaExists, rec.Body.String())
	}
}

func TestPostMetaDataWithInvalidArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books")

	// Assertions
	if assert.NoError(t, PostMetaInfo(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostProfileData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postProfileData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("300")

	// Assertions
	if assert.NoError(t, PostBookProfile(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postProfileData, rec.Body.String())
	}
}

func TestAfterPostProfileData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("300")

	// Assertions
	if assert.NoError(t, GetBookProfile(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, postProfileData, rec.Body.String())
	}
}

func TestPostProfileDataMultipleTimes(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postProfileData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("300")

	// Assertions
	if assert.NoError(t, PostBookProfile(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, profileExists, rec.Body.String())
	}
}

func TestPostProfileDataWithInvalidArgument(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidPostData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("300")

	// Assertions
	if assert.NoError(t, PostBookProfile(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidFormat, rec.Body.String())
	}
}

func TestPostProfileDataWithInvalidISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postProfileData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("baz")

	// Assertions
	if assert.NoError(t, PostBookProfile(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidISBN, rec.Body.String())
	}
}

func TestPostProfileDataWithMissingISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postProfileDataMissingMeta))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("500")

	// Assertions
	if assert.NoError(t, PostBookProfile(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, metaNotFound, rec.Body.String())
	}
}

func TestPostProfileDataWithInconsistentISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postProfileData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN")
	c.SetParamNames("ISBN")
	c.SetParamValues("200")

	// Assertions
	if assert.NoError(t, PostBookProfile(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, inconsistentISBN, rec.Body.String())
	}
}

func TestGetForumTitles(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/forum")
	c.SetParamNames("ISBN")
	c.SetParamValues("100")

	// Assertions
	if assert.NoError(t, GetForumTitles(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, forumTitlesTestData, rec.Body.String())
	}
}

func TestGetEmptyForumTitles(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/forum")
	c.SetParamNames("ISBN")
	c.SetParamValues("200")

	// Assertions
	if assert.NoError(t, GetForumTitles(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, emptyData, rec.Body.String())
	}
}

func TestGetForumTitlesWithInvalidISBN(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/forum")
	c.SetParamNames("ISBN")
	c.SetParamValues("foo")

	// Assertions
	if assert.NoError(t, GetForumTitles(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidISBN, rec.Body.String())
	}
}

func TestGetForumTitlesMissingBookData(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:ISBN/forum")
	c.SetParamNames("ISBN")
	c.SetParamValues("500")

	// Assertions
	if assert.NoError(t, GetForumTitles(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFound, rec.Body.String())
	}
}

func TestGetForumMessages(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/forum/:forumID")
	c.SetParamNames("forumID")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, GetForumMessages(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, forumMessagesTestData, rec.Body.String())
	}
}

func TestGetEmptyForumMessages(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/forum/:forumID")
	c.SetParamNames("forumID")
	c.SetParamValues("2")

	// Assertions
	if assert.NoError(t, GetForumMessages(c)) {
		res := rec.Result()
		assert.Equal(t, jsonHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, emptyData, rec.Body.String())
	}
}

func TestGetForumMessagesWithInvalidForumID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/forum/:forumID")
	c.SetParamNames("forumID")
	c.SetParamValues("foo")

	// Assertions
	if assert.NoError(t, GetForumMessages(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, invalidforumID, rec.Body.String())
	}
}

func TestGetForumMessagesMissingForumTitle(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/forum/:forumID")
	c.SetParamNames("forumID")
	c.SetParamValues("5")

	// Assertions
	if assert.NoError(t, GetForumMessages(c)) {
		res := rec.Result()
		assert.Equal(t, plainTextHeader, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, notFound, rec.Body.String())
	}
}
