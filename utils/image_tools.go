package utils

import (
	"image"
	"image/color"
	"image/draw"
	"bytes"
	"image/png"
)

func ConvertImgPalette(img image.Image, palette color.Palette) image.PalettedImage{
	ret := image.NewPaletted(img.Bounds(), palette)
	draw.Draw(ret, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return ret
}

func ImageToBuffer(img image.Image) (*bytes.Buffer, error) {
	buf := bytes.Buffer{}
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	err := enc.Encode(&buf, img)
	if err != nil{
		return nil, err
	}
	return &buf, nil
}

func Base32Image(img image.Image, remove_padding bool) (string, error) {
	buf, err := ImageToBuffer(img)
	if err != nil{
		return "", err
	}

	b32, err := Base32(buf, remove_padding)
	if err != nil{
		return "", err
	}

	return b32, nil
}
