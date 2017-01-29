// an example of a basic RestAPI without a database. Data is stored in a global variable.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux" // a helper to make endpoints easier to create
	"log"
	"net/http"
	"strconv"
)

// Pokemon data structure
type Pokemon struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	EvolvesFrom string `json:"evolvesfrom, omitempty"`
	EvolvesInto string `json:"evolvesinto, omitempty"` //omitempty param means that if property is null  it will be excluded from the JSON rather than showing as empty`
}

// define a method on the pokemon struct that allows it's ID to be set
func (p *Pokemon) SetId(id string) {
	p.ID = id
}

// define a method on a pokemon to return it's name
func (p *Pokemon) ReturnName() string {
	return p.Name
}

// instead of using a database this example uses a global variable (slice of type Pokemon) to hold data used by app
var pokemondata []Pokemon

// The json package provides Decoder and Encoder types to support the common operation of readin and writing streams of JSON data.
// The NewDecoder and NewEncoder functions wrap the io.Reader and io.Writer interface types
// Encoder and Decoder types can be used for reading and writing to HTTP connections, WebSockets or files

// The GetAPokemon returns the pokemon that matches requested ID to the front end
func GetAPokemon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)              // using the mux library can get any parameters which were passed in with the request.
	for _, item := range pokemondata { // loop over the global slice and look for any ids that match the id found in the request parameter
		if item.ID == params["id"] { // if a match is found...
			json.NewEncoder(w).Encode(item) // use the JSON encoder on it
			return
		}
	}
  // previous for loop has not found the pokemon
	msg := fmt.Sprintf("404 Pokemon Not Found")
  http.Error(w, msg, http.StatusNotFound)
}

// The GetAllPokemon returns all pokemon to the front end
func GetAllPokemons(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(pokemondata)
}

// AddNewPokemon adds a new pokemon to saved pokemondata
func AddNewPokemon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	var newpokemon Pokemon
	newpokemon.ID = string(params["id"])
	pokemondata = append(pokemondata, newpokemon)
	json.NewEncoder(w).Encode(pokemondata)
}

func DeleteAPokemon(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	for index, item := range pokemondata {
		if item.ID == params["id"] {
			pokemondata = append(pokemondata[:index], pokemondata[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(pokemondata)
}

func UpdateAPokemon(w http.ResponseWriter, r *http.Request) {
	//get the params attached to the request
	params := mux.Vars(r)
	// create a variable of type pokemon to hold the pokemonupdate
	var pokemonupdate Pokemon
	// decode the body of the request and store at address of pokemonupdate
	_ = json.NewDecoder(r.Body).Decode(&pokemonupdate)
	// set the id in the pokemonupdate, this is the same as the param id passed from request
	pokemonupdate.SetId(params["id"])
	// the params id is a string, the index of pokemondata is an int
	// so convert to params id to int
	i, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println("error converting string to int")
	}
	// params start counting at 1, pokemon data starts counting at 0
	// so subtract 1 from the pokemondata index to replace with pokemonupdate
	pokemondata[i-1] = pokemonupdate

}

func main() {

	// add some default pokemons
	pokemondata = append(pokemondata, Pokemon{ID: "1", Name: "Lampent", Type: "Ghost/Fire", EvolvesFrom: "Litwick", EvolvesInto: "Chandelure"})
	pokemondata = append(pokemondata, Pokemon{ID: "2", Name: "Pikachu", Type: "Electric", EvolvesFrom: "Pichu", EvolvesInto: "Raichu"})
	pokemondata = append(pokemondata, Pokemon{ID: "3", Name: "Roselia", Type: "Grass/Poison", EvolvesFrom: "Budew", EvolvesInto: "Roserade"})
	// set up a new Gorilla mux
	router := mux.NewRouter()
	// Gorilla mux supports regex, [0-9]+ one or more 0 - 9 digits
	router.HandleFunc("/pokemon/", GetAllPokemons).Methods("GET")
	router.HandleFunc("/pokemon/{id:[0-9]+}", GetAPokemon).Methods("GET")
	router.HandleFunc("/pokemon/add/", AddNewPokemon).Methods("POST") // add doesnt need an id it assigns the next available id for new pokemon
	router.HandleFunc("/pokemon/update/{id:[0-9]+}", UpdateAPokemon).Methods("POST")
	router.HandleFunc("/pokemon/{id:[0-9]+}", DeleteAPokemon).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

	// To view all pokemons
	// http://localhost:8080/pokemon/

	// To view a specific pokemon
	// http://localhost:8080/pokemon/2

	// To add a new pokemon from command line
	// curl -X POST -d '{"name":"Munchlax","type":"Normal","evolvesinto":"Snorlax"}' http://localhost:8080/pokemon/add/

	// To delete an existing pokemon from command line
	//  curl -X DELETE http://localhost:8080/pokemon/4

	// To update an existing pokemon from commandline
	// curl -X POST -d '{"name":"Munchlax","type":"Normal","evolvesinto":"Snorlax"}' http://localhost:8080/pokemon/update/1

}
