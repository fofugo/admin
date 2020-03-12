package benchmark

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

func AddBenchmarkContextHandler(c echo.Context) error {
	type request struct {
		No      int    `form:"no"`
		Context string `form:"context"`
	}
	input := request{}
	if err := c.Bind(&input); err != nil {
		panic(err)
	}
	db := c.Get("db").(*sql.DB)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	if input.No != 0 {
		if input.Context != "" {
			if _, err := db.Exec("UPDATE benchmark_no SET context = ? WHERE no = ?", input.Context, input.No); err != nil {
				panic(err)
			}
		}
	}
	return c.Render(http.StatusOK, "benchmark", nil)
}
