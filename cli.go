package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pb "github.com/inigofu/temac-user-service/proto/auth"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
)

func main() {

	cmd.Init()

	// Create new greeter client
	client := pb.NewAuthService("temac.auth", microclient.DefaultClient)

	var user pb.User
	configFile, err := os.Open("user.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&user)
	log.Print("user", user)

	ruser, err := client.Create(context.TODO(), &user)
	if err != nil {
		log.Println("Could not create: %v", err)
	} else {
		log.Printf("Created: %s", ruser.User.Idcode)
	}
	rauth, err := client.Auth(context.TODO(), &user)
	if err != nil {
		log.Fatalf("Could not auth: %v", err)
	}
	// let's just exit because
	log.Println("autg with token", rauth.Token.Token)

	var menu pb.Menu
	configFile, err = os.Open("menu.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser = json.NewDecoder(configFile)
	jsonParser.Decode(&menu)
	log.Println("menu", menu)
	ctx := metadata.NewContext(context.TODO(), map[string]string{
		"Authorization": rauth.Token.Token,
	})
	rmenu, err := client.CreateMenu(ctx, &menu)
	if err != nil {
		log.Println("Could not create: %v", err)
	} else {
		log.Printf("Created menu: %s", rmenu)
	}

	var role pb.Role
	configFile, err = os.Open("role.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser = json.NewDecoder(configFile)
	jsonParser.Decode(&role)
	log.Println("role", role)

	rrole, err := client.CreateRole(ctx, &role)
	if err != nil {
		log.Println("Could not create: %v", err)
	} else {
		log.Printf("Created role: %s", rrole)
	}
	temprole := make([]*pb.Role, 1)
	temprole[0] = &pb.Role{Idcode: rrole.Role.Idcode}
	user = *ruser.User
	user.Roles = temprole
	log.Printf("Updating user: %s", user)
	ruser, err = client.UpdateUser(ctx, &user)
	if err != nil {
		log.Println("Could not update: %v", err)
	} else {
		log.Printf("Created use: %s", ruser)
	}
	var form []pb.Form
	configFile, err = os.Open("form.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser = json.NewDecoder(configFile)
	jsonParser.Decode(&form)
	log.Println("form", form)
	for _, element := range form {
		rform, err := client.CreateForm(ctx, &element)
		if err != nil {
			log.Println("Could not create form: %v", err)
		} else {
			log.Printf("Created form: %s", rform)
		}
	}
	log.Printf("Procedure finished")
	os.Exit(0)
}
