package main

// imports

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PlayerRequest struct {
	Nickname string
	Life     int
	Attack   int
}

type Enemy struct {
	Nickname string
	Life     int
	Attack   int
}

type Battle struct {
	ID         string
	Enemy      string
	Player     string
	DiceThrown int
	Winner     string
}

type Response struct {
	Message string `json:"message"`
}

var players []PlayerRequest
var enemies []Enemy
var battles []Battle

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal Server Error"})
		return
	}

	if playerRequest.Nickname == "" || playerRequest.Life == 0 || playerRequest.Attack == 0 {
		json.NewEncoder(w).Encode(Response{Message: "Player nickname, life and attack is required"})
		return
	}

	if playerRequest.Attack > 10 || playerRequest.Attack <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Player attack must be between 1 and 10"})
		return
	}

	if playerRequest.Life > 10 || playerRequest.Life <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Player life must be between 1 and 10"})
		return
	}

	for _, player := range players {
		if player.Nickname == playerRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Message: "Player nickname already exists"})
			return
		}
	}

	player := PlayerRequest{
		Nickname: playerRequest.Nickname,
		Life:     playerRequest.Life,
		Attack:   playerRequest.Attack,
	}
	players = append(players, player)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

func LoadPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]

	for i, player := range players {
		if player.Nickname == nickname {
			players = append(players[:i], players[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Player nickname not found"})
}

func LoadPlayerByNickname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]

	for _, player := range players {
		if player.Nickname == nickname {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(player)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Player nickname not found"})
}

func SavePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nickname := mux.Vars(r)["nickname"]

	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal Server Error"})
		return
	}

	if playerRequest.Nickname == "" {
		json.NewEncoder(w).Encode(Response{Message: "Player nickname is required"})
		return
	}

	indexPlayer := -1
	for i, player := range players {
		if player.Nickname == nickname {
			indexPlayer = i
		}
		if player.Nickname == playerRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Message: "Player nickname already exists"})
			return
		}
	}

	if indexPlayer != -1 {
		players[indexPlayer].Nickname = playerRequest.Nickname
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(players[indexPlayer])
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{Message: "Player nickname not found"})
}

func AddEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var enemy Enemy
	if err := json.NewDecoder(r.Body).Decode(&enemy); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal Server Error"})
		return
	}

	if enemy.Nickname == "" {
		json.NewEncoder(w).Encode(Response{Message: "Enemy nickname is required"})
		return
	}

	for _, e := range enemies {
		if e.Nickname == enemy.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Message: "Enemy nickname already exists"})
			return
		}
	}

	rand.Seed(time.Now().UnixNano())
	enemy.Life = rand.Intn(10) + 1
	enemy.Attack = rand.Intn(10) + 1

	enemies = append(enemies, enemy)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enemy)
}

func LoadEnemies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enemies)
}

// Pega nickname para deletar

func DeleteEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]

	for i, enemy := range enemies {
		if enemy.Nickname == nickname {
			enemies = append(enemies[:i], enemies[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Enemy nickname not found"})
}

func LoadEnemyByNickname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]

	for _, enemy := range enemies {
		if enemy.Nickname == nickname {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(enemy)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Enemy nickname not found"})
}

//Salva Enemy e valida se ele existe

func SaveEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nickname := mux.Vars(r)["nickname"]

	var enemyRequest Enemy
	if err := json.NewDecoder(r.Body).Decode(&enemyRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal Server Error"})
		return
	}

	if enemyRequest.Nickname == "" {
		json.NewEncoder(w).Encode(Response{Message: "Enemy nickname is required"})
		return
	}

	indexEnemy := -1
	for i, enemy := range enemies {
		if enemy.Nickname == nickname {
			indexEnemy = i
		}
		if enemy.Nickname == enemyRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Message: "Enemy nickname already exists"})
			return
		}
	}

	if indexEnemy != -1 {
		enemies[indexEnemy].Nickname = enemyRequest.Nickname
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(enemies[indexEnemy])
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{Message: "Enemy nickname not found"})
}

func addBatle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var battleRequest struct {
		Enemy  string `json:"enemy"`
		Player string `json:"player"`
	}
	if err := json.NewDecoder(r.Body).Decode(&battleRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal Server Error"})
		return
	}

	var player *PlayerRequest
	var enemy *Enemy

	for i := range players {
		if players[i].Nickname == battleRequest.Player {
			player = &players[i]
			break
		}
	}

	for i := range enemies {
		if enemies[i].Nickname == battleRequest.Enemy {
			enemy = &enemies[i]
			break
		}
	}

	if player == nil || enemy == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Player or Enemy not found"})
		return
	}

	if player.Life <= 0 || enemy.Life <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Player or Enemy has no life remaining"})
		return
	}

	battleID := uuid.New().String()
	var winner string

	for player.Life > 0 && enemy.Life > 0 {
		rand.Seed(time.Now().UnixNano())
		dice := rand.Intn(6) + 1

		if dice <= 3 {
			player.Life -= enemy.Attack
			winner = enemy.Nickname
		} else {
			enemy.Life -= player.Attack
			winner = player.Nickname
		}

		battle := Battle{
			ID:         battleID,
			Enemy:      battleRequest.Enemy,
			Player:     battleRequest.Player,
			DiceThrown: dice,
			Winner:     winner,
		}
		battles = append(battles, battle)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(battles)
}

func LoadBattle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(battles)
}

func main() {
	r := mux.NewRouter()

	//Rotas do Player
	r.HandleFunc("/player", AddPlayer).Methods("POST")

	r.HandleFunc("/player", LoadPlayers).Methods("GET")

	r.HandleFunc("/player/{nickname}", DeletePlayer).Methods("DELETE")

	r.HandleFunc("/player/{nickname}", LoadPlayerByNickname).Methods("GET")

	r.HandleFunc("/player/{nickname}", SavePlayer).Methods("PUT")

	//Rotas do Enemies
	r.HandleFunc("/enemy", AddEnemy).Methods("POST")

	r.HandleFunc("/enemy", LoadEnemies).Methods("GET")

	r.HandleFunc("/enemy/{nickname}", DeleteEnemy).Methods("DELETE")

	r.HandleFunc("/enemy/{nickname}", LoadEnemyByNickname).Methods("GET")

	r.HandleFunc("/enemy/{nickname}", SaveEnemy).Methods("PUT")

	//Rotas de Batalha
	r.HandleFunc("/battle", addBatle).Methods("POST")

	r.HandleFunc("/battle", LoadBattle).Methods("GET")

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
