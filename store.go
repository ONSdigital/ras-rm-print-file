package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type PrintFileRequest struct {
	printFile *PrintFile
	filename string
	created time.Time
	Status     Status
}

type Status struct {
	TemplateComplete bool
	UploadedGCS      bool
	UploadedSFTP     bool
}

type Store interface {
	Init() error
	store(filename string, p *PrintFile) (*PrintFileRequest, error)
	update(pfr *PrintFileRequest) error
}

type store struct {
	ctx    context.Context
	client *datastore.Client
}

func (s *store) Init() error {
	var err error
	s.ctx = context.Background()
	log.Info("initialising google datastore connection")
	s.client, err = datastore.NewClient(s.ctx, viper.GetString("PROJECT_ID"))
	if err != nil {
		log.WithError(err).Error("error creating gcp client")
		return fmt.Errorf("datastore.NewClient: %v", err)
	}
	log.Info("connected to GCS")
	return nil
}

func (s *store) store(filename string, p *PrintFile) (*PrintFileRequest, error) {
	// store the initial response and the name of the file
	// we're meant to create
	key := datastore.NameKey("PrintFileRequest", filename, nil)
	pfr := &PrintFileRequest {
		p,
		filename,
		time.Now(),
		Status {
			TemplateComplete: false,
			UploadedGCS:      false,
			UploadedSFTP:     false,
		},
	}

	_, err := s.client.RunInTransaction(s.ctx, func(tx *datastore.Transaction) error {
		// We first check that there is no entity stored with the given key.
		var empty PrintFileRequest
		if err := tx.Get(key, &empty); err != datastore.ErrNoSuchEntity {
			return err
		}
		// If there was no matching entity, store it now.
		_, err := tx.Put(key, &pfr)
		return err
	})

	if err != nil {
		log.WithError(err).Error("unable to store entry")
		return nil, fmt.Errorf("unable to to store entry: %v", err)
	}
	return pfr, nil
}

func (s *store) update(pfr *PrintFileRequest) error {
	key := datastore.NameKey("PrintFileRequest", pfr.filename, nil)
	tx, err := s.client.NewTransaction(s.ctx)
	if err != nil {
		log.WithError(err).Error("unable to start transaction")
		return err
	}
	if _, err := tx.Put(key, pfr); err != nil {
		log.WithError(err).Error("unable to update entity")
		return err
	}
	if _, err := tx.Commit(); err != nil {
		log.WithError(err).Error("unable to commit entity to database")
		return err
	}
	return nil
}
