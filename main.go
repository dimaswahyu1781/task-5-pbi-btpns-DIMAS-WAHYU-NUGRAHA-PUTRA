package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Tabel
type Photo struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserID   string `json:"user_id"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/bptn?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Koneksi Gagal! error: ", err)
	} else {
		log.Println("Koneksi Berhasil!")
	}

	db.AutoMigrate(&Photo{})

	handleRequests()
}

// routing
func handleRequests() {
	log.Println("Start the development server localhost:2023")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/photos", createPhoto).Methods("POST")
	myRouter.HandleFunc("/users", createUser).Methods("POST")
	myRouter.HandleFunc("/photos/{id}", getPhotoByID).Methods("GET")
	myRouter.HandleFunc("/photosbyuserid/{id}", getPhotosByUserID).Methods("GET")
	myRouter.HandleFunc("/photos/{id}", updatePhoto).Methods("PUT")
	myRouter.HandleFunc("/photos/{id}", deletePhoto).Methods("DELETE")
	myRouter.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":2023", myRouter))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
}

// createUser function
func createUser(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	payloads, _ := ioutil.ReadAll(r.Body)
	// Unmarshal the request body into a user struct
	var user User
	json.Unmarshal(payloads, &user)
	// Check if the username is already taken
	var dbUser User
	db.Where("username = ?", user.Username).First(&dbUser)
	if dbUser.ID != 0 {
		// Write an error response with status code 409 Conflict
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Username already exists"))
		return
	}
	// Create the user in the database
	db.Create(&user)
	// Write a success response with status code 201 Created and the created user in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"user": %s}`, user)))
}

// getUserByID function
func getUserByID(w http.ResponseWriter, r *http.Request) {
	// Get the user id from the URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]
	// Find the user in the database by id
	var user User
	db.First(&user, userID)
	// Check if the user exists
	if user.ID == 0 {
		// Write an error response with status code 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	// Write a success response with status code 200 OK and the found user in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"user": %s}`, user)))
}

// getUsers function
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Find all users in the database
	var users []User
	db.Find(&users)
	// Write a success response with status code 200 OK and the users in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"users": %s}`, users)))
}

// updateUser function
func updateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user id from the URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]
	// Read the request body
	payloads, _ := ioutil.ReadAll(r.Body)
	// Unmarshal the request body into a user struct
	var userUpdates User
	json.Unmarshal(payloads, &userUpdates)
	// Find the user in the database by id
	var user User
	db.First(&user, userID)
	// Check if the user exists
	if user.ID == 0 {
		// Write an error response with status code 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	// Update the user in the database with the new data
	db.Model(&user).Update(userUpdates)
	// Write a success response with status code 200 OK and the updated user in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"user": %s}`, user)))
}

// deleteUser function
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the user id from the URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]
	// Find the user in the database by id
	var user User
	db.First(&user, userID)
	// Check if the user exists
	if user.ID == 0 {
		// Write an error response with status code 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	// Delete the user from the database
	db.Delete(&user)
	// Write a success response with status code 200 OK and a message in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User berhasil dihapus!"}`))
}

// GenerateToken function
func GenerateToken(user User) (string, error) {
	// Define the secret key for signing the token
	var secretKey = []byte("secret")
	// Create a new token object with the signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})
	// Sign the token with the secret key and return the string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken function
func VerifyToken(tokenString string) (User, error) {
	// Define the secret key for verifying the token
	var secretKey = []byte("secret")
	// Parse the token string with the secret key and the claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is the same as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return secretKey, nil
	})
	if err != nil {
		return User{}, err
	}
	// Check if the token is valid and get the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the id and username from the claims
		id := int(claims["id"].(float64))
		username := claims["username"].(string)
		// Create a user struct with the id and username
		user := User{
			ID:       id,
			Username: username,
		}
		// Return the user struct and nil error
		return user, nil
	} else {
		// Return an empty user struct and an error
		return User{}, fmt.Errorf("invalid token")
	}
}

// AuthMiddleware function
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header value from the request
		authHeader := r.Header.Get("Authorization")
		// Check if the header value is empty or does not start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// Write an error response with status code 401 Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing or invalid Authorization header"))
			return
		}
		// Get the token string by removing the "Bearer " prefix
		tokenString := authHeader[7:]
		// Verify the token string and get the user data
		user, err := VerifyToken(tokenString)
		if err != nil {
			// Write an error response with status code 401 Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		// Create a new context with the user data
		ctx := context.WithValue(r.Context(), "user", user)
		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// createPhoto function
func createPhoto(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	payloads, _ := ioutil.ReadAll(r.Body)
	// Unmarshal the request body into a photo struct
	var photo Photo
	json.Unmarshal(payloads, &photo)
	// Check if the user id is valid
	var user User
	db.First(&user, photo.UserID)
	if user.ID == 0 {
		// Write an error response with status code 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user id"))
		return
	}
	// Create the photo in the database
	db.Create(&photo)
	// Write a success response with status code 201 Created and the created photo in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"photo": %s}`, photo)))
}

// getPhotoByID function
func getPhotoByID(w http.ResponseWriter, r *http.Request) {
	// Get the photo id from the URL parameter
	vars := mux.Vars(r)
	photoID := vars["id"]
	// Find the photo in the database by id
	var photo Photo
	db.First(&photo, photoID)
	// Check if the photo exists
	if photo.ID == 0 {
		// Write an error response with status code 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Photo not found"))
		return
	}
	// Write a success response with status code 200 OK and the found photo in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"photo": %s}`, photo)))
}

// getPhotosByUserID function
func getPhotosByUserID(w http.ResponseWriter, r *http.Request) {
	// Get the user id from the URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]
	// Find all photos in the database that belong to the user id
	var photos []Photo
	db.Where("user_id = ?", userID).Find(&photos)
	// Write a success response with status code 200 OK and the photos in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"photos": %s}`, photos)))
}

// updatePhoto function
func updatePhoto(w http.ResponseWriter, r *http.Request) {
	// Get the photo id from the URL parameter
	vars := mux.Vars(r)
	photoID := vars["id"]
	// Read the request body
	payloads, _ := ioutil.ReadAll(r.Body)
	// Unmarshal the request body into a photo struct
	var photoUpdates Photo
	json.Unmarshal(payloads, &photoUpdates)
	// Find the photo in the database by id
	var photo Photo
	db.First(&photo, photoID)
	// Check if the photo exists
	if photo.ID == 0 {
		// Write an error response with status code 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Photo not found"))
		return
	}
	// Update the photo in the database with the new data
	db.Model(&photo).Update(photoUpdates)
	// Write a success response with status code 200 OK and the updated photo in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"photo": %s}`, photo)))
}

// deletePhoto function
func deletePhoto(w http.ResponseWriter, r *http.Request) {
	// Get the photo id from the URL parameter
	vars := mux.Vars(r)
	photoID := vars["id"]
	// Find the photo in the database by id
	var photo Photo
	db.First(&photo, photoID)
	// Check if the photo exists
	if photo.ID == 0 {
		// Write an error response with status code 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Photo not found"))
		return
	}
	// Delete the photo from the database
	db.Delete(&photo)
	// Write a success response with status code 200 OK and a message in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Photo berhasil dihapus!"}`))
}
