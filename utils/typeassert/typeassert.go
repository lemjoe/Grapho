package typeassert

import "log"

func InterfaceToString(i interface{}) string {
	val, ok := i.(string)
	if !ok {
		log.Printf("Failed to cast interface to string: %v", i)
		return ""
	}
	return val
}
