package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	sessionCreate = `{"drawingAreaWidth":1000,"drawingAreaHeight":1000}`
	sessionGet    = `{"id":"662c5ec6-7416-11e9-8c23-1681be663d3e"}`
	sessionJSON   = `{"id":"662c5ec6-7416-11e9-8c23-1681be663d3e","drawingAreaWidth":1000,"drawingAreaHeight":1000}`
)

func сallCreation(e *echo.Echo, drawingAreas map[string]*drawingArea) (*drawingArea, error) {
	req := httptest.NewRequest(http.MethodPost, "/sessions", strings.NewReader(sessionCreate))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()

	createSession(e.NewContext(req, res), drawingAreas)

	// assert.Equal(t, http.StatusCreated, res.Code)
	if http.StatusCreated != res.Code {
		return nil, errors.New("creation function return wrong code")
	}

	result := new(drawingArea)
	err := json.Unmarshal([]byte(res.Body.String()), &result)
	// assert.Equal(t, err, nil)
	if nil != err {
		return nil, errors.New("creation function return wrong data format")
	}

	return result, nil
}

func TestCreateSession(t *testing.T) {
	e := echo.New()

	drawingAreas := make(map[string]*drawingArea)
	area, err := сallCreation(e, drawingAreas)

	assert.Equal(t, err, nil)
	if nil != err {
		return
	}

	assert.Equal(t, area.Width, int32(1000))
	assert.Equal(t, area.Height, int32(1000))
}

func TestSessionCreatedAsUnique(t *testing.T) {
	e := echo.New()
	drawingAreas := make(map[string]*drawingArea)

	areaOne, err := сallCreation(e, drawingAreas)

	assert.Equal(t, err, nil)
	if nil != err {
		return
	}

	areaTwo, err := сallCreation(e, drawingAreas)

	assert.Equal(t, err, nil)
	if nil != err {
		return
	}

	assert.NotEqual(t, areaOne.AreaID, areaTwo.AreaID)
}

func TestGetUnexsistedSession(t *testing.T) {
	exptectedError := echo.NewHTTPError(http.StatusNoContent, "This session does not exsists")
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sessions/what_a_session",
		strings.NewReader(sessionGet))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()

	ctx := e.NewContext(req, res)
	ctx.SetParamNames("id")
	ctx.SetParamValues("what_a_session")

	drawingAreas := make(map[string]*drawingArea)
	result := getSession(ctx, drawingAreas)
	assert.Equal(t, exptectedError, result)
}

func TestGetCreatedSession(t *testing.T) {
	// Setup
	e := echo.New()
	addReq := httptest.NewRequest(http.MethodPost, "/sessions", strings.NewReader(sessionCreate))
	addReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	addRes := httptest.NewRecorder()

	drawingAreas := make(map[string]*drawingArea)
	createSession(e.NewContext(addReq, addRes), drawingAreas)

	// Assertions
	assert.Equal(t, http.StatusCreated, addRes.Code)
	addResult := new(drawingArea)
	err := json.Unmarshal([]byte(addRes.Body.String()), &addResult)

	assert.Equal(t, err, nil)
	if nil != err {
		return
	}

	getSessionURL := fmt.Sprintf("/sessions/%s", addResult.AreaID)

	getReq := httptest.NewRequest(http.MethodPost, getSessionURL, strings.NewReader(""))
	getReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	getRes := httptest.NewRecorder()

	ctx := e.NewContext(getReq, getRes)
	ctx.SetParamNames("id", addResult.AreaID)
	ctx.SetParamValues(addResult.AreaID)

	getSession(ctx, drawingAreas)

	assert.Equal(t, http.StatusOK, getRes.Code)

	getResult := new(drawingArea)
	err = json.Unmarshal([]byte(getRes.Body.String()), &getResult)

	assert.Equal(t, err, nil)
	if nil != err {
		return
	}

	assert.Equal(t, getResult.AreaID, addResult.AreaID)
	assert.Equal(t, getResult.Width, addResult.Width)
	assert.Equal(t, getResult.Height, addResult.Height)
}
