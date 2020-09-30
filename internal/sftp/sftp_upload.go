package sftp

import (
	"errors"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"os"
)

type SFTPUpload struct {
	conn *ssh.Client
}

func (s *SFTPUpload) Init() error {
	var err error
	addr := createSFTPAddress()
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //TODO remove this and check the key
		User:            viper.GetString("SFTP_USERNAME"),
		Auth: []ssh.AuthMethod{
			ssh.Password(viper.GetString("SFTP_PASSWORD")),
		},
	}
	s.conn, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		log.WithError(err).Error("unable to initialise the SFTP connection")
		return err
	}
	log.Info("connected to SFTP server")
	return nil
}

func createSFTPAddress() string {
	host := viper.GetString("SFTP_HOST")
	port := viper.GetString("SFTP_PORT")
	log.WithField("host", host).WithField("port", port).Info("initialising sftp connection")

	addr := host + ":" + port
	return addr
}

func (s *SFTPUpload) Close() error {
	log.Info("closing connection to SFTP")
	if s.conn == nil {
		return errors.New("please initialise connection")
	}
	err := s.conn.Close()
	log.Info("sftp connection closed")
	return err
}

func (s *SFTPUpload) UploadFile(filename string, contents []byte) error {
	log.WithField("filename", filename).Info("uploading to SFTP server")
	if s.conn == nil {
		return errors.New("please initialise connection")
	}
	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(s.conn)
	if err != nil {
		log.Error(err)
		return err
	}
	defer client.Close()

	path := filepath(filename)

	log.Info("creating file")
	f, err := client.Create(path)
	if err != nil {
		log.WithError(err).Error("unable to create file")
		return err
	}
	log.Info("writing contents")
	if _, err := f.Write(contents); err != nil {
		log.WithError(err).Error("unable to write file contents")
		return err
	}
	f.Close()

	// check it's there
	log.Info("confirming file exists")
	fi, err := client.Lstat(filename)
	if err != nil {
		log.WithError(err).Error("unable to write file contents")
		return err
	}
	log.WithField("file", fi.Name()).Info("upload complete")
	return nil
}

func filepath(filename string) string {
	dir := viper.GetString("SFTP_DIRECTORY")
	path := dir + string(os.PathSeparator) + filename
	return path
}
