package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"strings"
)

const (
	CommonType   = "Common"
	UncommonType = "Uncommon"
	RareType     = "Rare"
	MythicalType = "Mythical"
	LegendType   = "Legend"
)

type ImageLayer struct {
	Image image.Image
	XPos  int
	YPos  int
}

type BgProperty struct {
	Width   int
	Length  int
	BgColor color.Color
}

type ImageGenerator struct {
	layers    []*ImageLayer
	metadatas []*MetaplexMetadata
}

func NewImageGenerator() *ImageGenerator {
	return &ImageGenerator{
		layers:    make([]*ImageLayer, 0),
		metadatas: make([]*MetaplexMetadata, 0),
	}
}

func (imageGenerator ImageGenerator) AddLayer(layerPath string) {
	/*
		layerFiles = GetFiles(STROKES_PATH, strokeType)

		img, e := OpenImage(layer)
		if e != nil {
			panic(e)
		}

		imageGenerator.layers = append(imageGenerator.layers, &ImageLayer{
			Image: img,
			XPos:  0,
			YPos:  0,
		})
	*/
}

func (imageGenerator ImageGenerator) GenerateImages(imageNumber int) {
	//backgrounds := GetFiles(BACKGROUNDS_PATH, "All")

	for i := 0; i < imageNumber; i++ {
		layersPercents := make([]int, 0)
		layersTypes := make([]string, 0)
		//layersImages := make([]string, 0)

		// gen layers percents
		for j := 0; j < len(imageGenerator.layers); j++ {
			layersPercents = append(layersPercents, GenerateRandomNumber(0, 100))
		}

		// calc layers percents
		for k := 0; k < len(imageGenerator.layers); k++ {
			if layersPercents[k] >= 0 && layersPercents[k] < 50 {
				layersTypes[k] = CommonType
			} else if layersPercents[k] >= 50 && layersPercents[k] < 78 {
				layersTypes[k] = UncommonType
			} else if layersPercents[k] >= 78 && layersPercents[k] < 93 {
				layersTypes[k] = RareType
			} else if layersPercents[k] >= 93 && layersPercents[k] < 98 {
				layersTypes[k] = MythicalType
			} else {
				layersTypes[k] = LegendType
			}
		}
	}

}

func GetObjectName(object string) string {
	s := strings.Split(object, "/")
	objectName := s[len(s)-1]
	objectName = strings.Replace(objectName, ".png", "", 4)
	return objectName
}

func OpenImage(path string) (image.Image, error) {
	p := filepath.FromSlash(path)

	var file, err = os.OpenFile(p, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	imageFile, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return imageFile, err
}

func GetObjectRarity(object string) string {

	if strings.Contains(object, "Uncommon") {
		return "Uncommon"
	}

	if strings.Contains(object, "Common") {
		return "Common"
	}

	if strings.Contains(object, "Rare") {
		return "Rare"
	}

	if strings.Contains(object, "Mythical") {
		return "Mythical"
	}

	if strings.Contains(object, "Legend") {
		return "Legend"
	}

	return ""
}

func OverlapImages(imageLayers []ImageLayer, width, height int) (*image.RGBA, error) {
	// create empty RGBA rectagle with traits image size
	generatedImage := image.NewRGBA(image.Rect(0, 0, width, length))

	for _, layer := range imageLayers {
		// set pointer on X Y coordinate
		offset := image.Pt(layer.XPos, layer.YPos)

		draw.Draw(generatedImage, // draw on image buffer
			layer.Image.Bounds().Add(offset), // set bounds
			layer.Image, // set RGB pixel matrix
			image.Point{},
			draw.Over)
	}

	return generatedImage, nil
}
