package model

import (
	"encoding/json"
	"strings"
)

const (
	FINHELP_KEY_ADVANCED_PARSING = "finhelp_advanced_parsing"
)

// Структура данных финансовых документов
type FinHelp struct {
	Id             string `csv:"ID" db:"id"`
	Type           string `csv:"TYPE" db:"type"`
	Subtype        string `csv:"SUBTYPE" db:"subtype"`
	Bvkey          string `csv:"BVKEY" db:"bvkey"`
	Cp             string `csv:"CP" db:"cp"`
	P              string `csv:"P" db:"p"`
	Num            string `csv:"REG_NUM" db:"num"`
	U_o            string `csv:"U_O" db:"u_o"`
	Gap            string `csv:"GAP" db:"gap"`
	Status         string `csv:"STATUS" db:"status"`
	Tip            string `csv:"TIP" db:"tip"`
	Namemax        string `csv:"NAMEMAX" db:"namemax"`
	Namemax1       string `csv:"NAMEMAX1" db:"namemax1"`
	Name           string `csv:"NAME" db:"name"`
	Namer          string `csv:"NAMER" db:"namer"`
	Uf_old         string `csv:"UF_OLD" db:"uf_old"`
	Ust_f          string `csv:"UST_F" db:"ust_f"`
	Ust_fi         string `csv:"UST_FI" db:"ust_fi"`
	Ust_fb         string `csv:"UST_FB" db:"ust_fb"`
	Kol_i          string `csv:"KOL_I" db:"kol_i"`
	Kol_b          string `csv:"KOL_B" db:"kol_b"`
	Kolf           string `csv:"KOLF" db:"kolf"`
	Kolfd          string `csv:"KOLFD" db:"kolfd"`
	Data_reg       string `csv:"DATA_REG" db:"data_reg"`
	Mesto          string `csv:"MESTO" db:"mesto"`
	Data_preg      string `csv:"DATA_PREG" db:"data_preg"`
	Data_izmud     string `csv:"DATA_IZMUD" db:"data_izmud"`
	Regn           string `csv:"REGN" db:"regn"`
	Priz           string `csv:"PRIZ" db:"priz"`
	Qq             string `csv:"QQ" db:"qq"`
	Lic_gold       string `csv:"LIC_GOLD" db:"lic_gold"`
	Lic_gold1      string `csv:"LIC_GOLD1" db:"lic_gold1"`
	Lic_rub        string `csv:"LIC_RUB" db:"lic_rub"`
	Ogran          string `csv:"OGRAN" db:"ogran"`
	Ogran1         string `csv:"OGRAN1" db:"ogran1"`
	Ogran2         string `csv:"OGRAN2" db:"ogran2"`
	Lic_val        string `csv:"LIC_VAL" db:"lic_val"`
	Vv             string `csv:"VV" db:"vv"`
	Rei            string `csv:"REI" db:"rei"`
	No             string `csv:"NO" db:"no"`
	Opp            string `csv:"OPP" db:"opp"`
	Ko             string `csv:"KO" db:"ko"`
	Kp             string `csv:"KP" db:"kp"`
	Adres          string `csv:"ADRES" db:"adres"`
	Adres1         string `csv:"ADRES1" db:"adres1"`
	Telefon        string `csv:"TELEFON" db:"telefon"`
	Fax            string `csv:"FAX" db:"fax"`
	Prik           string `csv:"PRIK" db:"prik"`
	Fio_pr_pr      string `csv:"FIO_PR_PR" db:"fio_pr_pr"`
	Fio_gl_b       string `csv:"FIO_GL_B" db:"fio_gl_b"`
	Fio_zam_p      string `csv:"FIO_ZAM_P" db:"fio_zam_p"`
	Fio_zam_g      string `csv:"FIO_ZAM_G" db:"fio_zam_g"`
	Data_otz       string `csv:"DATA_OTZ" db:"data_otz"`
	Pric_otz       string `csv:"PRIC_OTZ" db:"pric_otz"`
	Data_pri       string `csv:"DATA_PRI" db:"data_pri"`
	Regn1          string `csv:"REGN1" db:"regn1"`
	Date_prb       string `csv:"DATE_PRB" db:"date_prb"`
	Type_prb       string `csv:"TYPE_PRB" db:"type_prb"`
	Date_nam       string `csv:"DATE_NAM" db:"date_nam"`
	Date_adr       string `csv:"DATE_ADR" db:"date_adr"`
	Date_fiz       string `csv:"DATE_FIZ" db:"date_fiz"`
	Numsrf         string `csv:"NUMSRF" db:"numsrf"`
	Okpo           string `csv:"OKPO" db:"okpo"`
	Dat_zayv       string `csv:"DAT_ZAYV" db:"dat_zayv"`
	Egr            string `csv:"EGR" db:"egr"`
	Prik_br        string `csv:"PRIK_BR" db:"prik_br"`
	Data_br        string `csv:"DATA_BR" db:"data_br"`
	Kliring        string `csv:"KLIRING" db:"kliring"`
	Cb_date        string `csv:"CB_DATE" db:"cb_date"`
	Fiz_end        string `csv:"FIZ_END" db:"fiz_end"`
	Num_ssv        string `csv:"NUM_SSV" db:"num_ssv"`
	Data_ssv       string `csv:"DATA_SSV" db:"data_ssv"`
	Cnlic          string `csv:"CNLIC" db:"cnlic"`
	Type_ko        string `csv:"TYPE_KO" db:"type_ko"`
	Name_en        string `csv:"NAME_EN" db:"name_en"`
	Ann            string `csv:"ANN" db:"ann"`
	Tip_lic        string `csv:"TIP_LIC" db:"tip_lic"`
	Prix           string `csv:"PRIX" db:"prix"`
	Data_vib       string `csv:"DATA_VIB" db:"data_vib"`
	Ssv            string `csv:"SSV" db:"ssv"`
	Inn            string `csv:"INN" db:"inn"`
	Ogrn           string `csv:"OGRN" db:"ogrn"`
	Reg_chartered  string `csv:"REG_CHARTERED" db:"reg_chartered"`
	Post_chartered string `csv:"POST_CHARTERED" db:"post_chartered"`
	Reg_city       string `csv:"REG_CITY" db:"reg_city"`
	Reg_address    string `csv:"REG_ADDRESS" db:"reg_address"`
	Ds             string `csv:"DS" db:"ds"`
	De             string `csv:"DE" db:"de"`
	Search         string `csv:"SEARCH" db:"search"`
}

