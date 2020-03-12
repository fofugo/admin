package benchmark

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v10"
)

type BenchmarkTemplate struct {
	No         int
	Context    string
	Benchmarks []Benchmark
}
type Benchmark struct {
	Id      int
	Title   string
	Content string
}

func BenchmarkHandler(c echo.Context) error {
	type request struct {
		No int `validate:"required" param:"no"`
	}
	var benchmarkTemplate BenchmarkTemplate
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
	if err := db.QueryRow("SELECT context FROM benchmark_no WHERE no=?", input.No).Scan(&benchmarkTemplate.Context); err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT id, title, content FROM benchmark WHERE no=?", input.No)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		benchmark := Benchmark{}
		rows.Scan(&benchmark.Id, &benchmark.Title, &benchmark.Content)
		benchmarkTemplate.Benchmarks = append(benchmarkTemplate.Benchmarks, benchmark)
	}

	benchmarkTemplate.No = input.No
	return c.Render(http.StatusOK, "benchmark", benchmarkTemplate)
}
