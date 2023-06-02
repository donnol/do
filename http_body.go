package do

import (
	"io"
	"mime/multipart"
)

// MultipartBody new a multipart writer with body, mark with fieldname  and name, write data to it. Return form data content type.
func MultipartBody(body io.Writer, fieldname, filename string, data []byte) (fileContentType string, err error) {
	bodyWriter := multipart.NewWriter(body)
	defer bodyWriter.Close()

	// write field
	fieldWriter, err := bodyWriter.CreateFormField(fieldname)
	if err != nil {
		return
	}
	_, err = fieldWriter.Write([]byte(filename))
	if err != nil {
		return
	}

	// write file
	fileWriter, err := bodyWriter.CreateFormFile(fieldname, filename)
	if err != nil {
		return
	}
	_, err = fileWriter.Write(data)
	if err != nil {
		return
	}

	fileContentType = bodyWriter.FormDataContentType()

	return
}
