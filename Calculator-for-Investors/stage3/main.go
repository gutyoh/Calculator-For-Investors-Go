package main

/*
[Calculator for Investors - Stage 3/4: Talking numbers](https://hyperskill.org/projects/264/stages/1337/implement)
-------------------------------------------------------------------------------
[Debugging Go code in GoLand](https://hyperskill.org/learn/step/23118)
[CRUD Operations — READ](https://hyperskill.org/learn/step/24151)
[CRUD Operations — UPDATE] - TODO
[CRUD Operations — DELETE] - TODO
*/

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

const (
	welcomeMsg        = "Welcome to the Investor Program!\n\n"
	notImplementedMsg = "Not implemented!\n\n"
	byeMsg            = "Have a nice day!\n"
	invalidOptionMsg  = "Invalid option!\n\n"
	companyNamePrompt = "Enter company name"
	optionPrompt      = "\nEnter an option:"
)

type Company struct {
	Ticker string `gorm:"primaryKey"`
	Name   *string
	Sector *string
}

type Financial struct {
	Ticker          string `gorm:"primaryKey"`
	Ebitda          *float64
	Sales           *float64
	NetProfit       *float64
	MarketPrice     *float64
	NetDebt         *float64
	Assets          *float64
	Equity          *float64
	CashEquivalents *float64
	Liabilities     *float64
}

func (*Financial) TableName() string {
	return "financial"
}

func (f *Financial) UpdateValue(field string) {
	newValue := getFloatInput(fmt.Sprintf("Enter %s", field))
	switch field {
	case "ebitda":
		f.Ebitda = &newValue
	case "sales":
		f.Sales = &newValue
	case "net profit":
		f.NetProfit = &newValue
	case "market price":
		f.MarketPrice = &newValue
	case "net debt":
		f.NetDebt = &newValue
	case "assets":
		f.Assets = &newValue
	case "equity":
		f.Equity = &newValue
	case "cash equivalents":
		f.CashEquivalents = &newValue
	case "liabilities":
		f.Liabilities = &newValue
	}
}

func (f *Financial) calculatePE() float64 {
	if f.NetProfit != nil && f.MarketPrice != nil {
		return *f.MarketPrice / *f.NetProfit
	}
	return 0
}

func (f *Financial) calculatePS() float64 {
	if f.Sales != nil && f.MarketPrice != nil {
		return *f.MarketPrice / *f.Sales
	}
	return 0
}

func (f *Financial) calculatePB() float64 {
	if f.Assets != nil && f.MarketPrice != nil {
		return *f.MarketPrice / *f.Assets
	}
	return 0
}

func (f *Financial) calculateNDEBITDA() float64 {
	if f.NetDebt != nil && f.Ebitda != nil {
		ndEbitda := *f.NetDebt / *f.Ebitda
		return ndEbitda
	}
	return 0
}

func (f *Financial) calculateROE() float64 {
	if f.NetProfit != nil && f.Equity != nil {
		return *f.NetProfit / *f.Equity
	}
	return 0
}

func (f *Financial) calculateROA() float64 {
	if f.NetProfit != nil && f.Assets != nil {
		return *f.NetProfit / *f.Assets
	}
	return 0
}

func (f *Financial) calculateLA() float64 {
	if f.Liabilities != nil && f.Assets != nil {
		return *f.Liabilities / *f.Assets
	}
	return 0
}

func (f *Financial) CalculateIndicators() FinancialIndicators {
	return FinancialIndicators{
		PE:       f.calculatePE(),
		PS:       f.calculatePS(),
		PB:       f.calculatePB(),
		NDEBITDA: f.calculateNDEBITDA(),
		ROE:      f.calculateROE(),
		ROA:      f.calculateROA(),
		LA:       f.calculateLA(),
	}
}

type FinancialIndicators struct {
	PE       float64
	PS       float64
	PB       float64
	NDEBITDA float64
	ROE      float64
	ROA      float64
	LA       float64
}

func (fi FinancialIndicators) print() {
	fmt.Printf("P/E = %.2f\n", fi.PE)
	fmt.Printf("P/S = %.2f\n", fi.PS)
	fmt.Printf("P/B = %.2f\n", fi.PB)
	if fi.NDEBITDA != 0 {
		fmt.Printf("ND/EBITDA = %.2f\n", fi.NDEBITDA)
	} else {
		fmt.Printf("ND/EBITDA = None\n")
	}
	fmt.Printf("ROE = %.2f\n", fi.ROE)
	fmt.Printf("ROA = %.2f\n", fi.ROA)
	fmt.Printf("L/A = %.2f\n\n", fi.LA)
}

