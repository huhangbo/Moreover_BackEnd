package util

func ArrayToString(tmpArray []string) string {
	var result string
	for i := 0; i < len(tmpArray); i++ {
		result += tmpArray[i] + ","
	}
	return result
}

func StringToArray(tmpString string) []string {
	var result []string
	var tmp string
	for _, item := range tmpString {
		if item == ',' {
			result = append(result, tmp)
			tmp = ""
			continue
		}
		tmp += string(item)
	}
	return result
}
