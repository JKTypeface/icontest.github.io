package main

import (
	"bytes"
	"testing"
)

func TestUnmask(t *testing.T) {
	before, err := loadImage("samples/before.png")

	if err != nil {
		t.Error(err)
	}

	after, err := loadImage("samples/after.png")

	if err != nil {
		t.Error(err)
	}

	before = Unmask(before)

	got, err := imageBytes(before)

	if err != nil {
		t.Error(err)
	}

	expected, err := imageBytes(after)

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(got, expected) {
		t.Errorf("Unmasking failed")
	}
}
