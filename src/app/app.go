package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/geo/s2"
	"github.com/labstack/echo"
)

var e = createMux()

func Start() {
	e.Start(":8080")
}

func init() {
	e.GET("/map", Map)
	e.GET("/cell/id/:id", GetCellById)
	e.GET("/cell/token/:token", GetCellByToken)
	e.GET("/cell/lat/:lat/lng/:lng", GetCellByLatLng)
	e.GET("/parents/id/:id", GetParentsCellById)
}

func Map(c echo.Context) error {
	return c.Render(http.StatusOK, "map.html", "")
}

func GetCellById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusOK, errorJSON(err))
	}
	cellID := s2.CellID(id)
	return c.JSON(http.StatusOK, cellJSON(cellID))
}

func GetCellByToken(c echo.Context) error {
	token := c.Param("token")
	cellID := s2.CellIDFromToken(token)
	return c.JSON(http.StatusOK, cellJSON(cellID))
}

func GetCellByLatLng(c echo.Context) error {
	lat, err := strconv.ParseFloat(c.Param("lat"), 64)
	if err != nil {
		return c.JSON(http.StatusOK, errorJSON(err))
	}
	lng, err := strconv.ParseFloat(c.Param("lng"), 64)
	if err != nil {
		return c.JSON(http.StatusOK, errorJSON(err))
	}
	cellID := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
	return c.JSON(http.StatusOK, cellJSON(cellID))
}

func GetParentsCellById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusOK, errorJSON(err))
	}
	cellID := s2.CellID(id)

	var parents []s2.CellID
	for lv, min := cellID.Level()-1, 0; min <= lv; lv-- {
		parents = append(parents, cellID.Parent(lv))
	}
	return c.JSON(http.StatusOK, parentsJSON(parents))
}

func cellJSON(cellID s2.CellID) map[string]string {
	return map[string]string{
		"cell_id": fmt.Sprintf("%d", cellID),
		"bits":    fmt.Sprintf("%b", cellID),
		"level":   fmt.Sprintf("%d", cellID.Level()),
	}
}

func parentsJSON(ps []s2.CellID) map[int]map[string]string {
	res := make(map[int]map[string]string)
	for _, p := range ps {
		res[p.Level()] = cellJSON(p)
	}
	return res
}

func errorJSON(err error) map[string]string {
	return map[string]string{
		"err": fmt.Sprintf("%v", err),
	}
}
