package connector

import (
	_ "github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

func TestOneDrive_Folder(t *testing.T) {

	c := OneDrive{
		param: &OneDriveParameters{
			Folder:    "",
			Recursive: false,
			Token:     oauth2.Token{},
		},
	}
	// root without recursive
	c.param.Recursive = false
	c.param.Folder = ""
	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), false)

	assert.Equal(t, c.isFilesAnalysing(""), true)
	assert.Equal(t, c.isFilesAnalysing("folder"), false)

	// root with recursive
	c.param.Recursive = true
	c.param.Folder = ""
	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFilesAnalysing(""), true)
	assert.Equal(t, c.isFilesAnalysing("folder"), true)
	assert.Equal(t, c.isFilesAnalysing("folder/chapter1"), true)

	// given folder name with recursive
	c.param.Recursive = true
	c.param.Folder = "folder"

	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFolderAnalysing("folder2/chapter1"), false)
	assert.Equal(t, c.isFolderAnalysing("docs/folder/chapter1"), false)

	assert.Equal(t, c.isFilesAnalysing(""), false)
	assert.Equal(t, c.isFilesAnalysing("folder"), true)
	assert.Equal(t, c.isFilesAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFilesAnalysing("folder2/chapter1"), false)
	assert.Equal(t, c.isFilesAnalysing("docs/folder/chapter1"), false)

	// given folder name without recursive
	c.param.Folder = "folder"
	c.param.Recursive = false
	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), false)
	assert.Equal(t, c.isFolderAnalysing("folder2/chapter1"), false)

	assert.Equal(t, c.isFilesAnalysing(""), false)
	assert.Equal(t, c.isFilesAnalysing("folder"), true)
	assert.Equal(t, c.isFilesAnalysing("folder/chapter1"), false)
	assert.Equal(t, c.isFilesAnalysing("folder2/chapter1"), false)

	// given subfolder name without recursive
	c.param.Folder = "folder/chapter1"
	c.param.Recursive = false
	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1/page1"), false)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter2"), false)
	assert.Equal(t, c.isFolderAnalysing("folder2/chapter1"), false)

	assert.Equal(t, c.isFilesAnalysing(""), false)
	assert.Equal(t, c.isFilesAnalysing("folder"), false)
	assert.Equal(t, c.isFilesAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFilesAnalysing("folder/chapter1/page1"), false)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter2"), false)
	assert.Equal(t, c.isFilesAnalysing("folder2/chapter1"), false)

	// given subfolder name with recursive
	c.param.Folder = "folder/chapter1"
	c.param.Recursive = true
	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1/page1"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter2"), false)
	assert.Equal(t, c.isFolderAnalysing("folder2/chapter1"), false)

	assert.Equal(t, c.isFilesAnalysing(""), false)
	assert.Equal(t, c.isFilesAnalysing("folder"), false)
	assert.Equal(t, c.isFilesAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFilesAnalysing("folder/chapter1/page1"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter2"), false)
	assert.Equal(t, c.isFilesAnalysing("folder2/chapter1"), false)
}
