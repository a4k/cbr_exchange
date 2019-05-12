// Parse Files from dbf, csv

package rufr

import (
	"bytes"
	"fmt"

	"encoding/json"
	"mime"

	"io"

	"net/http"

	"strings"
	"encoding/csv"
	"../utils/godbf"
	l4g "../utils/log4go"
	"../model"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	pop3 "../utils"
	"strconv"
	"time"
	"sort"
	"math/rand"
	"regexp"
)

type AdvancedParsing struct {
	CurrentTable  string `json:"current_table"`
}


const (

	LOG_ERR_25001 = "Ошибка во время импорта %v"
	LOG_ERR_25002 = "Не найден файл при начале парсинга %s"
	LOG_ERR_25003 = "Ошибки при парсинге файла  %s %s"
	LOG_ERR_25004 = "Ошибка создания папки для бэкапа файла из письма %s"
	LOG_ERR_25005 = "Ошибка создания бэкапа файла из письма %s"
	LOG_ERR_25006 = "Ошибка копирования бэкапа файла из письма %s"

	LOG_INFO_25007 =	"Начало импорта файлов"
	LOG_INFO_25008 =	"Получение писем"
	LOG_INFO_25009 =	"Обработка писем.."
	LOG_INFO_25010 =	"Новое письмо по запросу TOP не найдено"
	LOG_INFO_25011 =	"Найдены новые сообщения"
	LOG_INFO_25012 =	"Импорт завершен."
	LOG_INFO_25013 =	"Новых писем не найдено. Завершено."
	LOG_INFO_25014 =	"Новых писем не найдено. Завершено."
	LOG_INFO_25015 =	"Парсинг файла %s"
	LOG_INFO_25016 =	"[Org] Micro Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25017 =	"Парсинг файла завершен %s"
	LOG_INFO_25018 =	"Ошибки во время ClearAllPhones force Update"
	LOG_INFO_25019 =	"Ошибки во время ClearAllEmails force Update"
	LOG_INFO_25020 =	"Ошибки во время ClearAllFaxs force Update"
	LOG_INFO_25021 =	"Ошибки во время ClearAllWebsites force Update"
	LOG_INFO_25022 =	"Ошибки во время ClearAllInsurance force Update"
	LOG_INFO_25023 =	"Начинается парсинг файла %s"
	LOG_INFO_25024 =	"В файле найдены невалидные данные. Файл :: %s"
	LOG_INFO_25025 =	"[Org] Insurance Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25026 =	"[Org_insurance_lics] Добавлено: %d, Обновлено: %d, Ошибок: %d,"
	LOG_INFO_25027 =	"[Org_emails] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25028 =	"[Org_faxs] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25029 =	"[Org_phones] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25030 =	"[Org_websites] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25031 =	"Парсинг файла завершен %s"
	LOG_INFO_25032 =	"Не найден файл .DBF по адресу %s"
	LOG_INFO_25033 =	"Ошибки во время force Update"
	LOG_INFO_25034 =	"Org_phones очищена"
	LOG_INFO_25035 =	"Ошибки во время force Update"
	LOG_INFO_25036 =	"Org_emails очищена"
	LOG_INFO_25037 =	"Ошибки во время force Update"
	LOG_INFO_25038 =	"Org_lics очищена"
	LOG_INFO_25039 =	"Ошибка при чтении файла .DBF"
	LOG_INFO_25040 =	"Ошибка при чтении файла .DBF"
	LOG_INFO_25041 =	"Начинается парсинг файла %s"
	LOG_INFO_25042 =	"[Org] Bank Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25043 =	"[Org_faxs] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25044 =	"[Org_lics] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25045 =	"Парсинг файла завершен %s"
	LOG_INFO_25046 =	"Не найден файл .DBF по адресу %s"
	LOG_INFO_25047 =	"Ошибки во время force Update"
	LOG_INFO_25048 =	"Org_phones очищена"
	LOG_INFO_25049 =	"Ошибки во время force Update"
	LOG_INFO_25050 =	"Org_emails очищена"
	LOG_INFO_25051 =	"Ошибки во время force Update"
	LOG_INFO_25052 =	"Org_lics очищена"
	LOG_INFO_25053 =	"Ошибка при чтении файла .DBF"
	LOG_INFO_25054 =	"Ошибка при чтении файла .DBF"
	LOG_INFO_25055 =	"Начинается парсинг файла %s"
	LOG_INFO_25056 =	"[Org] Barx Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_INFO_25057 =	"Парсинг файла %s завершен"
	LOG_INFO_25058 =	"Получение доступа к почте pop3"
	LOG_INFO_25059 =	"Вход в почту выполнен без ошибок"
	LOG_INFO_25060 =	"Ошибки при получении данных о письме %s"
	LOG_INFO_25061 =	"Не удалось получить информацию из письма %s"
	LOG_INFO_25062 =	"Номер сообщения не может преобразоваться в uint64 %s"
	LOG_INFO_25063 =	"Найдено руфр письмо %s"
	LOG_INFO_25064 =	"Начинается парсинг письма %s"
	LOG_INFO_25065 =	"Создана папка для бэкапа файла из письма %s"
	LOG_INFO_25066 =	"В письме не найдено прикреплений"
	LOG_INFO_25067 =	"Не определен файл %s"
	LOG_INFO_25068 =	"Не определен файл %s"
	LOG_INFO_25069 =	"Не определен файл %s"
	LOG_INFO_25070 =	"Прикрепленный файл имеет не стандартный тип :: %s"
	LOG_INFO_25071 =	"Не найден Content-Type у письма"
	LOG_INFO_25072 =	"Идет процесс распаковки архива"
	LOG_INFO_25073 =	"В архиве не найдено файлов"
	LOG_INFO_25074 =	"Архив распакован"
	LOG_INFO_25075 =	"Не удалось получить тип файла"

	LOG_WARN_25076 = "Импорт завершен. Возникли ошибки во время импорта"
	LOG_WARN_25077 = "Ошибки при получении данных о письме %s"
	LOG_WARN_25078 = "[Org] Micro Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25079 = "[Org] Insurance Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25080 = "[Org_insurance_lics] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25081 = "[Org_emails] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25082 = "[Org_faxs] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25083 = "[Org_phones] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25084 = "[Org_websites] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25085 = "Ошибки при парсинге файла %s %s"
	LOG_WARN_25086 = "В файле не найдено строк для парсинга %s"
	LOG_WARN_25087 = "Найдены невалидные данные в файле %s"
	LOG_WARN_25088 = "Найдены ошибки при парсинге полей из файла. Количество ошибок: %d, Файл: %s"
	LOG_WARN_25089 = "[Org] Bank Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25090 = "[Org_faxs] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25091 = "[Org_lics] Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25092 = "Ошибки при парсинге файла %s %s"
	LOG_WARN_25093 = "Найдены невалидные данные в файле %s"
	LOG_WARN_25094 = "Найдены ошибки при парсинге полей из файла. Количество ошибок: %d, Файл: %s"
	LOG_WARN_25095 = "[Org] Barx Добавлено: %d, Обновлено: %d, Ошибок: %d"
	LOG_WARN_25096 = "Ошибки при парсинге файла %s %s"
	LOG_WARN_25097 = "Вложенность более одного раза не допускается"
	LOG_WARN_25098 = "Вложенность более одного раза не допускается"
	LOG_WARN_25099 = "Вложенность более одного раза не допускается"
	LOG_WARN_25100 = "Вложенность более одного раза не допускается"
	LOG_WARN_25101 = "Вложенность более одного раза не допускается"
	LOG_WARN_25102 = "Вложенность более одного раза не допускается"
	LOG_WARN_25103 = "Не удалось распаковать архив"
	LOG_WARN_25104 = "Не удалось распаковать архив"
	LOG_WARN_25105 = "Не удалось распаковать архив"
	LOG_WARN_25106 = "Не удалось открыть файл для определения его типа"
	LOG_WARN_25108 = "Во время подключения к pop3 возникли ошибки"
	LOG_WARN_25109 = "Во время авторизации возникли ошибки "
	LOG_WARN_25110 = "Ошибка при получении писем из базы данных "
	LOG_WARN_25111 = "Ошибка при получении писем из базы данных "
	LOG_WARN_25112 = "Ошибки при получении всех UIDL писем "
	LOG_WARN_25113 = "Не удалось получить сообщение в базе данных для сохранения "
	LOG_WARN_25114 = "Не удалось сохранить сообщение в базе данных"
	LOG_WARN_25115 = "При получении сообщений из хранилища возникли ошибки"
	LOG_WARN_25116 = "Не удалось посчитать количество сообщений "
	LOG_WARN_25117 = "Не найдено писем в почтовом ящике."
)


var (
	GET_MAILS_FROM_DB	= false // Были ли получены письма из базы данных
	all_mails = make(MailsSlice, 0, 0) // Сообщения из базы данных
	storage_mails = make(StorageMailsSlice, 0, 0) // Временные сообщения для записи в бд
	chunk	= 100 // Количество получение писем из бд за один раз
)


// Сообщения из почты
type MailsSlice []*model.Popmail
type StorageMailsSlice []*model.Popmail

// Результат добавления/обновления в Таблицы
type tableResult struct {
	Saved 				map[string]int64
	Updated 			map[string]int64
	Errors 				map[string]int64
}

type reader struct {
	file    *os.File
	headers []string
	reader  *csv.Reader
}

