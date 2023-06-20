package main

/*
[Calculator for Investors - Stage 2/4: Store it](https://hyperskill.org/projects/264/stages/1336/implement)
-------------------------------------------------------------------------------
[Maps](https://hyperskill.org/learn/step/16999)
[Parsing data from strings](https://hyperskill.org/learn/step/17935)
[Design principles](https://hyperskill.org/learn/step/8956)
[Single Responsibility Principle](https://hyperskill.org/learn/step/8963)
[Function decomposition](https://hyperskill.org/learn/topic/1893)
[Advanced Input](https://hyperskill.org/learn/topic/2027)
[Errors](https://hyperskill.org/learn/topic/1795)
[Public and private scopes](https://hyperskill.org/learn/topic/1894)
[Structs](https://hyperskill.org/learn/topic/1891)
[Methods](https://hyperskill.org/learn/topic/1928)
[CSV](https://hyperskill.org/learn/step/13164)
[Reading files](https://hyperskill.org/learn/step/16702)
[Debugging Go code](https://hyperskill.org/learn/step/23076)
[Introduction to GORM](https://hyperskill.org/learn/step/20695)
[Declaring GORM Models](https://hyperskill.org/learn/step/28639)
[Relationships between models](https://hyperskill.org/learn/step/29207)
[Migrations](https://hyperskill.org/learn/step/22043)
[CRUD Operations — Create](https://hyperskill.org/learn/step/22859)
*/

import (
	"encoding/csv"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
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

func (Financial) TableName() string {
	return "financial"
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
	fmt.Println("Database created successfully!")
}
