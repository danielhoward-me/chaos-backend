package worker

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

func convertToJpg(pngImage []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(pngImage))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
