/*
 * Таблица Org_faxs
*/
package sqlstore

import (
	"net/http"

	"../../model"
	"../../store"
	"strings"
)

type SqlOrgFaxsStore struct {
	SqlStore
}

func NewSqlOrgFaxsStore(sqlStore SqlStore) store.OrgFaxsStore {
	s := &SqlOrgFaxsStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Org_faxs{}, "org_faxs").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(36)
		table.ColMap("Number").SetMaxSize(128)
		table.ColMap("Org_id").SetMaxSize(128)
	}

	return s
}

func (jss SqlOrgFaxsStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_faxs_type", "org_faxs", "Id")
}

// Сохранение
func (jss SqlOrgFaxsStore) Save(faxs *model.Org_faxs) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		faxs.PreSave()
		if err := jss.GetMaster().Insert(faxs); err != nil {
			result.Err = model.NewAppError("SqlOrgFaxsStore.Save", "store.sql_faxs.save.app_error", nil, "id="+faxs.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = faxs
		}
	})
}

// Поиск в базе данных
func (us SqlOrgFaxsStore) GetByOrgId(org_id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var ffaxs = &model.Org_faxs{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_faxs WHERE Org_id = :Org_id ORDER BY Org_id DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Org_id": org_id})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgId", store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			ffaxs.Id = Id
			result.Data = ffaxs
		}
	})
}

// Обновление
func (us SqlOrgFaxsStore) Update(rfaxs *model.Org_faxs) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if sqlResult, err := us.GetReplica().Exec("UPDATE org_faxs SET 	Org_id = :Org_id, Number = :Number WHERE Id = :Id",
			map[string]interface{}{
				"Org_id":        rfaxs.Org_id,
				"Number":    	 rfaxs.Number,
				"Id":            rfaxs.Id,
			}); err != nil {
			result.Err = model.NewAppError("SqlOrgFaxsStore.Update", "store.sql_OrgFaxs.update.app_error", nil, "OrgFaxs, "+err.Error(), http.StatusInternalServerError)
		} else {
			rows, err := sqlResult.RowsAffected()

			if err != nil {
				result.Err = model.NewAppError("SqlOrgFaxsStore.Update", "store.sql_OrgFaxs.update.app_error", nil, "id="+rfaxs.Id+", "+err.Error(), http.StatusInternalServerError)
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

func (s SqlOrgFaxsStore) Clear() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				org_faxs`); err != nil {
			result.Err = model.NewAppError("SqlDomainStore.Clear", "store.sql_job.delete.app_error", nil, "" + err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = true
		}
	})
}