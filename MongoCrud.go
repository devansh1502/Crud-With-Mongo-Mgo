// Package controller RestApiOrder
//
// The purpose of this application is to store and retrieve order
//
//
//
//     BasePath: /api
//     Version: 0.0.1
//     Contact: Devansh<devanshgupta1502@gmail.com>
//
//     Consumes:
//       - application/json
//
//     Produces:
//       - application/json
//
//
//
// swagger:meta
package main

//go:generate swagger generate spec -m -o ./swagger.json

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var conn *mgo.Session
var proColl *mgo.Collection

var p []Product

var findValFunc = findVal

var dberror bool

//Product model
// swagger:model
type Product struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	OrderID  int           `bson:"orderid"`
	Name     string        `bson:"name"`
	Price    float32       `bson:"price"`
	Quantity int           `bson:"quantity"`
	Status   bool          `bson:"status"`
}

//Connection Establish
func init() {
	conn, _ = mgo.Dial("mongodb://localhost:27017")
	proColl = conn.DB("demosession1").C("product")
	index := mgo.Index{
		Key:      []string{"orderid"},
		Unique:   true,
		DropDups: true,
	}
	err := proColl.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
}

// swagger:operation GET /result Orders getOrders
//
// Get Orders
//
// This method Retrieves the full list of Orders.
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: order response
//     schema:
//       "$ref": "#/definitions/Product"
//   '400':
//     description: Bad Request
//   '405':
//     description: Method Not Allowed, likely url is not correct
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation

//GetResult Method to retrive full Orders
func GetResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := proColl.Find(bson.M{}).All(&p)
	if err != nil || dberror {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)

}

// swagger:operation GET /result/{id} Orders getSingleOrders
//
// Get Single Order
//
// This method retrieves a single result from the list
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   description: id to retrieve specific order details
//   required: true
//   type: string
// responses:
//   '200':
//     description: Orders response
//     schema:
//       "$ref": "#/definitions/Product"
//   '405':
//     description: Method Not Allowed, likely url is not correct
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation

//GetSingleResult Retrieves a single entry from Product
func GetSingleResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := getSingleValFromDB(r)
	if err != nil {
		if err == mgo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(body)
}

func getSingleValFromDB(r *http.Request) (Product, error) {
	var body Product
	params := mux.Vars(r)
	orderid := params["orderid"]
	i, err := strconv.ParseInt(orderid, 10, 64)
	if err != nil {
		return body, err
	}
	err = proColl.Find(bson.M{"orderid": i}).One(&body)
	return body, err
}

// swagger:operation POST /result Orders postOrders
//
// Post Orders
//
// Save Orders
//
// ---
// produces:
// - application/json
// parameters:
// - name: order
//   in: body
//   description: orders to create data
//   required: true
//   schema:
//     "$ref": "#/definitions/Product"
// responses:
//   '201':
//     description: Orders response
//     schema:
//       "$ref": "#/definitions/Product"
//   '409':
//     description: Conflict
//   '405':
//     description: Method Not Allowed, likely url is not correct
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation

//PostResult This method is to add a new entry in the DB.
func PostResult(w http.ResponseWriter, r *http.Request) {
	var body Product
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&body)
	err := proColl.Insert(body)
	if err != nil || dberror {
		if mgo.IsDup(err) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(body)
	w.Write(b)
}

// swagger:operation DELETE /result/{id} Orders deleteOrders
//
// Delete Orders
//
// This method will delete a specific entry from Orders
//
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   description: id to delete the order
//   required: true
//   type: string
// responses:
//   '201':
//     description: Orders response
//     schema:
//       "$ref": "#/definitions/Product"
//   '409':
//     description: Conflict
//   '405':
//     description: Method Not Allowed, likely url is not correct
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation

//DeleteResult This method deletes a particular entry from DB.
func DeleteResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params)
	orderid := params["orderid"]
	i, _ := strconv.Atoi(orderid)
	proColl.Remove(bson.M{"orderid": i})
}

// swagger:operation PUT /result/{id} Orders updateOrders
//
// Update Orders
//
// This method will update a specific entry from Orders
//
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   description: id to update the order
//   required: true
//   type: string
// responses:
//   '201':
//     description: Orders response
//     schema:
//       "$ref": "#/definitions/Product"
//   '409':
//     description: Conflict
//   '405':
//     description: Method Not Allowed, likely url is not correct
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation

//UpdateResult This method updates a particular document entry in DB.
func UpdateResult(w http.ResponseWriter, r *http.Request) {
	var body Product
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&body)
	params := mux.Vars(r)
	orderid := params["orderid"]
	i, _ := strconv.Atoi(orderid)
	proColl.Update(bson.M{"orderid": i}, &body)
}

func getVals() error {
	err := findValFunc()
	if err != nil {
		return err
	}
	return nil
}

func findVal() error {
	return nil
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	//Route Handlers / Endpoints
	r.HandleFunc("/api/result", GetResult).Methods("GET")
	r.HandleFunc("/api/result/{orderid}", GetSingleResult).Methods("GET")
	r.HandleFunc("/api/result", PostResult).Methods("POST")
	r.HandleFunc("/api/result/{orderid}", DeleteResult).Methods("DELETE")
	r.HandleFunc("/api/result/{orderid}", UpdateResult).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", r))
}
