package utils

// HandleError is a helper function to handle errors.
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
