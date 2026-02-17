package ftyp

import (
	"archive/zip"
	"errors"
	"io"
	"strings"
)

type Ftyp int

const (
	ZIP Ftyp = iota
	ODT
	ODS
	DOCX
	XLSX
	PPTX
)

// Determines whether a zip archive is odt, ods, docx, xlsx or pptx by looking
// for "mimetype" and "[Content_Types].xml" in the archive.
//
// Defaults to ZIP if an error occurs in the process.
func WhatZip(s string) Ftyp {
	r, err := zip.OpenReader(s)
	if err != nil {
		return ZIP
	}
	defer r.Close()

	for _, f := range r.File {
		switch f.Name {
		case "mimetype":
			rc, err := f.Open()
			if err != nil {
				return ZIP
			}
			defer rc.Close()

			buf := make([]byte, 256)
			n, err := rc.Read(buf)
			if err != nil && !errors.Is(err, io.EOF) {
				return ZIP
			}

			switch string(buf[:n]) {
			case "application/vnd.oasis.opendocument.text":
				return ODT
			case "application/vnd.oasis.opendocument.spreadsheet":
				return ODS
			}

			return ZIP
		case "[Content_Types].xml":
			for _, f = range r.File {
				name := f.Name
				switch {
				case strings.HasPrefix(name, "word/"):
					return DOCX
				case strings.HasPrefix(name, "xl/"):
					return XLSX
				case strings.HasPrefix(name, "ppt/"):
					return PPTX
				default:
					continue
				}
			}

			return ZIP
		default:
			continue
		}
	}

	return ZIP
}