func (p *AdvancedParsing) ToJson() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func AdvancedParsingFromJson(data io.Reader) *AdvancedParsing {
	var o *AdvancedParsing
	json.NewDecoder(data).Decode(&o)
	return o
}

func (p *AdvancedParsing) IsValid() bool {
	switch p.CurrentTable {
	case "ParseFiles":
	default:
		return false
	}

	return true
}

func (worker *Worker) runAdvancedParsing(lastDone string) (bool, string, *model.AppError) {
	var progress *AdvancedParsing
	a := worker.app
	if len(lastDone) == 0 {
		// Haven't started the rufr yet.
		progress = new(AdvancedParsing)
		progress.CurrentTable = "ParseFiles"
	} else {
		progress = AdvancedParsingFromJson(strings.NewReader(lastDone))
		if !progress.IsValid() {
			return false, "", model.NewAppError("FinHelpWorker.runAdvancedParsing", "finhelp.worker.run_advanced_parsing.invalid_progress",
				map[string]interface{}{"progress": progress.ToJson()}, "", http.StatusInternalServerError)
		}
	}

	if progress.CurrentTable == "ParseFiles" {
		// Run a ParseFiles finhelp batch
		a.CLog(25007, LOG_INFO_25007)
		defer func() {
			if err := recover(); err != nil {
				a.CLogErr(25001, LOG_ERR_25001, err)
			}
		}()


		// Получение писем с почты
		isMail := worker.getMails()
		if isMail {
			// Не было ошибок
			return true, progress.ToJson(), nil
		} else {
			// Произошли ошибки
			a.CLogWarning(25076, LOG_WARN_25076)
			return false, progress.ToJson(), nil
		}
	}

	return false, progress.ToJson(), nil
}


// Получение писем с почты
func (worker *Worker) getMails() (bool) {
	a := worker.app
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()


	c := worker.pop3Auth()

	a.CLog(25008, LOG_INFO_25008)

	// Получаем количество писем
	worker.getCountMessages(c)

	// Проверка на получение писем из БД
	worker.checkMailsInDB()

	// Проверяем на наличие новых сообщений
	mailArray := worker.checkNewMails(c)

	if len(mailArray) > 0 {
		// Есть новые сообщения
		a.CLog(25009, LOG_INFO_25009)

		for _, mail := range mailArray {
			// Получаем информацию о письме
			newMail, err := c.Top(mail.Seq, 25)
			if err != nil {
				a.CLogWarning(25077, LOG_WARN_25077, err.Error())
			} else {
				if len(newMail) > 0 {
					worker.getInfoAboutMail(mail, newMail)
				} else {
					a.CLog(25010, LOG_INFO_25010)
				}
			}
		}
		if len(storage_mails) > 0 {
			a.CLog(25011, LOG_INFO_25011)

			// Обрабатываем письма из storage
			res := worker.getMessagesFromStorage(c)

			// Очищаем storage
			storage_mails = storage_mails[:0]

			a.CLog(25012, LOG_INFO_25012)
			if res {
				return true
			} else {
				return false
			}
		} else {
			a.CLog(25013, LOG_INFO_25013)
		}
	} else {
		a.CLog(25014, LOG_INFO_25014)
	}
	return true
}

// Парсинг Микрофин. документа
func parseMicrofin(worker *Worker, target string) (bool) {
	// check for exist
	a := worker.app
	if _, err := os.Stat(target); err != nil {
		a.CLogErr(25002, LOG_ERR_25002, target)
		return false
	}

	saved := 0
	updated := 0
	errors := 0
	a.CLog(25015, LOG_INFO_25015, target)

	err := ParseEach(target, model.Micro{}, func(v interface{}) {
		//Cast interface to our type
		fmicro := v.(model.Micro)
		var rfin *model.FinHelp
		var err *model.AppError

		// Проверяем на валидность данных
		isValidRow := 0
		if(len(fmicro.Num) > 1) {
			isValidRow = 1
		} else if(len(fmicro.Inn) > 1) {
			isValidRow = 2
		}

		if isValidRow > 0 {
			if isValidRow == 1 {
				rfin, err = worker.app.GetOneByNum(fmicro.Num)
				if err != nil {
					errors++
				}
			} else {
				rfin, err = worker.app.GetOneByInn(fmicro.Inn)
				if err != nil {
					errors++
				}
			}
			if len(rfin.Id) < 3 {
				// save
				rfin = model.NormalizeMicro(fmicro)
				rfin.Type = "micro"
				rfin, err = worker.app.SaveInOrg(rfin)
				if err != nil {
					errors++
				} else {
					saved++;
				}
			} else {
				// update
				Id := rfin.Id
				rfin = model.NormalizeMicro(fmicro)
				rfin.Id = Id
				rfin, err = worker.app.UpdateInOrg(rfin)
				if err != nil {
					errors++
				} else {
					updated++;
				}
			}
		} else {
			errors++
		}

	})


	if errors > 0 {
		a.CLogWarning(25078, LOG_WARN_25078, saved, updated, errors)
	} else {
		a.CLog(25016, LOG_INFO_25016, saved, updated, errors)
	}
	if err != nil {
		a.CLogErr(25003, LOG_ERR_25003, target, err.Error())
		return false
	} else {
		a.CLog(25017, LOG_INFO_25017, target)
		return true
	}
}

// Сохранение Записи Insurance
func (worker *Worker) saveRowInsurance(rfin *model.FinHelp, row model.Insurance, tresult *tableResult) (*tableResult) {
	id := rfin.Id
	// Добавление в таблицу Org

	record2 := model.NormalizeInsurance(row)

	if len(id) > 2 {
		// Обновление в Org
		record2.Id = id
		_, err := worker.app.UpdateInOrg(record2)
		if err != nil {
			tresult.Errors["Org"]++
		} else {
			tresult.Updated["Org"]++
			row.Id = id
		}

	} else {
		// Сохранение в Org
		record2.Type = "insurance"
		record2, err := worker.app.SaveInOrg(record2)
		if err != nil {
			tresult.Errors["Org"]++
		} else {
			tresult.Saved["Org"]++
			row.Id = record2.Id
		}
	}

	record := row
	record.Org_id = row.Id

	// Добавление в таблицу Org_insurance_lics
	tres := addInOrgInsuranceLics(worker, record, tresult)
	if tres != nil {
		tresult = tres
	}
	// Добавление в таблицу Org_emails
	tres = addInOrgEmails(worker, record, tresult)
	if tres != nil {
		tresult = tres
	}
	// Добавление в таблицу Org_faxs
	tres = addInOrgFaxs(worker, record, tresult)
	if tres != nil {
		tresult = tres
	}
	// Добавление в таблицу Org_phones
	tres = addInOrgPhones(worker, record, tresult)
	if tres != nil {
		tresult = tres
	}
	// Добавление в таблицу Org_websites
	tres = addInOrgWebsites(worker, record, tresult)
	if tres != nil {
		tresult = tres
	}

	return tresult
}

