package main

/*
[Calculator for Investors - Stage 1/4: What's on the menu?](https://hyperskill.org/projects/264/stages/1335/implement)
-------------------------------------------------------------------------------
[Primitive types](https://hyperskill.org/learn/topic/1807)
[Input/Output](https://hyperskill.org/learn/topic/1506)
[Slices](https://hyperskill.org/learn/topic/1672)
[Control statements](https://hyperskill.org/learn/topic/1728)
[Loops](https://hyperskill.org/learn/topic/1531)
[Functions](https://hyperskill.org/learn/topic/1750)
*/

import (
	"fmt"
	"os"
)

const (
	notImplementedMsg = "Not implemented!\n\n"
	byeMsg            = "Have a nice day!\n"
	invalidOptionMsg  = "Invalid option!\n\n"
	optionPrompt      = "\nEnter an option:"
)

func printMainMenu() {
	fmt.Println("MAIN MENU")
	mainMenuItems := []string{
		"Exit",
		"CRUD operations",
		"Show top ten companies by criteria",
	}
	for idx, option := range mainMenuItems {
		fmt.Printf("%d %s\n", idx, option)
	}
	fmt.Println(optionPrompt)
	var mainMenuOption string
	fmt.Scanln(&mainMenuOption)

	switch mainMenuOption {
	case "0":
		fmt.Print(byeMsg)
		os.Exit(0)
	case "1":
		printCrudMenu()
	case "2":
		printTopTenMenu()
	default:
		fmt.Print(invalidOptionMsg)
	}
}

func printCrudMenu() {
	fmt.Println("\nCRUD MENU")
	crudMenuItems := []string{
		"Back",
		"Create a company",
		"Read a company",
		"Update a company",
		"Delete a company",
		"List all companies",
	}
	for idx, option := range crudMenuItems {
		fmt.Printf("%d %s\n", idx, option)
	}
	fmt.Println(optionPrompt)
	var crudMenuOption string
	fmt.Scanln(&crudMenuOption)

	switch crudMenuOption {
	case "0":
		return
	default:
		fmt.Print("Not implemented!\n\n")
	}
}

func printTopTenMenu() {
	fmt.Println("\nTOP TEN MENU")
	topTenMenuItems := []string{
		"Back",
		"List by ND/EBITDA",
		"List by ROE",
		"List by ROA",
	}
	for idx, option := range topTenMenuItems {
		fmt.Printf("%d %s\n", idx, option)
	}
	fmt.Println(optionPrompt)
	var topTenMenuOption string
	fmt.Scanln(&topTenMenuOption)

	switch topTenMenuOption {
	case "0":
		return
	default:
		fmt.Print(notImplementedMsg)
	}
}

func main() {
	for {
		printMainMenu()
	}
}
