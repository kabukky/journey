package methods

import (
	"github.com/kabukky/journey/structure"
	"strings"
)

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
