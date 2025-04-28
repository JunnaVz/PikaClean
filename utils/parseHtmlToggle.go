package utils

func ParseHtmlToggle(rawBool string) bool {
	// convert "on"/"off" to true/false
	return rawBool == "on"
}
