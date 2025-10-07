package domain

import "github.com/h2non/bimg"

type WatermarkOptions struct {
	Path     string
	Opacity  float32
	Position bimg.Gravity
}

type TransformationOptions struct {
	Width      int
	Height     int
	Format     string
	Quality    int
	Crop       string
	TargetType bimg.ImageType
	Watermark  WatermarkOptions
}
