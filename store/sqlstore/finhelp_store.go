/*
 * Таблица org
*/

package sqlstore

import (
	"net/http"

	"../../model"
	"../../store"
	"strings"
)

type SqlFinHelpStore struct {
	SqlStore
}

func NewSqlFinHelpStore(sqlStore SqlStore) store.FinHelpStore {
	s := &SqlFinHelpStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.FinHelp{}, "org").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(36)
		table.ColMap("Type").SetMaxSize(128)
		table.ColMap("Subtype").SetMaxSize(128)
		table.ColMap("Bvkey").SetMaxSize(128)
		table.ColMap("Cp").SetMaxSize(128)
		table.ColMap("P").SetMaxSize(128)
		table.ColMap("Num").SetMaxSize(128)
		table.ColMap("U_o").SetMaxSize(128)
		table.ColMap("Gap").SetMaxSize(128)
		table.ColMap("Status").SetMaxSize(128)
		table.ColMap("Tip").SetMaxSize(128)
		table.ColMap("Namemax").SetMaxSize(2048)
		table.ColMap("Namemax1").SetMaxSize(2048)
		table.ColMap("Name").SetMaxSize(2048)
		table.ColMap("Namer").SetMaxSize(2048)
		table.ColMap("Uf_old").SetMaxSize(128)
		table.ColMap("Ust_f").SetMaxSize(128)
		table.ColMap("Ust_fi").SetMaxSize(128)
		table.ColMap("Ust_fb").SetMaxSize(128)
		table.ColMap("Kol_i").SetMaxSize(128)
		table.ColMap("Kol_b").SetMaxSize(128)
		table.ColMap("Kolf").SetMaxSize(128)
		table.ColMap("Kolfd").SetMaxSize(128)
		table.ColMap("Data_reg").SetMaxSize(128)
		table.ColMap("Mesto").SetMaxSize(256)
		table.ColMap("Data_preg").SetMaxSize(128)
		table.ColMap("Data_izmud").SetMaxSize(128)
		table.ColMap("Regn").SetMaxSize(128)
		table.ColMap("Priz").SetMaxSize(128)
		table.ColMap("Qq").SetMaxSize(256)
		table.ColMap("Lic_gold").SetMaxSize(128)
		table.ColMap("Lic_gold1").SetMaxSize(128)
		table.ColMap("Lic_rub").SetMaxSize(128)
		table.ColMap("Ogran").SetMaxSize(128)
		table.ColMap("Ogran1").SetMaxSize(128)
		table.ColMap("Ogran2").SetMaxSize(128)
		table.ColMap("Lic_val").SetMaxSize(128)
		table.ColMap("Vv").SetMaxSize(128)
		table.ColMap("Rei").SetMaxSize(128)
		table.ColMap("No").SetMaxSize(128)
		table.ColMap("Opp").SetMaxSize(128)
		table.ColMap("Ko").SetMaxSize(128)
		table.ColMap("Kp").SetMaxSize(128)
		table.ColMap("Adres").SetMaxSize(2048)
		table.ColMap("Adres1").SetMaxSize(2048)
		table.ColMap("Telefon").SetMaxSize(256)
		table.ColMap("Fax").SetMaxSize(128)
		table.ColMap("Fio_pr_pr").SetMaxSize(128)
		table.ColMap("Fio_gl_b").SetMaxSize(128)
		table.ColMap("Fio_zam_p").SetMaxSize(128)
		table.ColMap("Fio_zam_g").SetMaxSize(128)
		table.ColMap("Data_otz").SetMaxSize(128)
		table.ColMap("Pric_otz").SetMaxSize(128)
		table.ColMap("Data_pri").SetMaxSize(128)
		table.ColMap("Regn1").SetMaxSize(128)
		table.ColMap("Date_prb").SetMaxSize(128)
		table.ColMap("Type_prb").SetMaxSize(128)
		table.ColMap("Date_nam").SetMaxSize(128)
		table.ColMap("Date_adr").SetMaxSize(128)
		table.ColMap("Date_fiz").SetMaxSize(128)
		table.ColMap("Numsrf").SetMaxSize(128)
		table.ColMap("Okpo").SetMaxSize(128)
		table.ColMap("Dat_zayv").SetMaxSize(128)
		table.ColMap("Egr").SetMaxSize(128)
		table.ColMap("Prik_br").SetMaxSize(128)
		table.ColMap("Data_br").SetMaxSize(128)
		table.ColMap("Kliring").SetMaxSize(128)
		table.ColMap("Cb_date").SetMaxSize(128)
		table.ColMap("Fiz_end").SetMaxSize(128)
		table.ColMap("Num_ssv").SetMaxSize(128)
		table.ColMap("Data_ssv").SetMaxSize(128)
		table.ColMap("Cnlic").SetMaxSize(128)
		table.ColMap("Type_ko").SetMaxSize(128)
		table.ColMap("Name_en").SetMaxSize(2048)
		table.ColMap("Ann").SetMaxSize(128)
		table.ColMap("Tip_lic").SetMaxSize(128)
		table.ColMap("Prix").SetMaxSize(128)
		table.ColMap("Data_vib").SetMaxSize(128)
		table.ColMap("Ssv").SetMaxSize(128)
		table.ColMap("Inn").SetMaxSize(128)
		table.ColMap("Ogrn").SetMaxSize(128)
		table.ColMap("Reg_chartered").SetMaxSize(2048)
		table.ColMap("Post_chartered").SetMaxSize(2048)
		table.ColMap("Reg_city").SetMaxSize(512)
		table.ColMap("Reg_address").SetMaxSize(2048)
		table.ColMap("Ds").SetMaxSize(128)
		table.ColMap("De").SetMaxSize(128)
		table.ColMap("Search").SetMaxSize(2048)
	}

	return s
}

