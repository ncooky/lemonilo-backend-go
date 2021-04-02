package models

import (
	"database/sql"
)

type Member struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Phone  int64  `json:"phone"`
	Status int    `json:"status"`
}

type MemberCollection struct {
	Member []Member `json:"members"`
}

type ErrorResponse struct {
	Error error `json:"errors"`
}

func GetMember(db *sql.DB) MemberCollection {
	sql := "SELECT * FROM member"
	rows, err := db.Query(sql)
	// Exit jika terjadi error
	if err != nil {
		panic(err)
	}
	// clean rows
	defer rows.Close()

	result := MemberCollection{}
	for rows.Next() {
		member := Member{}
		err2 := rows.Scan(&member.ID, &member.Name, &member.Phone, &member.Status)
		// Exit jika error
		if err2 != nil {
			panic(err2)
		}
		result.Member = append(result.Member, member)
	}
	return result
}

func PutMember(db *sql.DB, name string, phone int64, status int) (int64, error) {
	sql := "INSERT INTO member(name, phone, status) VALUES(?,?,?)"

	// membuat prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit jika error
	if err != nil {
		return 0, err
		// panic(err)
	}
	// memastikan statement ditutup setelah selesai
	defer stmt.Close()

	result, err2 := stmt.Exec(name, phone, status)
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
