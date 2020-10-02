package payload

import (
	mocks "github.com/ONSdigital/ras-rm-print-file/mocks/pkg"
	"testing"
)

func TestSave(t *testing.T) {

	data := []byte("test")
	filename := "test.csv"

	gcsUpload := new(mocks.Upload)
	gcsUpload.On("Init").Return(nil)
	gcsUpload.On("UploadFile", "test.json", data).Return(nil)
	gcsUpload.On("Close").Return(nil)

	payload := Payload{
		gcsUpload,
	}
	payload.Save(filename, data)

	gcsUpload.AssertCalled(t, "UploadFile", "test.json", data)
	gcsUpload.AssertExpectations(t)
}
