package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	_ "image/png"

	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ccolor "github.com/fatih/color"
	"github.com/schollz/progressbar/v3"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/h2non/bimg"
	gim "github.com/ozankasikci/go-image-merge"
)

// Collection traits folder and output folder
var COLLECTION_NAME = "ALIENS_FINAL"
var generatedFolder = COLLECTION_NAME + "_GENERATED"
var originalFolder = COLLECTION_NAME + "_GENERATED/" + "originalFolder"
var metadataFolder = COLLECTION_NAME + "_GENERATED/" + "metadataFolder"
var resizedFolder = COLLECTION_NAME + "_GENERATED/" + "resizedFolder"
var mergedFolder = COLLECTION_NAME + "_GENERATED/" + "mergedFolder"

// NUMBER OF ELEMENTS
var COMMON_START_RANGE = 0
var COMMON_END_RANGE = 5001 //

var bar = progressbar.Default(int64(COMMON_END_RANGE))

// image dimensions
var width = 900
var length = 1080

// rarity
var commonRangeFrom = 0
var commonRangeTo = 49

var uncommonRangeFrom = 49
var uncommonRangeTo = 74

var rareRangeFrom = 74
var rareRangeTo = 89

var mythicalRangeFrom = 89
var mythicalRangeTo = 98

// -----------------------------------------------------------------

// needed for generate true random number
var src cryptoSource
var rnd *rand.Rand

// mosaic size
var mosaicXCounter = 20
var mosaicYCounter = 20

var traitsCounter = 0

var metadatas = make([]*OpenseaMetadata, 0)
var layersFiles = make([][]string, 0)

// colors
// Create SprintXxx functions to mix strings with other non-colorized strings:
var yellow = ccolor.New(ccolor.FgYellow).SprintFunc()
var cyan = ccolor.New(ccolor.FgCyan).SprintFunc()
var green = ccolor.New(ccolor.FgGreen).SprintFunc()

/* -------------------------------------------------------------------------------------------------
MAIN
*/

