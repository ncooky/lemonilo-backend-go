package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ncooky/lemonilo-backend-go/models"

	"github.com/labstack/echo"
)

type H map[string]interface{}

func GetMembers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.GetMember(db))
	}
}

func PutMember(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var member models.Member

		c.Bind(&member)

		id, err := models.PutMember(db, member.Name, member.Phone, member.Status)

		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"created": id,
			})
		} else {
			fmt.Println(err)
			return c.JSON(http.StatusUnprocessableEntity, H{
				"error": err,
			})
		}

	}
}

func EditMember(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		var member models.Member
		c.Bind(&member)

		_, err := models.EditMember(db, member.ID, member.Name, member.Phone, member.Status)

		if err == nil {
			return c.JSON(http.StatusOK, H{
				"updated": member,
			})
		} else {
			return c.JSON(http.StatusUnprocessableEntity, H{
				"error": err,
			})
		}
	}
}

func DeleteMember(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		_, err := models.DeleteMember(db, id)

		if err == nil {
			return c.JSON(http.StatusOK, H{
				"deleted": id,
			})
		} else {
			return c.JSON(http.StatusUnprocessableEntity, H{
				"error": err,
			})
		}

	}
}
