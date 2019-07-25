package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type sessionParameters struct {
	Width  int32 `json:"drawingAreaWidth"`
	Height int32 `json:"drawingAreaHeight"`
}

type drawingArea struct {
	AreaID string `json:"id"`
	Width  int32  `json:"drawingAreaWidth"`
	Height int32  `json:"drawingAreaHeight"`
}

func createSession(c echo.Context, drawingAreas map[string]*drawingArea) error {
	params := new(sessionParameters)
	if err := c.Bind(params); err != nil || 0 == params.Width || 0 == params.Height {
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide area width and height")
	}

	idString := "None"
	for found := true; found; _, found = drawingAreas[idString] {
		id, err := uuid.NewUUID()
		if err != nil {
			// ToDo. Add logs here!
			return echo.NewHTTPError(http.StatusInternalServerError, "Somthing gone wrong, need admin help")
		}
		idString = id.String()
	}

	area := &drawingArea{
		AreaID: idString,
		Width:  params.Width,
		Height: params.Height,
	}

	drawingAreas[area.AreaID] = area

	return c.JSON(http.StatusCreated, area)
}

func getSession(c echo.Context, drawingAreas map[string]*drawingArea) error {
	id := c.Param("id")

	area, err := drawingAreas[id]
	if !err {
		return echo.NewHTTPError(http.StatusNoContent, "This session does not exsists")
	}

	return c.JSON(http.StatusOK, area)
}

func main() {
	// var drawingAreas map[string]*drawingArea
	drawingAreas := make(map[string]*drawingArea)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Go ask for a /session over POST")
	})

	e.POST("/sessions", func(c echo.Context) error {
		return createSession(c, drawingAreas)
	})
	e.GET("/sessions/:id", func(c echo.Context) error {
		return getSession(c, drawingAreas)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
