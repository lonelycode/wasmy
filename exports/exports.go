package exports

import (
	"fmt"

	"github.com/lonelycode/wasmy/shared_types"
)

func PrintHello(args *shared_types.Args) (interface{}, error) {
	val := fmt.Sprintf("From Host: Hello Mr. %s", args.Args[0].(string))
	fmt.Println(val)

	return val, nil
}
