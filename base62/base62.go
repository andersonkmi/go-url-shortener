// Package base62 provides encoding of numeric IDs into base62 strings
// using the alphabet [0-9A-Za-z].
package base62

// IdToBase62 converts a non-negative id into its base62 representation.
// The zero value is encoded as "0".
func IdToBase62(id int64) string {
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
