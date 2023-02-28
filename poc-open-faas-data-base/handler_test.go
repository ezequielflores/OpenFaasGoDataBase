package function

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testParam struct {
	testName  string
	pathParam string
	method    string
}

func mockFindCharacter(w http.ResponseWriter, _ *http.Request) {
	log.Println("find character controller mock ok")
	w.WriteHeader(http.StatusOK)

}

func mockCreateCharacter(w http.ResponseWriter, _ *http.Request) {
	log.Println("create character controller mock ok")
	w.WriteHeader(http.StatusOK)

}

func init() {
	routes["/characters/"] = mockFindCharacter
	routes["/characters"] = mockCreateCharacter
}

func TestUrlOk(t *testing.T) {

	testCase := []testParam{
		{testName: "find character", pathParam: "/api/v1/starwar/characters/1", method: http.MethodGet},
		{testName: "create character", pathParam: "/api/v1/starwar/characters", method: http.MethodPost},
	}

	for _, param := range testCase {
		t.Run(param.testName, func(t *testing.T) {
			req := httptest.NewRequest(
				param.method,
				param.pathParam,
				nil,
			)
			rec := httptest.NewRecorder()

			Handle(rec, req)
			resp := rec.Result()

			assert.EqualValues(t, resp.StatusCode, http.StatusOK)
		})
	}

}

func TestNotFound(t *testing.T) {

	testCase := []testParam{
		{testName: "find with letter path param", pathParam: "/api/v1/starwar/characters/a", method: http.MethodGet},
		{testName: "find with invalid url", pathParam: "/api/v1/starwar/characters/1/ad", method: http.MethodGet},
		{testName: "find with incomplete url", pathParam: "/api/v1/starwar/characters/", method: http.MethodGet},
		{testName: "find with url invalid", pathParam: "/api/v1/starwar/characters/_", method: http.MethodGet},
		{testName: "find with invalid star url", pathParam: "/hola/mundo/api/v1/starwar/characters/_", method: http.MethodGet},
		{testName: "find with invalid method", pathParam: "/api/v1/starwar/characters/1", method: http.MethodPut},
		{testName: "create with invalid method", pathParam: "/api/v1/starwar/characters", method: http.MethodPut},
	}

	for _, param := range testCase {
		t.Run(param.testName, func(t *testing.T) {
			req := httptest.NewRequest(
				param.method,
				param.pathParam,
				nil,
			)
			rec := httptest.NewRecorder()

			Handle(rec, req)
			resp := rec.Result()
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(resp.Body)
			if err != nil {
				log.Println("error reading response body from controller mock")
				return
			}
			assert.EqualValues(t, resp.StatusCode, http.StatusNotFound)
			assert.EqualValues(t, buf.String(), fmt.Sprintf("Url: %s Not Found", param.pathParam))
		})
	}

}
