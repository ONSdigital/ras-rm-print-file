package database

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"fmt"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type DataStore struct {
	ctx    context.Context
	client *datastore.Client
}

func (s *DataStore) Init() error {
	var err error
	s.ctx = context.Background()
	log.Info("initialising google datastore connection")
	s.client, err = datastore.NewClient(s.ctx, viper.GetString("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.WithError(err).Error("error creating gcp client")
		return fmt.Errorf("datastore.NewClient: %v", err)
	}
	log.Info("connected to google datastore")
	return nil

}

func (s *DataStore) Add(filename string, p *pkg.PrintFile) (*pkg.PrintFileRequest, error) {
	if s.client == nil {
		return nil, errors.New("please initialise the connection")
	}
	// DataStore the initial response and the name of the file
	// we're meant to create
	key := datastore.NameKey("PrintFileRequest", filename, nil)
	pfr := createPrintFileRequest(filename, p)

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
		log.WithError(err).Error("unable to DataStore entry")
		return nil, fmt.Errorf("unable to to DataStore entry: %v", err)
	}
	return pfr, nil
}

func createPrintFileRequest(filename string, p *pkg.PrintFile) *pkg.PrintFileRequest {
	pfr := &pkg.PrintFileRequest{
		PrintFile: p,
		Filename:  filename,
		Created:   time.Now(),
		Updated:   time.Now(),
		Attempts:  1,
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
	key := datastore.NameKey("PrintFileRequest", pfr.Filename, nil)
	tx, err := s.client.NewTransaction(s.ctx)
	if err != nil {
		log.WithError(err).Error("unable to start transaction")
		return err
	}
	if _, err := tx.Put(key, pfr); err != nil {
		log.WithError(err).Error("unable to Update entity")
		return err
	}
	if _, err := tx.Commit(); err != nil {
		log.WithError(err).Error("unable to commit entity to database")
		return err
	}
	return nil
}

func (s *DataStore) FindIncomplete() ([]*pkg.PrintFileRequest, error) {
	if s.client == nil {
		return nil, errors.New("please initialise the connection")
	}
	log.Debug("about to execute query on datastore")
	var pfr []*pkg.PrintFileRequest

	query := datastore.NewQuery("PrintFileRequest").Filter("Status.Completed =", false)
	keys, err := s.client.GetAll(s.ctx, query, &pfr)
	incomplete := len(keys)
	log.WithField("incomplete", incomplete).Info("found incomplete requests")
	for _, v := range keys {
		log.WithField("id", v).Debug("request found to be incomplete")
	}
	if err != nil {
		log.WithError(err).Error("unable to query datastore")
		return nil, err
	}
	return pfr, nil
}
