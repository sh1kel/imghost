package main

import (
	"mime/multipart"
	"fmt"
	//"path"
	"io/ioutil"
	//"github.com/rainycape/unidecode"
	"strings"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"github.com/nfnt/resize"
	"io"
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

type uploadedFile struct {
	Name string `json:"name"`
	User string	`json:"user"`
	Size int64  `json:"size"`
}

type uploadedFiles struct {
	dir   string
	items []uploadedFile
}

func saveFile(file multipart.File, handle *multipart.FileHeader, extenstion string, userName string) (string, int64, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("%v", err)
		return "", 0, err
	}
	if _, err := os.Stat(uploadDir + userName); os.IsNotExist(err) {
		os.Mkdir(uploadDir + userName, 0775)
	}
	hash := md5.New()
	_, err = io.Copy(hash, bytes.NewReader(data))
	if err != nil {
		return "", 0, err
	}
	//filename := path.Base(handle.Filename)
	//filename = unidecode.Unidecode(filename)
	filename := hex.EncodeToString(hash.Sum(nil)) + "." + extenstion
	err = ioutil.WriteFile(uploadDir + userName + "/" + filename, data, 0664)
	if err != nil {
		fmt.Printf("%v", err)
		return "", 0, err
	}
	fSize, err := os.Stat(uploadDir + userName + "/" + filename)
	if err != nil {
		fmt.Printf("%v", err)
		return "", 0, err
	}
	go createThumbnail(filename, userName)
	return filename, fSize.Size(), nil
}

func scanUploads(dir string, userName string) *uploadedFiles {
	f := new(uploadedFiles)
	f.scan(dir + userName + "/.", userName)
	return f
}

func (f *uploadedFiles) scan(dir string, user string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}
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
		f.add(info.Name(), user, info.Size())
		return nil
	})
}

func (f *uploadedFiles) add(name string, user string, size int64) uploadedFile {
	uf := uploadedFile{
		Name: name,
		User: user,
		Size: size,
	}
	f.items = append(f.items, uf)

	return uf
}

func createThumbnail(fileName string, userName string) {
	if _, err := os.Stat(uploadDir + userName + "/thumbnail"); os.IsNotExist(err) {
		os.Mkdir(uploadDir + userName + "/thumbnail", 0775)
	}
	file, err := os.Open(uploadDir + userName + "/" + fileName)
	if err != nil {
		fmt.Printf("Can not open source file: %v\n", err)

		return
	}
	defer file.Close()

	name := strings.ToLower(fileName)

	out, err := os.OpenFile(uploadDir + userName + "/thumbnail/" + fileName,
		os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("Can not create file: %v\n", err)
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
