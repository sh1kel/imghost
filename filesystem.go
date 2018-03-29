package main

import (
	"mime/multipart"
	"fmt"
	"path"
	"io/ioutil"
	"github.com/rainycape/unidecode"
	"strings"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"github.com/nfnt/resize"
)

type uploadedFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type uploadedFiles struct {
	dir   string
	items []uploadedFile
}

func saveFile(file multipart.File, handle *multipart.FileHeader) (string, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}
	filename := path.Base(handle.Filename)
	filename = unidecode.Unidecode(filename)
	err = ioutil.WriteFile(uploadDir + filename, data, 0664)
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}
	go createThumbnail(filename)
	return filename, nil
}

func scanUploads(dir string) *uploadedFiles {
	f := new(uploadedFiles)
	f.scan(dir)
	return f
}

func (f *uploadedFiles) scan(dir string) {
	f.dir = dir
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error access path: %v\n", err)
		}
		if info.IsDir() && info.Name() == "thumbnail" {
			return filepath.SkipDir
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		f.add(info.Name(), info.Size())
		return nil
	})
}

func (f *uploadedFiles) add(name string, size int64) uploadedFile {
	uf := uploadedFile{
		Name: name,
		Size: size,
	}
	f.items = append(f.items, uf)

	return uf
}

func createThumbnail(fileName string) {
	file, err := os.Open(uploadDir + fileName)
	if err != nil {
		return
	}
	defer file.Close()

	name := strings.ToLower(fileName)

	out, err := os.OpenFile(thumbsDir + fileName,
		os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return
	}
	defer out.Close()

	if strings.HasSuffix(name, ".jpg") {
		img, err := jpeg.Decode(file)
		if err != nil {
			return
		}

		resized := resize.Thumbnail(180, 180, img, resize.Lanczos3)
		jpeg.Encode(out, resized,
			&jpeg.Options{Quality: jpeg.DefaultQuality})

	} else if strings.HasSuffix(name, ".png") {
		img, err := png.Decode(file)
		if err != nil {
			return
		}

		resized := resize.Thumbnail(180, 180, img, resize.Lanczos3)
		png.Encode(out, resized)
	}
}
