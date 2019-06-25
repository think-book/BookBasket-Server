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

	// POST送信完了確認データ(詳細情報)
	profileDataAfterPost = `{"ISBN":300, "title":"epic book", "story":"funny"}`

	// ダメなPOST
	invalidPostData = `{"foo":"bar"}`

	// エラーメッセージ
	notFound = `Not Found`

	// エラーメッセージ
	invalidFormat = `Invalid Post Format`

	// エラーメッセージ
	invalidISBN = `ISBN must be an integer`

	// エラーメッセージ
	metaExists = `Meta info already exists`

	// エラーメッセージ
	profileExists = `Book profile already exists`

	// エラーメッセージ
	metaNotFound = `Book Meta Data Not Found`

	// エラーメッセージ
	inconsistentISBN = `ISBN is inconsistent`
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
