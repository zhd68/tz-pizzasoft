package generatorid

func GenerateId() func() string {
	id := 702
	return func() string {
		id++
		return base10ToBase26(id)
	}
}

func base10ToBase26(n int) string {
	stringID := ""
	for n > 0 {
		mod := byte('a' + (n-1)%26)
		stringID = string(mod) + stringID
		n = (n - 1) / 26
	}
	return stringID

}
