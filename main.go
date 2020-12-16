package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
)

// CleanParagraph is a parameter in the Word XML document.
type CleanParagraph struct {
	Text     string
	FileName string
	FilePath string
}

// WordDocument is a parameter in the Word XML document.
type WordDocument struct {
	XMLName xml.Name `xml:"document" json:"-"`
	Body    Body     `xml:"body" json:"body"`
}

// Body is a parameter in the Word XML document.
type Body struct {
	XMLName    xml.Name    `xml:"body" json:"-"`
	Paragraphs []Paragraph `xml:"p" json:"paragraph"`
}

// Paragraph is a parameter in the Word XML document.
type Paragraph struct {
	XMLName             xml.Name            `xml:"p" json:"-"`
	ParagraphParameters ParagraphParameters `xml:"pPr" json:"paragraphParameters"`
	Runs                []Run               `xml:"r" json:"runs"`
}

// ParagraphParameters is a parameter in the Word XML document.
type ParagraphParameters struct {
	XMLName             xml.Name            `xml:"pPr" json:"-"`
	PStyle              PStyle              `xml:"pStyle" json:"pStyle"`
	NumberingParameters NumberingParameters `xml:"numPr" json:"numberingParameters"`
	RunParameters       RunParameters       `xml:"rPr" json:"runParameters"`
}

// NumberingParameters is a parameter in the Word XML document.
type NumberingParameters struct {
	XMLName     xml.Name    `xml:"numPr" json:"-"`
	NumberingID NumberingID `xml:"numId" json:"numberingID"`
}

// NumberingID is a parameter in the Word XML document.
type NumberingID struct {
	XMLName        xml.Name `xml:"numId" json:"-"`
	NumberingIDInt int      `xml:"val,attr" json:"numberingIDString"`
}

// PStyle is a parameter in the Word XML document.
type PStyle struct {
	XMLName      xml.Name `xml:"pStyle" json:"-"`
	PStyleString string   `xml:"val,attr" json:"PStyleString,omitempty"`
}

// Run is a parameter in the Word XML document.
type Run struct {
	XMLName       xml.Name      `xml:"r" json:"-"`
	RunParameters RunParameters `xml:"rPr" json:"runParameters"`
	TextPart      string        `xml:"t" json:"textPart"`
}

// RunParameters is a parameter in the Word XML document.
type RunParameters struct {
	XMLName           xml.Name          `xml:"rPr" json:"-"`
	Italics           Italics           `xml:"i" json:"italics,omitempty"`
	Bold              Bold              `xml:"b" json:"bold,omitempty"`
	TextSize          TextSize          `xml:"sz" json:"textSize,omitempty"`
	Language          Language          `xml:"lang" json:"language,omitempty"`
	VerticalAlignment VerticalAlignment `xml:"vertAlign" json:"verticalAlignment,omitempty"`
}

// Italics is a parameter in the Word XML document.
type Italics struct {
	XMLName xml.Name `xml:"i" json:"-"`
}

// Bold is a parameter in the Word XML document.
type Bold struct {
	XMLName xml.Name `xml:"b" json:"-"`
}

// TextSize is a parameter in the Word XML document.
type TextSize struct {
	XMLName        xml.Name `xml:"sz" json:"-"`
	TextSizeString int      `xml:"val,attr" json:"textSizeString,omitempty"`
}

// Language is a parameter in the Word XML document.
type Language struct {
	XMLName        xml.Name `xml:"lang" json:"-"`
	LanguageString string   `xml:"val,attr" json:"languageString,omitempty"`
}

// VerticalAlignment is a parameter in the Word XML document.
type VerticalAlignment struct {
	XMLName                 xml.Name `xml:"vertAlign" json:"-"`
	VerticalAlignmentString string   `xml:"val,attr" json:"verticalAlignmentString,omitempty"`
}

