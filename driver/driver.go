package driver

//-------------------------------------------//
// Constructs and handles the Context object //
//-------------------------------------------//

import (
	"github.com/osingaatje/seshat/context"
	"github.com/osingaatje/seshat/queryable/tryout"
)

// Add your own FindQueries to this method!
func FindQueries(c *context.Ctx) {
	tryout.FindQueries(c)
}

// should be used for testing
func NewContext() *context.Ctx {
	ctx := context.New()
	FindQueries(ctx)
	return ctx
}

// should be used for UIs, CLIs, finished programs
func Bootstrap(f func(*context.Ctx)) int {
	ctx := context.New()
	FindQueries(ctx)

	var returnCode int = 0

	// TODO add this error handling later on maybe?
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
