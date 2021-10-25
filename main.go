package main

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	AddressName string
	UserID      int
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
			{AddressName: "Jl. Sehahtera"},
			{AddressName: "Jl. Cinta"},
		},
	}

	err := db.Create(&user).Error // pass pointer of data to Create

	if err != nil {
		panic(err)
	}
}

type Result struct {
	Name        string
	Number      string
	AddressName string
}

func FindUserByID(db *gorm.DB, id int, username string) {
	var user User
	//result := db.Last(&user)
	//result := db.First(&user, id)
	result := db.Where(&User{Username: username}).First(&user)

	fmt.Println("row affected: ", result.RowsAffected)
	fmt.Println(user)
	fmt.Println(user.Address)
	fmt.Println(user.CreditCard.Number)
}

func FindAllUserInRelatedField(db *gorm.DB) {
	var user []UserFieldAPI

	db.Model(&User{}).Limit(10).Find(&user)

	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(user)

}

func FindUserInformation(db *gorm.DB) {
	var result []Result

	//err := db.Joins("CreditCard").Where(&User{Username: "jin"}).Find(&result).Error
	// db.Model(&User{}).Select("users.name, credit_cards.number,addresses.name as addressName").Joins(
	// 	"left join credit_cards on credit_cards.user_id = users.id").Joins(
	// 	"left join addresses on addresses.user_id = users.id").Scan(&result)
	rows, err := db.Model(&User{}).Select("users.name, credit_cards.number,addresses.name as AddressName").Joins(
		"left join credit_cards on credit_cards.user_id = users.id").Joins(
		"left join addresses on addresses.user_id = users.id").Not("addresses.name = ?", "").Rows()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		// fmt.Println(rows)
		// fmt.Println(rows.Columns())
		// res := db.ScanRows(rows, &result)
		// fmt.Println(rows.Close().Error())

		// result = append(result, )
	}

	fmt.Println(result)

}

func CreateAnimals(db *gorm.DB) error {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&User{Name: "Tes"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Logger.Info(context.TODO(), "create")

	if err := tx.Create(&User{Name: "Yu"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=book_store port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	db.Logger.Info(context.Background(), "init")

	db.AutoMigrate(&User{}, &User_Book{}, &Company{}, &Book{}, &CreditCard{}, &Address{})

	//InsertUser(db)
	//FindUserByID(db, 3, "fadellh")
	// FindAllUserInRelatedField(db)
	// CreateAnimals(db)
	FindUserInformation(db)

	if err != nil {
		panic("Database Failed")
	}
}
