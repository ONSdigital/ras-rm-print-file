package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DataStore struct {
	ctx    context.Context
	client *datastore.Client
}

func (s *DataStore) Init() error {
	var err error
	s.ctx = context.Background()
	logger.Info("initialising google datastore connection")
	s.client, err = datastore.NewClient(s.ctx, viper.GetString("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		logger.Error("error creating gcp client",
			zap.Error(err))
		return fmt.Errorf("datastore.NewClient: %v", err)
	}
	logger.Info("connected to google datastore")
	return nil

}

func (s *DataStore) Close() error {
	if s.client == nil {
		return errors.New("please initialise the connection")
	}
	logger.Info("closing connection to datastore")
	err := s.client.Close()
	logger.Info("GCS connection closed")
	return err
}


func (s *DataStore) Add(printFilename string, dataFilename string) (*pkg.PrintFileRequest, error) {
	if s.client == nil {
		return nil, errors.New("please initialise the connection")
	}
	// DataStore the initial response and the name of the file
	// we're meant to create
	key := datastore.NameKey("PrintFileRequest", printFilename, nil)
	pfr := createPrintFileRequest(printFilename, dataFilename)

	_, err := s.client.RunInTransaction(s.ctx, func(tx *datastore.Transaction) error {
		// We first check that there is no entity stored with the given key.
		var empty pkg.PrintFileRequest
		if err := tx.Get(key, &empty); err != datastore.ErrNoSuchEntity {
			return err
		}
		// If there was no matching entity, DataStore it now.
		_, err := tx.Put(key, pfr)
		return err
	})

	if err != nil {
		logger.Error("unable to DataStore entry",
			zap.Error(err))
		return nil, fmt.Errorf("unable to to DataStore entry: %v", err)
	}
	return pfr, nil
}

func createPrintFileRequest(printFilename string, dataFilename string) *pkg.PrintFileRequest {
	pfr := &pkg.PrintFileRequest{
		DataFilename:  dataFilename,
		PrintFilename: printFilename,
		Created:       time.Now(),
		Updated:       time.Now(),
		Attempts:      1,
		Status: pkg.Status{
			Templated:    false,
			UploadedGCS:  false,
			UploadedSFTP: false,
		},
	}
	return pfr
}

func (s *DataStore) Update(pfr *pkg.PrintFileRequest) error {
	if s.client == nil {
		return errors.New("please initialise the connection")
	}
	pfr.Updated = time.Now()
	key := datastore.NameKey("PrintFileRequest", pfr.PrintFilename, nil)
	tx, err := s.client.NewTransaction(s.ctx)
	if err != nil {
		logger.Error("unable to start transaction",
			zap.Error(err))
		return err
	}
	if _, err := tx.Put(key, pfr); err != nil {
		logger.Error("unable to Update entity",
			zap.Error(err))
		return err
	}
	if _, err := tx.Commit(); err != nil {
		logger.Error("unable to commit entity to database",
			zap.Error(err))
		return err
	}
	return nil
}

func (s *DataStore) Delete(pfr *pkg.PrintFileRequest) error {
	if s.client == nil {
		return errors.New("please initialise the connection")
	}
	pfr.Updated = time.Now()
	key := datastore.NameKey("PrintFileRequest", pfr.PrintFilename, nil)
	tx, err := s.client.NewTransaction(s.ctx)
	if err != nil {
		logger.Error("unable to start transaction",
			zap.Error(err))
		return err
	}
	if err := tx.Delete(key); err != nil {
		logger.Error("unable to delete entity",
			zap.Error(err))
		return err
	}
	if _, err := tx.Commit(); err != nil {
		logger.Error("unable to commit entity to database",
			zap.Error(err))
		return err
	}
	return nil
}

func (s *DataStore) FindIncomplete() ([]*pkg.PrintFileRequest, error) {
	logger.Debug("finding all incomplete print file requests")
	return s.find(false)
}

func (s *DataStore) find(complete bool) ([]*pkg.PrintFileRequest, error) {
	if s.client == nil {
		return nil, errors.New("please initialise the connection")
	}
	logger.Debug("about to execute query on datastore")
	var pfr []*pkg.PrintFileRequest

	query := datastore.NewQuery("PrintFileRequest").Filter("Status.Completed =", complete)
	keys, err := s.client.GetAll(s.ctx, query, &pfr)
	if err != nil {
		logger.Error("unable to query datastore",
			zap.Error(err))
		return nil, err
	}
	results := len(keys)
	logger.Info("found  requests", zap.Bool("complete", complete), zap.Int("results", results))
	for _, v := range keys {
		logger.Debug("print file request found",
			zap.Bool("complete", complete),
			zap.Any("id", v))
	}
	return pfr, nil
}

func (s *DataStore) FindComplete() ([]*pkg.PrintFileRequest, error) {
	logger.Debug("finding all complete print file requests")
	return s.find(true)
}

