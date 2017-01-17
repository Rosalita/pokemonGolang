// an example of a basic RestAPI without a database. Data is stored in a global variable.
package main

import(
  "encoding/json"
  "log"
  "net/http"
  "github.com/gorilla/mux" // a helper to make endpoints easier to create
  "fmt"
)

//data structures
type Person struct {
  ID          string    `json:"id,omitempty"` //omitempty param means that if property is null  it will be excluded from the JSON rather than showing as empty
  Firstname   string    `json:"firstname, omitempty"`
  Lastname    string    `json:"lastname, omitempty"`
  Address     *Address  `json:"address, omitempty"` //nested json must be a pointer or omitempty will fail
}

type Address struct{
  City       string   `json:"city, omitempty"`
  State      string   `json:"state, omitempty"`
}

// instead of using a database this example usesa global variable (slice of type Person) to hold data used by app
var people []Person

// The json package provides Decoder and Encoder types to support the common operation of readin and writing streams of JSON data.
// The NewDecoder and NewEncoder functions wrap the io.Reader and io.Writer interface types
// Encoder and Decoder types can be used for reading and writing to HTTP connections, WebSockets or files

// The GetPeopleEndpoint is probably the easiest to understand because it returns all data to frontend.
// The GetPersonEndpoing returns a full person variable to the frontend
func GetPersonEndpoint(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r) // using the mux library can get any parameters which were passed in with the request.
  fmt.Printf("params are: %v", params)
  for _, item := range people { // loop over the global slice and look for any ids that match the id found in the request parameter
    if item.ID == params["id"] { // if a match is found...
      json.NewEncoder(w).Encode(item) // use the JSON encoder on it
      return
    }
  }
  json.NewEncoder(w).Encode(&Person{})
}

// Return all people
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request){
  json.NewEncoder(w).Encode(people)
}

func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r)
  var person Person
  _ = json.NewDecoder(r.Body).Decode(&person)
  person.ID = params["id"]
  people = append(people, person)
  json.NewEncoder(w).Encode(people)
}

func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r)
  for index, item := range people{
    if item.ID == params["id"]{
      people = append(people[:index], people[index+1:]...)
      break
    }
  }
  json.NewEncoder(w).Encode(people)
}

func main(){
  router := mux.NewRouter()
  people = append(people, Person{ID: "1", Firstname: "Rosie", Lastname: "Hamilton", Address: &Address{City: "Newcastle", State: "Tyne and Wear"} })
  people = append(people, Person{ID: "2", Firstname: "Jane", Lastname: "Doe"})
  router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
  router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
  router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
  router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
  log.Fatal(http.ListenAndServe(":8080", router))
}
