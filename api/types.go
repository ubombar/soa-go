package api

import (
	"github.com/ubombar/soa/internal/datetime"
)

// This is an interface that has a Kind() string function. It is used for
// note headers.
type Kinder interface {
	Kind() string
}

type QuestionHeader struct {
	Created  datetime.Date `buffer:"created"`  // creation date
	Question string        `buffer:"question"` // actual question string
	From     string        `buffer:"from"`     // this denotes where this note is spawed from
	Tags     []string      `buffer:"tags"`     // tags of the note
}

func (h QuestionHeader) Kind() string {
	return "question"
}

type QuestionBuffer struct {
	Header *QuestionHeader
}

type LiteratureHeader struct {
	Created datetime.Date `buffer:"created"` // creation date
	PDF     string        `buffer:"pdf"`     // path to the pdf file
	Tags    []string      `buffer:"tags"`    // tags of the note
}

func (h LiteratureHeader) Kind() string {
	return "literature"
}

type MeetingHeader struct {
	Created datetime.Date `buffer:"created"` // creation date
}

func (h MeetingHeader) Kind() string {
	return "meeting"
}

type PermanentHeader struct {
	Created datetime.Date `buffer:"created"` // creation date
}

func (h PermanentHeader) Kind() string {
	return "permanent"
}

type ZoteroAttachementResponse struct {
	JSONRPC string                  `json:"jsonrpc"`
	Result  []ZoteroAttachementItem `json:"result"`
	ID      *string                 `json:"id"` // or simply string if `null` never appears.
}

type ZoteroAttachementItem struct {
	Open        string             `json:"open"`
	Path        string             `json:"path"`
	Annotations []ZoteroAnnotation `json:"annotations"`
}

type AnnotationType string

const (
	Ink       AnnotationType = "ink"
	Note      AnnotationType = "note"
	Highlight AnnotationType = "highlight"
	Text      AnnotationType = "text"
	Image     AnnotationType = "image"
	Underline AnnotationType = "underline"
)

type AnnotationColor string

const (
	// used colors in personal annotation
	ColorYellow AnnotationColor = "#ffd400" // important highlights
	ColorRed    AnnotationColor = "#ff6666" // disagreements, needs fact check
	ColorGreen  AnnotationColor = "#5fb236" // unknown words/concepts
	ColorBlue   AnnotationColor = "#2ea8e5" // agreements and good ideas
	ColorPurple AnnotationColor = "#a28ae5" // very important highlights, facts

	// unused colors, to prevent errors still added here
	ColorMagenta AnnotationColor = "#e56eee" // unused
	ColorOrage   AnnotationColor = "#f19837" // unused
	ColorGray    AnnotationColor = "#aaaaaa" // unused
)

var AnnotationColorNames = map[AnnotationColor]string{
	ColorYellow:  "Yellow",
	ColorRed:     "Red",
	ColorGreen:   "Green",
	ColorBlue:    "Blue",
	ColorPurple:  "Purple",
	ColorMagenta: "Magenta",
	ColorOrage:   "Orage",
	ColorGray:    "Gray",
}

type ZoteroAnnotation struct {
	Key                  string                   `json:"key"`
	Version              int                      `json:"version"`
	ItemType             string                   `json:"itemType"`   // always annotation
	ParentItem           string                   `json:"parentItem"` // mostly itemid of pdf
	AnnotationType       AnnotationType           `json:"annotationType"`
	AnnotationAuthorName string                   `json:"annotationAuthorName"` // mostly empty
	AnnotationText       string                   `json:"annotationText"`
	AnnotationComment    string                   `json:"annotationComment"`
	AnnotationColor      AnnotationColor          `json:"annotationColor"`
	AnnotationPageLabel  string                   `json:"annotationPageLabel"`
	AnnotationSortIndex  string                   `json:"annotationSortIndex"`
	AnnotationPosition   ZoteroAnnotationPosition `json:"annotationPosition"`
	Tags                 []string                 `json:"tags"`
	Relations            map[string]interface{}   `json:"relations"` // Or a more specific type if known
	DateAdded            datetime.DateTime        `json:"dateAdded"`
	DateModified         datetime.DateTime        `json:"dateModified"`
}

type ZoteroAnnotationPosition struct {
	PageIndex int         `json:"pageIndex"` // index of the page, first page -> 0
	Rects     [][]float64 `json:"rects"`
}

type Citation []ZoteroCitationEntry

type ZoteroCitationEntry struct {
	ID             int               `json:"id"`
	Locator        string            `json:"locator"`
	SuppressAuthor bool              `json:"suppressAuthor"`
	Prefix         string            `json:"prefix"`
	Suffix         string            `json:"suffix"`
	Label          string            `json:"label"`
	CitationKey    string            `json:"citationKey"`
	ItemType       string            `json:"itemType"`
	Title          string            `json:"title"`
	Item           ZoteroItemDetails `json:"item"`
}

type ZoteroItemDetails struct {
	Version          int                    `json:"version"`
	ItemType         string                 `json:"itemType"`
	Title            string                 `json:"title"`
	AbstractNote     string                 `json:"abstractNote,omitempty"`
	Date             string                 `json:"date,omitempty"`
	ShortTitle       string                 `json:"shortTitle,omitempty"`
	Language         string                 `json:"language,omitempty"`
	LibraryCatalog   string                 `json:"libraryCatalog"`
	URL              string                 `json:"url,omitempty"`
	AccessDate       string                 `json:"accessDate,omitempty"`
	Volume           string                 `json:"volume,omitempty"`
	Pages            string                 `json:"pages,omitempty"`
	PublicationTitle string                 `json:"publicationTitle,omitempty"`
	DOI              string                 `json:"DOI,omitempty"`
	Issue            string                 `json:"issue,omitempty"`
	ISSN             string                 `json:"ISSN,omitempty"`
	Creators         []ZoteroCreator        `json:"creators"`
	Tags             []interface{}          `json:"tags"`
	Collections      []string               `json:"collections"`
	Relations        map[string]interface{} `json:"relations"`
	DateAdded        string                 `json:"dateAdded"`
	DateModified     string                 `json:"dateModified"`
	URI              string                 `json:"uri"`
	CitationKey      string                 `json:"citationKey"`
	ItemID           int                    `json:"itemID"`
	ItemKey          string                 `json:"itemKey"`
	LibraryID        int                    `json:"libraryID"`
	Attachments      []interface{}          `json:"attachments"` // this is usually not visible
}

type ZoteroCreator struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	CreatorType string `json:"creatorType"`
}
