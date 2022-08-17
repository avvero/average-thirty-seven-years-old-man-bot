package brain

type Intention interface {
	Express(text string) (has bool, response string)
}
