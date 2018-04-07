package response

const (


	/**
		Datastore
		1500~
	**/
	datastoreError = 1505
)

type Response struct {
	Code int
	Description string
}

func DatastoreError() interface{} {
	return Response{
		Code :        datastoreError,
		Description : "datastore error",
	}
}