func main() {
	// const directoryPath = "C:\\Users\\Jesper\\Dropbox\\Other\\Interesting\\Books"
	const directoryPath = "C:\\Users\\Jesper\\Dropbox\\Other\\Interesting\\Books\\Test"
	fileMap := getFiles(directoryPath)
	paragraphMap := map[string]string{}
	var paragraphArray = []CleanParagraph{}

	// Loops through all files.
	for filePath := range fileMap {
		fileName := getFileName(filePath)

		// Converts Word document to a collection of XML files.
		Unzip(filePath, "..\\zip_output")

		// Reads text from the main XML file.
		textData, _ := ioutil.ReadFile("C:\\Users\\Jesper\\Dropbox\\Factorial\\zip_output\\word\\document.xml")

		// Parses data and stores it in the defined Go data structures.
		var wordDocument WordDocument
		xml.Unmarshal(textData, &wordDocument)
		json.Marshal(wordDocument)

		paragraphCache := ""
		additionState := false
		startIndex := 0
		// Loops through all paragraphs in a Word document.
		for paragraphIndex, paragraph := range wordDocument.Body.Paragraphs[startIndex:] {
			paragraphText := ""
			lastTextSize := wordDocument.Body.Paragraphs[max(0, startIndex+paragraphIndex-1)].ParagraphParameters.RunParameters.TextSize.TextSizeString
			textSize := paragraph.ParagraphParameters.RunParameters.TextSize.TextSizeString

			lastNumberingID := wordDocument.Body.Paragraphs[max(0, startIndex+paragraphIndex-1)].ParagraphParameters.NumberingParameters.NumberingID.NumberingIDInt
			numberingID := paragraph.ParagraphParameters.NumberingParameters.NumberingID.NumberingIDInt
			nextNumberingID := wordDocument.Body.Paragraphs[min(len(wordDocument.Body.Paragraphs)-1, startIndex+paragraphIndex+1)].ParagraphParameters.NumberingParameters.NumberingID.NumberingIDInt

			// Obtains all text from paragraph, adding superscript and subscripts.
			for _, run := range paragraph.Runs {
				if run.RunParameters.VerticalAlignment.VerticalAlignmentString == "superscript" {
					paragraphText = paragraphText + "^"
				}
				if run.RunParameters.VerticalAlignment.VerticalAlignmentString == "subscript" {
					paragraphText = paragraphText + "_"
				}
				paragraphText = paragraphText + run.TextPart
			}

			// Greater than or equal to as the previous paragraph can either be a list or not.
			if numberingID-lastNumberingID >= 1 && lastTextSize <= 24 {
				additionState = true
				// Smaller than or equal to as the previous paragraph can either be a list or not.
			} else if numberingID-lastNumberingID <= -1 {
				additionState = false
				if textSize <= 12 && len(paragraphCache) != 0 {
					paragraphMap[paragraphCache] = getFileName(filePath)
					cleanParagraph := CleanParagraph{Text: paragraphCache, FileName: fileName, FilePath: filePath}
					paragraphArray = append(paragraphArray, cleanParagraph)
				}
				paragraphCache = ""
			}

			if additionState == true && textSize <= 24 {
				paragraphCache = paragraphCache + "\n- " + paragraphText
			} else {
				paragraphCache = paragraphText
				if nextNumberingID-numberingID < 1 && textSize <= 24 && len(paragraphCache) != 0 {
					paragraphMap[paragraphCache] = getFileName(filePath)
					cleanParagraph := CleanParagraph{Text: paragraphCache, FileName: fileName, FilePath: filePath}
					paragraphArray = append(paragraphArray, cleanParagraph)
					paragraphCache = ""
				}
			}
		}
	}

	// Selects a random paragraph from the collection of all paragraphs.
	// randomParagraph := getParagraph(paragraphArray)
	randomParagraph := paragraphArray[29]

	// Wraps text and adds proper page breaks for bullet points.
	paragraphText := formatText(randomParagraph.Text, 13)

	drawInterface(randomParagraph.FileName, paragraphText, randomParagraph, paragraphArray)
}
