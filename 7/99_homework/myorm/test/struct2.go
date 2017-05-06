package test

//myorm:users
type User2 struct {
	ID            uint   `myorm:"primary_key"`     // первичный ключ, в него мы пишем LastInsertId
	Login         string `myorm:"column:username"` // поле называется username в таблице
	Info          string `myorm2:"null"`           // поле может иметь тип null
	Balance       int
	Status        int
	SomeInnerFlag bool `myorm:"-"` //поля нет в таблице, игнорируем его
}

var aaa string

type MyInt int

//myorm:abba
func abba() {

}
