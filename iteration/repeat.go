package iteration

func Repeat(character string, repeatCount int) string {
	//return strings.Repeat(character, repeatCount)
	var repeated string
	for i := 0; i < repeatCount; i++ {
		repeated += character
	}
	return repeated
}
