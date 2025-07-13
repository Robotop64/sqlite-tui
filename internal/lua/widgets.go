package lua

type Widget interface {
	SetData(data interface{})
	SetActions(actions []interface{})
}
