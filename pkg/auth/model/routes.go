package model

import (
	"encoding/json"
	"fmt"
	"fullstack-shop/pkg/kernel/model/repo/fake"
	"fullstack-shop/pkg/kernel/model/subj"
	"github.com/nicourrrn/littleLogger"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//file, _ := os.Create("token.log")
	//logger, _ := littleLogger.NewLogger(file)

	if r.Method != "POST" {
		http.Error(w, "Method POST allowed only", http.StatusMethodNotAllowed)
		return
	}
	var regData RegisterReq
	err := json.NewDecoder(r.Body).Decode(&regData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // TODO не использовать err
		return
	}
	repo, err := fake.Load()
	if err != nil {
		log.Println(err)
	}
	user, err := repo.GetByEmail(regData.Email)
	if err != nil {
		log.Println(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(regData.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	accLifeTime, _ := strconv.Atoi(os.Getenv("ACCESS_MIN"))
	accToken, err := GenToken(user.ID,
		time.Now().Add(time.Minute*time.Duration(accLifeTime)).Unix(),
		os.Getenv("ACCESS_KEY"))

	cookie := http.Cookie{
		Name:    "access_token",
		Value:   accToken,
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, &cookie)
	user.AccessToken = accToken
	repo.Save()

	refLifeTime, _ := strconv.Atoi(os.Getenv("REFRESH_MIN"))
	refToken, err := GenToken(user.ID,
		time.Now().Add(time.Minute*time.Duration(refLifeTime)).Unix(),
		os.Getenv("REFRESH_KEY"))
	json.NewEncoder(w).Encode(map[string]string{
		"refresh_token": refToken,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	// TODO remove this block
	file, _ := os.Create("token.log")
	logger, _ := littleLogger.NewLogger(file)

	if r.Method != "POST" {
		http.Error(w, "Method POST allowed only", http.StatusMethodNotAllowed)
		return
	}
	var regData RegisterReq
	err := json.NewDecoder(r.Body).Decode(&regData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // TODO не использовать err
		return
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Server error", http.StatusServiceUnavailable)
		return
	}
	//repo := fake.NewUserRepo()
	repo, err := fake.Load()
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Server error", http.StatusServiceUnavailable)
		return
	}
	user := subj.NewUser(regData.Email, string(passHash))
	user.LastVisit = time.Now()
	repo.Users = append(repo.Users, user)
	defer func() {
		err = repo.Save()
		if err != nil {
			log.Println(err)
		}
	}()
	accLifeTime, _ := strconv.Atoi(os.Getenv("ACCESS_MIN"))
	accToken, err := GenToken(user.ID,
		time.Now().Add(time.Minute*time.Duration(accLifeTime)).Unix(),
		os.Getenv("ACCESS_KEY"))

	cookie := http.Cookie{
		Name:    "access_token",
		Value:   accToken,
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, &cookie)
	user.AccessToken = accToken
	repo.Save()

	refLifeTime, _ := strconv.Atoi(os.Getenv("REFRESH_MIN"))
	refToken, err := GenToken(user.ID,
		time.Now().Add(time.Minute*time.Duration(refLifeTime)).Unix(),
		os.Getenv("REFRESH_KEY"))
	json.NewEncoder(w).Encode(map[string]string{
		"refresh_token": refToken,
	})

}

func GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method GET allowed only", http.StatusMethodNotAllowed)
		return
	}
	token, _ := r.Cookie("access_token")
	claim, _ := GetClaim(token.Value, os.Getenv("ACCESS_KEY"))
	users, err := fake.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	user, err := users.GetById(claim.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if user.AccessToken != token.Value {
		http.Error(w, "Not actual token", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"login":        user.Login,
		"email":        user.Email,
		"home":         user.HomeAddress,
		"access_token": user.AccessToken,
	})

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method GET allowed only", http.StatusMethodNotAllowed)
	}
	token, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accClaim, err := GetClaim(token.Value, os.Getenv("ACCESS_KEY"))
	if err != nil && !strings.HasPrefix(err.Error(), "token is expired") {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users, err := fake.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	user, err := users.GetById(accClaim.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	fmt.Println(user.AccessToken)
	fmt.Println(token.Value)
	if user.AccessToken != token.Value {
		http.Error(w, "Not actual token", http.StatusUnauthorized)
		return
	}

	var data map[string]string
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	refToken := ParseBearer(data["refresh_token"])
	if refToken == "" {
		http.Error(w, "...", http.StatusUnauthorized)
		return
	}
	refClaim, err := GetClaim(refToken, os.Getenv("REFRESH_KEY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if refClaim.ID != accClaim.ID {
		http.Error(w, "Tokens not equal", http.StatusUnauthorized)
		return
	}

	accLifeTime, _ := strconv.Atoi(os.Getenv("ACCESS_MIN"))
	accToken, err := GenToken(user.ID,
		time.Now().Add(time.Minute*time.Duration(accLifeTime)).Unix(),
		os.Getenv("ACCESS_KEY"))

	cookie := http.Cookie{
		Name:    "access_token",
		Value:   accToken,
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, &cookie)
	user.AccessToken = accToken
	user.LastVisit = time.Now()
	users.Save()

	refLifeTime, _ := strconv.Atoi(os.Getenv("REFRESH_MIN"))
	refToken, err = GenToken(user.ID,
		time.Now().Add(time.Minute*time.Duration(refLifeTime)).Unix(),
		os.Getenv("REFRESH_KEY"))
	json.NewEncoder(w).Encode(map[string]string{
		"refresh_token": refToken,
	})

}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		//if cookie.Expires.After(time.Now()){
		//	http.Error(w, "You need update your token", http.StatusUnauthorized)
		//	return
		//}
		claim, err := GetClaim(cookie.Value, os.Getenv("ACCESS_KEY"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		err = claim.Valid()
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Todo покрытить тестами
		users, err := fake.Load()
		user, err := users.GetById(claim.ID)
		if user.LastVisit.Add(time.Minute * 30).After(time.Now()) {
			user.AccessToken = ""
			users.Save()
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method POST allowed only", http.StatusMethodNotAllowed)
		return
	}
	token, _ := r.Cookie("access_token")
	claim, _ := GetClaim(token.Value, os.Getenv("ACCESS_KEY"))
	users, err := fake.Load()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	user, err := users.GetById(claim.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if user.AccessToken != token.Value {
		http.Error(w, "Not actual token", http.StatusUnauthorized)
		return
	}
	user.AccessToken = ""
	users.Save()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "You was logout")
}
