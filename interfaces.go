package main

type Store interface {
	Init() error
	Add(filename string, p *PrintFile) (*PrintFileRequest, error)
	Update(pfr *PrintFileRequest) error
}

type Upload interface {
	Init() error
	Close()
	UploadFile(filename string, contents []byte) error
}
