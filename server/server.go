package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json: "id"`
	Name  string `json: "name"`
	Email string `json: "email"`
}

//CreateUser insert a user in the database.
func CreateUser(rw http.ResponseWriter, r *http.Request) {
	requisitionBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.Write([]byte("Failed to read the requisition body!"))
		return
	}
	var user user
	if err = json.Unmarshal(requisitionBody, &user); err != nil {
		rw.Write([]byte("Failed to convert user to struct!"))
	}

	db, err := database.Connect()
	if err != nil {
		rw.Write([]byte("Failed to connect with database!"))
		return
	}
	defer db.Close()

	//prepare statement
	statement, err := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if err != nil {
		rw.Write([]byte("Failed to create the statement!"))
		return
	}
	defer statement.Close()

	insertion, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		rw.Write([]byte("Failed to execute the statement!"))
		return
	}

	enteredId, err := insertion.LastInsertId()
	if err != nil {
		rw.Write([]byte("Failed to obtain the entered ID!"))
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fmt.Sprintf("User entered successfully! ID: %d", enteredId)))
}

//SearchUsers brings all the users saved on the database.
func SearchUsers(rw http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		rw.Write([]byte("Failed to connect with database!"))
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from usuarios")
	if err != nil {
		rw.Write([]byte("Failed to search users!"))
		return
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var user user
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			rw.Write([]byte("Failed scanning the user!"))
			return
		}
		users = append(users, user)
	}
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(users); err != nil {
		rw.Write([]byte("Failed converting users to JSON!"))
		return
	}
}

//SearchUser brings a specific user saved on the database.
func SearchUser(rw http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		rw.Write([]byte("Failed to convert parameter to int"))
		return
	}
	db, err := database.Connect()
	if err != nil {
		rw.Write([]byte("Failed to connect with database!"))
		return
	}
	row, err := db.Query("select * from usuarios where id = ?", ID)
	if err != nil {
		rw.Write([]byte("Failed to search the user!"))
		return
	}
	var user user
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			rw.Write([]byte("Failed to scan the user!"))
			return
		}
	}
	if err := json.NewEncoder(rw).Encode(user); err != nil {
		rw.Write([]byte("Failed convert user to JSON!"))
	}
}

//UpdateUser update the information of a user.
func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		rw.Write([]byte("Error to convert parameter to integer!"))
		return
	}
	requisitionBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.Write([]byte("Error reading the requisition body!"))
		return
	}
	var user user
	if err := json.Unmarshal(requisitionBody, &user); err != nil {
		rw.Write([]byte("Error to convert user to struct!"))
		return
	}
	db, err := database.Connect()
	if err != nil {
		rw.Write([]byte("Error to connect with database!"))
		return
	}
	defer db.Close()
	statement, err := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if err != nil {
		rw.Write([]byte("Error statement!"))
		return
	}
	defer statement.Close()
	if _, err := statement.Exec(user.Name, user.Email, ID); err != nil {
		rw.Write([]byte("Failed update user!"))
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		rw.Write([]byte("Failed to convert parameter to integer!"))
		return
	}
	db, err := database.Connect()
	if err != nil {
		rw.Write([]byte("Failed to connect to the database!"))
		return
	}
	defer db.Close()
	statement, err := db.Prepare("delete from usuarios where id = ?")
	if err != nil {
		rw.Write([]byte("Failed create the statement!"))
		return
	}
	defer statement.Close()
	if _, err := statement.Exec(ID); err != nil {
		rw.Write([]byte("Failed to delete the user!"))
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
