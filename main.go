package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("Home")
		uid := uuid.NewString()[:8]
		c.Cookie(&fiber.Cookie{
			Name:  "client_id",
			Value: uid,
		})
		return c.Render("index", fiber.Map{"ClientId": uid})
	})

	clients := make(map[string]*websocket.Conn)

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed")) // true
		// log.Println(c.Params("id"))       // 123
		// log.Println(c.Query("v"))         // 1.0
		client_id := c.Cookies("client_id")
		clients[client_id] = c
		log.Println("Client Online: ", client_id)
		log.Println("Clients: ", len(clients))

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				delete(clients, client_id)
				break
			}
			log.Printf("recv: %s", msg)
			m := make(map[string]interface{})
			erra := json.Unmarshal(msg, &m)
			if erra != nil {
				fmt.Println("Error Unmarshaling")
			}
			// fmt.Println(m["HEADERS"].(map[string]interface{})["HX-Target"])
			for cid, conn := range clients {
				var mesg = []byte(``)
				now := time.Now()
				if cid == client_id {
					mesg = []byte(`<div hx-swap-oob="beforeend:#messages"><div class="grid my-2"><span class="font-bold text-size-xl">` + client_id + ` (You):</span><span class="text-size-3 text-slate-500">` + now.Format("2006-01-02 15:04:05") + `</span><span class="font-italic text-size-lg">` + m["text"].(string) + `</span></div></div>`)
				} else {
					mesg = []byte(`<div hx-swap-oob="beforeend:#messages"><div class="grid my-2"><span class="font-bold text-size-xl">` + client_id + `:</span><span class="text-size-3 text-slate-500">` + now.Format("2006-01-02 15:04:05") + `</span><span class="font-italic text-size-lg">` + m["text"].(string) + `</span></div></div>`)
				}
				log.Println("Sending to : ", cid)
				if err = conn.WriteMessage(mt, mesg); err != nil {
					log.Println("write:", err)
					break
				}
			}

		}

	}))

	log.Fatal(app.Listen("0.0.0.0:3000"))
}
