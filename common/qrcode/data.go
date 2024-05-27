package qrcode

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

//	func GenerateQrCode(value string) ([]byte, error) {
//		byteArr, err := qrcode.Encode(value, qrcode.Medium, 250)
//		return byteArr, err
//	}
func GenerateQrCode(value string) ([]byte, error) {
	qrc, err := qrcode.NewWith(value, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w2 := standard.NewWithWriter(wr, standard.WithQRWidth(30), standard.WithLogoImageFilePNG("public/image/pitel-logo.png"))
	if err = qrc.Save(w2); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func ByteToString(b []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(b))
}
