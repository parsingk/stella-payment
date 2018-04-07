package datastore

import (
	"google.golang.org/appengine/datastore"
	"reflect"
	"golang.org/x/net/context"
	"log"
)

/**
	데이터 스토어 unindex 예제
**/
type Person struct {
	Name string
	Age  int `datastore:",noindex"`
}

/**
	데이터 스토어 Query.
	ID 말고도 여러 Filter, limit, offset 등의 처리를 하여 데이터를 가져올 수 있다.

	** ex ) Query("table").Project("name", "height").Distinct().Filter("height >", 20).Order("-height").Limit(10)

	어차피 쓸건 filter, order, limit 정도 일 것 같다..
	해당 패키지 안에서만 appengine/datastore패키지를 갖도록 하기 위해 만들어놈.. 두 datastore 패키지를 import할 필욘 없으니까~ (appengine context 처리는 추후 변경해야될듯)
**/
func Query (kind string) *datastore.Query {
	return datastore.NewQuery(kind)
}

/**
	데이터 스토어 Query Run.
	해당 쿼리를 실질적으로 구동(?) 시켜 데이터를 가져온다.
	두 가지 방식이 있다. Run함수를 호출 하는것과 GetAll 함수를 호출하는것.
	Run 함수는 Iterator를 반환한다. Next()함수를 가지고 있다...

	사실 Run보다 GetAll을 쓰겠..
**/
func Run (ctx context.Context, query *datastore.Query, fn func(interface{})) {

	results := query.Run(ctx) // return *datastore.Iterator

	for {
		var x interface{}
		_, error := results.Next(&x) // next로 돌린다.

		if error == datastore.Done {
			break
		}

		if error != nil {
			break
		}

		fn(x)
	}
}

/**
	데이터 스토어 Query GetAll.
	vo를 받아서 바로 해당 vo에 데이터를 담는다.
	키 배열을 return 한다.
**/
func SelectAllFromQuery(ctx context.Context, query *datastore.Query, vo interface{}) ([]*datastore.Key, error) {

	return query.GetAll(ctx, vo)
	//isError(error)
}

/**
	데이터 스토어 Select.
	ID를 알고 있어야 되며 ID를 모를 경우에는 Query문을 사용한다.
	원래 datastore에서 제공하는 메서드는 id는 string과 int64형 두개를 받는데 ( int64형은 데이터스토어 무작위 id, string은 커스텀 id를 넣었을 경우 ),
	아래 함수는 int64만을 위한 메서드이다.
**/
func SelectOne(ctx context.Context, kind string, id int64, vo interface{}, parent ...*datastore.Key) (*datastore.Key, error) {
	parentKey := getParentKey(parent)

	var key *datastore.Key = nil

	key = datastore.NewKey(ctx, kind, "", id, parentKey)

	error := datastore.Get(ctx, key, vo)

	log.Printf("%v", key)

	//isError(error)
	return key, error
}

/**
	데이터 스토어 Insert.
	해당 Kind(table)와 parent 정보를 가지고 키를 생성.

	** 사실상 datastore 라이브러리의 Put 메서드는 데이터가 있으면 update 기능을 없으면 insert 기능을 한다.

	생성된 키는 insert한 데이터의 ID를 담고 있으며, 데이터를 insert 한 후에는 키를 return 한다.
**/
func Insert(ctx context.Context, kind string, data interface{}, parent ...*datastore.Key) (*datastore.Key, error) {
	parentKey := getParentKey(parent)

	inCompleteKey := datastore.NewIncompleteKey(ctx, kind, parentKey)

	return datastore.Put(ctx, inCompleteKey, data)
	//isError(error)
}

/**
	데이터 스토어 멀티 Insert.
	여러 row의 데이터를 넣을 수 있다.
	datastore 라이브러리에서 PutMulti는 지원하는데 키 생성은 multi가 없는듯 하다...

	** 문제는 parent인데 키 생성시에 parent가 들어간다. 여기서 parent를 배열로 여러개를 받아야되는지 고민을 했..지만 하나만 받게함. 어차피 해당 테이블의 parent아닌감.

	insert가 완료되면 key 배열을 return 한다.
**/
func InsertMulti (ctx context.Context, kind string, data interface{}, parent ...*datastore.Key) ([]*datastore.Key, error) {
	parentKey := getParentKey(parent)

	reflectedStruct := reflect.ValueOf(data)

	inCompleteKeys := multiNewIncompleteKey(ctx, kind, reflectedStruct.Len(), parentKey)
	keys, error := datastore.PutMulti(ctx, inCompleteKeys, data)

	//isError(error)

	return keys, error
}

/**
	데이터 스토어 Update.
	datastore Put 메서드는 해당 Entity가 있으면 update를 치므로 ID를 받아서 Key를 생성하도록 하였다.
	인자로 받은 data로 update한다.
	update가 완료되면 key를 return 한다.
**/
func Update (ctx context.Context, kind string, data interface{}, id int64, parent ...*datastore.Key) (*datastore.Key, error) {
	parentKey := getParentKey(parent)

	inCompleteKey := datastore.NewKey(ctx, kind, "", id, parentKey)
	key, error := datastore.Put(ctx, inCompleteKey, data)

	//isError(error)

	return key, error
}

/**
	데이터 스토어 Delete.
	SELECT한 Key를 받아 Entity를 지운다.
**/
func DeleteByKey (ctx context.Context, key *datastore.Key) error {
	error := datastore.Delete(ctx, key)

	//isError(error)
	return error
}

/**
	데이터 스토어 Delete.
	ID를 받아 Entity를 지운다.
	해당 함수를 Transaction 안에서 쓸 경우에는 Transaction Option에 xg = true로 설정해야 한다.
**/
func DeleteById (ctx context.Context, kind string, id int64, parent ...*datastore.Key) error {
	parentKey := getParentKey(parent)

	inCompleteKey := datastore.NewKey(ctx, kind, "", id, parentKey)
	error := datastore.Delete(ctx, inCompleteKey)

	//isError(error)
	return error
}

/**
	데이터 스토어 Transaction
**/
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