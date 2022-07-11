package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"postgres/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int64  `json:"status"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading in .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	/* check for connection*/
	err = db.Ping()
	if err != nil {
		//panic("Failure to open database connection")
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database")
	return db
}

func CreateCountry(w http.ResponseWriter, r *http.Request) {
	var country models.Country

	err := json.NewDecoder(r.Body).Decode(&country)
	if err != nil {

		log.Fatalf("unable to decode the requested body . %v", err)
	}
	insertID := insertCountry(country)
	if country.Name == "" {
		//panic("Name field is missing")
		res := response{
			Message: "Name Field is missing",
			Status:  400,
		}

		json.NewEncoder(w).Encode(res)

	} else if country.Type == "" {
		//panic("Type field is missing")
		res := response{
			Message: "Type Field is Missing",
			Status:  400,
		}
		json.NewEncoder(w).Encode(res)
	} else {
		res := response{
			ID:      insertID,
			Message: "COUNTRY INSERTED SUCCESSFULLY",
			Status:  200,
		}
		json.NewEncoder(w).Encode(res)
	}

}

func GetCountry(w http.ResponseWriter, r *http.Request) {
	country, err := getallcountry()
	if err != nil {
		log.Fatalf("unable to get all country list %v", err)

	}

	json.NewEncoder(w).Encode(country)

}

func UpdateCountry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert the string into int %v", err)
	}

	var country models.Country

	err = json.NewDecoder(r.Body).Decode(&country)
	if err != nil {
		log.Fatalf("unable to decode the request bode %v", err)
	}

	updatedrows := updatecountry(int64(id), country)

	msg := fmt.Sprintf("COUNTRY UPDATED SUCCESSFULLY.TOTAL ROWS AFFECTED %v", updatedrows)

	res := response{
		ID:      updatedrows,
		Message: msg,
		Status:  200,
	}
	json.NewEncoder(w).Encode(res)

}

func DeleteCountry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to decode the request bode %v", err)
	}
	deleteddata := deleteCountry(int64(id))

	msg := fmt.Sprintf("COUNTRY DELETED  SUCCESSFULLY.TOTAL ROWS AFFECTED %v", deleteddata)

	res := response{
		ID:      deleteddata,
		Message: msg,
		Status:  200,
	}
	json.NewEncoder(w).Encode(res)

}
func insertCountry(country models.Country) int64 {

	db := CreateConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO category (name, type) VALUES ($1, $2) RETURNING categoryid`

	var id int64

	err := db.QueryRow(sqlStatement, country.Name, country.Type).Scan(&id)
	if err != nil {
		log.Fatalf("unable to execute the query %v", err)
	}

	// if country.Name == "" {
	// 	panic("Name field is missing")
	// } else if country.Type == "" {
	// 	panic("Type field is missing")
	//} else {
	fmt.Printf("inserted single record %v", id)
	return id
	//}

}

func getallcountry() ([]models.Country, error) {
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	var countries []models.Country

	// create the select sql query
	sqlStatement := `SELECT * FROM category`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var Country models.Country

		// unmarshal the row object to stock
		err = rows.Scan(&Country.CountryID, &Country.Name, &Country.Type)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the stock in the stocks slice
		countries = append(countries, Country)

	}

	// return empty stock on error
	return countries, err
}
func updatecountry(id int64, country models.Country) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `UPDATE category SET name=$2, type=$3  WHERE categoryid=$1`
	res, err := db.Exec(sqlStatement, id, country.Name, country.Type)

	if err != nil {
		log.Fatalf("unable to execute the query %v", err)
	}

	rowsaffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Errror while checking the affected rows %v", err)
	}
	fmt.Printf("TOTAL ROWS AFFECTED %v", rowsaffected)
	return rowsaffected

}

func deleteCountry(id int64) int64 {

	db := CreateConnection()
	defer db.Close()

	query := `DELETE FROM category  WHERE categoryid =$1`

	res, err := db.Exec(query, id)

	if err != nil {
		log.Fatalf("unable to execute the query %v", err)

	}

	rowsaffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("ERROR WHILE CHECKING THE AFFECTED ROWS %v", rowsaffected)
	}

	fmt.Printf("total rows affected %v", rowsaffected)
	return rowsaffected

}
