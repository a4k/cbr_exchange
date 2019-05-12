/*
 * Таблица Org_phones
*/
package sqlstore

import (
	"net/http"

	"../../model"
	"../../store"
	"strings"
)

type SqlOrgPhonesStore struct {
	SqlStore
}

func NewSqlOrgPhonesStore(sqlStore SqlStore) store.OrgPhonesStore {
	s := &SqlOrgPhonesStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Org_phones{}, "org_phones").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(36)
		table.ColMap("Number").SetMaxSize(256)
		table.ColMap("Org_id").SetMaxSize(128)
	}

	return s
}

func (jss SqlOrgPhonesStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_phones_type", "org_phones", "Id")
}

// Сохранение
func (jss SqlOrgPhonesStore) Save(phones *model.Org_phones) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		phones.PreSave()
		if err := jss.GetMaster().Insert(phones); err != nil {
			result.Err = model.NewAppError("SqlOrgPhonesStore.Save", "store.sql_phones.save.app_error", nil, "id="+phones.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = phones
		}
	})
}

// Поиск в базе данных
func (us SqlOrgPhonesStore) GetByOrgId(org_id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var fphones = &model.Org_phones{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_phones WHERE Org_id = :Org_id ORDER BY Org_id DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Org_id": org_id})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgId", store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			fphones.Id = Id
			result.Data = fphones
		}
	})
}

// Обновление
func (us SqlOrgPhonesStore) Update(rphones *model.Org_phones) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := us.GetMaster().Update(rphones); err != nil {
			result.Err = model.NewAppError("SqlOrgPhonesStore.Update", "store.sql_OrgPhones.update.app_error", nil, "", http.StatusInternalServerError)
		}
		if sqlResult, err := us.GetReplica().Exec("UPDATE org_phones SET 	Org_id = :Org_id, Number = :Number WHERE Id = :Id",
			map[string]interface{}{
				"Org_id":        rphones.Org_id,
				"Number":    	 rphones.Number,
				"Id":            rphones.Id,
			}); err != nil {
			result.Err = model.NewAppError("SqlOrgPhonesStore.Update", "store.sql_OrgPhones.update.app_error", nil, "OrgPhones, "+err.Error(), http.StatusInternalServerError)
		} else {
			rows, err := sqlResult.RowsAffected()

			if err != nil {
				result.Err = model.NewAppError("SqlOrgPhonesStore.Update", "store.sql_OrgPhones.update.app_error", nil, "id="+rphones.Id+", "+err.Error(), http.StatusInternalServerError)
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

func (s SqlOrgPhonesStore) Clear() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				org_phones`); err != nil {
			result.Err = model.NewAppError("SqlDomainStore.Clear", "store.sql_job.delete.app_error", nil, "" + err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = true
		}
	})
}