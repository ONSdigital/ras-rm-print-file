package main

import (
	"github.com/pkg/sftp"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	log "github.com/sirupsen/logrus"
)

type SFTPUpload struct {
	conn *ssh.Client
}

func (s *SFTPUpload) Init() error {
	var err error
	host := viper.GetString("SFTP_HOST")
	port := viper.GetString("SFTP_PORT")
	addr := host + ":" + port
	config := &ssh.ClientConfig{
		User: viper.GetString("SFTP_USERNAME"),
		Auth: []ssh.AuthMethod{
			ssh.Password(viper.GetString("SFTP_PASSWORD")),
		},
	}
	s.conn, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		log.WithError(err).Error("unable to initialise the SFTP connection")
		return err
	}
	return nil
}

func (s *SFTPUpload) Close() {
	s.conn.Close()
}

func (s *SFTPUpload) UploadFile(name string, contents []byte) error {
	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(s.conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	f, err := client.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(contents); err != nil {
		log.Fatal(err)
	}
	f.Close()

	// check it's there
	fi, err := client.Lstat(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)


	return nil
}