func (jss SqlFinHelpStore) CreateIndexesIfNotExists() {
	jss.CreateIndexIfNotExists("idx_finhelp_type", "org", "Id")
}

// Сохранение
func (jss SqlFinHelpStore) Save(fin *model.FinHelp) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		fin.PreSave()
		if err := jss.GetMaster().Insert(fin); err != nil {
			result.Err = model.NewAppError("SqlFinHelpStore.Save", "store.sql_finhelp.save.app_error", nil, "id="+fin.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = fin
		}
	})
}

// Получение записи по полю из базы данных
func (us SqlFinHelpStore) GetByInn(Inn string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		Inn = strings.ToLower(Inn)
		var fin = &model.FinHelp{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org WHERE Inn = :Inn ORDER BY Inn DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Inn": Inn})
		if err != nil {
			result.Err = model.NewAppError("SqlFinHelpStore.GetByInn", store.MISSING_INN_ERROR, nil, "Inn="+Inn+", "+err.Error(), http.StatusInternalServerError)
		} else {
			fin.Id = Id
			result.Data = fin
		}
	})
}
// Получение записи по полю из базы данных
func (us SqlFinHelpStore) GetByBvkey(Bvkey string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var fin = &model.FinHelp{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org WHERE Bvkey = :Bvkey ORDER BY Bvkey DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Bvkey": Bvkey})
		if err != nil {
			result.Err = model.NewAppError("SqlFinHelpStore.GetByBvkey", store.MISSING_INN_ERROR, nil, "Bvkey="+Bvkey+", "+err.Error(), http.StatusInternalServerError)
		} else {
			fin.Id = Id
			result.Data = fin
		}
	})
}
func (us SqlFinHelpStore) GetByNum(Num string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		Num = strings.ToLower(Num)
		var fin = &model.FinHelp{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org WHERE Num = :Num ORDER BY Num DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Num": Num})
		if err != nil {
			result.Err = model.NewAppError("SqlFinHelpStore.GetByNum", store.MISSING_INN_ERROR, nil, "Num="+Num+", "+err.Error(), http.StatusInternalServerError)
		} else {
			fin.Id = Id
			result.Data = fin
		}
	})
}
func (us SqlFinHelpStore) GetByRegn(Regn string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		Regn = strings.ToLower(Regn)
		var fin = &model.FinHelp{}
		Id, err := us.GetReplica().SelectStr("SELECT Id FROM org WHERE Regn = :Regn ORDER BY Regn DESC OFFSET 0 ROWS FETCH NEXT (1) ROWS ONLY", map[string]interface{}{"Regn": Regn})
		if err != nil {
			result.Err = model.NewAppError("SqlFinHelpStore.GetByInn", store.MISSING_INN_ERROR, nil, "Regn="+Regn+", "+err.Error(), http.StatusInternalServerError)
		} else {
			fin.Id = Id
			result.Data = fin
		}
	})
}


