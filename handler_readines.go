package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name string
		Age  int
	}

	names := [5]string{"nurbek", "adam", "amir", "ulug", "samir"}
	ages := [5]int{10, 20, 30, 40, 50}
	responseSlice := make([]User, 0)

	for i := 0; i < len(names); i++ {
		user := User{
			Name: names[i],
			Age:  ages[i],
		}
		responseSlice = append(responseSlice, user)
	}
	respondWithJSON(w, 200, responseSlice)
}
