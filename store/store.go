package store

import (
	"time"

	l4g "../utils/log4go"

	"../model"
)

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

func Do(f func(result *StoreResult)) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		result := StoreResult{}
		f(&result)
		storeChannel <- result
		close(storeChannel)
	}()
	return storeChannel
}

func Must(sc StoreChannel) interface{} {
	r := <-sc
	if r.Err != nil {
		l4g.Close()
		time.Sleep(time.Second)
		panic(r.Err)
	}

	return r.Data
}

type Store interface {
	Job() JobStore
	Popmail() PopmailStore

	Close()
	DropAllTables()
	TotalMasterDbConnections() int
	TotalReadDbConnections() int
	TotalSearchDbConnections() int
	FinHelp() FinHelpStore
	Insurance() InsuranceStore

	OrgEmails() OrgEmailsStore
	OrgFaxs() OrgFaxsStore
	OrgLics() OrgLicsStore
	OrgPhones() OrgPhonesStore
	OrgWebsites() OrgWebsitesStore
}


type PopmailStore interface {
	Save(popmail *model.Popmail) StoreChannel
	SaveOrUpdate(popmail *model.Popmail) StoreChannel
	Update(popmail *model.Popmail) StoreChannel
	Get(id string) StoreChannel
	GetByUidl(uidl string) StoreChannel
	GetList(listId []string) StoreChannel
	GetAllPopmail() StoreChannel
	GetChunk(offset int, limit int) StoreChannel
	GetPopmailCount() StoreChannel
	Delete(id string) StoreChannel
	Clear() StoreChannel
	DeleteList(listId []string) StoreChannel
}

type JobStore interface {
	Save(job *model.Job) StoreChannel
	UpdateOptimistically(job *model.Job, currentStatus string) StoreChannel
	UpdateStatus(id string, status string) StoreChannel
	UpdateStatusOptimistically(id string, currentStatus string, newStatus string) StoreChannel
	Get(id string) StoreChannel
	GetAllPage(offset int, limit int) StoreChannel
	GetAllByType(jobType string) StoreChannel
	GetAllByTypePage(jobType string, offset int, limit int) StoreChannel
	GetAllByStatus(status string) StoreChannel
	GetNewestJobByStatusAndType(status string, jobType string) StoreChannel
	GetCountByStatusAndType(status string, jobType string) StoreChannel
	Delete(id string) StoreChannel
}

type FinHelpStore interface {
	Save(fin *model.FinHelp) StoreChannel
	UpdateMicro(fin *model.FinHelp) StoreChannel
	UpdateBank(fin *model.FinHelp) StoreChannel
	GetByInn(Inn string) StoreChannel
	GetByBvkey(Bvkey string) StoreChannel
	GetByRegn(Regn string) StoreChannel
	GetByNum(Num string) StoreChannel
}

type InsuranceStore interface {
	Save(insurance *model.Org_insurance_lics) StoreChannel
	Update(finsurance *model.Org_insurance_lics) StoreChannel
	Clear() StoreChannel
	GetByOrgId(org_id string) StoreChannel
	GetByOrgIdAndType(org_id string, lic_type_code string) StoreChannel
}

type OrgEmailsStore interface {
	Save(email *model.Org_emails) StoreChannel
	Update(email *model.Org_emails) StoreChannel
	Clear() StoreChannel
	GetByOrgId(org_id string) StoreChannel
}

type OrgFaxsStore interface {
	Save(fax *model.Org_faxs) StoreChannel
	Update(fax *model.Org_faxs) StoreChannel
	Clear() StoreChannel
	GetByOrgId(org_id string) StoreChannel
}

type OrgLicsStore interface {
	Save(lic *model.Org_lics) StoreChannel
	Update(lic *model.Org_lics) StoreChannel
	Clear() StoreChannel
	GetByOrgId(org_id string) StoreChannel
	GetByOrgIdAndType(org_id string, rtype string) StoreChannel
}

type OrgPhonesStore interface {
	Save(phone *model.Org_phones) StoreChannel
	Update(phone *model.Org_phones) StoreChannel
	Clear() StoreChannel
	GetByOrgId(org_id string) StoreChannel
}

type OrgWebsitesStore interface {
	Save(website *model.Org_websites) StoreChannel
	Update(website *model.Org_websites) StoreChannel
	Clear() StoreChannel
	GetByOrgId(org_id string) StoreChannel
}
