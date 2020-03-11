package echo

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

func AddHandler(c echo.Context) error {
	type Input struct {
		Id      int    `form:"id"`
		No      int    `form:"no"`
		Title   string `form:"title"`
		Content string `form:"content"`
	}
	input := Input{}
	if err := c.Bind(&input); err != nil {
		panic(err)
	}
	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	if input.Id != 0 {
		if input.Title != "" {
			if _, err := db.Exec("UPDATE echo_board SET title = ? WHERE id = ?", input.Title, input.Id); err != nil {
				panic(err)
			}
		}
		if input.Content != "" {
			if _, err := db.Exec("UPDATE echo_board SET content = ? WHERE id = ?", input.Content, input.Id); err != nil {
				panic(err)
			}
		}
	} else {
		if input.No != 0 {
			if _, err := db.Exec("INSERT INTO echo_board (no,title,content) VALUES (?,?,?)", input.No, input.Title, input.Content); err != nil {
				panic(err)
			}
		}
	}

	return c.Render(http.StatusOK, "echo", nil)
}