// Парсинг страховых компаний документа
func parseInsurance(worker *Worker, target string) (bool) {
	// check for exist
	if _, err := os.Stat(target); err != nil {
		return false
	}
	a := worker.app
	cfg := a.Config()
	var forceUpdate *bool
	forceUpdate = cfg.ExchangeSettings.ForceUpdate

	if *forceUpdate == true {
		if _, err1 := worker.app.ClearAllPhones(); err1 != nil {
			a.CLog(25018, LOG_INFO_25018)
			return false
		}
		if _, err1 := worker.app.ClearAllEmails(); err1 != nil {
			a.CLog(25019, LOG_INFO_25019)
			return false
		}
		if _, err1 := worker.app.ClearAllFaxs(); err1 != nil {
			a.CLog(25020, LOG_INFO_25020)
			return false
		}
		if _, err1 := worker.app.ClearAllWebsites(); err1 != nil {
			a.CLog(25021, LOG_INFO_25021)
			return false
		}
		if _, err1 := worker.app.ClearAllInsurance(); err1 != nil {
			a.CLog(25022, LOG_INFO_25022)
			return false
		}
	}

	// Результат добавления записей в таблицы
	tresult := &tableResult{
		Saved: map[string]int64{
			"Org": 0,
			"Org_insurance_lics": 0,
			"Org_emails": 0,
			"Org_faxs": 0,
			"Org_phones": 0,
			"Org_websites": 0,
		},
		Updated: map[string]int64{
			"Org": 0,
			"Org_insurance_lics": 0,
			"Org_emails": 0,
			"Org_faxs": 0,
			"Org_phones": 0,
			"Org_websites": 0,
		},
		Errors: map[string]int64{
			"Org": 0,
			"Org_insurance_lics": 0,
			"Org_emails": 0,
			"Org_faxs": 0,
			"Org_phones": 0,
			"Org_websites": 0,
		},
	}

	isAllValid := true

	a.CLog(25023, LOG_INFO_25023, target)

	err := ParseEach(target, model.Insurance{}, func(v interface{}) {
		//Cast interface to our type
		finsurance := v.(model.Insurance)
		var rfin *model.FinHelp
		var err *model.AppError

		// Проверяем на валидность данных
		isValidRow := 0
		if(len(finsurance.Num) > 1) {
			isValidRow = 1
		} else if(len(finsurance.Inn) > 1) {
			isValidRow = 2
		}
		if(len(finsurance.Name) < 3) {
			isValidRow = 0
		}

		if isValidRow > 0 {
			if isValidRow == 1 {
				rfin, err = worker.app.GetOneByNum(finsurance.Num);
			} else {
				rfin, err = worker.app.GetOneByInn(finsurance.Inn);
			}
			if err != nil {
				isAllValid = false
			} else {
				tres := worker.saveRowInsurance(rfin, finsurance, tresult)
				if tres != nil {
					tresult = tres
				}
			}
		} else {
			isAllValid = false
		}
	})

	if !isAllValid {
		a.CLog(25024, LOG_INFO_25024, target)
	}

	if tresult.Errors["Org"] > 0 {
		a.CLogWarning(25079, LOG_WARN_25079, tresult.Saved["Org"], tresult.Updated["Org"], tresult.Errors["Org"])
	} else {
		a.CLog(25025, LOG_INFO_25025, tresult.Saved["Org"], tresult.Updated["Org"], tresult.Errors["Org"])
	}
	if tresult.Errors["Org_insurance_lics"] > 0 {
		a.CLogWarning(25080, LOG_WARN_25080, tresult.Saved["Org_insurance_lics"],
			tresult.Updated["Org_insurance_lics"], tresult.Errors["Org_insurance_lics"])
	} else {
		a.CLog(25026, LOG_INFO_25026, tresult.Saved["Org_insurance_lics"],
			tresult.Updated["Org_insurance_lics"], tresult.Errors["Org_insurance_lics"])
	}
	if tresult.Errors["Org_emails"] > 0 {
		a.CLogWarning(25081, LOG_WARN_25081, tresult.Saved["Org_emails"], tresult.Updated["Org_emails"], tresult.Errors["Org_emails"])
	} else {
		a.CLog(25027, LOG_INFO_25027, tresult.Saved["Org_emails"], tresult.Updated["Org_emails"], tresult.Errors["Org_emails"])
	}
	if tresult.Errors["Org_faxs"] > 0 {
		a.CLogWarning(25082, LOG_WARN_25082, tresult.Saved["Org_faxs"], tresult.Updated["Org_faxs"], tresult.Errors["Org_faxs"])
	} else {
		a.CLog(25028, LOG_INFO_25028, tresult.Saved["Org_faxs"], tresult.Updated["Org_faxs"], tresult.Errors["Org_faxs"])
	}
	if tresult.Errors["Org_phones"] > 0 {
		a.CLogWarning(25083, LOG_WARN_25083, tresult.Saved["Org_phones"], tresult.Updated["Org_phones"], tresult.Errors["Org_phones"])
	} else {
		a.CLog(25029, LOG_INFO_25029, tresult.Saved["Org_phones"], tresult.Updated["Org_phones"], tresult.Errors["Org_phones"])
	}
	if tresult.Errors["Org_websites"] > 0 {
		a.CLogWarning(25084, LOG_WARN_25084, tresult.Saved["Org_websites"], tresult.Updated["Org_websites"], tresult.Errors["Org_websites"])
	} else {
		a.CLog(25030, LOG_INFO_25030, tresult.Saved["Org_websites"], tresult.Updated["Org_websites"], tresult.Errors["Org_websites"])
	}

	if err != nil {
		a.CLogWarning(25085, LOG_WARN_25085, target, err.Error())
		return false
	} else {
		a.CLog(25031, LOG_INFO_25031, target)
		return true
	}
}

// Добавление или обновление страховой компании в Org_insurance_lics
func addInOrgInsuranceLics(worker *Worker, frecord model.Insurance, tresult *tableResult) (*tableResult) {
	var record *model.Org_insurance_lics
	record, err := worker.app.GetInsuranceLicsByOrgIdAndType(frecord.Id, frecord.Lic_type_code)
	if err != nil {
		tresult.Errors["Org_insurance_lics"]++
		return tresult
	}
	isValid := false
	if len(frecord.Name) > 2 && len(frecord.Num) > 1 && record != nil {
		isValid = true
	}
	if isValid {
		record2 := model.NormalizeOrigInsurance(frecord)
		if len(record.Id) > 2 {
			// update
			Id := record.Id
			record2.Id = Id
			record2.Org_id = frecord.Id
			record2, err = worker.app.UpdateInOrgInsuranceLics(record2)
			if err != nil {
				tresult.Errors["Org_insurance_lics"]++
			} else {
				tresult.Updated["Org_insurance_lics"]++
			}
		} else {
			// save
			record2.Id = ""
			record2.Org_id = frecord.Id
			record2, err = worker.app.SaveInOrgInsuranceLics(record2)
			if err != nil {
				tresult.Errors["Org_insurance_lics"]++
			} else {
				tresult.Saved["Org_insurance_lics"]++
			}
		}
	} else {
		tresult.Errors["Org_insurance_lics"]++
	}
	return tresult
}

// Добавление или обновление страховой компании в Org_emails
func addInOrgEmails(worker *Worker, frecord model.Insurance, tresult *tableResult) (*tableResult) {
	var record *model.Org_emails
	record, err := worker.app.GetEmailsByOrgId(frecord.Id);
	if err != nil {
		tresult.Errors["Org_emails"]++
		return tresult
	}
	if record != nil {
		if len(record.Id) > 2 {
			// update
			Id := record.Id
			record2 := model.NormalizeOrgEmails(frecord)
			record2.Id = Id
			record2.Org_id = frecord.Id
			record2, err = worker.app.UpdateInOrgEmails(record2)
			if err != nil {
				tresult.Errors["Org_emails"]++
			} else {
				tresult.Updated["Org_emails"]++
			}
		} else {
			// save
			record2 := model.NormalizeOrgEmails(frecord)
			record2.Id = ""
			record2.Org_id = frecord.Id
			record2, err = worker.app.SaveInOrgEmails(record2)
			if err != nil {
				tresult.Errors["Org_emails"]++
			} else {
				tresult.Saved["Org_emails"]++
			}
		}
	}
	return tresult
}

// Добавление или обновление страховой компании в Org_Phones
func addInOrgPhones(worker *Worker, frecord model.Insurance, tresult *tableResult) (*tableResult) {
	var record *model.Org_phones
	record, err := worker.app.GetPhonesByOrgId(frecord.Id)
	if err != nil {
		tresult.Errors["Org_phones"]++
		return tresult
	}
	if record != nil {
		if len(record.Id) > 2 {
			// update
			Id := record.Id
			record2 := model.NormalizeOrgPhones(frecord)
			record2.Id = Id
			record2.Org_id = frecord.Id
			record2, err = worker.app.UpdateInOrgPhones(record2)
			if err != nil {
				tresult.Errors["Org_phones"]++
			} else {
				tresult.Updated["Org_phones"]++
			}
		} else {
			// save
			record2 := model.NormalizeOrgPhones(frecord)
			record2.Id = ""
			record2.Org_id = frecord.Id
			record2, err = worker.app.SaveInOrgPhones(record2)
			if err != nil {
				tresult.Errors["Org_phones"]++
			} else {
				tresult.Saved["Org_phones"]++
			}
		}
	}
	return tresult
}

// Добавление или обновление страховой компании в Org_Faxs
func addInOrgFaxs(worker *Worker, frecord model.Insurance, tresult *tableResult) (*tableResult) {
	var record *model.Org_faxs
	record, err := worker.app.GetFaxsByOrgId(frecord.Id)
	if err != nil {
		tresult.Errors["Org_faxs"]++
		return tresult
	}
	if record != nil {
		if len(record.Id) > 2 {
			// update
			Id := record.Id
			record2 := model.NormalizeOrgFaxs(frecord)
			record2.Id = Id
			record2.Org_id = frecord.Id
			record2, err = worker.app.UpdateInOrgFaxs(record2)
			if err != nil {
				tresult.Errors["Org_faxs"]++
			} else {
				tresult.Updated["Org_faxs"]++
			}
		} else {
			// save
			record2 := model.NormalizeOrgFaxs(frecord)
			record2.Id = ""
			record2.Org_id = frecord.Id
			record2, err = worker.app.SaveInOrgFaxs(record2)
			if err != nil {
				tresult.Errors["Org_faxs"]++
			} else {
				tresult.Saved["Org_faxs"]++
			}
		}
	}
	return tresult
}

// Добавление или обновление страховой компании в Org_Websites
func addInOrgWebsites(worker *Worker, frecord model.Insurance, tresult *tableResult) (*tableResult) {
	var record *model.Org_websites
	record, err := worker.app.GetWebsitesByOrgId(frecord.Id)
	if err != nil {
		tresult.Errors["Org_websites"]++
		return tresult
	}
	if record != nil {
		if len(record.Id) > 2 {
			// update
			Id := record.Id
			record2 := model.NormalizeOrgWebsites(frecord)
			record2.Id = Id
			record2.Org_id = frecord.Id
			record2, err = worker.app.UpdateInOrgWebsites(record2)
			if err != nil {
				tresult.Errors["Org_websites"]++
			} else {
				tresult.Updated["Org_websites"]++
			}
		} else {
			// save
			record2 := model.NormalizeOrgWebsites(frecord)
			record2.Id = ""
			record2.Org_id = frecord.Id
			record2, err = worker.app.SaveInOrgWebsites(record2)
			if err != nil {
				tresult.Errors["Org_websites"]++
			} else {
				tresult.Saved["Org_websites"]++
			}
		}
	}
	return tresult
}

