package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/hahaclassic/cmdrouter"
)

// panic: runtime error: index out of range [1] with length 0
func panicFunc(_ context.Context) error {
	var a []int
	fmt.Println("we are thrown into the root router, where the recover middleware is installed.")
	fmt.Println(a[0])
	return nil
}

func main() {
	ctx := context.Background()

	// handlers
	panicOption := cmdrouter.Option{
		Name:    "Panic",
		Handler: panicFunc,
	}

	errOption := cmdrouter.Option{
		Name: "Error",
		Handler: func(ctx context.Context) error {
			fmt.Println("An error has occurred. But there is no logging?")
			return errors.New("some error")
		},
	}

	router := cmdrouter.NewCmdRouter("lvl 0")
	router.AddMiddleware(cmdrouter.DefaultLoggerMiddleware,
		cmdrouter.DefaultRecoverMiddleware)
	router.PathShow(true)

	// crash
	crashGroup := router.Group("lvl 1 [bad]")
	crashGroup.AddOptions(panicOption, errOption)
	crashGroupLevel2 := crashGroup.Group("lvl 2 [bad]")
	crashGroupLevel2.AddOptions(panicOption, errOption)

	// correct processing
	group := router.GroupWithMiddleware("lvl 1")
	group.AddOptions(panicOption, errOption)
	groupLevel2 := group.GroupWithMiddleware("lvl 2")
	groupLevel2.AddOptions(panicOption, errOption)

	router.Run(ctx)
}