// Структура микрофинансовых
type Micro struct {
	Id             string `csv:"ID"`      // Id
	Namemax        string `csv:"NAME"` // Name
	Inn            string `csv:"INN"`
	Ogrn           string `csv:"OGRN"`
	Reg_chartered  string `csv:"REG_CHARTERED"`
	Post_chartered string `csv:"POST_CHARTERED"`
	Reg_city       string `csv:"REG_CITY"`
	Reg_address    string `csv:"REG_ADDRESS"`
	Name           string `csv:"SHORT_NAME"`    // Short_name
	Num            string `csv:"REG_NUM"` // Reg_num
	Ds             string `csv:"DS"`
	De             string `csv:"DE"`
	Search         string `csv:"NAME"`
}

// Структура страховых компаний
type Insurance struct {
	Id             string `csv:"ID"` // Id
	Org_id         string `csv:"#SBJ_ID"`
	Namemax        string `csv:"NAME"`
	Namemax1       string `csv:"NAME"`
	Name           string `csv:"SHORT_NAME"`
	Inn            string `csv:"INN"`
	Ogrn           string `csv:"OGRN"`
	Num            string `csv:"REG_NUM"`
	Reg_city       string `csv:"REG_CITY"`
	Reg_settlement string `csv:"REG_SETTLEMENT"`
	Reg_address    string `csv:"REG_ADDRESS"`
	Post_chartered string `csv:"POST_CHARTERED"`
	Reg_chartered  string `csv:"REG_CHARTERED"`
	Email          string `csv:"EMAIL"`
	Phones         string `csv:"PHONES"`
	Website        string `csv:"WEBSITE"`
	Description    string `csv:"DESCRIPTION"`
	Lic_type_code  string `csv:"LIC_TYPE_CODE"`
	Status         string `csv:"STATUS"`
	Lic_number     string `csv:"LIC_NUMBER"`
	Lic_date       string `csv:"LIC_DATE"`
	Lic_expire     string `csv:"LIC_EXPIRE"`
	Lic_end        string `csv:"LIC_END"`
	End_reason     string `csv:"END_REASON"`
	Search         string `csv:"SHORT_NAME"`
}

// Структура страховых компаний
type Org_insurance_lics struct {
	Id            string `csv:"ID" db:"id"`
	Org_id        string `csv:"SBJ_ID" db:"org_id"`
	Description   string `csv:"DESCRIPTION" db:"description"`
	Lic_type_code string `csv:"LIC_TYPE_CODE" db:"lic_type_code"`
	Status        string `csv:"STATUS" db:"status"`
	Lic_number    string `csv:"LIC_NUMBER" db:"lic_number"`
	Lic_date      string `csv:"LIC_DATE" db:"lic_date"`
	Lic_expire    string `csv:"LIC_EXPIRE" db:"lic_expire"`
	Lic_end       string `csv:"LIC_END" db:"lic_end"`
	End_reason    string `csv:"END_REASON" db:"end_reason"`
}