// Добавление или обновление банковской компании в Org_Faxs
func addBankInOrgFaxs(worker *Worker, frecord *model.FinHelp, tresult *tableResult) (*tableResult) {
	var record *model.Org_faxs
	record, err := worker.app.GetFaxsByOrgId(frecord.Id)
	if err != nil {
		tresult.Errors["Org_faxs"]++
		return tresult
	}
	if record != nil {
		if len(record.Id) > 2 {
			// update
			record.Number = frecord.Fax
			record.Org_id = frecord.Id
			record, err = worker.app.UpdateInOrgFaxs(record)
			if err != nil {
				tresult.Errors["Org_faxs"]++
			} else {
				tresult.Updated["Org_faxs"]++
			}
		} else {
			// save
			record.Number = frecord.Fax
			record.Org_id = frecord.Id
			record, err = worker.app.SaveInOrgFaxs(record)
			if err != nil {
				tresult.Errors["Org_faxs"]++
			} else {
				tresult.Saved["Org_faxs"]++
			}
		}
	}
	return tresult
}

// Добавление или обновление банковской компании в Org_lics
func addBankInOrgLics(worker *Worker, frecord *model.FinHelp, tresult *tableResult) (*tableResult) {
	var record *model.Org_lics

	record = &model.Org_lics{}
	record.Org_id = frecord.Id

	aNull := "0"
	aNotNull := "1"
	a1 := aNull
	a2 := aNull
	a3 := aNull
	a4 := aNull
	a5 := aNull
	a6 := aNull
	a7 := aNull
	a8 := aNull

	if len(frecord.Pric_otz) < 1 {
		// not limit
		if len(frecord.Date_fiz) > 0 && len(frecord.Ogran) < 1 {
			a1 = aNotNull
		}
		if len(frecord.Lic_val) > 0 && len(frecord.Ogran2) < 1 {
			a2 = aNotNull
		}
		a3 = aNotNull
		a4 = aNotNull
		if len(frecord.Date_fiz) > 0 && len(frecord.Lic_val ) > 0 && len(frecord.Ogran) < 1 {
			a5 = aNotNull
		}
		if len(frecord.Lic_val ) > 0 {
			a6 = aNotNull
		}
		if len(frecord.Date_fiz ) > 0 && len(frecord.Lic_gold1 ) > 0 {
			a7 = aNotNull
		}
		if len(frecord.Lic_val ) > 0 {
			a8 = aNotNull
		}
	}
	var records []*model.Org_lics
	records = append(records, &model.Org_lics{
		Name: "Депозиты в рублях",
		Type: "deposits_in_rub",
		Allow: a1,
	})
	records = append(records, &model.Org_lics{
		Name: "Обмен валюты",
		Type: "exchenge_currency",
		Allow: a2,
	})
	records = append(records, &model.Org_lics{
		Name: "Банковские переводы в рублях",
		Type: "bank_transfers_in_rub",
		Allow: a3,
	})
	records = append(records, &model.Org_lics{
		Name: "Кредиты в рублях",
		Type: "credit_in_rub",
		Allow: a4,
	})
	records = append(records, &model.Org_lics{
		Name: "Депозиты в валюте",
		Type: "deposits_in_currency",
		Allow: a5,
	})
	records = append(records, &model.Org_lics{
		Name: "Кредиты в валюте",
		Type: "credit_in_currency",
		Allow: a6,
	})
	records = append(records, &model.Org_lics{
		Name: "Депозиты в драгоценных металлах",
		Type: "deposits_in_precious_metals",
		Allow: a7,
	})
	records = append(records, &model.Org_lics{
		Name: "Банковские переводы в валюте",
		Type: "bank_transfers_in_currency",
		Allow: a8,
	})

	for _, record := range records {
		record.Org_id = frecord.Id
		tres := addRecordInOrgLics(worker, *record, tresult)
		if tres != nil {
			tresult = tres
		}
	}
	return tresult
}

// Добавление или обновление записи в Org_lics
func addRecordInOrgLics(worker *Worker, frecord model.Org_lics, tresult *tableResult) (*tableResult) {
	var record *model.Org_lics
	record, err := worker.app.GetLicsByOrgIdAndType(frecord.Org_id, frecord.Type)
	if err != nil {
		tresult.Errors["Org_lics"]++
		return tresult
	}

	if record != nil {
		if len(record.Id) > 2 {
			// update
			record.Type = frecord.Type
			record.Name = frecord.Name
			record.Allow = frecord.Allow
			record.Org_id = frecord.Org_id
			record, err = worker.app.UpdateInOrgLics(record)
			if err != nil {
				tresult.Errors["Org_lics"]++
			} else {
				tresult.Updated["Org_lics"]++
			}
		} else {
			// save
			record.Type = frecord.Type
			record.Name = frecord.Name
			record.Allow = frecord.Allow
			record.Org_id = frecord.Org_id
			record, err = worker.app.SaveInOrgLics(record)
			if err != nil {
				tresult.Errors["Org_lics"]++
			} else {
				tresult.Saved["Org_lics"]++
			}
		}
	}
	return tresult
}

// Сохранение Записи Bank
func (worker *Worker) saveRowBank(rfin *model.FinHelp, tresult *tableResult) (*tableResult) {
	id := rfin.Id

	// Добавление в таблицу Org
	if len(id) > 2 {
		// Обновление
		_, err := worker.app.UpdateBankInOrg(rfin)
		if err != nil {
			tresult.Errors["Org"]++
		} else {
			tresult.Updated["Org"]++
		}
	} else {
		// Сохранение
		rfin2, err := worker.app.SaveInOrg(rfin)
		if err != nil {
			tresult.Errors["Org"]++
		} else {
			tresult.Saved["Org"]++
			rfin.Id = rfin2.Id
		}
	}

	// Добавление в таблицу Org_faxs
	tres := addBankInOrgFaxs(worker, rfin, tresult)
	if tres != nil {
		tresult = tres
	}
	// Добавление в таблицу Org_lics всех лицензий
	tres = addBankInOrgLics(worker, rfin, tresult)
	if tres != nil {
		tresult = tres
	}

	return tresult
}

// Сохранение Записи Barx
func (worker *Worker) saveRowBarx(rfin *model.FinHelp, tresult *tableResult) (*tableResult) {
	id := rfin.Id

	// Добавление в таблицу Org
	if len(id) > 2 {
		// Обновление
		_, err := worker.app.UpdateBankInOrg(rfin)
		if err != nil {
			tresult.Errors["Org"]++
		} else {
			tresult.Updated["Org"]++
		}
	} else {
		// Сохранение
		rfin2, err := worker.app.SaveInOrg(rfin)
		if err != nil {
			tresult.Errors["Org"]++
		} else {
			tresult.Saved["Org"]++
			rfin.Id = rfin2.Id
		}
	}

	return tresult
}