func (us SqlFinHelpStore) UpdateBank(record *model.FinHelp) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := us.GetReplica().Exec(`
UPDATE org
SET Type           = :Type,
    Subtype        = :Subtype,
    Bvkey          = :Bvkey,
    Cp             = :Cp,
    P              = :P,
    Num            = :Num,
    U_o            = :U_o,
    Gap            = :Gap,
    Status         = :Status,
    Tip            = :Tip,
    Namemax        = :Namemax,
    Namemax1       = :Namemax1,
    Name           = :Name,
    Namer          = :Namer,
    Uf_old         = :Uf_old,
    Ust_f          = :Ust_f,
    Ust_fi         = :Ust_fi,
    Ust_fb         = :Ust_fb,
    Kol_i          = :Kol_i,
    Kol_b          = :Kol_b,
    Kolf           = :Kolf,
    Kolfd          = :Kolfd,
    Data_reg       = :Data_reg,
    Mesto          = :Mesto,
    Data_preg      = :Data_preg,
    Data_izmud     = :Data_izmud,
    Regn           = :Regn,
    Priz           = :Priz,
	Prik           = :Prik,
    Qq             = :Qq,
    Lic_gold       = :Lic_gold,
    Lic_gold1      = :Lic_gold1,
    Lic_rub        = :Lic_rub,
    Ogran          = :Ogran,
    Ogran1         = :Ogran1,
    Ogran2         = :Ogran2,
    Lic_val        = :Lic_val,
    Vv             = :Vv,
    Rei            = :Rei,
    No             = :No,
    Opp            = :Opp,
    Ko             = :Ko,
    Kp             = :Kp,
    Adres          = :Adres,
    Adres1         = :Adres1,
    Telefon        = :Telefon,
    Fax            = :Fax,
    Fio_pr_pr      = :Fio_pr_pr,
    Fio_gl_b       = :Fio_gl_b,
    Fio_zam_p      = :Fio_zam_p,
    Fio_zam_g      = :Fio_zam_g,
    Data_otz       = :Data_otz,
    Pric_otz       = :Pric_otz,
    Data_pri       = :Data_pri,
    Date_prb       = :Date_prb,
    Type_prb       = :Type_prb,
    Date_nam       = :Date_nam,
    Date_adr       = :Date_adr,
    Date_fiz       = :Date_fiz,
    Numsrf         = :Numsrf,
    Okpo           = :Okpo,
    Dat_zayv       = :Dat_zayv,
    Egr            = :Egr,
    Data_br        = :Data_br,
    Prik_br        = :Prik_br,
    Kliring        = :Kliring,
    Cb_date        = :Cb_date,
    Fiz_end        = :Fiz_end,
    Num_ssv        = :Num_ssv,
    Data_ssv       = :Data_ssv,
    Cnlic          = :Cnlic,
    Type_ko        = :Type_ko,
    Name_en        = :Name_en,
    Ann            = :Ann,
    Tip_lic        = :Tip_lic,
    Prix           = :Prix,
    Data_vib       = :Data_vib,
    Ssv            = :Ssv,
    Inn            = :Inn,
    Ogrn           = :Ogrn,
    Reg_chartered  = :Reg_chartered,
    Post_chartered = :Post_chartered,
    Reg_city       = :Reg_city,
    Reg_address    = :Reg_address,
    Ds             = :Ds,
    De             = :De,
    Search         = :Search
WHERE Id = :Id`,
			map[string]interface{}{
				"Type": record.Type,
				"Subtype": record.Subtype,
				"Bvkey": record.Bvkey,
				"Cp": record.Cp,
				"P": record.P,
				"Num": record.Num,
				"U_o": record.U_o,
				"Gap": record.Gap,
				"Status": record.Status,
				"Tip": record.Tip,
				"Namemax": record.Namemax,
				"Namemax1": record.Namemax1,
				"Name": record.Name,
				"Namer": record.Namer,
				"Uf_old": record.Uf_old,
				"Ust_f": record.Ust_f,
				"Ust_fi": record.Ust_fi,
				"Ust_fb": record.Ust_fb,
				"Kol_i": record.Kol_i,
				"Kol_b": record.Kol_b,
				"Kolf": record.Kolf,
				"Kolfd": record.Kolfd,
				"Data_reg": record.Data_reg,
				"Mesto": record.Mesto,
				"Data_preg": record.Data_preg,
				"Data_izmud": record.Data_izmud,
				"Regn": record.Regn,
				"Priz": record.Priz,
				"Qq": record.Qq,
				"Lic_gold": record.Lic_gold,
				"Lic_gold1": record.Lic_gold1,
				"Lic_rub": record.Lic_rub,
				"Ogran": record.Ogran,
				"Ogran1": record.Ogran1,
				"Ogran2": record.Ogran2,
				"Lic_val": record.Lic_val,
				"Vv": record.Vv,
				"Rei": record.Rei,
				"No": record.No,
				"Opp": record.Opp,
				"Ko": record.Ko,
				"Kp": record.Kp,
				"Adres": record.Adres,
				"Adres1": record.Adres1,
				"Telefon": record.Telefon,
				"Fax": record.Fax,
				"Prik": record.Prik,
				"Fio_pr_pr": record.Fio_pr_pr,
				"Fio_gl_b": record.Fio_gl_b,
				"Fio_zam_p": record.Fio_zam_p,
				"Fio_zam_g": record.Fio_zam_g,
				"Data_otz": record.Data_otz,
				"Pric_otz": record.Pric_otz,
				"Data_pri": record.Data_pri,
				"Date_prb": record.Date_prb,
				"Type_prb": record.Type_prb,
				"Date_nam": record.Date_nam,
				"Date_adr": record.Date_adr,
				"Date_fiz": record.Date_fiz,
				"Numsrf": record.Numsrf,
				"Okpo": record.Okpo,
				"Dat_zayv": record.Dat_zayv,
				"Egr": record.Egr,
				"Prik_br": record.Prik_br,
				"Data_br": record.Data_br,
				"Kliring": record.Kliring,
				"Cb_date": record.Cb_date,
				"Fiz_end": record.Fiz_end,
				"Num_ssv": record.Num_ssv,
				"Data_ssv": record.Data_ssv,
				"Cnlic": record.Cnlic,
				"Type_ko": record.Type_ko,
				"Name_en": record.Name_en,
				"Ann": record.Ann,
				"Tip_lic": record.Tip_lic,
				"Prix": record.Prix,
				"Data_vib": record.Data_vib,
				"Ssv": record.Ssv,
				"Inn": record.Inn,
				"Ogrn": record.Ogrn,
				"Reg_chartered": record.Reg_chartered,
				"Post_chartered": record.Post_chartered,
				"Reg_city": record.Reg_city,
				"Reg_address": record.Reg_address,
				"Ds": record.Ds,
				"De": record.De,
				"Search": record.Search,
				"Id": record.Id,
			}); err != nil {
				result.Err = model.NewAppError("update", "import", nil, err.Error(), 500)
			result.Data = false
		} else {

			result.Data = true

		}

	})
}


func (us SqlFinHelpStore) UpdateMicro(micro *model.FinHelp) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := us.GetReplica().Exec(`
UPDATE org
SET Namemax        = :Namemax,
    Inn            = :Inn,
    Ogrn           = :Ogrn,
    Reg_chartered  = :Reg_chartered,
    Post_chartered = :Post_chartered,
    Reg_city       = :Reg_city,
    Reg_address    = :Reg_address,
    Name           = :Name,
    Num            = :Num,
    Ds             = :Ds,
    De             = :De,
    Search         = :Search
WHERE Id = :Id`,
			map[string]interface{}{
				"Namemax":        micro.Namemax,
				"Inn":            micro.Inn,
				"Ogrn":           micro.Ogrn,
				"Reg_chartered":  micro.Reg_chartered,
				"Post_chartered": micro.Post_chartered,
				"Reg_city":       micro.Reg_city,
				"Reg_address":    micro.Reg_address,
				"Name":           micro.Name,
				"Num":            micro.Num,
				"Ds":             micro.Ds,
				"De":             micro.De,
				"Search":         micro.Search,
				"Id":             micro.Id,
			}); err == nil {
			result.Data = true
		} else {

			result.Data = false

		}

	})
}