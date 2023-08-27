package main

import (
	"fmt"
	"tinyorm"

	_ "github.com/lib/pq"
)

func main() {
	e, err1 := tinyorm.NewEngine("postgres", "user=postgres dbname=test password=123456 sslmode=disable")
	if err1 != nil {
		panic(err1)
	}
	defer e.Close()
	s := e.NewSession()
	s.Raw("DROP TABLE IF EXISTS users").Exec()
	s.Raw("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50))").Exec()
	s.Raw("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50))").Exec()
	res, err2 := s.Raw("INSERT INTO users VALUES ($1, $2)", 1, "tom").Exec()
	if err2 != nil {
		panic(err2)
	}
	count, _ := res.RowsAffected()
	fmt.Printf("insert %d rows\n", count)
}
