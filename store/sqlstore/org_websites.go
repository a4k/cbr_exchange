/*
 * Таблица Org_websites
*/
package sqlstore

import (
	"net/http"

	"../../model"
	"../../store"
	"strings"
)

type SqlOrgWebsitesStore struct {
	SqlStore
}

func NewSqlOrgWebsitesStore(sqlStore SqlStore) store.OrgWebsitesStore {
	s := &SqlOrgWebsitesStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Org_websites{}, "org_websites").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(36)
		table.ColMap("Website").SetMaxSize(128)
		table.ColMap("Org_id").SetMaxSize(128)
	}

	return s
}

func (jss SqlOrgWebsitesStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_websites_type", "org_websites", "Id")
}

// Сохранение
func (jss SqlOrgWebsitesStore) Save(websites *model.Org_websites) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		websites.PreSave()
		if err := jss.GetMaster().Insert(websites); err != nil {
			result.Err = model.NewAppError("SqlOrgWebsitesStore.Save",
				"store.sql_websites.save.app_error",
				nil, "id="+websites.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = websites
		}
	})
}

// Поиск в базе данных
func (us SqlOrgWebsitesStore) GetByOrgId(org_id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var fwebsites = &model.Org_websites{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_websites WHERE Org_id = :Org_id " +
			"ORDER BY Org_id DESC OFFSET 0 " +
			"ROWS FETCH NEXT (1) " +
			"ROWS ONLY", map[string]interface{}{"Org_id": org_id})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgId",
				store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", " +
				""+err.Error(), http.StatusInternalServerError)
		} else {
			fwebsites.Id = Id
			result.Data = fwebsites
		}
	})
}

// Обновление
func (us SqlOrgWebsitesStore) Update(rwebsites *model.Org_websites) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if sqlResult, err := us.GetReplica().Exec("UPDATE org_websites SET 	Org_id = :Org_id, Website = :Website WHERE Id = :Id",
			map[string]interface{}{
				"Org_id":        rwebsites.Org_id,
				"Website":    	 rwebsites.Website,
				"Id":            rwebsites.Id,
			}); err != nil {
			result.Err = model.NewAppError("SqlOrgWebsitesStore.Update", "store.sql_OrgWebsites.update.app_error", nil, "OrgWebsites, "+err.Error(), http.StatusInternalServerError)
		} else {
			rows, err := sqlResult.RowsAffected()

			if err != nil {
				result.Err = model.NewAppError("SqlOrgWebsitesStore.Update", "store.sql_OrgWebsites.update.app_error", nil, "id="+rwebsites.Id+", "+err.Error(), http.StatusInternalServerError)
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

func (s SqlOrgWebsitesStore) Clear() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				org_websites`); err != nil {
			result.Err = model.NewAppError("SqlDomainStore.Clear", "store.sql_job.delete.app_error", nil, "" + err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = true
		}
	})
}