package tests

type Person struct {
	Name string `tinyorm:"PRIMARY KEY"`
	Age  int
}
