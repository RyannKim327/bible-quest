package web_components

import (
	"fmt"
	"html"
	"strings"
)

type spriteStyle struct {
	hair     string
	robe     string
	shade    string
	trim     string
	longHair bool
	beard    bool
	veil     bool
	halo     bool
	crown    bool
	staff    bool
	key      bool
	scroll   bool
}

type pixelCharacter struct {
	Name  string
	Role  string
	Verse string
	Style spriteStyle
}

var pixelCharacters = []struct {
	Key string
	pixelCharacter
}{
	{"jesus", pixelCharacter{"Jesus", "The Savior", "John 3:16", spriteStyle{hair: "#6b4429", robe: "#fffaf0", shade: "#e8dcc4", trim: "#c0563f", longHair: true, beard: true, halo: true}}},
	{"abraham", pixelCharacter{"Abraham", "Father of faith", "Genesis 15:6", spriteStyle{hair: "#e8e2d5", robe: "#e7c98a", shade: "#c9a45f", trim: "#8a5a2b", beard: true, staff: true}}},
	{"moses", pixelCharacter{"Moses", "The lawgiver", "Exodus 3:10", spriteStyle{hair: "#e8e2d5", robe: "#5b86b8", shade: "#3f6697", trim: "#e0b23c", beard: true, staff: true}}},
	{"david", pixelCharacter{"David", "The shepherd king", "Psalm 23:1", spriteStyle{hair: "#7a4a25", robe: "#7a5a97", shade: "#5c4075", trim: "#e0b23c", crown: true}}},
	{"peter", pixelCharacter{"Peter", "The rock", "1 Peter 5:7", spriteStyle{hair: "#4a3423", robe: "#4f8a72", shade: "#376754", trim: "#e0b23c", beard: true, key: true}}},
	{"paul", pixelCharacter{"Paul", "The missionary", "Philippians 4:13", spriteStyle{hair: "#3f3a34", robe: "#d9c7a3", shade: "#b8a37d", trim: "#4a6fa5", beard: true, scroll: true}}},
	{"mary", pixelCharacter{"Mary", "The mother", "Luke 1:46", spriteStyle{hair: "#8a5a2b", robe: "#6f9bd1", shade: "#4f79aa", trim: "#fffaf0", veil: true}}},
	{"daniel", pixelCharacter{"Daniel", "The faithful one", "Daniel 6:23", spriteStyle{hair: "#2f2a26", robe: "#c98b5a", shade: "#a86e42", trim: "#e0b23c"}}},
}

var basePalette = map[rune]string{
	'o': "#2b2118",
	's': "#e8b98a",
	'e': "#2b2118",
	'x': "#4a3423",
	'g': "#e0a52a",
	'w': "#fffdf5",
	'a': "#9a6b3f",
}

func baseRows() []string {
	return []string{
		"................",
		".....oooooo.....",
		"....ohhhhhho....",
		"...ohhhhhhhho...",
		"...ohssssssho...",
		"...ohsessesho...",
		"...ohssssssho...",
		"....osssssso....",
		"......osso......",
		"...oorrrrrroo...",
		"..orrrrrrrrrro..",
		"..orrRRRRrrrro..",
		"..orrRRRRrrrro..",
		"..orrRttRRrrro..",
		"..orrrrrrrrrro..",
		"....orrrrrro....",
		"....orrrrrro....",
		"...oxo....oxo...",
		"..oxxo....oxxo..",
		"................",
	}
}

func setPixel(rows []string, x, y int, value rune) {
	if y < 0 || y >= len(rows) || x < 0 || x >= len(rows[y]) {
		return
	}
	cells := []rune(rows[y])
	cells[x] = value
	rows[y] = string(cells)
}

func buildRows(style spriteStyle) []string {
	rows := baseRows()

	if style.veil {
		for y := 2; y <= 7; y++ {
			rows[y] = strings.ReplaceAll(rows[y], "h", "v")
		}
		setPixel(rows, 4, 7, 'v')
		setPixel(rows, 11, 7, 'v')
		setPixel(rows, 4, 8, 'v')
		setPixel(rows, 11, 8, 'v')
	}

	if style.longHair {
		setPixel(rows, 4, 7, 'h')
		setPixel(rows, 11, 7, 'h')
		setPixel(rows, 4, 8, 'h')
		setPixel(rows, 11, 8, 'h')
	}

	if style.beard {
		for x := 6; x <= 9; x++ {
			setPixel(rows, x, 6, 'h')
		}
		for x := 5; x <= 10; x++ {
			setPixel(rows, x, 7, 'h')
		}
	}

	if style.halo {
		rows[0] = "....gggggggg...."
	}

	if style.crown {
		rows[0] = "....g..g..g....."
		rows[1] = "....gggggggg...."
	}

	if style.staff {
		for y := 8; y <= 17; y++ {
			setPixel(rows, 1, y, 'a')
		}
		setPixel(rows, 2, 11, 's')
		setPixel(rows, 2, 12, 's')
	}

	if style.key {
		for y := 10; y <= 13; y++ {
			setPixel(rows, 15, y, 'g')
		}
		setPixel(rows, 14, 10, 'g')
		setPixel(rows, 14, 12, 'g')
	}

	if style.scroll {
		for y := 10; y <= 13; y++ {
			setPixel(rows, 14, y, 'w')
			setPixel(rows, 15, y, 'w')
		}
		setPixel(rows, 13, 10, 'o')
	}

	return rows
}

func pixelSVG(rows []string, palette map[rune]string, class string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, `<svg class="%s" viewBox="0 0 %d %d" width="%d" height="%d" shape-rendering="crispEdges" aria-hidden="true" focusable="false">`, html.EscapeString(class), len(rows[0]), len(rows), len(rows[0]), len(rows))

	for y, row := range rows {
		for x := 0; x < len(row); {
			value := rune(row[x])
			if value == '.' {
				x++
				continue
			}

			run := 1
			for x+run < len(row) && rune(row[x+run]) == value {
				run++
			}

			if color, ok := palette[value]; ok {
				fmt.Fprintf(&builder, `<rect x="%d" y="%d" width="%d" height="1" fill="%s"/>`, x, y, run, color)
			}
			x += run
		}
	}

	builder.WriteString("</svg>")
	return builder.String()
}

func spritePalette(style spriteStyle) map[rune]string {
	palette := make(map[rune]string, len(basePalette)+3)
	for key, value := range basePalette {
		palette[key] = value
	}
	palette['h'] = style.hair
	palette['v'] = "#a9c7e8"
	palette['r'] = style.robe
	palette['R'] = style.shade
	palette['t'] = style.trim
	return palette
}

func CharacterSprite(key, class string) string {
	for _, character := range pixelCharacters {
		if character.Key == key {
			return pixelSVG(buildRows(character.Style), spritePalette(character.Style), class)
		}
	}
	return ""
}

func PixelBird(class string) string {
	rows := []string{
		"..b...b..",
		".bbb.bbb.",
		"..bbbbb..",
		"...bbb...",
	}
	return pixelSVG(rows, map[rune]string{'b': "#2b2118"}, class)
}

func PixelCross(class string) string {
	rows := []string{
		"....gg....",
		"....gg....",
		"gggggggggg",
		"gggggggggg",
		"....gg....",
		"....gg....",
	}
	return pixelSVG(rows, map[rune]string{'g': "#c0563f"}, class)
}
