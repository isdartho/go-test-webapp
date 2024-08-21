package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("Home")
		return c.Render("index", fiber.Map{"data": "test"})
	})

	// app.Use("/ws", func(c *fiber.Ctx) error {
	// 	// IsWebSocketUpgrade returns true if the client
	// 	// requested upgrade to the WebSocket protocol.
	// 	log.Println("Websocket requested!")
	// 	if websocket.IsWebSocketUpgrade(c) {
	// 		c.Locals("allowed", true)
	// 		return c.Next()
	// 	}
	// 	return fiber.ErrUpgradeRequired
	// })

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed")) // true
		// log.Println(c.Params("id"))       // 123
		// log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			m := make(map[string]interface{})
			erra := json.Unmarshal(msg, &m)
			if erra != nil {
				fmt.Println("Error Unmarshaling")
			}
			// fmt.Println(m["HEADERS"].(map[string]interface{})["HX-Target"])
			mesg := []byte(`<div hx-swap-oob="beforeend:#messages"><p><i style="color: green" class="fa fa-circle"></i> ` + m["text"].(string) + `</p></div>`)
			if err = c.WriteMessage(mt, mesg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	app.Listen(":3000")
}
