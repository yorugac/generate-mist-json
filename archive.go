package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type ArchiveReader interface {
	// tar | zip
	Type() string

	// Without extension
	BareName() string

	// Read filenames in root directory of the archive
	TopFiles() []string

	// Close all associated streams
	Close() error
}

type ZipArchiveReader struct {
	zipr     *zip.ReadCloser
	bareName string
}

func NewZipArchiveReader(name string) (ArchiveReader, error) {
	zipr, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}

	return &ZipArchiveReader{zipr, name[:len(name)-4]}, nil
}

func (a *ZipArchiveReader) Type() string { return "zip" }

func (a *ZipArchiveReader) BareName() string { return a.bareName }

func (a *ZipArchiveReader) TopFiles() []string {
	filenames := make([]string, 2)
	for _, file := range a.zipr.File {
		if !file.FileInfo().IsDir() {
			filenames = append(filenames, file.Name)
		}
	}
	return filenames
}

func (a *ZipArchiveReader) Close() error {
	a.zipr.Close()
	return nil
}

type TarballArchiveReader struct {
	gzipr    *gzip.Reader
	tarr     *tar.Reader
	bareName string
}

func NewTarballArchiveReader(file *os.File) (ArchiveReader, error) {
	gzipr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	tarr := tar.NewReader(gzipr)
	name := file.Name()
	return &TarballArchiveReader{gzipr, tarr, name[:len(name)-7]}, nil
}

func (a *TarballArchiveReader) Type() string { return "tar" }

func (a *TarballArchiveReader) BareName() string { return a.bareName }

func (a *TarballArchiveReader) TopFiles() (filenames []string) {
	for {
		header, err := a.tarr.Next()
		if err == io.EOF {
			return
		}
		if !header.FileInfo().IsDir() {
			filenames = append(filenames, header.Name)
		}
	}
}

func (a *TarballArchiveReader) Close() error {
	a.gzipr.Close()
	return nil
}

func OpenArchive(filename string, file *os.File) (ArchiveReader, error) {
	if strings.HasSuffix(filename, ".zip") {
		return NewZipArchiveReader(filename)
	}
	if strings.HasSuffix(filename, ".tar.gz") {
		return NewTarballArchiveReader(file)
	}
	return nil, fmt.Errorf("unsupported archive" + filename)
}

// InvestigateArchive looks into an existing archive.
func InvestigateArchive(filename, binaryName string) (binaryNames [2]string, archiveType, md5String string, err error) {
	log.Println("Investigating archive", filename)
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	archive, err := OpenArchive(filename, file)
	if archive == nil {
		return
	}

	archiveType = archive.Type()
	files := archive.TopFiles()
	for _, f := range files {
		delimIndex := strings.LastIndex(f, "/")
		if delimIndex > 0 && len(f) > delimIndex+4 && f[delimIndex+1:delimIndex+5] == binaryName {
			binaryNames[0] = f[delimIndex+1:]
			binaryNames[1] = f
			break
		}
	}

	archive.Close()
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return
	}
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return
	}
	hashInBytes := hash.Sum(nil)[:16]
	md5String = hex.EncodeToString(hashInBytes)
	return
}