// Структура почтовых адресов
type Org_emails struct {
	Id     string `csv:"ID" db:"id"` // Id
	Email  string `csv:"EMAIL"  db:"email"`
	Org_id string `csv:"ORG_ID" db:"org_id"` // Org_id
}

// Структура номеров факса
type Org_faxs struct {
	Id     string `csv:"ID" db:"id"` // Id
	Number string `csv:"Number" db:"number"`
	Org_id string `csv:"ORG_ID" db:"org_id"` // Org_id
}

// Структура лицензий
type Org_lics struct {
	Id     string `csv:"ID" db:"id"` // Id
	Type   string `csv:"TYPE"  db:"type"`
	Name   string `csv:"NAME"  db:"name"`
	Allow  string `csv:"Allow"  db:"allow"`
	Org_id string `csv:"ORG_ID" db:"org_id"` // Org_id
}

// Структура номеров телефона
type Org_phones struct {
	Id     string `csv:"ID" db:"id"` // Id
	Number string `csv:"PHONES" db:"number"`
	Org_id string `csv:"ORG_ID" db:"org_id"` // Org_id
}

// Структура вебсайтов
type Org_websites struct {
	Id      string `csv:"ID" db:"id"` // Id
	Website string `csv:"WEBSITE" db:"website"`
	Org_id  string `csv:"ORG_ID" db:"org_id"` // Org_id
}

func (u *FinHelp) DeepCopy() *FinHelp {
	copyFinHelp := *u

	return &copyFinHelp
}

// ToJson convert a User to a json string
func (u *FinHelp) ToJson() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// IsValid validates the user and returns an error if it isn't configured
// correctly.
func (u *FinHelp) IsValid() *AppError {

	return nil
}

// PreSave will set the Id and Username if missing.  It will also fill
// in the CreateAt, UpdateAt times.  It will also hash the password.  It should
// be run before saving the user to the db.
func (u *FinHelp) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}

func (u *Org_insurance_lics) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}
func (u *Org_emails) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}
func (u *Org_faxs) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}
func (u *Org_phones) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}
func (u *Org_lics) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}

func (u *Org_websites) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}

// Приводим данные к единому виду
func NormalizeMicro(micro Micro) *FinHelp {
	a := &FinHelp{}
	a.Type = "micro"
	a.Namemax = micro.Namemax
	a.Inn = micro.Inn
	a.Ogrn = micro.Ogrn
	a.Reg_chartered = micro.Reg_chartered
	a.Post_chartered = micro.Post_chartered
	a.Reg_city = micro.Reg_city
	a.Reg_address = micro.Reg_address
	a.Name = micro.Name
	a.Num = strings.Replace(micro.Num,"-","",-1)
	a.Ds = micro.Ds
	a.De = micro.De
	a.Search = micro.Namemax
	return a
}

func NormalizeInsurance(insurance Insurance) *FinHelp {
	a := &FinHelp{}
	a.Type = "insurance"
	a.Namemax = insurance.Namemax
	a.Namemax1 = insurance.Namemax1
	a.Name = insurance.Name
	a.Inn = insurance.Inn
	a.Ogrn = insurance.Ogrn
	a.Num = insurance.Num
	a.Reg_city = insurance.Reg_city
	a.Reg_address = insurance.Reg_address
	a.Post_chartered = insurance.Post_chartered
	a.Reg_chartered = insurance.Reg_chartered
	a.Search = insurance.Namemax
	return a
}

func NormalizeOrigInsurance(insurance Insurance) *Org_insurance_lics {
	a := &Org_insurance_lics{}
	a.Id = insurance.Id
	a.Org_id = insurance.Org_id
	a.Description = insurance.Description
	a.Lic_type_code = insurance.Lic_type_code
	a.Status = insurance.Status
	a.Lic_number = insurance.Lic_number
	a.Lic_date = insurance.Lic_date
	a.Lic_expire = insurance.Lic_expire
	a.Lic_end = insurance.Lic_end
	a.End_reason = insurance.End_reason
	return a
}

func NormalizeOrgEmails(insurance Insurance) *Org_emails {
	a := &Org_emails{}
	a.Id = insurance.Id
	a.Org_id = insurance.Org_id
	a.Email = insurance.Email
	return a
}
func NormalizeOrgFaxs(insurance Insurance) *Org_faxs {
	a := &Org_faxs{}
	a.Id = insurance.Id
	a.Org_id = insurance.Org_id
	a.Number = insurance.Phones
	return a
}
func NormalizeOrgPhones(insurance Insurance) *Org_phones {
	a := &Org_phones{}
	a.Id = insurance.Id
	a.Org_id = insurance.Org_id
	a.Number = insurance.Phones
	return a
}
func NormalizeOrgWebsites(insurance Insurance) *Org_websites {
	a := &Org_websites{}
	a.Id = insurance.Id
	a.Org_id = insurance.Org_id
	a.Website = insurance.Website
	return a
}
