package base62

// Convert ID to base62 string
func idToBase62(id int64) string {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if id == 0 {
		return string(base62Chars[0])
	}

	var result []byte
	base := int64(len(base62Chars))

	for id > 0 {
		remainder := id % base
		result = append([]byte{base62Chars[remainder]}, result...)
		id = id / base
	}

	return string(result)
}
