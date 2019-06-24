package channels

// SetStringChan is used with a goroutine to send a string value to a channel.
func SetStringChan(sChan chan<- string, s string) {
	sChan <- s
}
