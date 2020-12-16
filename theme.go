package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var (
	beige  = &color.RGBA{R: 213, G: 177, B: 155, A: 255}
	beige2 = &color.RGBA{R: 193, G: 157, B: 135, A: 255}
	blue   = &color.RGBA{R: 0, G: 153, B: 153, A: 255}
	grey   = &color.Gray{Y: 123}
	grey2  = &color.Gray{Y: 160}
	orange = &color.RGBA{R: 198, G: 123, B: 0, A: 255}
	purple = &color.RGBA{R: 128, G: 0, B: 128, A: 255}
	yellow = &color.RGBA{R: 225, G: 241, B: 167, A: 255}
)

type customTheme struct {
}

func (customTheme) BackgroundColor() color.Color {
	return grey
}

func (customTheme) ButtonColor() color.Color {
	return grey2
}

func (customTheme) DisabledButtonColor() color.Color {
	return color.White
}

func (customTheme) HyperlinkColor() color.Color {
	return orange
}

func (customTheme) TextColor() color.Color {
	return color.White
}

func (customTheme) DisabledTextColor() color.Color {
	return color.Black
}

func (customTheme) IconColor() color.Color {
	return color.White
}

func (customTheme) DisabledIconColor() color.Color {
	return color.Black
}

func (customTheme) PlaceHolderColor() color.Color {
	return grey
}

func (customTheme) PrimaryColor() color.Color {
	return orange
}

func (customTheme) HoverColor() color.Color {
	return grey
}

func (customTheme) FocusColor() color.Color {
	return grey
}

func (customTheme) ScrollBarColor() color.Color {
	return grey
}

func (customTheme) ShadowColor() color.Color {
	return &color.RGBA{0xcc, 0xcc, 0xcc, 0xcc}
}

func (customTheme) TextSize() int {
	return 24
}

func (customTheme) TextFont() fyne.Resource {
	return theme.DefaultTextBoldFont()
}

func (customTheme) TextBoldFont() fyne.Resource {
	return theme.DefaultTextBoldFont()
}

func (customTheme) TextItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

func (customTheme) TextBoldItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

func (customTheme) TextMonospaceFont() fyne.Resource {
	return theme.DefaultTextMonospaceFont()
}

func (customTheme) Padding() int {
	return 10
}

func (customTheme) IconInlineSize() int {
	return 20
}

func (customTheme) ScrollBarSize() int {
	return 10
}

func (customTheme) ScrollBarSmallSize() int {
	return 5
}

func newCustomTheme() fyne.Theme {
	return &customTheme{}
}
