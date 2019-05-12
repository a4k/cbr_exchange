/*
 * Таблица Org_emails
*/
package sqlstore

import (
	"net/http"

	"../../model"
	"../../store"
	"strings"
)

type SqlOrgEmailsStore struct {
	SqlStore
}

func NewSqlOrgEmailsStore(sqlStore SqlStore) store.OrgEmailsStore {
	s := &SqlOrgEmailsStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Org_emails{}, "org_emails").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(36)
		table.ColMap("Email").SetMaxSize(128)
		table.ColMap("Org_id").SetMaxSize(128)
	}

	return s
}

func (jss SqlOrgEmailsStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_emails_type", "org_emails", "Id")
}

// Сохранение
func (jss SqlOrgEmailsStore) Save(emails *model.Org_emails) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		emails.PreSave()
		if err := jss.GetMaster().Insert(emails); err != nil {
			result.Err = model.NewAppError("SqlOrgEmailsStore.Save", "store.sql_emails.save.app_error", nil, "id="+emails.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = emails
		}
	})
}

// Поиск в базе данных
func (us SqlOrgEmailsStore) GetByOrgId(org_id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var femails = &model.Org_emails{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_emails WHERE Org_id = :Org_id ORDER BY Org_id DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Org_id": org_id})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgId", store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			femails.Id = Id
			result.Data = femails
		}
	})
}

// Обновление
func (us SqlOrgEmailsStore) Update(remails *model.Org_emails) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if sqlResult, err := us.GetReplica().Exec("UPDATE org_emails SET 	Org_id = :Org_id, Email = :Email WHERE Id = :Id",
			map[string]interface{}{
				"Org_id":        remails.Org_id,
				"Email":    remails.Email,
				"Id":            remails.Id,
			}); err != nil {
			result.Err = model.NewAppError("SqlOrgEmailsStore.Update", "store.sql_OrgEmails.update.app_error", nil, "OrgEmails, "+err.Error(), http.StatusInternalServerError)
		} else {
			rows, err := sqlResult.RowsAffected()

			if err != nil {
				result.Err = model.NewAppError("SqlOrgEmailsStore.Update", "store.sql_OrgEmails.update.app_error", nil, "id="+remails.Id+", "+err.Error(), http.StatusInternalServerError)
			} else {
				if rows == 1 {
					result.Data = true
				} else {
					result.Data = false
				}
			}
		}

	})
}

func (s SqlOrgEmailsStore) Clear() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				org_emails`); err != nil {
			result.Err = model.NewAppError("SqlDomainStore.Clear", "store.sql_job.delete.app_error", nil, "" + err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = true
		}
	})
}