func main() {

	rnd = rand.New(src)

	// GET TRAITS FOLDERS
	traitsFolders := GetTraitsFolders(COLLECTION_NAME)
	CreateGeneratedFolders()

	//usedBodyType := "Common"
	usedBodyName := "Reptilian"

	// get traits folder
	for folderIndex := 0; folderIndex < len(traitsFolders); folderIndex++ {
		folders := GetFiles(traitsFolders[folderIndex], "All")
		layersFiles = append(layersFiles, folders)
		traitsCounter += len(folders)
	}

	// PRINT

	fmt.Println()
	fmt.Printf("%s %s\n", cyan("Traits Folders length:"), yellow(len(traitsFolders)))
	fmt.Println()

	fmt.Println()
	fmt.Printf("%s %s\n", cyan("Traits Number:"), yellow(traitsCounter))
	fmt.Println()
	// ----------------------------------------------------------------------------------------

	//
	// -----------------------------------------------------
	for imagesCounter := COMMON_START_RANGE; imagesCounter <= COMMON_END_RANGE; imagesCounter++ {
		layers := make([]ImageLayer, 0)
		counterString := strconv.Itoa(imagesCounter)

		// set generic metadata for all nfts
		metadata := new(OpenseaMetadata)
		metadata.Name = "PXT #" + counterString
		metadata.Description = "The X-Collective reveals itself, expanding the Project X ecosystem with choice and user-driven utility behind a limited release of NFTs"
		metadata.Image = "https://ipfs.io/ipfs/QmbuUiiF2BBxWnrixBQotoy8e8HoC69j5DFkAJYRfAPDn3?filename=" + counterString + ".png"

		// ---------------------------------------------------------------------------------------------

		// loop about all traits folders
		for layersIndex := 0; layersIndex < len(layersFiles); layersIndex++ {
			fmt.Println("LAYER INDEX ", layersIndex)

			// generate random number for calc object rarity
			rarityPercent := GenerateRandomNumber(0, 99)
			traitType := GetType(rarityPercent)
			trait := GetObjectFromType(layersFiles[layersIndex], traitType)
			//traitTypeString := GetTypeFromFileName(trait)

			addLayer := true

			//fmt.Println(traitTypeString)

			if trait != "" {

				hasTrapCap := FindTrait(metadata.Attributes, "Trap Cap")
				hasSkateHelmet := FindTrait(metadata.Attributes, "Skate Helmet")
				hasCap := FindTrait(metadata.Attributes, "Cap")
				hasMask := FindTrait(metadata.Attributes, "Mask")
				hasSmiley := FindTrait(metadata.Attributes, "Smiley")
				hasSki := FindTrait(metadata.Attributes, "Ski")
				hasVR := FindTrait(metadata.Attributes, "VR")
				hasBeard := FindTrait(metadata.Attributes, "Beard")
				hasDaft := FindTrait(metadata.Attributes, "Daft")

				// ...
				// GET USED BODY
				if strings.Contains(trait, "B_") {
					//usedBodyType = traitType
					usedBodyName = strings.Split(trait, "/")[3]
					usedBodyName = strings.Replace(usedBodyName, ".png", "", 1)
				}

				if strings.Contains(trait, "C_") {
					if strings.Contains(trait, "Mythical") {
						hasSameBodyClothes := strings.Contains(trait, usedBodyName)
						for next := true; next; next = hasSameBodyClothes == false {
							trait = GetObjectFromType(layersFiles[layersIndex], traitType)
							hasSameBodyClothes = strings.Contains(trait, usedBodyName)
						}
					}
				}

				if strings.Contains(trait, "D_") {
					//trait = GetObjectFromType(layersFiles[layersIndex], usedBodyType)

					hasSame := strings.Contains(trait, usedBodyName)
					for next := true; next; next = hasSame == false {
						trait = GetObjectFromType(layersFiles[layersIndex], traitType)
						hasSame = strings.Contains(trait, usedBodyName)
					}

					/*
						hasHelmetPercent := GenerateRandomNumber(0, 99)
						hasCapPercent := GenerateRandomNumber(0, 99)
						hasTrapCapPercent := GenerateRandomNumber(0, 99)

						if hasHelmetPercent > 85 {
							condition := strings.Contains(trait, "Skate Helmet")
							for next := true; next; next = condition == true {
								trait = GetObjectFromType(layersFiles[layersIndex], usedBodyType)
								condition = strings.Contains(trait, "Skate Helmet")
							}
						}

						if hasTrapCapPercent > 75 {
							condition := strings.Contains(trait, "Trap Cap")
							for next := true; next; next = condition == true {
								trait = GetObjectFromType(layersFiles[layersIndex], usedBodyType)
								condition = strings.Contains(trait, "Trap Cap")
							}
						}

						if hasCapPercent > 75 {
							condition := strings.Contains(trait, "Cap")
							for next := true; next; next = condition == true {
								trait = GetObjectFromType(layersFiles[layersIndex], usedBodyType)
								condition = strings.Contains(trait, "Cap")
							}
						}
					*/
				}

				if hasTrapCap {
					if strings.Contains(trait, "H_") || strings.Contains(trait, "VR") {
						if strings.Contains(trait, "Ski Goggles") {
							addLayer = false
						}
					}
				}

				if hasSkateHelmet {
					if strings.Contains(trait, "H_") {
						if strings.Contains(trait, "Ski Goggles") || strings.Contains(trait, "VR") {
							addLayer = false
						}
					}
				}

				if hasCap {
					if strings.Contains(trait, "H_") {
						if strings.Contains(trait, "Ski Goggles") || strings.Contains(trait, "VR") {
							addLayer = false
						}
					}
				}

				if hasTrapCap {
					if strings.Contains(trait, "H_") {
						if strings.Contains(trait, "Crown") {
							addLayer = false
						}
					}
				}

				if hasSkateHelmet {
					if strings.Contains(trait, "H_") {
						if strings.Contains(trait, "Crown") {
							addLayer = false
						}
					}
				}

				if hasCap {
					if strings.Contains(trait, "H_") {
						if strings.Contains(trait, "Crown") {
							addLayer = false
						}
					}
				}

				if hasMask {
					if strings.Contains(trait, "I_") {
						addLayer = false
					}
				}

				if hasSki {
					if strings.Contains(trait, "I_") {
						addLayer = false
					}
				}

				if hasVR {
					if strings.Contains(trait, "I_") {
						addLayer = false
					}
				}

				if hasDaft {
					if strings.Contains(trait, "I_") {
						addLayer = false
					}
				}

				if hasSmiley {
					if strings.Contains(trait, "F_") {
						if strings.Contains(trait, "Rare") {
							addLayer = false
						}
					}
				}

				if hasSmiley {
					if strings.Contains(trait, "F_") {
						if strings.Contains(trait, "Uncommon") {
							hasSame := strings.Contains(trait, "Smiley")
							for next := true; next; next = hasSame == false {
								trait = GetObjectFromType(layersFiles[layersIndex], traitType)
								hasSame = strings.Contains(trait, "Smiley")
							}
						}
					}
				}

				if !hasSmiley {
					if strings.Contains(trait, "F_") {
						if strings.Contains(trait, "Smiley") {
							addLayer = false
						}
					}
				}

				if hasBeard {
					if strings.Contains(trait, "H_") {
						if strings.Contains(trait, "Headphones") {
							addLayer = false
						}
					}
				}

				if addLayer {
					AddLayer(&layers, trait)
					metadata.Attributes = append(metadata.Attributes, Attribute{
						TraitType: strings.Split(strings.Split(trait, "/")[1], "_")[1],
						Value:     GetObjectName(trait),
					})
				}
			}
		}

		//go func () {
		// overlap all layers
		overlapedImage, overlapImagesErr := OverlapImages(layers, width, length)
		if overlapImagesErr != nil {
			log.Printf("Error OverlapImages: %+v\n", overlapImagesErr)
		}

		fmt.Println(counterString)
		savePath := generatedFolder + "/originalFolder/" + counterString + ".png"

		saveImage(overlapedImage, 1, savePath)
		saveOpenSeaMetadataFile(metadata, counterString)

		fmt.Println()
		fmt.Printf("%s %s\n", cyan("Traits Folders length:"), green(savePath))
		fmt.Println()
		bar.Add(1)
		//}()
	}

	//j2, _ := json.Marshal(elements)
	j2, _ := json.MarshalIndent(metadatas, "", "    ")

	// to append to a file
	// create the file if it doesn't exists with O_CREATE, Set the file up for read write, add the append flag and set the permission
	f2, errr2 := os.OpenFile(generatedFolder+"/originalFolder/whole.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer f2.Close()
	if errr2 != nil {
		log.Fatal(errr2)
	}
	// write to file, f.Write()
	f2.Write(j2)

	ResizeImages(COMMON_END_RANGE, 128)
	MergeImages(COMMON_END_RANGE, mosaicXCounter, mosaicYCounter, mosaicXCounter*mosaicYCounter)
}

func encodePNG(img *image.RGBA) []byte {
	pngBuffer := new(bytes.Buffer)
	encodePngErr := png.Encode(pngBuffer, img)
	if encodePngErr != nil {
		fmt.Println(encodePngErr)
	}
	return pngBuffer.Bytes()

}

func saveOpenSeaMetadataFile(metadata *OpenseaMetadata, counter string) {
	meta, _ := json.MarshalIndent(metadata, "", "    ")
	// to append to a file
	// create the file if it doesn't exists with O_CREATE, Set the file up for read write, add the append flag and set the permission
	f3, errr3 := os.OpenFile(generatedFolder+"/originalFolder/"+counter+".json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer f3.Close()
	if errr3 != nil {
		log.Fatal(errr3)
	}
	// write to file, f.Write()
	f3.Write(meta)
	fmt.Println("----------------------------------------------------------------------------")
}

func saveMetadataFile(metadata *MetaplexMetadata, counter string) {
	metadata.Name = "PXT #" + counter
	metadata.Image = counter + ".png"
	metadata.Properties.Files = append(metadata.Properties.Files, File{
		URI:  counter + ".png",
		Type: "image/png",
	})
	meta, _ := json.MarshalIndent(metadata, "", "    ")
	// to append to a file
	// create the file if it doesn't exists with O_CREATE, Set the file up for read write, add the append flag and set the permission
	f3, errr3 := os.OpenFile(generatedFolder+"/originalFolder/"+counter+".json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer f3.Close()
	if errr3 != nil {
		log.Fatal(errr3)
	}
	// write to file, f.Write()
	f3.Write(meta)
	fmt.Println("----------------------------------------------------------------------------")
}

func saveImage(img *image.RGBA, compression int, savePath string) {

	buffer := encodePNG(img)

	options := bimg.Options{
		Width:   width,
		Height:  length,
		Quality: compression,
	}

	// encode image
	outImage, err3 := bimg.NewImage(buffer).Process(options)
	if err3 != nil {
		fmt.Fprintln(os.Stderr, err3)
	}

	if bimg.NewImage(outImage).Type() != "png" {
		panic(errors.New("Erro exporting"))
	}

	// save file
	r := bimg.Write(savePath, outImage)
	if r != nil {
		fmt.Println(r)
	}
}

func FindTraitfFromName(attributes []Attribute, traitName string) bool {
	for i := 0; i < len(attributes); i++ {
		if strings.Contains(attributes[i].Value, traitName) {
			return true
		}
	}
	return false
}

func FindTrait(attributes []Attribute, name string) bool {
	for i := 0; i < len(attributes); i++ {
		if strings.Contains(attributes[i].Value, name) {
			return true
		}
	}
	return false
}

func AddLayer(layers *[]ImageLayer, imagePath string) {
	img, e := openImage(imagePath)
	if e != nil {
		fmt.Println(e)
	}

	*layers = append(*layers, ImageLayer{
		Image: img,
		XPos:  0,
		YPos:  0,
	})
	//fmt.Println("INSERTED LAYER")
	//fmt.Println(imagePath)
}

func AddAttributes(metadata *MetaplexMetadata, name string, value string) {
	metadata.Attributes = append(metadata.Attributes, Attribute{
		TraitType: name,
		Value:     GetObjectName(value),
	})
}

var (
	colorB = [3]float64{248, 54, 0}
	colorA = [3]float64{254, 140, 0}
)

var (
	max = float64(0)
)

func linearGradient(x, y float64) (uint8, uint8, uint8) {
	d := x / max
	r := colorA[0] + d*(colorB[0]-colorA[0])
	g := colorA[1] + d*(colorB[1]-colorA[1])
	b := colorA[2] + d*(colorB[2]-colorA[2])
	return uint8(r), uint8(g), uint8(b)
}

func GenerateGradientImage(filename string, width, height int) *image.RGBA {
	var r1 = randomcolor.GetRandomColorInRgb()
	var r2 = randomcolor.GetRandomColorInRgb()
	max = float64(width)

	colorA = [3]float64{float64(r1.Red), float64(r1.Green), float64(r1.Blue)}
	colorB = [3]float64{float64(r2.Red), float64(r2.Green), float64(r2.Blue)}

	var w, h int = width, height
	dst := image.NewRGBA(image.Rect(0, 0, w, h)) //*NRGBA (image.Image interface)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b := linearGradient(float64(x), float64(y))
			c := color.RGBA{

				r,
				g,
				b,
				255,
			}
			dst.Set(x, y, c)
		}
	}

	//img, _ := os.Create(filename)
	//defer img.Close()
	//png.Encode(img, dst) //Encode writes the Image m to w in PNG format.

	return dst
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := make([]string, 0)
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GetTraitsFolders(root string) []string {
	var folders []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		dir, err2 := os.Stat(path)
		if err2 != nil {
			panic(err2)
		}
		if dir.IsDir() {

			if len(strings.Split(path, "/")) >= 2 {
				folders = append(folders, strings.Split(path, "/")[0]+"/"+strings.Split(path, "/")[1])
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	f := RemoveDuplicateStr(folders)
	//sort.Strings(f)

	return f
}

func CreateGeneratedFolders() {
	// clean folder
	// remove generated folder
	removeGeneratedFolderErr := os.RemoveAll(generatedFolder)
	if removeGeneratedFolderErr != nil {
		fmt.Println(removeGeneratedFolderErr)
	}

	//Create a folder/directory at a full qualified path
	err := os.Mkdir(generatedFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}

	//Create a folder/directory at a full qualified path
	err = os.Mkdir(originalFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}

	//Create a folder/directory at a full qualified path
	err = os.Mkdir(metadataFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}

	//Create a folder/directory at a full qualified path
	err = os.Mkdir(resizedFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}

	//Create a folder/directory at a full qualified path
	err = os.Mkdir(mergedFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func printMap(m map[string]int) {
	var maxLenKey int
	for k, _ := range m {
		if len(k) > maxLenKey {
			maxLenKey = len(k)
		}
	}

	for k, v := range m {
		fmt.Println(k + ": " + strings.Repeat(" ", maxLenKey-len(k)) + strconv.Itoa(v))
	}
}

func GetTypeFromFileName(fileName string) string {
	if strings.Contains(fileName, "Common") {
		return "Common"
	} else if strings.Contains(fileName, "Uncommon") {
		return "Uncommon"
	} else if strings.Contains(fileName, "Rare") {
		return "Rare"
	} else if strings.Contains(fileName, "Mythical") {
		return "Mythical"
	} else {
		return "Legendary"
	}
}

func GetType(traitPercent int) string {
	if traitPercent >= commonRangeFrom && traitPercent < commonRangeTo {
		return "Common"
	} else if traitPercent >= uncommonRangeFrom && traitPercent < uncommonRangeTo {
		return "Uncommon"
	} else if traitPercent >= rareRangeFrom && traitPercent < rareRangeTo {
		return "Rare"
	} else if traitPercent >= mythicalRangeFrom && traitPercent < mythicalRangeTo {
		return "Mythical"
	} else {
		return "Legendary"
	}
}

func openImage(path string) (image.Image, error) {
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

func GetObjectFromType(list []string, fileType string) string {
	objectsList := make([]string, 0)

	if len(list) > 0 {
		for i := 0; i < len(list); i++ {
			if strings.Contains(list[i], fileType) {
				objectsList = append(objectsList, list[i])
			}
		}

		if len(objectsList) > 0 {
			return objectsList[GenerateRandomNumber(0, len(objectsList)-1)]
		}
	}
	return ""
}

func GetFiles(root string, fileType string) []string {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".png") {
			t := GetObjectRarity(path)
			if fileType == "All" {
				files = append(files, path)
			} else {
				if t == fileType {
					files = append(files, path)
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	/*
		for _, file := range files {
			fmt.Println(file)
		}
	*/
	return files
}

func GenerateRandomNumber(min, max int) int {
	ns := make([]int, 0)

	if max == 0 {
		return 0
	}

	for i := 0; i < 100; i++ {
		//fmt.Println("max")
		//fmt.Println(max)
		n := rnd.Intn(max)
		ns = append(ns, n)
	}
	n2 := rnd.Intn(len(ns))

	return ns[n2]
}

func ResizeImages(amount int, width int) {
	for i := 0; i < amount; i++ {
		counterString := strconv.Itoa(i)
		ResizeImage("./"+originalFolder+"/"+counterString+".png", generatedFolder+"/resizedFolder/", width)
	}
}

func MergeImages(amount int, imageCountDX, imageCountDY int, divider int) {

	grids := make([]*gim.Grid, 0)

	for i := 0; i < amount; i++ {
		counterString := strconv.Itoa(i)
		currentFile := resizedFolder + "/" + counterString + ".png"

		g := &gim.Grid{
			ImageFilePath: currentFile,
		}
		grids = append(grids, g)

		fmt.Println(counterString)

		if i != 0 && i%divider == 0 {
			rgba, err := gim.New(grids, imageCountDX, imageCountDY).Merge()
			if err != nil {
				panic(err)
			}

			// save the output to jpg or png
			file, err2 := os.Create(generatedFolder + "/mergedFolder/" + counterString + ".png")
			if err2 != nil {
				panic(err2)
			}
			err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})
			err = png.Encode(file, rgba)
			grids = []*gim.Grid{}
		}
	}
}

func ResizeImage(path string, out string, width int) {
	buffer, err := bimg.Read(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	newImage, err := bimg.NewImage(buffer).Resize(width, width)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	size, err := bimg.NewImage(newImage).Size()
	if size.Width == width && size.Height == width {
		fmt.Println("The image size is valid " + path)
	}

	s := strings.Split(path, "/")
	bimg.Write(out+s[len(s)-1], newImage)
}

/*
	var colorInRGB = randomcolor.GetRandomColorInRgb()
	col := color.RGBA{uint8(colorInRGB.Red), uint8(colorInRGB.Green), uint8(colorInRGB.Blue), 0xaa}
	c := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	draw.Draw(c, c.Bounds(), image.NewUniform(col), image.ZP, draw.Src)

	layers = append(layers, ImageLayer{
		Image: c,
		XPos:  0,
		YPos:  0,
	})
*/

/*
	layers = append(layers, ImageLayer{
		Image: GenerateImage("t", 1024, 1024),
		XPos:  0,
		YPos:  0,
	})
*/

//AddElement(GetFiles("BabypunksOriginal/baseback", "All")[0])
