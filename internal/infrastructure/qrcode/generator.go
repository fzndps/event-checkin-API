package qrcode

import (
	"bytes"
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateQRCode akan menghasilkan QR dalam format PNG bytes
// Input data adalah token yang sudah dibuat sebelumnya
func (g *Generator) GenerateQRCode(data string, size int) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// 256x256 pixel
	if size == 0 {
		size = 256
	}

	// Generate PNG bytse
	pngBytes, err := qr.PNG(size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PNG: %w", err)
	}

	return pngBytes, nil
}

// Menghasilkan qr code dalam format base64 untuk embed langsung di html email
func (g *Generator) GenerateQRCodeBase64(data string, size int) (string, error) {
	pngBytes, err := g.GenerateQRCode(data, size)
	if err != nil {
		return "", err
	}

	base64STR := base64.StdEncoding.EncodeToString(pngBytes)

	return fmt.Sprintf("data:image/png;base64,%s", base64STR), nil
}

func (g *Generator) GenerateQRCodeBuffer(data string, size int) (*bytes.Buffer, error) {
	pngBytes, err := g.GenerateQRCode(data, size)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(pngBytes), nil
}

// Simpan qr code ke file
func (g *Generator) SaveQRCodeToFile(data string, filename string, size int) error {
	if size == 0 {
		size = 256
	}

	err := qrcode.WriteFile(data, qrcode.Medium, size, filename)
	if err != nil {
		return fmt.Errorf("failed to write QR code file: %w", err)
	}

	return nil
}
