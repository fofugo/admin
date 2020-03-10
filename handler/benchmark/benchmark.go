package benchmark

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v10"
)

type BenchmarkTemplate struct {
	Id         int
	Benchmarks []Benchmark
}
type Benchmark struct {
	Id      int
	Title   string
	Content string
}

func BenchmarkHandler(c echo.Context) error {
	type Input struct {
		Id int `validate:"required" param:"id"`
	}
	var benchmarkTemplate BenchmarkTemplate
	input := Input{}
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
	benchmarkTemplate.Id = input.Id
	rows, err := db.Query("SELECT id, title, content FROM benchmark WHERE benchmark_id=?", input.Id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		benchmark := Benchmark{}
		rows.Scan(&benchmark.Id, &benchmark.Title, &benchmark.Content)
		benchmarkTemplate.Benchmarks = append(benchmarkTemplate.Benchmarks, benchmark)
	}

	return c.Render(http.StatusOK, "benchmark", benchmarkTemplate)
}
