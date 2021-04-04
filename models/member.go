package models

import (
	"database/sql"
)

type Member struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    int64  `json:"phone"`
	Status   int    `json:"status"`
	Password string `json:"password"`
}

type MemberLogin struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    int64  `json:"phone"`
	Status   int    `json:"status"`
}

type MemberCollection struct {
	Member []Member `json:"members"`
}

type MemberData struct {
	Member []MemberLogin `json:"members"`
}

type ErrorResponse struct {
	Error error `json:"errors"`
}

func GetMember(db *sql.DB) MemberData {
	sql := "SELECT id, name, username, phone, status FROM member"
	rows, err := db.Query(sql)
	// Exit jika terjadi error
	if err != nil {
		panic(err)
	}
	// clean rows
	defer rows.Close()

	result := MemberData{}
	for rows.Next() {
		member := MemberLogin{}
		err2 := rows.Scan(&member.ID, &member.Name, &member.Username, &member.Phone, &member.Status)
		// Exit jika error
		if err2 != nil {
			panic(err2)
		}
		result.Member = append(result.Member, member)
	}
	return result
}

func PutMember(db *sql.DB, name string, username string, password string, phone int64, status int) (int64, error) {
	sql := "INSERT INTO member(name, username, password, phone, status) VALUES(?,?,?,?,?)"

	// membuat prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit jika error
	if err != nil {
		return 0, err
		// panic(err)
	}
	// memastikan statement ditutup setelah selesai
	defer stmt.Close()

	result, err2 := stmt.Exec(name, username, password, phone, status)
	// Exit jika error
	if err2 != nil {
		return 0, err2
		// panic(err2)
	}

	return result.LastInsertId()
}

func EditMember(db *sql.DB, memberId int, name string, phone int64, status int) (int64, error) {
	sql := "UPDATE member set name = ?, phone = ?, status = ? WHERE id = ?"

	stmt, err := db.Prepare(sql)

	if err != nil {
		return 0, err
		// panic(err)
	}

	result, err2 := stmt.Exec(name, phone, status, memberId)

	if err2 != nil {
		return 0, err2
		// panic(err2)
	}

	return result.RowsAffected()
}

func DeleteMember(db *sql.DB, id int) (int64, error) {
	sql := "DELETE FROM member WHERE id = ?"

	// buat prepare statement
	stmt, err := db.Prepare(sql)
	// Exit jika error
	if err != nil {
		return 0, err
		// panic(err)
	}

	result, err2 := stmt.Exec(id)
	// Exit jika error
	if err2 != nil {
		return 0, err2
		// panic(err2)
	}

	return result.RowsAffected()
}

func LoginMember(db *sql.DB, username string) (*Member, error) {
	sql := "SELECT * FROM member where username = ?"
	stmt, err := db.Prepare(sql)
	// Exit jika error
	if err != nil {
		return &Member{}, err
		// panic(err)
	}

	rows, err2 := stmt.Query(username)
	// Exit jika error
	defer rows.Close()

	result := Member{}
	for rows.Next() {
		member := Member{}
		err2 := rows.Scan(&member.ID, &member.Name, &member.Username, &member.Password, &member.Phone, &member.Status)
		// Exit jika error
		if err2 != nil {
			return &Member{}, err2
			// panic(err2)
		}
		result = member
	}
	if err2 != nil {
		return &Member{}, err2
		// panic(err2)
	}
	return &result, err2
}
