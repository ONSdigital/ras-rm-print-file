package main

type Store interface {
	Init() error
	store(filename string, p *PrintFile) (*PrintFileRequest, error)
	update(pfr *PrintFileRequest) error
}

type Upload interface {
	Init() error
	Close()
	UploadFile(filename string, contents []byte) error
}
