package main

import (
	"fmt"
	"strings"
)

// https://github.com/myano/jenni/wiki/IRC-String-Formatting

const (
	COLOR_WHITE       = "00"
	COLOR_BLACK       = "01"
	COLOR_BLUE        = "02"
	COLOR_NAVY        = COLOR_BLUE
	COLOR_GREEN       = "03"
	COLOR_RED         = "04"
	COLOR_BROWN       = "05"
	COLOR_MAROON      = COLOR_BROWN
	COLOR_PURPLE      = "06"
	COLOR_ORANGE      = "07"
	COLOR_YELLOW      = "08"
	COLOR_LIGHT_GREEN = "09"
	COLOR_LIME        = COLOR_LIGHT_GREEN
	COLOR_TEAL        = "10"
	COLOR_LIGHT_CYAN  = "11"
	COLOR_CYAN        = COLOR_LIGHT_CYAN
	COLOR_AQUA        = COLOR_LIGHT_CYAN
	COLOR_LIGHT_BLUE  = "12"
	COLOR_ROYAL       = COLOR_LIGHT_BLUE
	COLOR_PINK        = "13"
	COLOR_GREY        = "14"
	COLOR_LIGHT_GRAY  = "15"
	COLOR_SILVER      = COLOR_LIGHT_GRAY

	// format
	FMT_BOLD       = "\x02"
	FMT_COLORED    = "\x03"
	FMT_ITALIC     = "\x1D"
	FMT_UNDERLINED = "\x1F"
	FMT_SWAP_BG_FG = "\x16"
	FMT_RESET      = "\x0F"
)

var irc2rgb = map[string]string{
	"00": "#D73C2C",
	"01": "#870000",
	"02": "#ca2c68",
	"03": "#fa5c98",
	"04": "#7e349d",
	"05": "#3e005d",
	"06": "#0067b0",
	"07": "#22A7F0",
	"08": "#009b90",
	"09": "#106b60",
	"10": "#106b60",
	"11": "#006c11",
	"12": "#f9b32f",
	"13": "#e67e22",
	"14": "#5c6a79",
	"15": "#2c3a49",
}

func fmtUrlMeta(title string, mime string, err error, url string) string {
	if err != nil {
		return fmt.Sprintf("%s <- %s", err.Error(), url)
	}
	if strings.HasPrefix(strings.ToLower(mime), "text/html") {
		return fmt.Sprintf("Title: %s", title)
	}
	return fmt.Sprintf("Title: %s %s", title, mime)
}
