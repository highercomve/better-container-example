package utils

// A utility to convert the values to proper strings.
func Int8ToStr(arr interface{}) string {
	var b []byte

	switch arr := arr.(type) {
	case []int8:
		b = make([]byte, 0, len(arr))
		for _, v := range arr {
			if v == 0x00 {
				break
			}
			b = append(b, byte(v))
		}
	case []uint8:
		b = make([]byte, 0, len(arr))
		for _, v := range arr {
			if v == 0x00 {
				break
			}
			b = append(b, byte(v))
		}
	}

	return string(b)
}
