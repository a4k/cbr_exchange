/*
 * Таблица Org_insurance_lics
*/
package sqlstore

import (
	"net/http"

	"../../model"
	"../../store"
	"strings"
)

type SqlInsuranceStore struct {
	SqlStore
}

func NewSqlInsuranceStore(sqlStore SqlStore) store.InsuranceStore {
	s := &SqlInsuranceStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Org_insurance_lics{}, "org_insurance_lics").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(36)
		table.ColMap("Org_id").SetMaxSize(128)
		table.ColMap("Description").SetMaxSize(2048)
		table.ColMap("Lic_type_code").SetMaxSize(128)
		table.ColMap("Status").SetMaxSize(128)
		table.ColMap("Lic_number").SetMaxSize(128)
		table.ColMap("Lic_date").SetMaxSize(128)
		table.ColMap("Lic_expire").SetMaxSize(128)
		table.ColMap("Lic_end").SetMaxSize(128)
		table.ColMap("End_reason").SetMaxSize(128)
	}

	return s
}

func (jss SqlInsuranceStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_insurance_type", "org_insurance_lics", "Id")
}

// Сохранение
func (jss SqlInsuranceStore) Save(insurance *model.Org_insurance_lics) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		insurance.PreSave()
		if err := jss.GetMaster().Insert(insurance); err != nil {
			result.Err = model.NewAppError("SqlInsuranceStore.Save", "store.sql_insurance.save.app_error", nil, "id="+insurance.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = insurance
		}
	})
}

// Поиск в базе данных
func (us SqlInsuranceStore) GetByOrgId(org_id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var finsurance = &model.Org_insurance_lics{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_insurance_lics WHERE org_id = :Org_id ORDER BY Org_id DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Org_id": org_id})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgId", store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			finsurance.Id = Id
			result.Data = finsurance
		}
	})
}
// Поиск в базе данных
func (us SqlInsuranceStore) GetByOrgIdAndType(org_id string, lic_type_code string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		org_id = strings.ToLower(org_id)
		var finsurance = &model.Org_insurance_lics{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org_insurance_lics WHERE Org_id = :Org_id AND Lic_type_code = :Lic_type_code ORDER BY Lic_type_code DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY",
			map[string]interface{}{
				"Org_id": org_id,
				"Lic_type_code": lic_type_code,
			})
		if err != nil {
			result.Err = model.NewAppError("SqlOrgLicsStore.GetByOrgIdAndType", store.MISSING_LIC_NUMBER_ERROR, nil, "Org_id="+org_id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			finsurance.Id = Id
			result.Data = finsurance
		}
	})
}

// Обновление
func (us SqlInsuranceStore) Update(rinsurance *model.Org_insurance_lics) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if sqlResult, err := us.GetReplica().Exec("UPDATE org_insurance_lics SET 	Org_id = :Org_id, Description = :Description, " +
			"Lic_type_code = :Lic_type_code, Status = :Status, Lic_number = :Lic_number, Lic_date = :Lic_date, Lic_expire = :Lic_expire, " +
			"Lic_end = :Lic_end, End_reason = :End_reason WHERE Id = :Id",
			map[string]interface{}{
				"Org_id":        rinsurance.Org_id,
				"Description":   rinsurance.Description,
				"Lic_type_code": rinsurance.Lic_type_code,
				"Status":        rinsurance.Status,
				"Lic_number":    rinsurance.Lic_number,
				"Lic_date":      rinsurance.Lic_date,
				"Lic_expire":    rinsurance.Lic_expire,
				"Lic_end":       rinsurance.Lic_end,
				"End_reason":    rinsurance.End_reason,
				"Id":            rinsurance.Id,
			}); err != nil {
			result.Err = model.NewAppError("SqlInsuranceStore.Update", "store.sql_Insurance.update.app_error", nil, "Insurance, "+err.Error(), http.StatusInternalServerError)
		} else {
			rows, err := sqlResult.RowsAffected()

			if err != nil {
				result.Err = model.NewAppError("SqlInsuranceStore.Update", "store.sql_Insurance.update.app_error", nil, "id="+rinsurance.Id+", "+err.Error(), http.StatusInternalServerError)
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

func (s SqlInsuranceStore) Clear() store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := s.GetMaster().Exec(
			`DELETE FROM
 				org_insurance_lics`); err != nil {
			result.Err = model.NewAppError("SqlDomainStore.Clear", "store.sql_job.delete.app_error", nil, "" + err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = true
		}
	})
}