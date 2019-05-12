package sqlstore

import (
	"net/http"
	"../../model"
	"../../store"
)

type SqlPopmailStore struct {
	SqlStore
}

func NewSqlPopmailStore(sqlStore SqlStore) store.PopmailStore {
	s := &SqlPopmailStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Popmail{}, "popmail").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(128)
		table.ColMap("Uidl").SetMaxSize(128)
		table.ColMap("Type").SetMaxSize(64)
		table.ColMap("Date").SetMaxSize(64)
		table.ColMap("CreateAt").SetMaxSize(64)
	}

	return s
}

func (s SqlPopmailStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_popmail_type", "popmail", "Id")
}

func (s SqlPopmailStore) Save(system *model.Popmail) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		system.PreSave()
		system.CreateAt = model.GetMillis()
		if err := s.GetMaster().Insert(system); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.Save", "store.sql_system.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = system
		}
	})
}

func (s SqlPopmailStore) SaveOrUpdate(popmail *model.Popmail) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if err := s.GetReplica().SelectOne(&model.Popmail{}, "SELECT * FROM popmail WHERE ID = :Id", map[string]interface{}{"Name": popmail.Id}); err == nil {
			if _, err := s.GetMaster().Update(popmail); err != nil {
				result.Err = model.NewAppError("SqlPopmailStore.SaveOrUpdate", "store.sql_system.update.app_error", nil, "", http.StatusInternalServerError)
			}
		} else {
			if err := s.GetMaster().Insert(popmail); err != nil {
				result.Err = model.NewAppError("SqlPopmailStore.SaveOrUpdate", "store.sql_system.save.app_error", nil, "", http.StatusInternalServerError)
			}
		}
	})
}

func (s SqlPopmailStore) Delete(id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				popmail
 			WHERE
 				Id = :Id`, map[string]interface{}{"Id": id}); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.DeleteById", "store.sql_job.delete.app_error", nil, "id="+id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = id
		}
	})
}

func (s SqlPopmailStore) Clear() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				popmail`); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.DeleteById", "store.sql_job.delete.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = true
		}
	})
}

func (s SqlPopmailStore) DeleteList(listId []string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		inQuery := ""
		for i := 0; i < len(listId); i++ {
			inQuery += "id = \"" + listId[i] + "\""

			if (len(listId)-1 != i) {
				inQuery += " OR "
			}
		}

		if _, err := s.GetMaster().Exec(`DELETE FROM popmail WHERE ` + inQuery + ``); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.DeleteList", "store.sql_job.delete_list.app_error", nil, "id", http.StatusInternalServerError)
		}

	})
}

func (s SqlPopmailStore) Update(popmail *model.Popmail) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Update(popmail); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.Update", "store.sql_system.update.app_error", nil, "", http.StatusInternalServerError)
		}
	})
}

func (s SqlPopmailStore) Get(id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if obj, err := s.GetReplica().Get(model.Popmail{}, id); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.Get", "store.popmail_store.get.app_error", nil, "user_id="+id+", "+err.Error(), http.StatusInternalServerError)
		} else if obj == nil {
			result.Err = model.NewAppError("SqlPopmailStore.Get", "store.popmail_store.get.app_error", nil, "user_id="+id, http.StatusNotFound)
		} else {
			result.Data = obj.(*model.Popmail)
		}
	})
}

func (s SqlPopmailStore) GetList(listId []string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var popmail []*model.Popmail
		strFilter := ""
		for i := 0; i < len(listId); i++ {
			strFilter += "id = \"" + listId[i] + "\""

			if (len(listId)-1 != i) {
				strFilter += " OR "
			}
		}

		query := "SELECT * FROM popmail WHERE " + strFilter

		if _, err := s.GetReplica().Select(&popmail, query); err != nil {
			result.Err = model.NewAppError("SqlUserStore.GetAllPopmail", "store.sql_popmail.getAllPopmail.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = popmail
		}

	})
}

// Получить письмо по Uidl
func (us SqlPopmailStore) GetByUidl(uidl string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var record = &model.Popmail{}
		if err := us.GetReplica().SelectOne(record, "SELECT * FROM Popmail WHERE Uidl = :Uidl", map[string]interface{}{"Uidl": uidl}); err != nil {
			result.Err = model.NewAppError("SqlFinHelpStore.GetByUidl", store.MISSING_POPMAIL_ERROR, nil, "UIDL="+uidl+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = record
		}
	})
}

// Получить все письма
func (us SqlPopmailStore) GetAllPopmail() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var popmail []*model.Popmail

		query := "SELECT * FROM Popmail ORDER BY CreateAt DESC"

		if _, err := us.GetReplica().Select(&popmail, query); err != nil {
			result.Err = model.NewAppError("SqlUserStore.GetAllPopmail", "store.sql_popmail.getAllPopmail.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = popmail
		}
	})
}

// Чанк сообщений в порядке возрастания
func (s SqlPopmailStore) GetChunk(offset int, limit int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if limit > 1000 {
			limit = 1000
			result.Err = model.NewAppError("SqlPopmailStore.Get", "store.sql_popmail.get.limit.app_error", nil, "user_id=", http.StatusBadRequest)
			return
		}

		query := "SELECT * FROM Popmail"
		query += " ORDER BY uidl ASC OFFSET :Offset ROWS FETCH NEXT :Limit ROWS ONLY"

		var records []*model.Popmail
		if _, err := s.GetReplica().Select(&records, query, map[string]interface{}{"Limit": limit, "Offset": offset}); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.Get", "store.sql_popmail.get.finding.app_error", nil, "user_id=", http.StatusInternalServerError)
		} else {
			result.Data = records
		}
	})
}

// Количество сообщений
func (us SqlPopmailStore) GetPopmailCount() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if count, err := us.GetReplica().SelectInt("SELECT COUNT(Id) FROM Popmail"); err != nil {
			result.Err = model.NewAppError("SqlPopmailStore.GetPopmailsCount", "store.sql_popmail.get_popmail_count.app_error", nil, err.Error(), http.StatusInternalServerError)
			result.Data = 0
		} else {
			result.Data = count
		}
	})
}