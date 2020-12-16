package main

import (
	"archive/zip"
	"fmt"
	"image/color"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

func changeContent(c fyne.Canvas) {
	time.Sleep(time.Second * 2)

	c.SetContent(canvas.NewRectangle(color.Black))

	time.Sleep(time.Second * 2)
	c.SetContent(canvas.NewLine(color.Gray{0x66}))

	time.Sleep(time.Second * 2)
	circle := canvas.NewCircle(color.White)
	circle.StrokeWidth = 4
	circle.StrokeColor = color.RGBA{0xff, 0x33, 0x33, 0xff}
	c.SetContent(circle)

	time.Sleep(time.Second * 2)
	c.SetContent(canvas.NewImageFromResource(theme.FyneLogo()))
}

// Retrieves all Microsoft Word files recursively.
func getFiles(directoryPath string) map[string]struct{} {
	fileMap := map[string]struct{}{}
	_ = filepath.Walk(directoryPath, func(path string, info os.FileInfo, error error) error {
		if !info.IsDir() {
			splitPath := strings.Split(path, "\\")
			fileName := splitPath[len(splitPath)-1]
			splitName := strings.Split(fileName, ".")
			fileExtension := splitName[len(splitName)-1]
			if fileExtension == "docx" && fileName[0] != '~' { // The names of temporary Word files start with '~'.
				fileMap[path] = struct{}{}
			}
		}
		return error
	})
	return fileMap
}

// Returns a paragraph from the paragraph array.
func getParagraph(paragraphArray []CleanParagraph) CleanParagraph {
	rand.Seed(time.Now().UnixNano())
	var randomParagraphIndex = rand.Intn(len(paragraphArray))
	// print(paragraphArray[randomParagraphIndex].FileName) // DEBUG
	// print("\n") // DEBUG
	// print(paragraphArray[randomParagraphIndex].Text) // DEBUG
	return paragraphArray[randomParagraphIndex]
}

// Returns file name without extension.
func getFileName(filePath string) string {
	splitPath := strings.Split(filePath, "\\")
	fileName := splitPath[len(splitPath)-1]
	splitName := strings.Split(fileName, ".")
	shortFileName := strings.Join(splitName[0:len(splitName)-1], ".")
	return shortFileName
}

// Wraps text and adds proper page breaks for bullet points.
func formatText(text string, wordLimit int) string {
	wrappedText := textWrap(text, wordLimit)
	splitBulletedText := strings.Split(wrappedText, "-")
	bulletedText := splitBulletedText[0]
	for _, line := range splitBulletedText[1:] {
		bulletedText = bulletedText + "\n- " + strings.ReplaceAll(line, "\n", " ")
	}
	return bulletedText
}

// Returns the smallest of two integers.
func min(integerOne, integerTwo int) int {
	if integerOne > integerTwo {
		return integerTwo
	}
	return integerOne
}

// Returns the largest of two integers.
func max(integerOne, integerTwo int) int {
	if integerOne > integerTwo {
		return integerOne
	}
	return integerTwo
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// Wraps text according to a given word limit per line.
func textWrap(inputString string, limit int) string {
	if strings.TrimSpace(inputString) == "" {
		return inputString
	}

	// Convert string to slice.
	strSlice := strings.Fields(inputString)

	var result string = ""

	if len(strSlice) > limit {
		for len(strSlice) >= 1 {
			// Convert slice/array back to string but insert \r\n at specified limit.
			result = result + strings.Join(strSlice[:limit], " ") + "\r\n"

			// Discard the elements that were copied over to result.
			strSlice = strSlice[limit:]

			// Change the limit to cater for the last few words in.
			if len(strSlice) < limit {
				limit = len(strSlice)
			}
		}
	} else {
		return strings.Join(strSlice, " ")
	}
	return result
}
