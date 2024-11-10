package io

import (
	"testing"
)

var stringToTestAgainst = "9onI2wUvHXlehdrXCg7Gx2GbIjeTJxDotfdT9jRd6eiGDuDsYp8zrJgJxbVdmbmmRIAov7VPD3O7KEz12iw279HTyO9u0F6RetjzM5Kdkwa3gJPOQtJW8ZgMMZ2K7OIa07Glp3sJQO00JMaERJu0tVAwMDQn6ZPMSJblXa95PYtpYVOz1JXbxwn02SNHolSGgCCCsAeIXHHtd84BE0bQHhYsrZ7ejizQiaEyshdpBQHWIktBUF9ITbieTtlzgQtkUlcvZ4WUk923SBsOCY5RODTRfVpaPAndSAwrtdKO5mJ0H0ikKN1VM37FAW5dUwqlcDrm1JSfqj4rfzBXl4GNUJMxAla7kja92w2LH7LHFFMPqQaYzAXpOyDORNJelQ7xKya5mYbRRVSxF5CUOvI7mB4ZghjzgoGejXjg9P85HeliEzgZ7unlRv35KkiS1BcLGNrXjJKIYxxkmK4ADF7886by6lkhtx9KY8FhAF2PdeqDyZbjfqcrbOMTnJbCQKzGdgKfyT3UtcmCpLLjcl2ttsn9eg6WQXcTjTOlLOTikDH5VZ07oMI4Np9h4kqtxEIyfCYeQ8kmTZoZKQ7EZMFtP5npmQGQt9iLFmLI"

func BenchmarkSecret1(b *testing.B) {
	for range b.N {
		MakeSecret(stringToTestAgainst)
	}
}

func BenchmarkSecret2(b *testing.B) {
	for range b.N {
		makeSecret2(stringToTestAgainst)
	}
}

func BenchmarkSecret3(b *testing.B) {
	for range b.N {
		makeSecret3(stringToTestAgainst)
	}
}

func BenchmarkSecret4(b *testing.B) {
	for range b.N {
		makeSecret4(stringToTestAgainst)
	}
}
