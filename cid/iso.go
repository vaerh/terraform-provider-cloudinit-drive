package cid

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	isoCommands = map[string]func(isoLabel, src, dst string) []string{
		"genisoimage": func(isoLabel, src, dst string) []string {
			return []string{"-joliet", "-rock", "-volid", isoLabel, "-o", dst, src}
		},
		"mkisofs": func(isoLabel, src, dst string) []string {
			return []string{"-joliet", "-rock", "-volid", isoLabel, "-o", dst, src}
		},
		"hdiutil": func(isoLabel, src, dst string) []string {
			return []string{"makehybrid", "-hfs", "-joliet", "-iso", "-default-volume-name", isoLabel, "-o", dst, src}
		},
		"oscdimg": func(isoLabel, src, dst string) []string {
			return []string{"-j1", "-o", "-m", "-l" + isoLabel, src, dst}
		},
		"xorriso": func(isoLabel, src, dst string) []string {
			return []string{"-as", "genisoimage", "-rock", "-joliet", "-volid", isoLabel, "-output", dst, src}
		},
	}
)

type Iso struct {
	isoMaker string
	files    map[string]io.ReadCloser
	label    string
}

func NewISOWriter(isoMakerCommand string) (*Iso, error) {
	return &Iso{
		isoMaker: isoMakerCommand,
		files:    make(map[string]io.ReadCloser),
	}, nil
}

// data is io.Reader or io.ReadCloser
func (i *Iso) AddFile(data any, filePath string) {
	switch data := data.(type) {
	case io.ReadCloser:
		i.files[filePath] = data
	case io.Reader:
		i.files[filePath] = io.NopCloser(data)
	default:
		panic("Method not applicable for types other than io.Reader or io.ReadCloser")
	}
}

func (i *Iso) SetLabel(label string) {
	i.label = label
}

func (i *Iso) WriteTo(w io.Writer) (int64, error) {
	if len(i.files) == 0 {
		return -1, fmt.Errorf("no files specified, ISO file will not be made")
	}

	if i.label == "" {
		i.label = "config-2"
	}

	isoFilename, err := os.CreateTemp("", "cloudinit*.iso")
	if err != nil {
		return -1, fmt.Errorf("error creating temporary file for ISO: %s", err)
	}
	defer os.Remove(isoFilename.Name())

	rootFolder, err := os.MkdirTemp("", "cid*")
	if err != nil {
		return -1, fmt.Errorf("failed to create temp dir %v", err)
	}
	defer os.RemoveAll(rootFolder)

	// Creating an ISO structure.
	for filePath, reader := range i.files {
		err = os.MkdirAll(filepath.Dir(filepath.Join(rootFolder, filePath)), os.ModePerm)
		if err != nil {
			return -1, fmt.Errorf("error creating directory: %s", filePath)
		}

		dest, err := os.Create(filepath.Join(rootFolder, filePath))
		if err != nil {
			return -1, fmt.Errorf("error creating file for copy %s to ISO root", filePath)
		}
		defer dest.Close()

		_, err = io.Copy(dest, reader)
		if err != nil {
			return -1, fmt.Errorf("error copying %s to ISO root", filePath)
		}

		reader.Close()
	}

	// The process of creating an ISO file.
	cmd := exec.Command(i.isoMaker, isoCommands[i.isoMaker](i.label, rootFolder, isoFilename.Name())...)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return -1, fmt.Errorf("failed to exec %s, %s %s", i.isoMaker, data, err)
	}

	if !cmd.ProcessState.Success() {
		return -1, fmt.Errorf("failed to create ISO: %s", data)
	}

	if info, _ := isoFilename.Stat(); info.Size() == 0 {
		return -1, fmt.Errorf("empty iso created")
	}

	f, err := os.Open(isoFilename.Name())
	if err != nil {
		return -1, fmt.Errorf("failed to open new ISO: %s", err)
	}
	defer f.Close()

	return io.Copy(w, f)
}
