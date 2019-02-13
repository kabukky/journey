package methods

import (
	"github.com/kabukky/journey/structure"
	"strings"
)

// ProcessHelperArguments: Function to put all arguments into a neatly organized map (splitting argument.Name with format "name=argument" into map["name"]"argument")
// for easier lookup and use in helper functions.
func ProcessHelperArguments(arguments []structure.Helper) map[string]string {
	argumentsMap := make(map[string]string)
	for index, _ := range arguments {
		// Separate = arguments and put them in map
		argumentParts := strings.SplitN(arguments[index].Name, "=", 2)
		if len(argumentParts) > 1 {
			argumentsMap[argumentParts[0]] = argumentParts[1]
		} else {
			argumentsMap[arguments[index].Name] = ""
		}
	}
	return argumentsMap
}
