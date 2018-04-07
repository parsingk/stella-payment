// +build appengine

package main

import (
	"api"
	"google.golang.org/appengine"
	"net/http"
)

func main() {
	http.Handle("/", api.Server)
	appengine.Main()
}
