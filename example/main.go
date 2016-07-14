package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/captncraig/erroneous"
	"github.com/captncraig/erroneous/web"
	"github.com/pkg/errors"
)

func main() {
	//Initialize the logger (optional)
	erroneous.Use(&erroneous.ErrorLogger{
		Store: &erroneous.MemoryStore{
			RollupDuration: 5 * time.Minute,
		},
		MachineName: "ny-web42",
	})

	err := SomeFunction()
	if err != nil {
		erroneous.LogError(err)
	}

	for i := 0; i < 5; i++ {
		a()
	}

	http.Handle("/errors/", web.GetMux(erroneous.Default, ""))
	http.ListenAndServe(":9090", nil)
}

func a() {
	erroneous.LogError(SomeWrappingFunction())
}

func SomeWrappingFunction() error {
	return SomeOtherWrappingFunction()
}

func SomeOtherWrappingFunction() error {
	err := fmt.Errorf("Server not found")
	return errors.Wrap(err, "Error calling db")
}

func SomeFunction() error {
	return SomeOtherFunction()
}

func SomeOtherFunction() error {
	return fmt.Errorf("Something Bad Happened!")
}
