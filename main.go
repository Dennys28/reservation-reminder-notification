package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Inicializar la conexión con la base de datos
func initDB() {
	var err error
	dsn := "admin:Hola1244@tcp(54.147.30.4:3306)/reservation_db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error al verificar conexión con la base de datos: %v", err)
	}
	fmt.Println("Conexión exitosa con la base de datos")
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// Endpoint para enviar recordatorios
	r.POST("/send-reminder", func(c *gin.Context) {
		type ReminderRequest struct {
			ReservationID int    `json:"reservation_id"`
			Email         string `json:"email"`
		}

		var request ReminderRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: asegúrate de enviar 'reservation_id' y 'email'"})
			return
		}

		// Consultar la reservación en la base de datos
		var name, date string
		err := db.QueryRow("SELECT name, reservation_date FROM reservations WHERE id = ?", request.ReservationID).Scan(&name, &date)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Reservación no encontrada"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al buscar la reservación: %v", err)})
			}
			return
		}

		// Validar la fecha
		parsedDate, err := time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fecha de reservación con formato inválido"})
			return
		}

		// Simular el envío del recordatorio
		fmt.Printf("Enviando recordatorio a %s (%s) para la fecha %s\n", name, request.Email, parsedDate.Format("2006-01-02 15:04:05"))
		c.JSON(http.StatusOK, gin.H{"message": "Recordatorio enviado exitosamente"})
	})

	r.GET("/reservation-details", func(c *gin.Context) {
		reservationID := c.Query("reservation_id") // Obtener el parámetro desde la URL

		var name, email, date string
		err := db.QueryRow("SELECT name, email, reservation_date FROM reservations WHERE id = ?", reservationID).Scan(&name, &email, &date)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservación no encontrada"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reservation_id":   reservationID,
			"name":             name,
			"email":            email,
			"reservation_date": date,
		})
	})

	r.Run(":8080") // Escuchar en el puerto 8080
}
