package connector

import "cognix.ch/api/v2/core/proto"

var supportedMimeTypes = map[string]proto.FileType{
	mineURL: proto.FileType_URL,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": proto.FileType_XLSX,
	"application/vnd.ms-excel": proto.FileType_XLS,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": proto.FileType_DOCX,
	"application/msword": proto.FileType_DOC,
	"application/pdf":    proto.FileType_PDF,
	"text/plain":         proto.FileType_TXT,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": proto.FileType_PPTX,
	"application/vnd.ms-powerpoint":                                             proto.FileType_PPT,
	"application/vnd.ms-xpsdocument":                                            proto.FileType_XPS,
	"application/oxps":                                                          proto.FileType_XPS,
	"application/epub+zip":                                                      proto.FileType_EPUB,
	"application/hwp+zip":                                                       proto.FileType_HWPX,
	"text/markdown":                                                             proto.FileType_MD,
	"application/x-mobipocket-ebook":                                            proto.FileType_MOBI,
	"application/fb2":                                                           proto.FileType_FB2,
}

var supportedExtensions = map[string]string{
	"MD":   "text/markdown",
	"HWPX": "application/hwp+zip",
	"MOBI": "application/x-mobipocket-ebook",
	"FB2":  "application/fb2",
}
