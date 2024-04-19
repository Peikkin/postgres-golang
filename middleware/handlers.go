package middleware

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Peikkin/postgres-golang/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка загрузки .env файла")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось подключиться к базе данных: ")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Отсутствие ответа от базы данных")
	}

	log.Info().Msg("Подключение к базе данных выполнено!")

	return db
}

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения данных")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "Создание успешно!",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["id"]
	id, err := strconv.Atoi(params)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения параметра id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stock, err := getStock(int64(id))
	if err != nil {
		log.Error().Err(err).Msg("Ошибка получения данных о запасе")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stock, err := getAllStock()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка получения данных о запасах")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(stock)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["id"]
	id, err := strconv.Atoi(params)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения параметра id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения данных о запасе")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updateStock(int64(id), stock)

	res := response{
		ID:      int64(stock.ID),
		Message: "Обновление успешно!",
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["id"]
	id, err := strconv.Atoi(params)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения параметра id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deleteStock(int64(id))
	res := response{
		ID:      int64(id),
		Message: "Удаление успешно",
	}
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка выполнения запроса на создание")
	}
	log.Info().Msg("Запрос на создание выполнен")
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		log.Error().Err(err).Msg("Записи отсутствуют")
		return stock, nil
	case nil:
		log.Info().Msg("Запрос на получение данных выполнен")
		return stock, nil
	default:
		log.Error().Err(err).Msg("Ошибка выполнения запроса на получение данных о запасе")
	}

	return stock, err
}
func getAllStock() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка выполнения запроса на получение данных о запасах")
		return stocks, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Error().Err(err).Msg("Ошибка получения получения данных о запасах")
			return stocks, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, err
}
func updateStock(id int64, stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE id=$1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка обновления данных о запасе")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка обновления данных о запасе")
	}
	log.Info().Msg("Обновление данных о запасе выполнено")
	return rowsAffected
}
func deleteStock(id int64) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка удаления данных о запасе")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка удаления данных о запасе")
	}
	return rowsAffected
}
