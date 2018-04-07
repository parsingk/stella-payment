package datastore

import (
	"google.golang.org/appengine/datastore"
	"reflect"
	"golang.org/x/net/context"
	"log"
)

func Query (kind string) *datastore.Query {
	return datastore.NewQuery(kind)
}

func Run (ctx context.Context, query *datastore.Query, fn func(interface{})) {

	results := query.Run(ctx) // return *datastore.Iterator

	for {
		var x interface{}
		_, error := results.Next(&x)

		if error == datastore.Done {
			break
		}

		if error != nil {
			break
		}

		fn(x)
	}
}

func SelectAllFromQuery(ctx context.Context, query *datastore.Query, vo interface{}) ([]*datastore.Key, error) {

	return query.GetAll(ctx, vo)
	//isError(error)
}

func SelectOne(ctx context.Context, kind string, id int64, vo interface{}, parent ...*datastore.Key) (*datastore.Key, error) {
	parentKey := getParentKey(parent)

	var key *datastore.Key = nil

	key = datastore.NewKey(ctx, kind, "", id, parentKey)

	error := datastore.Get(ctx, key, vo)

	//isError(error)
	return key, error
}

func Insert(ctx context.Context, kind string, data interface{}, parent ...*datastore.Key) (*datastore.Key, error) {
	parentKey := getParentKey(parent)

	inCompleteKey := datastore.NewIncompleteKey(ctx, kind, parentKey)

	return datastore.Put(ctx, inCompleteKey, data)
	//isError(error)
}

func InsertMulti (ctx context.Context, kind string, data interface{}, parent ...*datastore.Key) ([]*datastore.Key, error) {
	parentKey := getParentKey(parent)

	reflectedStruct := reflect.ValueOf(data)

	inCompleteKeys := multiNewIncompleteKey(ctx, kind, reflectedStruct.Len(), parentKey)
	keys, error := datastore.PutMulti(ctx, inCompleteKeys, data)

	//isError(error)

	return keys, error
}

func Update (ctx context.Context, kind string, data interface{}, id int64, parent ...*datastore.Key) (*datastore.Key, error) {
	parentKey := getParentKey(parent)

	inCompleteKey := datastore.NewKey(ctx, kind, "", id, parentKey)
	key, error := datastore.Put(ctx, inCompleteKey, data)

	//isError(error)

	return key, error
}

func DeleteByKey (ctx context.Context, key *datastore.Key) error {
	error := datastore.Delete(ctx, key)

	//isError(error)
	return error
}

func DeleteById (ctx context.Context, kind string, id int64, parent ...*datastore.Key) error {
	parentKey := getParentKey(parent)

	inCompleteKey := datastore.NewKey(ctx, kind, "", id, parentKey)
	error := datastore.Delete(ctx, inCompleteKey)

	//isError(error)
	return error
}

func Transaction (ctx context.Context, fn func(ct context.Context) error, options *datastore.TransactionOptions) error {
	error := datastore.RunInTransaction(ctx, func (transactionCtx context.Context) error {
		return fn(transactionCtx)
	}, options)

	//isError(error)
	return error
}

func multiNewIncompleteKey(ctx context.Context, kind string, count int, parent *datastore.Key) []*datastore.Key {
	var keys []*datastore.Key

	for i := 0; i < count; i++ {
		key := datastore.NewKey(ctx, kind, "", 0, parent)
		keys[i] = key
	}

	return keys
}

/**
	For default parentKey
 */
func getParentKey(parent []*datastore.Key) *datastore.Key {
	var parentKey *datastore.Key = nil

	if len(parent) > 0 {
		parentKey = parent[0]
	}

	return parentKey
}

func isError(error error) bool {
	if error != nil {
		log.Println(error)
		panic(error)
	}

	return false
}