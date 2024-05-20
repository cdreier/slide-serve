package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_holder_generateSlides_hide(t *testing.T) {
	slideContent := `
# first slide

# second slide
@hidden

# third slide`

	h := holder{}
	h.slides = make([]slide, 0)
	h.generateSlides(slideContent)

	expectedSlides := 2
	got := len(h.slides)
	if got != expectedSlides {
		t.Errorf("wanted slides %d, got %d ", expectedSlides, got)
	}

	expectedLines := [2]int{1, 1}
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
			name: "test with hidden",
			slideContent: `
# first slide

@hidden
# second slide`,
			expectedSlides: 1,
		},
		{
			name: "test with hidden and append",
			slideContent: `
# first slide

@hidden
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

func Test_holder_generateSlides_zoom(t *testing.T) {
	slideContent := `
# start slide

# first slide
@zoom/1.5

# second slide
@zoom/2, center left

# third slide
@zoom/3, 50% 30%
`

	h := holder{}
	h.slides = make([]slide, 0)
	h.generateSlides(slideContent)

	expectedSlides := 4
	got := len(h.slides)
	if got != expectedSlides {
		t.Errorf("wanted slides %d, got %d ", expectedSlides, got)
	}

	assert.False(t, h.slides[0].zoom.enabled)
	assert.True(t, h.slides[1].zoom.enabled)
	assert.True(t, h.slides[2].zoom.enabled)
	assert.True(t, h.slides[3].zoom.enabled)

	assert.Equal(t, 1.5, h.slides[1].zoom.Scale)
	assert.Equal(t, 2.0, h.slides[2].zoom.Scale)
	assert.Equal(t, 3.0, h.slides[3].zoom.Scale)

	assert.Equal(t, "center center", h.slides[1].zoom.Origin)
	assert.Equal(t, "center left", h.slides[2].zoom.Origin)
	assert.Equal(t, "50% 30%", h.slides[3].zoom.Origin)

}
