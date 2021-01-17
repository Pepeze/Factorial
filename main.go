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
	XMLName    xml.Name `xml:"b" json:"-"`
	BoldString string   `xml:"val,attr" json:"boldString,omitempty"`
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
	// Initialization of global variables.
	// const directoryPath = "C:\\Users\\Jesper\\Dropbox\\Other\\Interesting\\Books"
	const directoryPath = "C:\\Users\\Jesper\\Dropbox\\Factorial\\Test"
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
		// Loops through all paragraphs in a Word document.
		for paragraphIndex, paragraph := range wordDocument.Body.Paragraphs {
			// Retrieves text sizes.
			lastTextSize := wordDocument.Body.Paragraphs[max(0, paragraphIndex-1)].ParagraphParameters.RunParameters.TextSize.TextSizeString
			textSize := paragraph.ParagraphParameters.RunParameters.TextSize.TextSizeString

			// Retrieves numbering IDs, used for determining if text is in a bulleted list or not.
			lastNumberingID := wordDocument.Body.Paragraphs[max(0, paragraphIndex-1)].ParagraphParameters.NumberingParameters.NumberingID.NumberingIDInt
			numberingID := paragraph.ParagraphParameters.NumberingParameters.NumberingID.NumberingIDInt
			nextNumberingID := wordDocument.Body.Paragraphs[min(len(wordDocument.Body.Paragraphs)-1, paragraphIndex+1)].ParagraphParameters.NumberingParameters.NumberingID.NumberingIDInt

			// Obtains all text from paragraph, adding superscript and subscripts.
			paragraphText := getParagraphText(paragraph.Runs)

			boldIndicator := false
			if len(paragraphText) != 0 {
				if isNumeric(string(paragraphText[0])) && paragraphText[len(paragraphText)-1:] == ":" {
					additionState = true
					boldIndicator = true
				} else {
					additionState = false
				}
			}

			// Greater than or equal to as the previous paragraph can either be a list or not.
			// Means that the current paragraph should be added with the previous one.
			if numberingID-lastNumberingID >= 1 && lastTextSize <= 24 {
				additionState = true
			// Smaller than or equal to as the previous paragraph can either be a list or not.
			// Means that we now stepped out of a bullet list and needs to add the paragraph to the array.
			} else if numberingID-lastNumberingID <= -1 && textSize <= 12 && len(paragraphCache) != 0 {
					paragraphMap[paragraphCache] = fileName
					cleanParagraph := CleanParagraph{Text: paragraphCache, FileName: fileName, FilePath: filePath}
					paragraphArray = append(paragraphArray, cleanParagraph)
					additionState = false
					paragraphCache = ""
			}
			// Special case for when entries from Days are added.
			if additionState == true && boldIndicator == true {
				paragraphCache += "\n " + paragraphText
			// Adds current bullet paragraph to the previous one through the cache.
			} else if additionState == true && textSize <= 24 {
				paragraphCache += "\n- " + paragraphText
			// Simply adds the paragraph text to the array if no bullet lists exist.
			} else {
				paragraphCache += paragraphText
				if nextNumberingID-numberingID < 1 && textSize <= 24 && len(paragraphCache) != 0 {
					paragraphMap[paragraphCache] = fileName
					cleanParagraph := CleanParagraph{Text: paragraphCache, FileName: fileName, FilePath: filePath}
					paragraphArray = append(paragraphArray, cleanParagraph)
					paragraphCache = ""
				}	
			} 
				
		}
	}

	// Selects a random paragraph from the collection of all paragraphs.
	randomParagraph := getParagraph(paragraphArray)

	drawInterface(randomParagraph, paragraphArray)
}
