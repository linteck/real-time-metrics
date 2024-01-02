package controllers

import (
	"encoding/json"
	"fmt"
	"stream/models"

	"context"
	"log"
	"time"

	echo "github.com/labstack/echo/v4"
	"nhooyr.io/websocket"
)

// Controller interface has two methods
type Controller interface {
	// Homecontroller renders initial home page
	HomeController(e echo.Context) error

	// StreamController responds with live cpu status over websocket
	StreamController(e echo.Context) error
}

type controller struct {
}

func NewController() Controller {
	return &controller{}
}

var model models.Model

// Initializes the models
func Init() {
	model = models.NewModel()
}

func (c *controller) HomeController(e echo.Context) error {
	return e.File("views/index.html")
}

func (c *controller) StreamController(e echo.Context) error {

	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	}

	conn, err := websocket.Accept(e.Response().Writer, e.Request(), opts)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.CloseNow()
	ctx, cancel := context.WithTimeout(e.Request().Context(), time.Second*100)
	defer cancel()

	status, err := model.GetLiveCpuUsage()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		// Write
		newVal := <-status

		jsonResponse, _ := json.Marshal(newVal)
		err := conn.Write(ctx, websocket.MessageText, jsonResponse)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
