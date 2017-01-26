// an example of a basic RestAPI without a database. Data is stored in a global variable.
package main

import(
  "encoding/json"
  "log"
  "net/http"
  "github.com/gorilla/mux" // a helper to make endpoints easier to create
)

// Pokemon data structure
type Pokemon struct {
  ID            string    `json:"id"`
  Name          string    `json:"name"`
  Type          string    `json:"type"`
  EvolvesFrom   string    `json:"evolvesfrom, omitempty"`
  EvolvesInto   string    `json:"evolvesinto, omitempty"` //omitempty param means that if property is null  it will be excluded from the JSON rather than showing as empty`
}


// instead of using a database this example uses a global variable (slice of type Pokemon) to hold data used by app
var pokemondata []Pokemon

// The json package provides Decoder and Encoder types to support the common operation of readin and writing streams of JSON data.
// The NewDecoder and NewEncoder functions wrap the io.Reader and io.Writer interface types
// Encoder and Decoder types can be used for reading and writing to HTTP connections, WebSockets or files

// The GetAPokemon returns the pokemon that matches requested ID to the front end
func GetAPokemon(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r) // using the mux library can get any parameters which were passed in with the request.
  for _, item := range pokemondata { // loop over the global slice and look for any ids that match the id found in the request parameter
    if item.ID == params["id"] { // if a match is found...
      json.NewEncoder(w).Encode(item) // use the JSON encoder on it
      return
    }
  }
}

// The GetAllPokemon returns all pokemon to the front end
func GetAllPokemons(w http.ResponseWriter, r *http.Request){
  json.NewEncoder(w).Encode(pokemondata)
}


// AddNewPokemon adds a new pokemon to saved pokemondata
func AddNewPokemon(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r)
  var newpokemon Pokemon
  _ = json.NewDecoder(r.Body).Decode(&newpokemon)
  newpokemon.ID = string(params["id"])
  pokemondata = append(pokemondata, newpokemon)
  json.NewEncoder(w).Encode(pokemondata)
}

func DeleteAPokemon(w http.ResponseWriter, r *http.Request){
  params := mux.Vars(r)
  for index, item := range pokemondata{
    if item.ID == params["id"]{
      pokemondata = append(pokemondata[:index], pokemondata[index+1:]...)
      break
    }
  }
  json.NewEncoder(w).Encode(pokemondata)
}

func main(){

  // add some default pokemons
  pokemondata = append(pokemondata, Pokemon{ID: "1", Name: "Lampent", Type: "Ghost/Fire", EvolvesFrom: "Litwick", EvolvesInto: "Chandelure"})
  pokemondata = append(pokemondata, Pokemon{ID: "2", Name: "Pikachu", Type: "Electric", EvolvesFrom: "Pichu", EvolvesInto: "Raichu"})
  pokemondata = append(pokemondata, Pokemon{ID: "3", Name: "Roselia", Type: "Grass/Poison", EvolvesFrom: "Budew", EvolvesInto: "Roserade"})

  // set up a new Gorilla mux
  router := mux.NewRouter()

  // Gorilla mux supports regex, d+ one or more digits
  router.HandleFunc("/pokemon/", GetAllPokemons).Methods("GET")
  router.HandleFunc("/pokemon/{id:[0-9]+}", GetAPokemon).Methods("GET")
  router.HandleFunc("/pokemon/{id:[0-9]+}", AddNewPokemon).Methods("POST")
  router.HandleFunc("/pokemon/{id:[0-9]+}", DeleteAPokemon).Methods("DELETE")
  log.Fatal(http.ListenAndServe(":8080", router))

// To view all pokemons
// http://localhost:8080/pokemon/

// To view a specific pokemon
// http://localhost:8080/pokemon/2

// To add a new pokemon from command line
// curl -X POST -d '{"name":"Munchlax","type":"Normal","evolvesinto":"Snorlax"}' http://localhost:8080/pokemon/4

// To delete an existing pokemon from command line
//  curl -X DELETE http://localhost:8080/pokemon/4


}
