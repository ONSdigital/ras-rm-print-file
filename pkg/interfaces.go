package pkg

type Store interface {
	Init() error
	Close() error
	Add(printFilename string, dataFilename string) (*PrintFileRequest, error)
	Update(pfr *PrintFileRequest) error
	FindIncomplete() ([]*PrintFileRequest, error)
}

type Download interface {
	Init() error
	Close() error
	DownloadFile(filename string) (*PrintFile, error)
}

type Upload interface {
	Init() error
	Close() error
	UploadFile(filename string, contents []byte) error
}

type Retry interface {
	Start() error
}

type Printer interface {
	Process(filename string, dataFilename string) error
	ReProcess(printFileRequest *PrintFileRequest) error
}
