package helpers

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "image/jpeg"

	"github.com/aniketDinda/zocket/models"
	compression "github.com/nurlantulemisov/imagecompression"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DownloadImages(product *models.Product) error {

	if err := os.Mkdir(product.ProductID.Hex(), os.ModePerm); err != nil {
		log.Fatal(err)
		return err
	}
	for index, imageURL := range product.ProductImages {
		response, e := http.Get(imageURL)
		if e != nil {
			log.Fatal(e)
		}
		defer response.Body.Close()

		filePath := fmt.Sprintf("%s/%d.png", product.ProductID.Hex(), index+1)
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func CompressImages(id string) ([]string, error) {
	var compressedUrls []string
	direcPath := fmt.Sprintf("%s_compress", id)
	if err := os.Mkdir(direcPath, os.ModePerm); err != nil {
		log.Fatal(err)
		return nil, err
	}
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	files, err := ioutil.ReadDir(id)
	if err != nil {
		fmt.Println("error reading directory:", err)
		return nil, err
	}
	for index, file := range files {
		filePath := fmt.Sprintf("%s/%s", id, file.Name())

		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf(err.Error())
		}

		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatalf(err.Error())
		}

		compressing, _ := compression.New(99)
		compressingImage := compressing.Compress(img)

		savePath := fmt.Sprintf("%s/%d.png", direcPath, index+1)
		f, err := os.Create(savePath)
		if err != nil {
			log.Fatalf("error creating file: %s", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatalf(err.Error())
			}
		}(f)

		err = png.Encode(f, compressingImage)
		if err != nil {
			log.Fatalf(err.Error())
		}

		fileLoc := fmt.Sprintf("%s/%s", path, savePath)
		compressedUrls = append(compressedUrls, fileLoc)
	}
	return compressedUrls, nil
}

func MapProductInputToProductModel(input *models.ProductInput) models.Product {
	return models.Product{
		ProductID:               primitive.NewObjectID(),
		ProductName:             input.ProductName,
		ProductDescription:      input.ProductDescription,
		ProductImages:           input.ProductImages,
		ProductPrice:            input.ProductPrice,
		CompressedProductImages: []string{},
		CreatedAt:               time.Time{},
		UpdatedAt:               time.Time{},
	}
}
