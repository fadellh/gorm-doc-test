package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	Name       string
	Username   string
	Email      string
	CreditCard CreditCard
	Address    []Address
	gorm.Model
}

type Company struct {
	ID   int
	Name string
}

type User_Book struct {
	UserID int
	BookID int
	gorm.Model
}

type Book struct {
	Title string
	gorm.Model
}

//each user can have many address
type Address struct {
	Name   string
	UserID int
	gorm.Model
}

//Each user have one credit card
type CreditCard struct {
	gorm.Model
	UserID uint
	Number string
}

type NewUser struct {
}

type UserFieldAPI struct {
	Name     string
	Username string
	Email    string
}

func InsertUser(db *gorm.DB) {
	user := User{
		Name:     "Fadel",
		Username: "fadellh",
		Email:    "fadellh@gmail.com",
		CreditCard: CreditCard{
			Number: "3456676556"},
		Address: []Address{
			{Name: "Jl. Sehahtera"},
			{Name: "Jl. Cinta"},
		},
	}

	err := db.Create(&user).Error // pass pointer of data to Create

	if err != nil {
		panic(err)
	}
}

func FindUserByID(db *gorm.DB, id int, username string) {
	var user User
	//result := db.Last(&user)
	//result := db.First(&user, id)
	result := db.Where(&User{Username: username}).First(&user)

	fmt.Println("row affected: ", result.RowsAffected)
	fmt.Println(user)
}

func FindAllUserInRelatedField(db *gorm.DB) {
	var user []UserFieldAPI

	db.Model(&User{}).Limit(10).Find(&user)

	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(user)

}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=book_store port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&User{}, &User_Book{}, &Company{}, &Book{}, &CreditCard{}, &Address{})

	//InsertUser(db)
	// FindUserByID(db, 3, "fadellh")
	FindAllUserInRelatedField(db)

	if err != nil {
		panic("Database Failed")
	}
}
