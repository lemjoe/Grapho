package typeassert

import (
	"github.com/lemjoe/Grapho/internal/service"
)

func InterfaceToString(i interface{}) string {
	logger := service.GetLogger()
	val, ok := i.(string)
	if !ok {
		logger.Errorf("Failed to cast interface to string: %v", i)
		return ""
	}
	return val
}
