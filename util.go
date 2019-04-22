package main

func verifyNoError(err error) {
	if err != nil {
		panic(err)
	}
}
