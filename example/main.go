package main

import (
	"fmt"
	"net/http"

	"github.com/captncraig/erroneous"
	"github.com/captncraig/erroneous/web"
	"github.com/pkg/errors"
)

func main() {
	//Initialize the logger (optional)
	erroneous.Use(&erroneous.ErrorLogger{Store: &erroneous.MemoryStore{}, MachineName: "ny-web42"})

	err := SomeFunction()
	if err != nil {
		erroneous.LogError(err)
	}

	err = SomeWrappingFunction()
	if err != nil {
		erroneous.LogError(err)
	}
	http.Handle("/errors", web.GetMux(erroneous.Default, ""))
	http.ListenAndServe(":9090", nil)
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
