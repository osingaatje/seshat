package driver

import (
	"github.com/osingaatje/seshat/context"
	"github.com/osingaatje/seshat/queryable/tryout"
)

// Add your own FindQueries to this method!
func FindQueries(c *context.Ctx) {
	tryout.FindQueries(c)
}

func Bootstrap(f func(*context.Ctx)) int {
	ctx := context.New()
	FindQueries(ctx)

	var returnCode int = 0

	// TODO add this later on
	//	defer func() {
	//		if r := recover(); r != nil {
	//			returnCode = 1
	//
	//			errorStr := "An error occurred while running Seshat: \n"
	//			errorStr += fmt.Sprintf("%s", r)
	//			fmt.Println(errorStr)
	//		}
	//	}()

	// run the code
	f(ctx)

	return returnCode
}
