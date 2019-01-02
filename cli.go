package main

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/inigofu/shippy-user-service/proto/auth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

func main() {

	// Creates a database connection and handles
	// closing it again before exit.
	host := "localhost:54321"
	username := "postgres"
	DBName := "postgres"
	password := "postgres"
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			username, password, host, DBName,
		),
	)
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	// db.LogMode(true)

	user := &pb.User{}
	//var roles []*pb.Role
	var menues []*pb.Menu
	var rolmenuesall []*pb.Menu
	email := "martin4"
	if err := db.Preload("Roles.Menues").Select("id").Where("email = ?", email).
		First(&user).Error; err != nil {
		fmt.Println("errp", err)
	}

	for _, role := range user.Roles {
		rolmenuesall = append(rolmenuesall, role.Menues...)
	}
	var rolmenues []string
	for _, role := range rolmenuesall {
		rolmenues = append(rolmenues, role.Id)
	}
	type Result struct {
		Children_id string
	}
	// fmt.Println(rolmenues)
	var results []Result
	var childrenid []string
	db.Raw("SELECT children_id FROM menu_childrens").Scan(&results)
	for _, result := range results {
		childrenid = append(childrenid, result.Children_id)
	}
	// (*sql.Row)
	// fmt.Println(childrenid)
	if err := db.Not(childrenid).Where(rolmenues).Preload("Children", "id in (?)", rolmenues).Find(&menues).Error; err != nil {
		fmt.Println("errp 2", err)
	}
	// fmt.Println(menues)

	key := []byte("mySuperSecretKeyLol")
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: 24,
			Issuer:    "shippy.user",
		},
	}
	fmt.Println("start token")
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	fmt.Println("token", tokenString)
	token2, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	fmt.Println(token2)
	// Validate the token and return the custom claims
	if claims, ok := token2.Claims.(*CustomClaims); ok && token2.Valid {
		fmt.Println(claims)
	} else {
		fmt.Println(err)
	}
}
