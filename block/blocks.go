package block

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Block represents an element in a Pachca form
type Block interface {
	// Init sets the block type before serialization
	Init()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Header displays a heading in a form
//
// https://dev.pachca.com/guides/forms#zagolovok-header
type Header struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// PlainText displays unformatted text in a form
//
// https://dev.pachca.com/guides/forms#obychnyy-tekst-plaintext
type PlainText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Markdown displays Markdown-formatted text in a form
//
// https://dev.pachca.com/guides/forms#formatirovannyy-tekst-markdown
type Markdown struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Divider visually separates sections of a form
//
// https://dev.pachca.com/guides/forms#razdelitel-divider
type Divider struct {
	Type string `json:"type"`
}

// Input accepts single-line or multiline text
//
// https://dev.pachca.com/guides/forms#tekstovoe-pole-input
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

// Select lets the user choose an option from a dropdown
//
// https://dev.pachca.com/guides/forms#vypadayuschiy-spisok-select
type Select struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Label      string  `json:"label"`
	Hint       string  `json:"hint,omitempty"`
	Options    Options `json:"options,omitempty"`
	IsRequired bool    `json:"required"`
}

// Radio lets the user choose one option from a group of radio buttons
//
// https://dev.pachca.com/guides/forms#radioknopki-radio
type Radio struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Label      string  `json:"label"`
	Hint       string  `json:"hint,omitempty"`
	Options    Options `json:"options,omitempty"`
	IsRequired bool    `json:"required"`
}

// Checkbox lets the user choose multiple options from a group of checkboxes
//
// https://dev.pachca.com/guides/forms#chekboksy-checkbox
type Checkbox struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Label      string  `json:"label"`
	Hint       string  `json:"hint,omitempty"`
	Options    Options `json:"options,omitempty"`
	IsRequired bool    `json:"required"`
}

// Date lets the user pick a calendar date
//
// https://dev.pachca.com/guides/forms#vybor-daty-date
type Date struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	Hint         string `json:"hint,omitempty"`
	InitialValue string `json:"initial_date,omitempty"`
	IsRequired   bool   `json:"required"`
}

// Time lets the user pick a time of day
//
// https://dev.pachca.com/guides/forms#vybor-vremeni-time
type Time struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Label        string `json:"label"`
	Hint         string `json:"hint,omitempty"`
	InitialValue string `json:"initial_time,omitempty"`
	IsRequired   bool   `json:"required"`
}

// Files lets the user attach files to a form
//
// https://dev.pachca.com/guides/forms#zagruzka-faylov-fileinput
type Files struct {
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Label      string   `json:"label"`
	Hint       string   `json:"hint,omitempty"`
	Filetypes  []string `json:"filetypes,omitempty"`
	MaxFiles   int      `json:"max_files,omitzero"`
	IsRequired bool     `json:"required"`
}

// Options holds the available choices for a selection block
type Options []*Option

// Option defines a choice and whether it is selected by default
type Option struct {
	Text        string `json:"text"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	IsSelected  bool   `json:"selected"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Init sets the JSON block type for Header
func (b *Header) Init() {
	b.Type = "header"
}

// Init sets the JSON block type for PlainText
func (b *PlainText) Init() {
	b.Type = "plain_text"
}

// Init sets the JSON block type for Markdown
func (b *Markdown) Init() {
	b.Type = "markdown"
}

// Init sets the JSON block type for Divider
func (b *Divider) Init() {
	b.Type = "divider"
}

// Init sets the JSON block type for Input
func (b *Input) Init() {
	b.Type = "input"
}

// Init sets the JSON block type for Select
func (b *Select) Init() {
	b.Type = "select"
}

// Init sets the JSON block type for Radio
func (b *Radio) Init() {
	b.Type = "radio"
}

// Init sets the JSON block type for Checkbox
func (b *Checkbox) Init() {
	b.Type = "checkbox"
}

// Init sets the JSON block type for Date
func (b *Date) Init() {
	b.Type = "date"
}

// Init sets the JSON block type for Time
func (b *Time) Init() {
	b.Type = "time"
}

// Init sets the JSON block type for Files
func (b *Files) Init() {
	b.Type = "file_input"
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AddOption appends a choice to the dropdown and returns the dropdown
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

// AddOptionIf appends a choice to the dropdown if cond is true and returns the dropdown
func (b *Select) AddOptionIf(cond bool, text, value string, selected bool) *Select {
	if !cond {
		return b
	}

	return b.AddOption(text, value, selected)
}

// AddOption appends a choice to the radio group and returns the radio group
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

// AddOptionIf appends a choice to the radio group if cond is true and returns
// the radio group
func (b *Radio) AddOptionIf(cond bool, text, value, desc string, selected bool) *Radio {
	if !cond {
		return b
	}

	return b.AddOption(text, value, desc, selected)
}

// AddOption appends a choice to the checkbox group and returns the checkbox group
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

// AddOptionIf appends a choice to the checkbox group if cond is true and returns
// the checkbox group
func (b *Checkbox) AddOptionIf(cond bool, text, value, desc string, selected bool) *Checkbox {
	if !cond {
		return b
	}

	return b.AddOption(text, value, desc, selected)
}

// Set sets the initial date from the year, month, and day and returns the date picker
func (b *Date) Set(year, month, day int) *Date {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%d-%02d-%02d", year, month, day)

	return b
}

// SetIf sets the initial date from the year, month, and day if cond is true
// and returns the date picker
func (b *Date) SetIf(cond bool, year, month, day int) *Date {
	if !cond {
		return b
	}

	return b.Set(year, month, day)
}

// SetWithDate sets the initial date from a time.Time value and returns the date picker
func (b *Date) SetWithDate(d time.Time) *Date {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%d-%02d-%02d", d.Year(), d.Month(), d.Day())

	return b
}

// SetWithDateIf sets the initial date from a time.Time value if cond is true
// and returns the date picker
func (b *Date) SetWithDateIf(cond bool, d time.Time) *Date {
	if !cond {
		return b
	}

	return b.SetWithDate(d)
}

// Set sets the initial time from the hour and minute and returns the time picker
func (b *Time) Set(hour, minute int) *Time {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%02d:%02d", hour, minute)

	return b
}

// SetIf sets the initial time from the hour and minute if cond is true and returns
// the time picker
func (b *Time) SetIf(cond bool, hour, minute int) *Time {
	if !cond {
		return b
	}

	return b.Set(hour, minute)
}

// SetWithDate sets the initial time from a time.Time value and returns the time picker
func (b *Time) SetWithDate(d time.Time) *Time {
	if b == nil {
		return b
	}

	b.InitialValue = fmt.Sprintf("%02d:%02d", d.Hour(), d.Minute())

	return b
}

// SetWithDateIf sets the initial time from a time.Time value if cond is true
// and returns the time picker
func (b *Time) SetWithDateIf(cond bool, d time.Time) *Time {
	if !cond {
		return b
	}

	return b.SetWithDate(d)
}

// ////////////////////////////////////////////////////////////////////////////////// //
