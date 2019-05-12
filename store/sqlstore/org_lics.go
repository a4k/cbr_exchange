/*
 * Таблица Org_lics
*/
package sqlstore

import (
	"net/http"

	"../../model"
	"../../store"
	"strings"
)

type SqlOrgLicsStore struct {
	SqlStore
}

func NewSqlOrgLicsStore(sqlStore SqlStore) store.OrgLicsStore {
	s := &SqlOrgLicsStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Org_lics{}, "org_lics").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(36)
		table.ColMap("Type").SetMaxSize(128)
		table.ColMap("Name").SetMaxSize(128)
		table.ColMap("Allow").SetMaxSize(128)
		table.ColMap("Org_id").SetMaxSize(128)
	}

	return s
}

func (jss SqlOrgLicsStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_lics_type", "org_lics", "Id")
}

// Сохранение
func (jss SqlOrgLicsStore) Save(lics *model.Org_lics) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		lics.PreSave()
		if err := jss.GetMaster().Insert(lics); err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.Save", "store.sql_lics.save.app_error", nil, "id="+lics.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = lics
		}
	})
}

// Поиск в базе данных
func (us SqlOrgLicsStore) GetByOrgId(org_id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var flics = &model.Org_lics{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_lics WHERE Org_id = :Org_id ORDER BY Org_id DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Org_id": org_id})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgId", store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			flics.Id = Id
			result.Data = flics
		}
	})
}

// Поиск в базе данных
func (us SqlOrgLicsStore) GetByOrgIdAndType(org_id string, rtype string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var flics = &model.Org_lics{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_lics WHERE Org_id = :Org_id and Type = :Type ORDER BY Type DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{
			"Org_id": org_id,
			"Type": rtype,
		})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgId", store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			flics.Id = Id
			result.Data = flics
		}
	})
}

// Обновление
func (us SqlOrgLicsStore) Update(rlics *model.Org_lics) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if sqlResult, err := us.GetReplica().Exec("UPDATE org_lics SET Org_id = :Org_id, Type = :Type, Name = :Name, Allow = :Allow WHERE Id = :Id",
			map[string]interface{}{
				"Org_id":       rlics.Org_id,
				"Type":    		rlics.Type,
				"Name":    		rlics.Name,
				"Allow":    	rlics.Allow,
				"Id":           rlics.Id,
			}); err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.Update", "store.sql_OrgLics.update.app_error", nil, "OrgLics, "+err.Error(), http.StatusInternalServerError)
		} else {
			rows, err := sqlResult.RowsAffected()

			if err != nil {
				result.Err = model.NewAppError("SqlOrgLicsStore.Update", "store.sql_OrgLics.update.app_error", nil, "id="+rlics.Id+", "+err.Error(), http.StatusInternalServerError)
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

func (s SqlOrgLicsStore) Clear() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				org_lics`); err != nil {
			result.Err = model.NewAppError("SqlDomainStore.Clear", "store.sql_job.delete.app_error", nil, "" + err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = true
		}
	})
}