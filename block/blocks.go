package block

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Block interface {
	Init()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Header is a header block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-header
type Header struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// PlainText is a plain text block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-plaintext
type PlainText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Markdown is a text block with markdown formatting
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-markdown
type Markdown struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Divider is a divider block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-divider
type Divider struct {
	Type string `json:"type"`
}

// Input is text input block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-input
type Input struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	Placeholder  string `json:"placeholder,omitempty"`
	Hint         string `json:"hint,omitempty"`
	InitialValue string `json:"initial_value,omitempty"`
	MinLength    int    `json:"min_length,omitzero"`
	MaxLength    int    `json:"max_length,omitzero"`
	IsMultiline  bool   `json:"multiline"`
	IsRequired   bool   `json:"required"`
}

// Select is select list block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-select
type Select struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Label      string  `json:"label"`
	Hint       string  `json:"hint,omitempty"`
	Options    Options `json:"options,omitempty"`
	IsRequired bool    `json:"required"`
}

// Radio is radio buttons block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-radio
type Radio struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Label      string  `json:"label"`
	Hint       string  `json:"hint,omitempty"`
	Options    Options `json:"options,omitempty"`
	IsRequired bool    `json:"required"`
}

// Checkbox is checkboxes block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-checkbox
type Checkbox struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Label      string  `json:"label"`
	Hint       string  `json:"hint,omitempty"`
	Options    Options `json:"options,omitempty"`
	IsRequired bool    `json:"required"`
}

// Date is date selection block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-date
type Date struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	Hint         string `json:"hint,omitempty"`
	InitialValue string `json:"initial_date,omitempty"`
	IsRequired   bool   `json:"required"`
}

// Time is time selection block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-time
type Time struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	Hint         string `json:"hint,omitempty"`
	InitialValue string `json:"initial_time,omitempty"`
	IsRequired   bool   `json:"required"`
}

// Files is file attachment block
//
// https://crm.pachca.com/dev/forms/views/blocks/#title-files
type Files struct {
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Label      string   `json:"label"`
	Hint       string   `json:"hint,omitempty"`
	Filetypes  []string `json:"filetypes,omitempty"`
	MaxFiles   int      `json:"max_files,omitzero"`
	IsRequired bool     `json:"required"`
}

// Options is a slice with block options
type Options []*Option

// Option is block option
type Option struct {
	Text        string `json:"text"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	IsSelected  bool   `json:"selected"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Init sets block type
func (b *Header) Init() {
	b.Type = "header"
}

// Init sets block type
func (b *PlainText) Init() {
	b.Type = "plain_text"
}

// Init sets block type
func (b *Markdown) Init() {
	b.Type = "markdown"
}

// Init sets block type
func (b *Divider) Init() {
	b.Type = "divider"
}

// Init sets block type
func (b *Input) Init() {
	b.Type = "input"
}

// Init sets block type
func (b *Select) Init() {
	b.Type = "select"
}

// Init sets block type
func (b *Radio) Init() {
	b.Type = "radio"
}

// Init sets block type
func (b *Checkbox) Init() {
	b.Type = "checkbox"
}

// Init sets block type
func (b *Date) Init() {
	b.Type = "date"
}

// Init sets block type
func (b *Time) Init() {
	b.Type = "time"
}

// Init sets block type
func (b *Files) Init() {
	b.Type = "file_input"
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AddOption adds a new option to select list
func (b *Select) AddOption(text, value string, selected bool) *Select {
	if b == nil {
		return b
	}

	b.Options = append(b.Options, &Option{
		Text:       text,
		Value:      value,
		IsSelected: selected,
	})

	return b
}

// AddOptionIf conditionally adds a new option to select list
func (b *Select) AddOptionIf(cond bool, text, value string, selected bool) *Select {
	if cond == false {
		return b
	}

	return b.AddOption(text, value, selected)
}

// AddOption adds a new option to radio group
func (b *Radio) AddOption(text, value, desc string, selected bool) *Radio {
	if b == nil {
		return b
	}

	b.Options = append(b.Options, &Option{
		Text:        text,
		Value:       value,
		Description: desc,
		IsSelected:  selected,
	})

	return b
}

// AddOptionIf conditionally adds a new option to radio group
func (b *Radio) AddOptionIf(cond bool, text, value, desc string, selected bool) *Radio {
	if cond == false {
		return b
	}

	return b.AddOption(text, value, desc, selected)
}

// AddOption adds a new option to checkbox group
func (b *Checkbox) AddOption(text, value, desc string, selected bool) *Checkbox {
	if b == nil {
		return b
	}

	b.Options = append(b.Options, &Option{
		Text:        text,
		Value:       value,
		Description: desc,
		IsSelected:  selected,
	})

	return b
}

// AddOptionIf conditionally adds a new option to checkbox group
func (b *Checkbox) AddOptionIf(cond bool, text, value, desc string, selected bool) *Checkbox {
	if cond == false {
		return b
	}

	return b.AddOption(text, value, desc, selected)
}

// Set sets initial date
func (b *Date) Set(year, month, day int) *Date {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%d-%02d-%02d", year, month, day)

	return b
}

// SetIf conditionally sets initial date
func (b *Date) SetIf(cond bool, year, month, day int) *Date {
	if cond == false {
		return b
	}

	return b.Set(year, month, day)
}

// SetWithDate sets initial date using time.Time
func (b *Date) SetWithDate(d time.Time) *Date {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%d-%02d-%02d", d.Year(), d.Month(), d.Day())

	return b
}

// SetWithDateIf conditionally sets initial date using time.Time
func (b *Date) SetWithDateIf(cond bool, d time.Time) *Date {
	if cond == false {
		return b
	}

	return b.SetWithDate(d)
}

// Set sets initial time
func (b *Time) Set(hour, minute int) *Time {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%02d:%02d", hour, minute)

	return b
}

// SetIf conditionally sets initial time
func (b *Time) SetIf(cond bool, hour, minute int) *Time {
	if cond == false {
		return b
	}

	return b.Set(hour, minute)
}

// SetWithDate sets initial time using time.Time
func (b *Time) SetWithDate(d time.Time) *Time {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%02d:%02d", d.Hour(), d.Minute())

	return b
}

// SetWithDateIf conditionally sets initial time using time.Time
func (b *Time) SetWithDateIf(cond bool, d time.Time) *Time {
	if cond == false {
		return b
	}

	return b.SetWithDate(d)
}

// ////////////////////////////////////////////////////////////////////////////////// //