// Парсинг Банков
func parseBank(worker *Worker, target string) (bool) {
	defer func() {
		if err := recover(); err != nil {
			panic("При парсинге сообщений возникли ошибки ")
		}
	}()
	a := worker.app

	// check for exist
	if _, err := os.Stat(target); err != nil {
		a.CLog(25032, LOG_INFO_25032, target)
		return false
	}
	cfg := a.Config()
	var forceUpdate *bool
	forceUpdate = cfg.ExchangeSettings.ForceUpdate

	if *forceUpdate == true {
		if _, err1 := worker.app.ClearAllPhones(); err1 != nil {
			a.CLog(25033, LOG_INFO_25033)
			return false
		}
		a.CLog(25034, LOG_INFO_25034)
		if _, err1 := worker.app.ClearAllLics(); err1 != nil {
			a.CLog(25035, LOG_INFO_25035)
			return false
		}
		a.CLog(25036, LOG_INFO_25036)
		if _, err1 := worker.app.ClearAllFaxs(); err1 != nil {
			a.CLog(25037, LOG_INFO_25037)
			return false
		}
		a.CLog(25038, LOG_INFO_25038)
	}

	dbfTable, err := godbf.NewFromFile(target, "CP866")
	if err != nil {
		a.CLog(25039, LOG_INFO_25039)
		return false
	}
	var BancList []model.FinHelp
	if dbfTable != nil {
		BancList = make([]model.FinHelp, dbfTable.NumberOfRecords())
	} else {
		a.CLog(25040, LOG_INFO_25040)
		return false
	}

	// Результат добавления записей в таблицы
	tresult := &tableResult{
		Saved: map[string]int64{
			"Org": 0,
			"Org_faxs": 0,
			"Org_lics": 0,
		},
		Updated: map[string]int64{
			"Org": 0,
			"Org_faxs": 0,
			"Org_lics": 0,
		},
		Errors: map[string]int64{
			"Org": 0,
			"Org_faxs": 0,
			"Org_lics": 0,
			"Model": 0,
		},
	}

	a.CLog(25041, LOG_INFO_25041, target)

	isValidRow := false
	fields := dbfTable.FieldNames() // Поля в файле

	// Проверка на валидность, что это действительно BANC.DBF
	for _, fieldName := range fields {
		if inStr(fieldName, "BVKEY") {
			isValidRow = true
			break
		}
	}

	for i := 0; i < dbfTable.NumberOfRecords(); i++ {
		if isValidRow {
			BancList[i] = model.FinHelp{}
			BancList[i].Type = "bank"
			BancList[i].Subtype = "archived"
			BancList[i].Bvkey = dbfTable.FieldValue(i, 0)
			BancList[i].Cp = dbfTable.FieldValue(i, 1)
			BancList[i].P = dbfTable.FieldValue(i, 2)
			BancList[i].Num = dbfTable.FieldValue(i, 3)
			BancList[i].U_o = dbfTable.FieldValue(i, 4)
			BancList[i].Gap = dbfTable.FieldValue(i, 5)
			BancList[i].Status = dbfTable.FieldValue(i, 6)
			BancList[i].Tip = dbfTable.FieldValue(i, 7)
			BancList[i].Namemax = dbfTable.FieldValue(i, 8)
			BancList[i].Namemax1 = dbfTable.FieldValue(i, 9)
			BancList[i].Name = dbfTable.FieldValue(i, 10)
			BancList[i].Namer = dbfTable.FieldValue(i, 11)
			BancList[i].Uf_old = strings.Replace(dbfTable.FieldValue(i, 12),".",",",-1)
			BancList[i].Ust_f = strings.Replace(dbfTable.FieldValue(i, 13),".",",",-1)
			BancList[i].Ust_fi = strings.Replace(dbfTable.FieldValue(i, 14),".",",",-1)
			BancList[i].Ust_fb = strings.Replace(dbfTable.FieldValue(i, 15),".",",",-1)
			BancList[i].Kol_i = dbfTable.FieldValue(i, 16)
			BancList[i].Kol_b = dbfTable.FieldValue(i, 17)
			BancList[i].Kolf = dbfTable.FieldValue(i, 18)
			BancList[i].Kolfd = dbfTable.FieldValue(i, 19)
			BancList[i].Data_reg = DbfDateFormat(dbfTable.FieldValue(i, 20))
			BancList[i].Mesto = dbfTable.FieldValue(i, 21)
			BancList[i].Data_preg = DbfDateFormat(dbfTable.FieldValue(i, 22))
			BancList[i].Data_izmud = DbfDateFormat(dbfTable.FieldValue(i, 23))
			BancList[i].Regn = dbfTable.FieldValue(i, 24)
			BancList[i].Priz = dbfTable.FieldValue(i, 25)
			BancList[i].Qq = DbfDateFormat(dbfTable.FieldValue(i, 26))
			BancList[i].Lic_gold = DbfDateFormat(dbfTable.FieldValue(i, 27))
			BancList[i].Lic_gold1 = DbfDateFormat(dbfTable.FieldValue(i, 28))
			BancList[i].Lic_rub = DbfDateFormat(dbfTable.FieldValue(i, 29))
			BancList[i].Ogran = dbfTable.FieldValue(i, 30)
			BancList[i].Ogran1 = dbfTable.FieldValue(i, 31)
			BancList[i].Lic_val = DbfDateFormat(dbfTable.FieldValue(i, 32))
			BancList[i].Vv = dbfTable.FieldValue(i, 33)
			BancList[i].Rei = dbfTable.FieldValue(i, 34)
			BancList[i].No = dbfTable.FieldValue(i, 35)
			BancList[i].Opp = dbfTable.FieldValue(i, 36)
			BancList[i].Ko = dbfTable.FieldValue(i, 37)
			BancList[i].Kp = dbfTable.FieldValue(i, 38)
			BancList[i].Adres = dbfTable.FieldValue(i, 39)
			BancList[i].Adres1 = dbfTable.FieldValue(i, 40)
			BancList[i].Telefon = dbfTable.FieldValue(i, 41)
			BancList[i].Fax = dbfTable.FieldValue(i, 42)
			BancList[i].Fio_pr_pr = dbfTable.FieldValue(i, 43)
			BancList[i].Fio_gl_b = dbfTable.FieldValue(i, 44)
			BancList[i].Fio_zam_p = dbfTable.FieldValue(i, 45)
			BancList[i].Fio_zam_g = dbfTable.FieldValue(i, 46)
			BancList[i].Data_otz = DbfDateFormat(dbfTable.FieldValue(i, 47))
			BancList[i].Pric_otz = dbfTable.FieldValue(i, 48)
			BancList[i].Data_pri = DbfDateFormat(dbfTable.FieldValue(i, 49))
			BancList[i].Date_prb = DbfDateFormat(dbfTable.FieldValue(i, 50))
			BancList[i].Type_prb = dbfTable.FieldValue(i, 51)
			BancList[i].Date_nam = DbfDateFormat(dbfTable.FieldValue(i, 52))
			BancList[i].Date_adr = DbfDateFormat(dbfTable.FieldValue(i, 53))
			BancList[i].Date_fiz = DbfDateFormat(dbfTable.FieldValue(i, 54))
			BancList[i].Numsrf = dbfTable.FieldValue(i, 55)
			BancList[i].Okpo = dbfTable.FieldValue(i, 56)
			BancList[i].Dat_zayv = DbfDateFormat(dbfTable.FieldValue(i, 57))
			BancList[i].Egr = dbfTable.FieldValue(i, 58)
			BancList[i].Kliring = dbfTable.FieldValue(i, 59)
			BancList[i].Cb_date = DbfDateFormat(dbfTable.FieldValue(i, 60))
			BancList[i].Fiz_end = DbfDateFormat(dbfTable.FieldValue(i, 61))
			BancList[i].Num_ssv = dbfTable.FieldValue(i, 62)
			BancList[i].Data_ssv = DbfDateFormat(dbfTable.FieldValue(i, 63))
			BancList[i].Cnlic = dbfTable.FieldValue(i, 64)
			BancList[i].Type_ko = dbfTable.FieldValue(i, 65)
			BancList[i].Tip_lic = dbfTable.FieldValue(i, 66)
			BancList[i].Search = dbfTable.FieldValue(i, 10)
			BancList[i].Ogrn = dbfTable.FieldValue(i, 58)
		} else {
			tresult.Errors["Model"]++
		}
	}

	isValidAll := true
	for i := 0; i < len(BancList); i++ {
		var rfin *model.FinHelp
		var err *model.AppError

		if rfin, err = worker.app.GetOneByRegn(BancList[i].Regn); err != nil {
			tresult.Errors["Org"]++
		}
		BancList[i].Id = rfin.Id
		rfin = &BancList[i]

		// проверяем на валидность данных
		isValid := false
		if len(rfin.Name) > 2 && len(rfin.Regn) > 0 {
			isValid = true
		}
		if isValid {
			if len(BancList[i].Pric_otz) == 0 {
				BancList[i].Subtype = "unlimited"
			} else  {
				BancList[i].Subtype = "withdraw"
			}
			/*else {
				BancList[i].Subtype = "archived"
			}*/

			// Сохранение строки Bank в таблицы
			tres := worker.saveRowBank(rfin, tresult)
			if tres != nil {
				tresult = tres
			}

		} else {
			isValidAll = false
		}
	}

	// Пустые строки в файле
	if len(BancList) < 1 {
		a.CLogWarning(25086, LOG_WARN_25086, target)
	}

	if !isValidAll {
		a.CLogWarning(25087, LOG_WARN_25087, target)
	}
	if tresult.Errors["Model"] > 0 {
		a.CLogWarning(25088, LOG_WARN_25088, tresult.Errors["Model"], target)
	}

	if tresult.Errors["Org"] > 0 {
		a.CLogWarning(25089, LOG_WARN_25082, tresult.Saved["Org"], tresult.Updated["Org"], tresult.Errors["Org"])
	} else {
		a.CLog(25042, LOG_INFO_25042, tresult.Saved["Org"], tresult.Updated["Org"], tresult.Errors["Org"])
	}
	if tresult.Errors["Org_faxs"] > 0 {
		a.CLogWarning(25090, LOG_WARN_25090, tresult.Saved["Org_faxs"], tresult.Updated["Org_faxs"], tresult.Errors["Org_faxs"])
	} else {
		a.CLog(25043, LOG_INFO_25043, tresult.Saved["Org_faxs"], tresult.Updated["Org_faxs"], tresult.Errors["Org_faxs"])
	}
	if tresult.Errors["Org_lics"] > 0 {
		a.CLogWarning(25091, LOG_WARN_25091, tresult.Saved["Org_lics"], tresult.Updated["Org_lics"], tresult.Errors["Org_lics"])
	} else {
		a.CLog(25044, LOG_INFO_25044, tresult.Saved["Org_lics"], tresult.Updated["Org_lics"], tresult.Errors["Org_lics"])
	}

	if err != nil {
		a.CLogWarning(25092, LOG_WARN_25092, target, err.Error())
		return false
	} else {
		a.CLog(25045, LOG_INFO_25045, target)
		return true
	}
}

