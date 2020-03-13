package echo

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v10"
)

type echoTemplate struct {
	No         int
	Context    string
	EchoBoards []echoBoard
}
type echoBoard struct {
	Id      int
	Title   string
	Content string
}

func EchoHandler(c echo.Context) error {
	type request struct {
		No int `validate:"required" param:"no"`
	}
	var echoTemplate echoTemplate
	input := request{}
	if err := c.Bind(&input); err != nil {
		panic(err)
	}
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		panic(err)
	}

	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}

	_ = db.QueryRow("SELECT context FROM echo_board_no WHERE no=?", input.No).Scan(&echoTemplate.Context)

	rows, err := db.Query("SELECT id, title, content FROM echo_board WHERE no = ?", input.No)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		echoBoard := echoBoard{}
		rows.Scan(&echoBoard.Id, &echoBoard.Title, &echoBoard.Content)
		echoTemplate.EchoBoards = append(echoTemplate.EchoBoards, echoBoard)
	}

	echoTemplate.No = input.No
	return c.Render(http.StatusOK, "echo", echoTemplate)
}
