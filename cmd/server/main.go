

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)


func main(){

mux := http.NewServeMux()
mux.HandleFunc("GET / ", rootHandler)
mux.HandleFunc("GET /users", getAllUser)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)
mux.HandleFunc("GET /health" , healthHandler)
fmt.Println("Server listening to :8080")

err := http.ListenAndServe(":8080",mux)
if err!= nil{
	fmt.Println("Server error ",err)
	return

}
}

func rootHandler(w http.ResponseWriter , r *http.Request){
	w.Write([]byte("Welcome to Go API"))
}
func healthHandler (w http.ResponseWriter ,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
func getAllUser (w http.ResponseWriter, r *http.Request){
	cacheMutex.RLock()
	users:= make([]User, 0 , len(userCache))

	for _, user := range userCache{
		users= append(users, user)
	}

	cacheMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
   err := json.NewEncoder(w).Encode(users)
   if err != nil{

http.Error(w , "failed to encode users", http.StatusInternalServerError)
return
   }

}

func getUser(w http.ResponseWriter , r *http.Request){
	id, err := strconv.Atoi( r.PathValue("id"))
   if err != nil{
	http.Error(w , "invalid user id ", http.StatusBadRequest)
	return
   }

   cacheMutex.RLock()
   user , ok := userCache[id]
   cacheMutex.RUnlock()
 if !ok {
	http.Error(w, "user not found",http.StatusNotFound)
	return

 } 
w.Header().Set("Content-Type","application/json")

 err = json.NewEncoder(w).Encode(user)
 if err !=nil{
	http.Error(w,"failed to encode user", http.StatusInternalServerError)
	return

 }

}


func createUser(w http.ResponseWriter , r *http.Request){
	defer r.Body.Close()
  var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err!= nil{
		http.Error(w , "invalid JSON body ",http.StatusBadRequest)
		return
	}
	if user.Name ==""{
		http.Error(w,"name is required",http.StatusBadRequest)
		return
	}
	// creating user mechanism 
	cacheMutex.Lock()
	user.ID =nextID
	userCache[nextID] = user
   nextID ++
   cacheMutex.Unlock()


   w.Header().Set("Content-Type","application/json")
   w.WriteHeader(http.StatusCreated)


   err = json.NewEncoder(w).Encode(user)
   if err != nil{
	http.Error(w, "failed to encode response", http.StatusInternalServerError)
	return
   }
	}

func deleteUser(w http.ResponseWriter  , r *http.Request){
 id , err := strconv.Atoi(r.PathValue("id"))
 if err != nil{
	http.Error(w, "invalid id ", http.StatusBadRequest)
	return

 }
 cacheMutex.Lock()
 defer cacheMutex.Unlock()

 if _,ok := userCache[id];!ok{
	http.Error(w, "user not found", http.StatusNotFound)
	return

 }
 delete(userCache, id)

 w.WriteHeader(http.StatusNoContent)


}