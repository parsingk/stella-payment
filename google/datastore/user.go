package datastore

import (
	"github.com/labstack/echo"
	ds "google/datastore/lib"
	"google"
	"google.golang.org/appengine/datastore"
)

const kindUser = "user"

func InsertUser (c echo.Context, data interface{}) (*datastore.Key, error) {
	ctx := google.GetContext(c)

	return ds.Insert(ctx, kindUser, data)
}