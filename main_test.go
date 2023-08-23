package tinyorm_test

import (
	"fmt"
	"testing"
	"tinyorm"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	e, _ := tinyorm.NewEngine("postgres", "user=postgres dbname=test password=123456 sslmode=disable")
	defer e.Close()
	s := e.NewSession()
	s.Raw("DROP TABLE IF EXISTS users").Exec()
	s.Raw("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50))").Exec()
	s.Raw("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50))").Exec()
	res, err := s.Raw("INSERT INTO users VALUES ($1, $2)", 1, "tom").Exec()
	if err != nil {
		panic(err)
	}
	count, _ := res.RowsAffected()
	fmt.Printf("insert %d rows\n", count)
}
