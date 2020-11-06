package sftp

import (
	"errors"
	"os"

	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/pkg/sftp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type SFTPUpload struct {
	conn *ssh.Client
}

func (s *SFTPUpload) Init() error {
	var err error
	addr := createSFTPAddress()
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO remove this and check the key
		User:            viper.GetString("SFTP_USERNAME"),
		Auth: []ssh.AuthMethod{
			ssh.Password(viper.GetString("SFTP_PASSWORD")),
		},
	}
	s.conn, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		logger.Error("unable to initialise the SFTP connection",
			zap.Error(err))
		return err
	}
	logger.Info("connected to SFTP server")
	return nil
}

func createSFTPAddress() string {
	host := viper.GetString("SFTP_HOST")
	port := viper.GetString("SFTP_PORT")
	logger.Info("initialising sftp connection",
		zap.String("host", host),
		zap.String("port", port))

	addr := host + ":" + port
	return addr
}

func (s *SFTPUpload) Close() error {
	logger.Info("closing connection to SFTP")
	if s.conn == nil {
		return errors.New("please initialise connection")
	}
	err := s.conn.Close()
	logger.Info("sftp connection closed")
	return err
}

func (s *SFTPUpload) UploadFile(filename string, contents []byte) error {
	logger.Info("uploading to SFTP server",
		zap.String("filename", filename))
	if s.conn == nil {
		return errors.New("please initialise connection")
	}
	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(s.conn)
	if err != nil {
		logger.Error("unable to create new SFTP connection",
			zap.Error(err))
		return err
	}
	defer client.Close()

	workingDir, err := client.Getwd()
	if err != nil {
		logger.Error("unable to get current working directory",
			zap.Error(err))
	}

	logger.Info("working dir",
		zap.String("workingDir", workingDir))
	path := filepath(workingDir, filename)

	logger.Info("creating file")

	// check the file is there
	fi, err := client.Lstat(path)
	if err != nil {
		logger.Info("file does not exist, creating",
			zap.String("filepath", path))
	} else if fi.Size() != 0 {
		logger.Info("file already exists and is not empty",
			zap.String("filepath", path))
		return nil
	}

	f, err := client.Create(path)
	if err != nil {
		logger.Info("unable to create file",
			zap.String("filepath", path))
		return err
	}
	logger.Info("writing contents")
	if _, err := f.Write(contents); err != nil {
		logger.Error("unable to write file contents",
			zap.String("filepath", path),
			zap.Error(err))
		return err
	}
	f.Close()

	// check it's there
	logger.Info("confirming file exists")
	fi, err = client.Lstat(path)
	if err != nil {
		logger.Warn("unable to confirm file exists",
			zap.String("filepath", path),
			zap.Error(err))
	}
	logger.Info("upload complete",
		zap.String("file", fi.Name()))
	return nil
}

func filepath(workingDir string, filename string) string {
	dir := viper.GetString("SFTP_DIRECTORY")
	ps := string(os.PathSeparator)
	path := workingDir + ps + dir + ps + filename
	return path
}
