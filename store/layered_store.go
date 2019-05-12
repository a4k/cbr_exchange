package store

import (
	"context"
)

type LayeredStoreDatabaseLayer interface {
	LayeredStoreSupplier
	Store
}

type LayeredStore struct {
	TmpContext context.Context

	FinHelpStore    FinHelpStore
	InsuranceStore  InsuranceStore
	OrgEmailsStore  OrgEmailsStore
	OrgFaxsStore  	OrgFaxsStore
	OrgLicsStore  	OrgLicsStore
	OrgPhonesStore  OrgPhonesStore
	OrgWebsitesStore  OrgWebsitesStore

	DatabaseLayer   LayeredStoreDatabaseLayer
	RedisLayer      *RedisSupplier
	LayerChainHead  LayeredStoreSupplier
}

func NewLayeredStore(db LayeredStoreDatabaseLayer) Store {
	store := &LayeredStore{
		TmpContext:      context.TODO(),
		DatabaseLayer:   db,
	}


	return store
}

type QueryFunction func(LayeredStoreSupplier) *LayeredStoreSupplierResult

func (s *LayeredStore) RunQuery(queryFunction QueryFunction) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := queryFunction(s.LayerChainHead)
		storeChannel <- result.StoreResult
	}()

	return storeChannel
}

func (s *LayeredStore) Popmail() PopmailStore {
	return s.DatabaseLayer.Popmail()
}

func (s *LayeredStore) FinHelp() FinHelpStore {
	return s.DatabaseLayer.FinHelp()
}

func (s *LayeredStore) Insurance() InsuranceStore {
	return s.DatabaseLayer.Insurance()
}

func (s *LayeredStore) OrgEmails() OrgEmailsStore {
	return s.DatabaseLayer.OrgEmails()
}


func (s *LayeredStore) OrgFaxs() OrgFaxsStore {
	return s.DatabaseLayer.OrgFaxs()
}


func (s *LayeredStore) OrgLics() OrgLicsStore {
	return s.DatabaseLayer.OrgLics()
}


func (s *LayeredStore) OrgPhones() OrgPhonesStore {
	return s.DatabaseLayer.OrgPhones()
}

func (s *LayeredStore) OrgWebsites() OrgWebsitesStore {
	return s.DatabaseLayer.OrgWebsites()
}

func (s *LayeredStore) Job() JobStore {
	return s.DatabaseLayer.Job()
}

func (s *LayeredStore) Close() {
	s.DatabaseLayer.Close()
}

func (s *LayeredStore) DropAllTables() {
	s.DatabaseLayer.DropAllTables()
}

func (s *LayeredStore) TotalMasterDbConnections() int {
	return s.DatabaseLayer.TotalMasterDbConnections()
}

func (s *LayeredStore) TotalReadDbConnections() int {
	return s.DatabaseLayer.TotalReadDbConnections()
}

func (s *LayeredStore) TotalSearchDbConnections() int {
	return s.DatabaseLayer.TotalSearchDbConnections()
}

type LayeredReactionStore struct {
	*LayeredStore
}
