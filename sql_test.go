package belajar_golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUE('joko','Joko')"
	_, err := db.ExecContext(ctx, script)

	if err != nil {
		panic(err)
	}

	fmt.Println("succes insert new customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, FROM customer"
	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name,)
		if err != nil {
			panic(err)
		}
		fmt.Println("id : ", id)
		fmt.Println("name : ", name)
		
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var createdAt time.Time
		var birthDate sql.NullTime
		var married bool

		err = rows.Scan(&id, &name, &email, &balance, &rating,  &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("==============")
		fmt.Println("id : ", id)
		fmt.Println("name : ", name)
		if email.Valid{
			fmt.Println("email : ", email)
		}
		fmt.Println("balance : ", balance)
		fmt.Println("rating : ", rating)
		if birthDate.Valid{
			fmt.Println("birthdate : ", birthDate)
		}
		fmt.Println("married : ", married)
		fmt.Println("create_at : ", createdAt)
		
	}
}

func TestSqlInjectionSave(t *testing.T){
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin';#"
	password := "salah"

	script := "SELECT username FROM user WHERE username = '"+username+"' AND password = '"+password+"' LIMIT 1"
	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next(){
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		rows.Scan(&username)
		fmt.Println("Sukses login", username)
	}else{
		fmt.Println("gagal login")
	}
}

func TestSqlInjection(t *testing.T){
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin"
	password := "admin"

	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"
	rows, err := db.QueryContext(ctx, script, username, password)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next(){
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		rows.Scan(&username)
		fmt.Println("Sukses login", username)
	}else{
		fmt.Println("gagal login")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "eko"
	password := "eko"

	script := "INSERT INTO user(username, password) VALUE(?,?)"
	_, err := db.ExecContext(ctx, script, username, password)

	if err != nil {
		panic(err)
	}

	fmt.Println("succes insert new customer")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "eko@gmail.com"
	comment := "Test Komen"

	script := "INSERT INTO comments(email, comment) VALUE(?,?)"
	result, err := db.ExecContext(ctx, script, email, comment)

	if err != nil {
		panic(err)
	}
	insertId, err := result.LastInsertId()
	if err != nil{
		panic(err)
	}

	fmt.Println("succes insert new comment with id", insertId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO comments(email, comment) VALUE(?,?)"
	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "eko" +strconv.Itoa(i)+"@gmail.com"
		comment := "komentar ke "+strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment id ", id)
	}
 
}

func TestTransaction(t *testing.T){
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "INSERT INTO comments(email, comment) VALUE(?,?)"
	
	//do transaction
	for i := 0; i < 10; i++ {
		email := "eko" +strconv.Itoa(i)+"@gmail.com"
		comment := "komentar ke "+strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment id ", id)
	}

	tx.Rollback()
	if err != nil {
		panic(err)
	}
}