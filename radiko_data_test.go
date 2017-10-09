package main

import (
    "testing"
)

func TestGetCount(t *testing.T) {
    want := "01"
    got := getCount("「1」")
    if got != want {
        t.Errorf("getCount(\"「1」\") == \"01\", want %q", want)
    }
}
