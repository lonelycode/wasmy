//go:generate msgp
package shared_types

//tinyjson:json
type Args struct {
	Args []interface{} `msg:"args"`
}

//tinyjson:json
type Payload struct {
	Data interface{}       `msg:"data"`
	Meta map[string]string `msg:"meta"`
}
