package pack_bps

import (
	"bytes"
	b64 "encoding/base64"
	"io"
	"log"

	"github.com/andybalholm/brotli"
)

func BpsPack(data string) string {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, brotli.BestCompression)
	_, _ = w.Write([]byte(data))
	_ = w.Close()
	wStr := b64.StdEncoding.EncodeToString(b.Bytes())
	return wStr
}

func BpsUnpack(w string) string {
	wBytes, _ := b64.StdEncoding.DecodeString(w)
	reader := brotli.NewReader(bytes.NewReader(wBytes))
	if plain, err := io.ReadAll(reader); err != nil {
		log.Fatalln(err)
	} else {
		return string(plain)
	}
	return ""
}
