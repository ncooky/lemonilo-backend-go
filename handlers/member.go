package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ncooky/lemonilo-backend-go/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

type H map[string]interface{}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func validate(m *models.Member) error {
	if m.Password == "" {
		return errors.New("Password Must be filled")
	}
	if m.Username == "" {
		return errors.New("Username must be filled")
	}
	return nil
}

func GetMembers(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.GetMember(db))
	}
}

func PutMember(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var member models.Member

		c.Bind(&member)
		err := validate(&member)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, H{
				"error": err.Error(),
			})
		}
		hashedPassword, err := Hash(member.Password)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, H{
				"error": "hash error",
			})
		}
		id, err := models.PutMember(db, member.Name, member.Username, string(hashedPassword), member.Phone, member.Status)

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

func Login(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var member models.Member

		c.Bind(&member)
		result, err := models.LoginMember(db, member.Username)

		err = VerifyPassword(result.Password, member.Password)
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			return c.JSON(http.StatusUnprocessableEntity, H{
				"error": err,
			})
		}

		loginData := models.MemberLogin{}
		loginData.ID = result.ID
		loginData.Username = result.Username
		loginData.Name = result.Name
		loginData.Phone = result.Phone
		loginData.Status = result.Status
		if err == nil {
			return c.JSON(http.StatusOK, loginData)
		} else {
			fmt.Println(err)
			return c.JSON(http.StatusUnprocessableEntity, H{
				"error": err,
			})
		}

	}
}
