package main

import (
	"strings"
	"testing"
)

func Test_holder_generateSlides_append(t *testing.T) {
	slideContent := `
# first slide

# second slide
@append

# third slide
@append`

	h := holder{}
	h.slides = make([]slide, 0)
	h.generateSlides(slideContent)

	expectedSlides := 3
	got := len(h.slides)
	if got != expectedSlides {
		t.Errorf("wanted slides %d, got %d ", expectedSlides, got)
	}

	expectedLines := [3]int{1, 2, 3}
	for i, s := range h.slides {
		lines := len(strings.Split(s.content, "\n")) - 1 // the -1 strips the newline at the end
		if lines != expectedLines[i] {
			t.Errorf("wanted slide-lines %d, got %d ", expectedLines[i], lines)
		}
	}

}

func Test_holder_generateSlides(t *testing.T) {
	tests := []struct {
		name           string
		slideContent   string
		expectedSlides int
	}{
		{
			name: "basic parsing with code",
			slideContent: `
# simple demo

for tests

@code/go
	func()test{
	  
	}

# end`,
			expectedSlides: 4,
		},
		{
			name: "all sides are generated with @append",
			slideContent: `
# first slide

# second slide
@append`,
			expectedSlides: 2,
		},
		{
			name: "test with hide",
			slideContent: `
# first slide

@hide
# second slide`,
			expectedSlides: 1,
		},
		{
			name: "test with hide and append",
			slideContent: `
# first slide

@hide
# second slide
@append

# third slide
@append`,
			expectedSlides: 2,
		},
	}
	h := holder{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.slides = make([]slide, 0)
			h.generateSlides(tt.slideContent)
			got := len(h.slides)
			if got != tt.expectedSlides {
				t.Errorf("wanted slides %d, got %d in test '%s'", tt.expectedSlides, got, tt.name)
			}
		})
	}
}
