package main

// Store dayly heatpump values into my database: heatpump.db
// Mo 17. Okt 15:02:58 CEST 2016

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const heatpumpDB = "/home/mhusmann/Documents/src/pyt/heizung/heatpump.db"

// const heatpumpDB = "heatpump.db"
const shortForm = "2006-01-02"

type entry struct {
	ID      int64
	Day     *string
	Ht      int64
	Nt      int64
	TarifID int64
}

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func lastEntry(db *sql.DB) entry {
	sqlReadMaxID := "SELECT MAX(id), day, ht," +
		"nt, tarifid FROM dayly"
	rows, err := db.Query(sqlReadMaxID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := entry{}
	for rows.Next() {
		err2 := rows.Scan(&result.ID, &result.Day,
			&result.Ht, &result.Nt, &result.TarifID)
		if err2 != nil {
			panic(err2)
		}
	}
	return result
}

func keyboardEntry(what string) string {
	fmt.Printf("Eingabe oder Abbruch (A) %s\t", what)
	var input string
	fmt.Scanln(&input)
	if strings.HasPrefix(input, "A") {
		fmt.Println("Programm Abbruch")
		os.Exit(1)
	}
	return strings.TrimSpace(input)
}

func isValid(day string) bool {
	_, err := time.Parse(shortForm, day)
	if err != nil {
		return false
	}
	return true
}

func getDay(day, last string) string {
	var today string
	for {
		today = keyboardEntry(day)
		if isValid(today) && today > last {
			break
		}
		fmt.Printf("%s Vermutlich falsches Datum\n", today)
	}
	return today
}

func getInt(num string, last int64) int64 {
	var validInt = regexp.MustCompile(`\d{5}`)
	var value int64
	var input string
	var err error
	for {
		input = keyboardEntry(num)
		value, err = strconv.ParseInt(input, 10, 32)
		if err != nil {
			panic(err)
		}
		if validInt.MatchString(input) && value >= last {
			break
		}
		fmt.Printf("%s Vermutlich falscher Wert\n", input)
	}
	return value
}
func storeData(db *sql.DB, newRow entry) bool {
	sqlStoreRow := `INSERT INTO dayly(id, day, ht, nt,
                        tarifid) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(sqlStoreRow)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(newRow.ID, newRow.Day,
		newRow.Ht, newRow.Nt, 1)
	if err2 != nil {
		panic(err2)
	}
	return true
}

func main() {
	t := time.Now()
	today := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	fmt.Printf("Eintrag für die DB Wärmepumpe am %s Datenbank ist %s\n",
		today, heatpumpDB)
	db := initDB(heatpumpDB)
	defer db.Close()
	res := lastEntry(db)
	fmt.Printf("Letzter Eintrag  ID: %d, DAY:  "+
		"%s HT: %d NT: %d TARIFID: %d\n",
		res.ID, *res.Day, res.Ht, res.Nt, res.TarifID)

	var newRow entry
	for {
		day := getDay("Day", *res.Day)
		ht := getInt("Ht", res.Ht)
		nt := getInt("Nt", res.Nt)
		newRow = entry{res.ID + 1, &day, ht, nt, 1}
		diffHt := newRow.Ht - res.Ht
		diffNt := newRow.Nt - res.Nt
		fmt.Printf("Neue Daten. Tag: %s, Diff Ht %d, Diff Nt %d, "+
			"Total: %d\n", *newRow.Day, diffHt, diffNt,
			diffHt+diffNt)
		r := keyboardEntry("ist das korrekt (j/n)")
		if r == "j" {
			break
		}
	}
	fmt.Println("Schreibe neue Werte in die Datenbank")
	storeData(db, newRow)
}
