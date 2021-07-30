package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func main() {
	engine := handlebars.New("./views", ".hbs")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(context.Background())

	app.Get("", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{}, "layouts/main")
	})

	app.Post("", func(c *fiber.Ctx) error {
		body := new(PwBody)

		if err := c.BodyParser(body); err != nil {
			return err
		}

		var id string

		err = db.QueryRow(
			context.Background(),
			"INSERT INTO pwtoshare (pw, created_at, days_limit, views_remaining) VALUES($1, $2, $3, $4) RETURNING (id)",
			body.Password,
			time.Now(),
			body.DaysLimit,
			body.ViewsRemaining,
		).Scan(&id)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert due to: %v", err)
		}

		return c.Render("index", fiber.Map{
			"id": id,
		}, "layouts/main")
	})

	app.Get("/secret/:id", func(c *fiber.Ctx) error {
		var pw string
		var createdAt time.Time
		var daysLimit int
		var viewsRemaining int
		var found = true

		err := db.QueryRow(
			context.Background(),
			"SELECT pw, created_at, days_limit, views_remaining FROM pwtoshare WHERE id=$1",
			c.Params("id")).Scan(&pw, &createdAt, &daysLimit, &viewsRemaining)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch secret: %v\n", err)
			found = false
		}

		var maxDate = createdAt.AddDate(0, 0, daysLimit)

		if maxDate.After(time.Now()) && viewsRemaining >= 0 {
			viewsRemaining = viewsRemaining - 1
			_, err = db.Exec(
				context.Background(),
				"UPDATE pwtoshare SET views_remaining=views_remaining-1 WHERE id=$1",
				c.Params("id"),
			)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to update due to: %v\n", err)
			}
		}

		return c.Render("secret", fiber.Map{
			"found":          found,
			"canShow":        viewsRemaining >= 0 && maxDate.After(time.Now()),
			"password":       pw,
			"viewsRemaining": viewsRemaining,
			"maxDate":        maxDate,
			"createdAt":      createdAt,
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}