func readCSVData(filename string) []map[string]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open CSV file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		log.Panicf("failed to read CSV data: %s", err)
	}
	data := make([]map[string]string, len(lines)-1)

	for i, line := range lines[1:] {
		row := make(map[string]string)
		for j, field := range line {
			row[lines[0][j]] = field
		}
		data[i] = row
	}
	return data
}

func createCompanyRecord(db *gorm.DB, row map[string]string) {
	company := Company{
		Ticker: row["ticker"],
		Name:   nullableString(row["name"]),
		Sector: nullableString(row["sector"]),
	}
	db.Create(&company)
}

func createFinancialRecord(db *gorm.DB, row map[string]string) {
	financial := Financial{
		Ticker:          row["ticker"],
		Ebitda:          nullableFloat(row["ebitda"]),
		Sales:           nullableFloat(row["sales"]),
		NetProfit:       nullableFloat(row["net_profit"]),
		MarketPrice:     nullableFloat(row["market_price"]),
		NetDebt:         nullableFloat(row["net_debt"]),
		Assets:          nullableFloat(row["assets"]),
		Equity:          nullableFloat(row["equity"]),
		CashEquivalents: nullableFloat(row["cash_equivalents"]),
		Liabilities:     nullableFloat(row["liabilities"]),
	}
	db.Create(&financial)
}

func insertIntoDB(db *gorm.DB, data []map[string]string, tableName string) {
	for _, row := range data {
		for key, value := range row {
			if value == "" {
				row[key] = "nil"
			}
		}

		switch tableName {
		case "companies":
			createCompanyRecord(db, row)
		case "financial":
			createFinancialRecord(db, row)
		}
	}
}

func nullableString(s string) *string {
	if s == "nil" {
		return nil
	}
	return &s
}

func nullableFloat(s string) *float64 {
	if s == "nil" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("failed to convert string to float64: %v", err)
	}
	return &f
}

func getUserInput(prompt string) string {
	var scanner = bufio.NewScanner(os.Stdin)
	fmt.Printf("%s:\n", prompt)
	scanner.Scan()
	return scanner.Text()
}

func getFloatInput(prompt string) float64 {
	var input float64
	fmt.Printf("%s (in the format '987654321'):\n", prompt)
	fmt.Scanln(&input)
	return input
}

func printOperationResult(operation string) {
	fmt.Printf("Company %s successfully!\n\n", operation)
}

func createCompany(db *gorm.DB) {
	ticker, name, sector := getCompanyData()
	ebitda, sales, netProfit, marketPrice, netDebt, assets, equity, cashEquivalents, liabilities := getFinancialData()
	company := Company{
		Ticker: ticker,
		Name:   &name,
		Sector: &sector,
	}
	db.Create(&company)

	financial := Financial{
		Ticker:          ticker,
		Ebitda:          &ebitda,
		Sales:           &sales,
		NetProfit:       &netProfit,
		MarketPrice:     &marketPrice,
		NetDebt:         &netDebt,
		Assets:          &assets,
		Equity:          &equity,
		CashEquivalents: &cashEquivalents,
		Liabilities:     &liabilities,
	}
	db.Create(&financial)
	printOperationResult("created")
}

func getFinancialData() (float64, float64, float64, float64, float64, float64, float64, float64, float64) {
	ebitda := getFloatInput("Enter ebitda")
	sales := getFloatInput("Enter sales")
	netProfit := getFloatInput("Enter net profit")
	marketPrice := getFloatInput("Enter market price")
	netDebt := getFloatInput("Enter net debt")
	assets := getFloatInput("Enter assets")
	equity := getFloatInput("Enter equity")
	cashEquivalents := getFloatInput("Enter cash equivalents")
	liabilities := getFloatInput("Enter liabilities")

	return ebitda, sales, netProfit, marketPrice, netDebt, assets, equity, cashEquivalents, liabilities
}

func getCompanyData() (string, string, string) {
	ticker := getUserInput("Enter ticker (in the format 'MOON')")
	name := getUserInput("Enter company (in the format 'Moon Corp')")
	sector := getUserInput("Enter industries (in the format 'Technology')")

	return ticker, name, sector
}