// Парсинг Банков
func parseBarx(worker *Worker, target string) (bool) {
	a := worker.app
	defer func() {
		if err := recover(); err != nil {
			panic("При парсинге сообщений возникли ошибки ")
		}
	}()

	// check for exist
	if _, err := os.Stat(target); err != nil {
		a.CLog(25046, LOG_INFO_25046, target)
		return false
	}
	cfg := a.Config()
	var forceUpdate *bool
	forceUpdate = cfg.ExchangeSettings.ForceUpdate

	if *forceUpdate == true {
		if _, err1 := worker.app.ClearAllPhones(); err1 != nil {
			a.CLog(25047, LOG_INFO_25047)
			return false
		}
		a.CLog(25048, LOG_INFO_25048)
		if _, err1 := worker.app.ClearAllLics(); err1 != nil {
			a.CLog(25049, LOG_INFO_25049)
			return false
		}
		a.CLog(25050, LOG_INFO_25050)
		if _, err1 := worker.app.ClearAllFaxs(); err1 != nil {
			a.CLog(25051, LOG_INFO_25051)
			return false
		}
		a.CLog(25052, LOG_INFO_25052)
	}

	dbfTable, err := godbf.NewFromFile(target, "CP866")
	if err != nil {
		a.CLog(25053, LOG_INFO_25053)
		return false
	}
	var BancList []model.FinHelp
	if dbfTable != nil {
		BancList = make([]model.FinHelp, dbfTable.NumberOfRecords())
	} else {
		a.CLog(25054, LOG_INFO_25054)
		return false
	}

	// Результат добавления записей в таблицы
	tresult := &tableResult{
		Saved: map[string]int64{
			"Org": 0,
		},
		Updated: map[string]int64{
			"Org": 0,
		},
		Errors: map[string]int64{
			"Org": 0,
			"Model": 0,
		},
	}

	a.CLog(25055, LOG_INFO_25055, target)

	isValidRow := false
	fields := dbfTable.FieldNames() // Поля в файле
	// Проверка на валидность, что это действительно B_ARX.DBF
	for _, fieldName := range fields {
		if inStr(fieldName, "KLIRING") {
			isValidRow = true
			break
		}
	}

	for i := 0; i < dbfTable.NumberOfRecords(); i++ {
		if isValidRow {
			BancList[i] = model.FinHelp{}
			BancList[i].Type = "bank"
			BancList[i].Data_vib = DbfDateFormat(dbfTable.FieldValue(i, 0))
			BancList[i].Prik = dbfTable.FieldValue(i, 2)
			BancList[i].Cp = dbfTable.FieldValue(i, 3)
			BancList[i].P = dbfTable.FieldValue(i, 4)
			BancList[i].Num = dbfTable.FieldValue(i, 5)
			BancList[i].U_o = dbfTable.FieldValue(i, 6)
			BancList[i].Gap = dbfTable.FieldValue(i, 7)
			BancList[i].Status = dbfTable.FieldValue(i, 8)
			BancList[i].Tip = dbfTable.FieldValue(i, 9)
			BancList[i].Namemax = dbfTable.FieldValue(i, 10)
			BancList[i].Namemax1 = dbfTable.FieldValue(i, 11)
			BancList[i].Name = dbfTable.FieldValue(i, 12)
			BancList[i].Namer = dbfTable.FieldValue(i, 13)
			BancList[i].Uf_old = strings.Replace(dbfTable.FieldValue(i, 14),".",",",-1)
			BancList[i].Ust_f = strings.Replace(dbfTable.FieldValue(i, 15),".",",",-1)
			BancList[i].Ust_fi = strings.Replace(dbfTable.FieldValue(i, 16),".",",",-1)
			BancList[i].Ust_fb = strings.Replace(dbfTable.FieldValue(i, 17),".",",",-1)
			BancList[i].Kol_i = dbfTable.FieldValue(i, 18)
			BancList[i].Kol_b = dbfTable.FieldValue(i, 19)
			BancList[i].Kolf = dbfTable.FieldValue(i, 20)
			BancList[i].Kolfd = dbfTable.FieldValue(i, 21)
			BancList[i].Data_reg = DbfDateFormat(dbfTable.FieldValue(i, 22))
			BancList[i].Mesto = dbfTable.FieldValue(i, 23)
			BancList[i].Data_preg = DbfDateFormat(dbfTable.FieldValue(i, 24))
			BancList[i].Data_izmud = DbfDateFormat(dbfTable.FieldValue(i, 25))
			BancList[i].Regn = dbfTable.FieldValue(i, 26)
			BancList[i].Priz = dbfTable.FieldValue(i, 27)
			BancList[i].Qq = DbfDateFormat(dbfTable.FieldValue(i, 28))
			BancList[i].Lic_gold = DbfDateFormat(dbfTable.FieldValue(i, 29))

			BancList[i].Lic_gold1 = DbfDateFormat(dbfTable.FieldValue(i, 30))
			BancList[i].Lic_rub = DbfDateFormat(dbfTable.FieldValue(i, 31))
			BancList[i].Ogran = dbfTable.FieldValue(i, 32)
			BancList[i].Ogran1 = dbfTable.FieldValue(i, 33)
			BancList[i].Lic_val = DbfDateFormat(dbfTable.FieldValue(i, 34))
			BancList[i].Vv = dbfTable.FieldValue(i, 35)
			BancList[i].Rei = dbfTable.FieldValue(i, 36)
			BancList[i].No = dbfTable.FieldValue(i, 37)
			BancList[i].Opp = dbfTable.FieldValue(i, 38)
			BancList[i].Ko = dbfTable.FieldValue(i, 39)
			BancList[i].Kp = dbfTable.FieldValue(i, 40)
			BancList[i].Adres = dbfTable.FieldValue(i, 41)
			BancList[i].Adres1 = dbfTable.FieldValue(i, 42)
			BancList[i].Telefon = dbfTable.FieldValue(i, 43)
			BancList[i].Fax = dbfTable.FieldValue(i, 44)
			BancList[i].Fio_pr_pr = dbfTable.FieldValue(i, 45)
			BancList[i].Fio_gl_b = dbfTable.FieldValue(i, 46)
			BancList[i].Fio_zam_p = dbfTable.FieldValue(i, 47)
			BancList[i].Fio_zam_g = dbfTable.FieldValue(i, 48)
			BancList[i].Data_otz = DbfDateFormat(dbfTable.FieldValue(i, 49))
			BancList[i].Pric_otz = dbfTable.FieldValue(i, 50)
			BancList[i].Data_pri = DbfDateFormat(dbfTable.FieldValue(i, 51))
			BancList[i].Regn1 = dbfTable.FieldValue(i, 52)
			BancList[i].Date_prb = DbfDateFormat(dbfTable.FieldValue(i, 53))
			BancList[i].Type_prb = dbfTable.FieldValue(i, 54)
			BancList[i].Date_nam = DbfDateFormat(dbfTable.FieldValue(i, 55))
			BancList[i].Date_adr = DbfDateFormat(dbfTable.FieldValue(i, 56))
			BancList[i].Date_fiz = DbfDateFormat(dbfTable.FieldValue(i, 57))
			BancList[i].Numsrf = dbfTable.FieldValue(i, 58)
			BancList[i].Okpo = dbfTable.FieldValue(i, 59)
			BancList[i].Dat_zayv = DbfDateFormat(dbfTable.FieldValue(i, 60))
			BancList[i].Egr = dbfTable.FieldValue(i, 61)
			BancList[i].Prik_br = dbfTable.FieldValue(i, 62)
			BancList[i].Data_br = DbfDateFormat(dbfTable.FieldValue(i, 63))
			BancList[i].Kliring = dbfTable.FieldValue(i, 64)
			BancList[i].Cb_date = DbfDateFormat(dbfTable.FieldValue(i, 65))
			BancList[i].Fiz_end = DbfDateFormat(dbfTable.FieldValue(i, 66))
			BancList[i].Num_ssv = dbfTable.FieldValue(i, 67)
			BancList[i].Data_ssv = DbfDateFormat(dbfTable.FieldValue(i, 68))
			BancList[i].Cnlic = dbfTable.FieldValue(i, 69)
			BancList[i].Type_ko = dbfTable.FieldValue(i, 70)
			BancList[i].Tip_lic = dbfTable.FieldValue(i, 71)
			BancList[i].Search = dbfTable.FieldValue(i, 12)
		} else {
			tresult.Errors["Model"]++
		}
	}

	isValidAll := true

	for i := 0; i < len(BancList); i++ {
		var rfin *model.FinHelp
		var err *model.AppError
		if rfin, err = worker.app.GetOneByRegn(BancList[i].Regn); err != nil {
			tresult.Errors["Org"]++
		}
		BancList[i].Id = rfin.Id
		rfin = &BancList[i]
		// проверяем на валидность данных
		isValid := false
		if len(rfin.Name) > 2 && len(rfin.Regn) > 0 {
			isValid = true
		}
		if isValid {
			BancList[i].Subtype = "archived"

			// Сохранение строки Bank в таблицы
			tres := worker.saveRowBarx(rfin, tresult)
			if tres != nil {
				tresult = tres
			}
		} else {
			isValidAll = false
		}
	}

	if !isValidAll {
		a.CLogWarning(25093, LOG_WARN_25093, target)
	}
	if tresult.Errors["Model"] > 0 {
		a.CLogWarning(25094, LOG_WARN_25094, tresult.Errors["Model"], target)
	}

	if tresult.Errors["Org"] > 0 {
		a.CLogWarning(25095, LOG_WARN_25095, tresult.Saved["Org"], tresult.Updated["Org"], tresult.Errors["Org"])
	} else {
		a.CLog(25056, LOG_INFO_25056, tresult.Saved["Org"], tresult.Updated["Org"], tresult.Errors["Org"])
	}

	if err != nil {
		a.CLogWarning(25096, LOG_WARN_25096, target, err.Error())
		return false
	} else {
		a.CLog(25057, LOG_INFO_25057, target)
		return true
	}
}

// Parse Files
func (r *reader) getHeaders() (err error) {
	r.headers, err = r.reader.Read()
	if err != nil {
		return err
	}
	return
}

