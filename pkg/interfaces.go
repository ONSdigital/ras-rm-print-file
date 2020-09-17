package pkg

type Store interface {
	Init() error
	Add(filename string, p *PrintFile) (*PrintFileRequest, error)
	Update(pfr *PrintFileRequest) error
	FindIncomplete() ([]*PrintFileRequest, error)
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
	Process(filename string, printFile *PrintFile) error
	ReProcess(printFileRequest *PrintFileRequest) error
}