func readCompany(db *gorm.DB) {
	name := getUserInput(companyNamePrompt)

	company, done := findCompany(db, name)
	if done {
		return
	}

	fmt.Printf("%s %s\n", company.Ticker, *company.Name)

	financial := findFinancial(db, company)
	indicators := financial.CalculateIndicators()
	indicators.print()
}

func updateCompany(db *gorm.DB) {
	name := getUserInput(companyNamePrompt)

	company, done := findCompany(db, name)
	if done {
		return
	}
	financial := findFinancial(db, company)

	financialValues := []string{
		"ebitda",
		"sales",
		"net profit",
		"market price",
		"net debt",
		"assets",
		"equity",
		"cash equivalents",
		"liabilities",
	}

	for _, value := range financialValues {
		financial.UpdateValue(value)
	}
	db.Save(&financial)

	fmt.Printf("Company updated successfully!\n\n")
}

func findFinancial(db *gorm.DB, company Company) Financial {
	var financial Financial
	result := db.First(&financial, "ticker = ?", company.Ticker)
	if result.Error != nil {
		log.Fatalf("cannot retrieve financial: %v", result.Error)
	}
	return financial
}

func deleteCompany(db *gorm.DB) {
	name := getUserInput(companyNamePrompt)

	company, done := findCompany(db, name)
	if done {
		return
	}
	db.Delete(&company)

	printOperationResult("deleted")
}

func findCompany(db *gorm.DB, name string) (Company, bool) {
	var companies []Company
	result := db.Find(&companies, "name LIKE ?", "%"+name+"%")
	if result.Error != nil {
		log.Fatalf("cannot retrieve companies: %v", result.Error)
	}
	if len(companies) == 0 {
		fmt.Printf("Company not found!\n\n")
		return Company{}, true
	}
	for i, company := range companies {
		fmt.Printf("%d %s\n", i, *company.Name)
	}

	var index int
	fmt.Println("Enter company number:")
	fmt.Scanln(&index)

	company := companies[index]
	return company, false
}

func listAllCompanies(db *gorm.DB) {
	var companies []Company
	db.Order("ticker").Find(&companies)
	fmt.Println("COMPANY LIST")
	for _, company := range companies {
		fmt.Println(company.Ticker, *company.Name, *company.Sector)
	}
	fmt.Printf("\n")
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
}

func handleCrudOption(db *gorm.DB, option string) {
	switch option {
	case "0":
		fmt.Printf("\n")
	case "1":
		createCompany(db)
	case "2":
		readCompany(db)
	case "3":
		updateCompany(db)
	case "4":
		deleteCompany(db)
	case "5":
		listAllCompanies(db)
	default:
		fmt.Print(invalidOptionMsg)
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
}

func handleTopTenOption(option string) {
	switch option {
	case "0":
		fmt.Printf("\n")
	case "1":
		fmt.Print(notImplementedMsg)
	case "2":
		fmt.Print(notImplementedMsg)
	default:
		fmt.Print(invalidOptionMsg)
	}
}

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
}

func processUserInput(db *gorm.DB) {
	for {
		printMainMenu()

		var option string
		fmt.Scanln(&option)
		switch option {
		case "0":
			fmt.Print(byeMsg)
			os.Exit(0)
		case "1":
			printCrudMenu()

			var crudOption string
			fmt.Scanln(&crudOption)
			handleCrudOption(db, crudOption)
		case "2":
			printTopTenMenu()

			var topTenOption string
			fmt.Scanln(&topTenOption)
			handleTopTenOption(topTenOption)
		default:
			fmt.Print(invalidOptionMsg)
		}
	}
}

func checkExistingRecords(db *gorm.DB) {
	var companyRecordCount, financialRecordCount int64
	db.Model(&Company{}).Count(&companyRecordCount)
	db.Model(&Financial{}).Count(&financialRecordCount)

	if companyRecordCount == 0 || financialRecordCount == 0 {
		data := readCSVData("companies.csv")
		insertIntoDB(db, data, "companies")
		data = readCSVData("financial.csv")
		insertIntoDB(db, data, "financial")
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("investor.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Company{}, &Financial{})
	if err != nil {
		log.Fatal(err)
	}

	checkExistingRecords(db)

	fmt.Print(welcomeMsg)
	processUserInput(db)
}
