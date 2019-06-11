package font

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
)

// LoadFont loads a font file and returns *truetype.Font.
func LoadFont(fontFile string) (*truetype.Font, error) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	return freetype.ParseFont(fontBytes)
}

func NewFreeTypeContext() *freetype.Context {
	return freetype.NewContext()
}
