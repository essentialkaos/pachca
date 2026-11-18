package data

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v14/timeutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Date holds the calendar date submitted through a modal form
type Date struct {
	time.Time
}

// Time holds the hour and minute submitted through a modal form
type Time struct {
	Hour   int
	Minute int
}

// File holds metadata for a file submitted through a modal form
type File struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

// Files holds the files submitted through a file input block
type Files []*File

// Options holds the selected values from a checkbox block
type Options []string

// ////////////////////////////////////////////////////////////////////////////////// //

// UnmarshalJSON accepts dates in YYYY-MM-DD or DD.MM.YYYY format and resets
// the date on null
func (d *Date) UnmarshalJSON(b []byte) error {
	data := string(b)

	if data == "null" {
		d.Time = time.Time{}
		return nil
	}

	date, err := timeutil.ParseWithAny(data, `"2006-01-02"`, `"02.01.2006"`)

	if err != nil {
		return err
	}

	d.Time = date

	return nil
}

// UnmarshalJSON decodes a 24-hour time string and leaves the receiver unchanged
// on null
func (t *Time) UnmarshalJSON(b []byte) error {
	data := string(b)

	if data == "null" {
		return nil
	}

	hour, min, ok := strings.Cut(strings.Trim(data, `"`), ":")

	if !ok {
		return fmt.Errorf("invalid time format")
	}

	h, err := strconv.Atoi(hour)

	if err != nil {
		return fmt.Errorf("can't parse time hour value: %w", err)
	}

	if h < 0 || h > 23 {
		return fmt.Errorf("invalid hour value \"%d\"", h)
	}

	m, err := strconv.Atoi(min)

	if err != nil {
		return fmt.Errorf("can't parse time minute value: %w", err)
	}

	if m < 0 || m > 59 {
		return fmt.Errorf("invalid minute value \"%d\"", m)
	}

	t.Hour, t.Minute = h, m

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// String formats the time as HH:MM and returns 00:00 for a nil receiver
func (t *Time) String() string {
	if t == nil {
		return "00:00"
	}

	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Has reports whether the given option is selected
func (o Options) Has(option string) bool {
	return slices.Contains(o, option)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Get returns the file at the given zero-based index or nil if it is out of range
func (f Files) Get(index int) *File {
	if index >= len(f) || index < 0 {
		return nil
	}

	return f[index]
}
