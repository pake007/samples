package tests

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"shorturl/models"
	_ "shorturl/routers"
	"testing"
	"time"
)

type ShortResult struct {
	UrlShort string
	UrlLong  string
}

var cacheExistMock func(key string) bool
var getCacheMock func(key string) interface{}
var putCacheMock func(key string, val interface{}, timeout time.Duration) error

type cacheManagerMock struct{}

func (m cacheManagerMock) IsExist(key string) bool {
	return cacheExistMock(key)
}

func (m cacheManagerMock) Get(key string) interface{} {
	return getCacheMock(key)
}

func (m cacheManagerMock) Put(key string, val interface{}, timeout time.Duration) error {
	return nil
}

// start test
func init() {
	_, file, _, _ := runtime.Caller(0)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(appPath)

	models.CacheCond = cacheManagerMock{}
}

func TestShortWithoutCache(t *testing.T) {
	// mock
	cacheExistMock = func(key string) bool {
		return false
	}
	// test
	r, _ := http.NewRequest("GET", "/v1/shorten", nil)
	q := r.URL.Query()
	q.Add("longurl", "http://www.beego.me/")
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("testing", "TestShort", "Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Have ShortResult", func() {
			var s ShortResult
			contents, _ := ioutil.ReadAll(w.Body)
			json.Unmarshal(contents, &s)
			So(s.UrlShort, ShouldNotBeBlank)
		})
	})
}

func TestShortWithCache(t *testing.T) {
	// mock
	cacheExistMock = func(key string) bool {
		return true
	}
	getCacheMock = func(key string) interface{} {
		return "something"
	}
	// test
	r, _ := http.NewRequest("GET", "/v1/shorten", nil)
	q := r.URL.Query()
	q.Add("longurl", "http://www.beego.me/")
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("testing", "TestShort", "Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Have ShortResult", func() {
			var s ShortResult
			contents, _ := ioutil.ReadAll(w.Body)
			json.Unmarshal(contents, &s)
			So(s.UrlShort, ShouldNotBeBlank)
		})
	})
}

func TestExpandWithoutCache(t *testing.T) {
	// mock
	cacheExistMock = func(key string) bool {
		return false
	}
	// test
	r, _ := http.NewRequest("GET", "/v1/expand", nil)
	q := r.URL.Query()
	q.Add("shorturl", "5laZF")
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("testing", "TestExpand", "Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Have LongResult If No Cache Exists", func() {
			var s ShortResult
			contents, _ := ioutil.ReadAll(w.Body)
			json.Unmarshal(contents, &s)
			So(s.UrlLong, ShouldBeBlank)
		})
	})
}

func TestExpandWithCache(t *testing.T) {
	// mock
	cacheExistMock = func(key string) bool {
		return true
	}
	getCacheMock = func(key string) interface{} {
		return "www.google.com"
	}
	// execute
	r, _ := http.NewRequest("GET", "/v1/expand", nil)
	q := r.URL.Query()
	q.Add("shorturl", "5laZF")
	r.URL.RawQuery = q.Encode()
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("testing", "TestExpand", "Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("Then Result Should Have LongResult If Cache Exists", func() {
			var s ShortResult
			contents, _ := ioutil.ReadAll(w.Body)
			json.Unmarshal(contents, &s)
			So(s.UrlLong, ShouldEqual, "www.google.com")
		})
	})
}