// ParseEach parses each line of csv file and passing interface{} to callback function
func ParseEach(file string, v interface{}, callback func(result interface{})) error {
	var err error
	r := new(reader)

	dataType := reflect.TypeOf(v)
	newData := reflect.New(dataType).Elem()
	if r.file, err = os.Open(file); err != nil {
		return err
	}
	defer r.file.Close()

	// Определение кодировки
	body, err := ioutil.ReadAll(r.file)
	if err != nil {
		l4g.Info("Ошибки во время определения кодировки: %v", err)
		return err
	}

	strBody := string(body)

	re := strings.Replace(strBody,`"`,`""`,-1)
	re2 := strings.Replace(re,`;""`,`;"`,-1)
	re3 := strings.Replace(re2,`"";`,`";`,-1)
	re4 := strings.Replace(re3,`""`+string(10),`"`+string(10),-1)
	newString := strings.Replace(re4,string(10)+`""`,string(10)+`"`,-1)
	r.reader = csv.NewReader(strings.NewReader(newString))
	r.reader.Comma = ';'
	r.reader.LazyQuotes = true
	r.reader.FieldsPerRecord = -1
	r.headers = make([]string, 0)
	err = r.getHeaders()

	if err != nil {
		return err
	}

	if len(r.headers) == 1 {
		// Возможно шапка сформирована запятыми
		strHeader := r.headers[0]
		keyHeader := strings.Split(strHeader, ",")
		r.headers = keyHeader
	}


	for {
		row, err := r.reader.Read()
		if err != nil {
			//log.Printf("ParseEach error read ", err)
			break
		}

		for i := 0; i < dataType.NumField(); i++ {
			f := dataType.Field(i)
			index := 0

			fieldName := f.Tag.Get("csv")
			for k, v := range r.headers {
				if v == fieldName {
					index = k
					break
				}
			}

			newField := newData.FieldByName(f.Name)

			if newField.IsValid() {
				if newField.CanSet() {
					if index < len(row) {
						value := reflect.ValueOf(DecodeStringWindows1251(row[index]))
						newField.Set(value)
					}
				}
			}
		}
		callback(newData.Interface())
	}
	return nil
}

// Авторизация в pop3
func (worker *Worker) pop3Auth() (*pop3.Client) {
	a := worker.app
	cfg := a.Config()

	// Dial
	c, err := pop3.Dial(*cfg.EmailSettings.Pop3Url, pop3.UseTLS(nil))
	if err != nil {
		a.CLogWarning(25108, LOG_WARN_25108)
		panic(LOG_WARN_25108 + err.Error())
	}

	a.CLog(25058, LOG_INFO_25058)

	// Login
	if err := c.Auth(*cfg.EmailSettings.Pop3Login, *cfg.EmailSettings.Pop3Password); err != nil {
		a.CLogWarning(25109, LOG_WARN_25109)
		panic(LOG_WARN_25109 + err.Error())
	}
	a.CLog(25059, LOG_INFO_25059)
	return c
}

// Получение писем из базы данных
func (worker *Worker) checkMailsInDB() {
	a := worker.app
	if !GET_MAILS_FROM_DB {
		// еще не были получены сообщения из базы данных
		// Получение писем из базы данных
		count, err := a.GetPopmailCount()
		if err != nil {
			a.CLogWarning(25110, LOG_WARN_25110)
			panic(LOG_WARN_25110 + err.Error())
		}
		// Сообщения получены
		GET_MAILS_FROM_DB = true
		if count > 0 {
			for i := 0; i < int(count); i += chunk {
				end := i + chunk

				if end > int(count) {
					end = int(count)
				}
				offset := end - i

				// Чанк писем из бд
				mails, err := a.GetChunkMails(i, offset)
				if err != nil {
					a.CLogWarning(25111, LOG_WARN_25111)
					panic(LOG_WARN_25111 + err.Error())
				}

				// Добавляем в стек
				all_mails = append(all_mails, mails...)
			}
		}
		sort.Sort(all_mails)
	}
}

// Провека на получение новых писем
func (worker *Worker) checkNewMails(c *pop3.Client) ([]*pop3.MessageInfo) {
	// Получаем ID всех сообщений
	a := worker.app
	allMails, err := c.UIDlAll()
	if err != nil {
		a.CLogWarning(25112, LOG_WARN_25112)
		panic(LOG_WARN_25112 + err.Error())
	}
	// Сортируем по uidl
	sort.SliceStable(all_mails, func(i, j int) bool {
		return all_mails[i].Uidl < all_mails[j].Uidl
	})

	// Сравниваем, нет ли новых писем
	var mailArray []*pop3.MessageInfo
	for _, mail := range allMails {
		isFound := false

		i := sort.Search(len(all_mails), func(i int) bool { return all_mails[i].Uidl >= mail.UID })
		if i < len(all_mails) && all_mails[i].Uidl == mail.UID {
			isFound = true
		} else {
			isFound = false
		}
		if !isFound {
			// Если не нашли сообщение, то оно новое
			mailArray = append(mailArray, mail)
		}
	}
	// Обратно сортируем по дате
	sort.Sort(all_mails)

	return mailArray
}

// Получаем дату последнего сообщения
func (worker *Worker) getLastMessageDate() (int64) {
	if len(all_mails) > 0 {
		return all_mails[len(all_mails)-1].Date
	}
	return 0
}

// Получаем информацию о сообщении
func (worker *Worker) getInfoAboutMail(mail *pop3.MessageInfo, newMail string) (bool) {
	// Дата последнего сообщения
	lastMailDate := worker.getLastMessageDate()
	a := worker.app
	var reader io.Reader
	reader = strings.NewReader(newMail)
	email, err := pop3.ShortParse(reader)
	if err != nil {
		a.CLog(25061, LOG_INFO_25061)
		panic(LOG_INFO_25061 + err.Error())
	}

	mailType := ""
	// Проверка на совпадение темы письма
	subject_to_parse := worker.app.Config().ExchangeSettings.Subject
	re := regexp.MustCompile(subject_to_parse)
	// проверка на совпадение со всеми кодировками
	dec := new(mime.WordDecoder)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "koi8-r":
			content, err := ioutil.ReadAll(input)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(content), nil
		default:
			return nil, fmt.Errorf("unhandled charset %q", charset)
		}
	}
	subject, subjErr := dec.Decode(email.Subject)
	if subjErr != nil {
		subject = email.Subject
	}
	l4g.Debug("Subject string before: \"",email.Subject,"\"")
	l4g.Debug("Subject string after: \"",subject,"\"")
	if re.MatchString(subject) {
		mailType = "exchange_mail"
	}
	if re.MatchString(DecodeStringKOI8R(subject)) {
		mailType = "exchange_mail"
	}
	if re.MatchString(DecodeStringCp866(subject)) {
		mailType = "exchange_mail"
	}
	if re.MatchString(DecodeStringWindows1251(subject)) {
		mailType = "exchange_mail"
	}
	newDate := email.Date.Unix()
	seq := strconv.FormatUint(uint64(mail.Seq), 10)

	if lastMailDate > newDate {
		// Это старое сообщение
		l4g.Debug("Last mail date in database: \"",lastMailDate,"\"")
		l4g.Debug("Current mail date: \"",newDate,"\"")
	} else {
		tmail := &model.Popmail{
			seq,
			mail.UID,
			mailType,
			newDate,
			0,
		}
		// Сохраняем для обработки сообщения
		storage_mails = append(storage_mails, tmail)
	}
	rmail := &model.Popmail{
		seq,
		mail.UID,
		mailType,
		newDate,
		0,
	}
	// Сохраняем в базу данных
	worker.saveMailInDB(rmail)
	return true
}

// Получение руфр сообщения
func (worker *Worker) getRufrMessageBySeq(c *pop3.Client, seq uint32) (*pop3.Email, error) {
	a := worker.app

	// Получаем данные о руфр сообщении
	newMail, err := c.Retr(seq)
	if err != nil {
		a.CLog(25060, LOG_INFO_25060, err.Error())
		return nil, err
	}

	// Получено сообщение
	var reader io.Reader
	reader = strings.NewReader(newMail)
	email, err := pop3.Parse(reader)
	if err != nil {
		a.CLog(25061, LOG_INFO_25061, err.Error())
		return nil, err
	}
	return &email, nil
}

// Сохранение сообщения в базе данных
func (worker *Worker) saveMailInDB(mail *model.Popmail) (bool) {
	a := worker.app

	// Проверяем существует ли сообщение в БД
	_, err := a.GetPopmailByUidl(mail.Uidl)
	if err != nil {
		mail.Id = ""
		_, err2 := a.SaveInPopmail(mail)
		if err2 != nil {
			a.CLog(25114, LOG_WARN_25114)
			panic(LOG_WARN_25114 + err.Error())
		}

	}
	return true
}

// Получение и обработка сообщений из временного хранилища
func (worker *Worker) getMessagesFromStorage(c *pop3.Client) (bool) {
	a := worker.app
	defer func() {
		if err := recover(); err != nil {
			a.CLog(25115, LOG_WARN_25115)
			panic(LOG_WARN_25115)
		}
	}()

	// Получаем сообщения из storage
	for _, mail := range storage_mails {
		// Сохраняем в сохраненные сообщения
		all_mails = append(all_mails, mail)

		seq, err := strconv.ParseUint(mail.Id, 10, 10)
		if err != nil {
			a.CLog(25062, LOG_INFO_25062, err.Error())
		} else {
			// Проверяем это руфр или нет
			if mail.Type == "exchange_mail" {
				email, err := worker.getRufrMessageBySeq(c, uint32(seq))
				if err != nil {

				} else {
					a.CLog(25063, LOG_INFO_25063, email.Subject)
					_, err := worker.parseRufrMail(*email)
					if err != nil {

					}
				}
			}
		}
	}
	return true
}

