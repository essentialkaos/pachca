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

	"github.com/essentialkaos/ek/v13/timeutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Date contains info from Date block of modal window
type Date struct {
	time.Time
}

// Time contains info from Time block of modal window
type Time struct {
	Hour   int
	Minute int
}

// File contains info from File block of modal window
type File struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

// Files is a slice of files
type Files []*File

// Options contains slice with selected options from Checkbox block
// of modal window
type Options []string

// ////////////////////////////////////////////////////////////////////////////////// //

// UnmarshalJSON parses JSON date
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

// UnmarshalJSON parses JSON date
func (t *Time) UnmarshalJSON(b []byte) error {
	data := string(b)

	if data == "null" {
		return nil
	}

	hour, min, ok := strings.Cut(strings.Trim(data, `"`), ":")

	if !ok {
		return fmt.Errorf("Invalid time format")
	}

	h, err := strconv.Atoi(hour)

	if err != nil {
		return fmt.Errorf("Can't parse time hour value: %w", err)
	}

	if h < 0 || h > 23 {
		return fmt.Errorf("Invalid hour value \"%d\"", h)
	}

	m, err := strconv.Atoi(min)

	if err != nil {
		return fmt.Errorf("Can't parse time minute value: %w", err)
	}

	if m < 0 || m > 59 {
		return fmt.Errorf("Invalid minute value \"%d\"", m)
	}

	t.Hour, t.Minute = h, m

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// String returns string representation of time
func (t *Time) String() string {
	if t == nil {
		return "00:00"
	}

	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Has returns true if slice contains given option
func (o Options) Has(option string) bool {
	return slices.Contains(o, option)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Get returns file with given index
func (f Files) Get(index int) *File {
	if index >= len(f) || index < 0 {
		return nil
	}

	return f[index]
}
