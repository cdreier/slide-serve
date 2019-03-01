package main

import (
	"testing"
)

func Test_holder_generateSlides(t *testing.T) {
	h := holder{
		slides: make([]slide, 0),
	}
	slideContent := `
# simple demo

for tests

@code/go
  func()test{
  
	}

# end
	`
	h.generateSlides(slideContent)
	expected := 4
	got := len(h.slides)
	if got != expected {
		t.Errorf("wanted slides %d, got %d", expected, got)
	}
}