// Получить количество сообщений
func (worker *Worker) getCountMessages (c *pop3.Client) {
	count, _, err := c.Stat()
	a := worker.app
	if err != nil {
		a.CLog(25116, LOG_WARN_25116)
		panic(LOG_WARN_25116 + err.Error())
	}
	if count < 1 {
		// Нету писем
		a.CLog(25117, LOG_WARN_25117)
		panic(LOG_WARN_25117 + err.Error())
		panic("Не найдено писем в почтовом ящике.")
	}
}

// Парсинг руфр письма
// Возвращает результат и ошибку
func (worker *Worker) parseRufrMail(email pop3.Email) (result bool, err error) {
	a := worker.app
	cfg := a.Config().ExchangeSettings
	dir := cfg.DirToDocTemp
	a.CLog(25064, LOG_INFO_25064, email.Subject)
	// Путь до бэкап папки

	ctime := time.Now().Format("20060102_150405.000000")
	numberRand := strconv.Itoa(rand.Intn(100))
	ctime = "tmp_" + ctime + "_" + numberRand
	dir = dir + "\\" + ctime

	success := false

	if len(email.Attachments) > 0 {
		// Найдены прикрепления
		for _, g := range (email.Attachments) {
			success = false
			filename := g.Filename
			data := g.Data

			// Создание tmp папки для бэкапа файлов из почты
			target := filepath.Join(dir, filename)
			canCreate := true
			// if a dir and doesn't exist
			if _, err := os.Stat(dir); err != nil {
				if err := os.MkdirAll(dir, os.ModePerm); err != nil {
					a.CLogErr(25004, LOG_ERR_25004, err.Error())
					canCreate = false
				} else {
					a.CLog(25065, LOG_INFO_25065, dir)
				}
			}
			if canCreate {
				img, err := os.Create(target)
				if err != nil {
					a.CLogErr(25005, LOG_ERR_25005, err.Error())
					canCreate = false
				}

				if canCreate {
					_, err = io.Copy(img, data)
					if err != nil {
						a.CLogErr(25006, LOG_ERR_25006, err.Error())
						canCreate = false
					}

					if canCreate {
						// Парсим прикрепления
						res := worker.parseAttachments(g.ContentType, dir, target, 0)
						if res {
							success = true
						}
					}
				}
			}
		}
		return success, nil
	} else {
		a.CLog(25066, LOG_INFO_25066)
		return false, nil
	}


	return false, nil
}

// Парсинг прикреплений сообщения
func (worker *Worker) parseAttachments(contentType string, dir string, target string, lvl int) (bool) {
	a := worker.app
	cfg := a.Config().ExchangeSettings
	docfile_bank := cfg.Docfile_bank
	docfile_barx := cfg.Docfile_barx
	docfile_strahov := cfg.Docfile_strahov
	docfile_microfin := cfg.Docfile_microfin
	if contentType != "" {
		switch contentType {
		case "application/arj", "application/x-arj":
			// Arj архив
			// Распаковываем
			if lvl > 0 {
				a.CLogWarning(25096, LOG_WARN_25097)
				return false
			}
			res := worker.getArchive(dir, target, "arj")
			return res
			break;
		case "application/x-gzip", "application/gzip":
			// Tar архив
			// Распаковываем
			if lvl > 0 {
				a.CLogWarning(25098, LOG_WARN_25098)
				return false
			}
			res := worker.getArchive(dir, target, "tar")
			return res
			break;
		case "text/plain", "text/plain; charset=utf-8", "application/octet-stream":
			// Отдельный случай
			if inStr(target, docfile_microfin) && len(docfile_microfin) > 1 {
				return parseMicrofin(worker, target)
			} else if inStr(target, docfile_strahov) && len(docfile_strahov) > 1 {
				return parseInsurance(worker, target)
			} else if inStr(target, docfile_barx) && len(docfile_barx) > 1 {
				return parseBarx(worker, target)
			} else if inStr(target, docfile_bank) && len(docfile_bank) > 1 {
				return parseBank(worker, target)
			} else {
				// Это архив
				if inStr(target, ".arj") || inStr(target, ".ARJ") || inStr(target, "arj") ||
					inStr(target, "ARJ") {
					if lvl > 0 {
						a.CLogWarning(25099, LOG_WARN_25099)
						return false
					}
					return worker.getArchive(dir, target, "arj")
				} else if inStr(target, ".zip") || inStr(target, ".ZIP") || inStr(target, "zip") ||
					inStr(target, "ZIP") {
					if lvl > 0 {
						a.CLogWarning(25100, LOG_WARN_25100)
						return false
					}
					return worker.getArchive(dir, target, "zip")
				} else if inStr(target, ".tar") || inStr(target, ".TAR") || inStr(target, "tar") ||
					inStr(target, "TAR") {
					if lvl > 0 {
						a.CLogWarning(25101, LOG_WARN_25101)
						return false
					}
					return worker.getArchive(dir, target, "tar")
				} else {
					a.CLog(25067, LOG_INFO_25067, target)
				}
			}
			return false
			break;
		case "application/zip", "application/x-zip-compressed", "application/x-zip",
		"application/x-compress", "application/x-compressed", "multipart/x-zip":
			// Zip архив
			// Распаковываем
			if lvl > 0 {
				a.CLogWarning(25102, LOG_WARN_25102)
				return false
			}
			res := worker.getArchive(dir, target, "zip")
			return res
			break;
		case "text/csv", "application/csv", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/comma-separated-values", "application/excel", "application/vnd.msexcel", "text/anytext":
			// .csv
			if inStr(target, docfile_microfin) && len(docfile_microfin) > 1 {
				return parseMicrofin(worker, target)
			} else if inStr(target, docfile_strahov) && len(docfile_strahov) > 1 {
				return parseInsurance(worker, target)
			}
			a.CLog(25068, LOG_INFO_25068, target)
			return false
			break;
		case "application/dbf", "application/x-dbf", "application/dbase", "application/x-dbase", "zz-application/zz-winassoc-dbf":
			// .dbf
			if inStr(target, docfile_barx) && len(docfile_barx) > 1 {
				return parseBarx(worker, target)
			} else if inStr(target, docfile_bank) && len(docfile_bank) > 1 {
				return parseBank(worker, target)
			}
			a.CLog(25069, LOG_INFO_25069, target)
			return false
			break;
		default:
			a.CLog(25070, LOG_INFO_25070, contentType)
			return false
			break;
		}
	} else {
		a.CLog(25071, LOG_INFO_25071)
		return false
	}
	return true
}

func (worker *Worker) getArchive(dir string, target string, ztype string) (bool) {
	a := worker.app
	a.CLog(25072, LOG_INFO_25072)
	var files []string
	var err error
	// Определяем тип
	if (ztype == "zip") {
		// zip
		files, err = a.Unzip(dir, target)
		if err != nil {
			a.CLogWarning(25103, LOG_WARN_25103)
			return false
		}
	} else if(ztype == "tar") {
		// tar
		files, err = a.Untar(dir, target)
		if err != nil {
			a.CLogWarning(25104, LOG_WARN_25104)
			return false
		}
	} else if(ztype == "arj") {
		// arj
		files, err = a.Unarj(dir, target)
		if err != nil {
			a.CLogWarning(25105, LOG_WARN_25105)
			return false
		}
	}
	// Проверка на пустоту архива
	if len(files) < 1 {
		a.CLog(25073, LOG_INFO_25073)
	} else {
		a.CLog(25074, LOG_INFO_25074)
	}
	for _, target2 := range files {
		// Определить файл
		f, err := os.Open(target2)
		if err != nil {
			a.CLogWarning(25106, LOG_WARN_25106)
		}
		defer f.Close()
		if err == nil {
			contentType2, err := GetFileContentType(f)
			if err != nil {
				a.CLog(25075, LOG_INFO_25075)
				//return false
			} else {
				res := worker.parseAttachments(contentType2, dir, target2, 1)
				if !res {
				}
			}
		}
	}

	return true
}

// В строке
func inStr(s string, substr string) (bool) {
	re := regexp.MustCompile(substr)
	return re.MatchString(s)
}

// Определение Content type у файла
func GetFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// Len is part of sort.Interface.
func (d MailsSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d MailsSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d MailsSlice) Less(i, j int) bool {
	return d[i].Date < d[j].Date
}

// Len is part of sort.Interface.
func (d StorageMailsSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d StorageMailsSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d StorageMailsSlice) Less(i, j int) bool {
	return d[i].Date < d[j].Date
}

// Разбор кодировки строки
func DecodeString(s string) (string) {
	return s;
}

// Разбор кодировки строки
func DecodeStringCp866(s string) (string) {
	I := strings.NewReader(s)
	O := transform.NewReader(I, charmap.CodePage866.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		l4g.Debug("error with decode string")
	}
	//l4g.Info(string(d))
	return string(d)
}
func DecodeStringKOI8R(s string) (string) {
	I := strings.NewReader(s)
	O := transform.NewReader(I, charmap.KOI8R.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		l4g.Debug("error with decode string")
	}
	//l4g.Info(string(d))
	return string(d)
}
func DecodeStringWindows1251(s string) (string) {
	I := strings.NewReader(s)
	O := transform.NewReader(I, charmap.Windows1251.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		l4g.Debug("error with decode string")
	}
	//l4g.Info(string(d))
	return string(d)
}
func DbfDateFormat (s string) (string) {
	fromString := []rune(s)
	if len(fromString) != 0 {
		toString := string(fromString[6:8]) + "." + string(fromString[4:6]) + "." + string(fromString[0:4])
		return toString
	}
	return s
}