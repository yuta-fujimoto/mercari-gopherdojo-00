package convertImage

import "os"

// valid image format(PGN is for output)
const (
	JPEG = ".jpg"
	PNG  = ".png"
	PGM  = ".pgm"
	GIF  = ".gif"
)

// ConvertImage at first opens all files and specifies formats to store them into Params
type Params struct {
	Infile  []*os.File
	Outfile []*os.File
	Inform  string
	Outform string
}
