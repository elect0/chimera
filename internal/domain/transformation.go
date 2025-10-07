package domain

import "github.com/h2non/bimg"

type TransformationOptions struct {
	Width      int
	Height     int
	Format     string
	Quality    int
	Crop       string
	TargetType bimg.ImageType
}
