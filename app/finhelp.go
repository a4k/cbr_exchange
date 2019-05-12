package app

import (
	"../model"
	"net/http"
	"archive/tar"
	"../utils"
	"io"
	"path/filepath"
	"os"
	"compress/gzip"
	"archive/zip"
	"os/exec"
	"io/ioutil"
	"reflect"
	"strings"
)

func (a *App) CreateFinSingle(fin *model.FinHelp) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().Save(fin); result.Err != nil {
		return nil, result.Err
	} else {
		rfin := result.Data.(*model.FinHelp)

		return rfin, nil
	}
}

// Сохранение микрофин. в общую таблицу
func (a *App) SaveInOrg(fin *model.FinHelp) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().Save(fin); result.Err != nil {
		return nil, result.Err
	} else {
		rfin := result.Data.(*model.FinHelp)

		return rfin, nil
	}
}

// Сохранение страхов. в общую таблицу
func (a *App) SaveInsuranceInOrg(insurance model.Insurance) (*model.FinHelp, *model.AppError) {
	insuranceFin := model.NormalizeInsurance(insurance)
	insuranceFin.Type = "insurance"
	if result := <-a.Srv.Store.FinHelp().Save(insuranceFin); result.Err != nil {
		return nil, result.Err
	} else {
		rinsurance := result.Data.(*model.FinHelp)

		return rinsurance, nil
	}
}

// Сохранение страхов. в общую таблицу
func (a *App) SaveInOrgInsuranceLics(insurance *model.Org_insurance_lics) (*model.Org_insurance_lics, *model.AppError) {
	if result := <-a.Srv.Store.Insurance().Save(insurance); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't save the Org_insurance_lics err=%v", result.Err))
		return nil, result.Err
	} else {
		rinsurance := result.Data.(*model.Org_insurance_lics)

		return rinsurance, nil
	}
}

// Сохранение страховой компании в Org_emails
func (a *App) SaveInOrgEmails(frecord *model.Org_emails) (*model.Org_emails, *model.AppError) {
	if result := <-a.Srv.Store.OrgEmails().Save(frecord); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't Save the Org_emails err=%v", result.Err))
		return nil, result.Err
	} else {
		frecord := result.Data.(*model.Org_emails)

		return frecord, nil
	}
}

// Сохранение страховой компании в Org_faxs
func (a *App) SaveInOrgFaxs(frecord *model.Org_faxs) (*model.Org_faxs, *model.AppError) {
	if result := <-a.Srv.Store.OrgFaxs().Save(frecord); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't Save the Org_faxs err=%v", result.Err))
		return nil, result.Err
	} else {
		frecord := result.Data.(*model.Org_faxs)

		return frecord, nil
	}
}

// Сохранение страховой компании в Org_lics
func (a *App) SaveInOrgLics(frecord *model.Org_lics) (*model.Org_lics, *model.AppError) {
	if result := <-a.Srv.Store.OrgLics().Save(frecord); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't Save the Org_lics err=%v", result.Err))
		return nil, result.Err
	} else {
		frecord := result.Data.(*model.Org_lics)

		return frecord, nil
	}
}

// Сохранение страховой компании в Org_phones
func (a *App) SaveInOrgPhones(frecord *model.Org_phones) (*model.Org_phones, *model.AppError) {
	if result := <-a.Srv.Store.OrgPhones().Save(frecord); result.Err != nil {
		return nil, result.Err
	} else {
		frecord := result.Data.(*model.Org_phones)

		return frecord, nil
	}
}

// Сохранение страховой компании в Org_websites
func (a *App) SaveInOrgWebsites(frecord *model.Org_websites) (*model.Org_websites, *model.AppError) {
	if result := <-a.Srv.Store.OrgWebsites().Save(frecord); result.Err != nil {
		return nil, result.Err
	} else {
		frecord := result.Data.(*model.Org_websites)

		return frecord, nil
	}
}

// Сохранение письма в Popmail
func (a *App) SaveInPopmail(frecord *model.Popmail) (*model.Popmail, *model.AppError) {
	if result := <-a.Srv.Store.Popmail().Save(frecord); result.Err != nil {
		return nil, result.Err
	} else {
		frecord := result.Data.(*model.Popmail)

		return frecord, nil
	}
}


// Обновление микрофин. в общую таблицу
func (a *App) UpdateInOrg(fin *model.FinHelp) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().UpdateMicro(fin); result.Err != nil {
		return nil, result.Err
	} else {
		return fin, nil
	}
}

// Обновление банк. в общую таблицу
func (a *App) UpdateBankInOrg(fin *model.FinHelp) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().UpdateBank(fin); result.Err != nil {
		return nil, result.Err
	} else {
		return fin, nil
	}
}

// Обновление страховой компании в Org_insurance_lics
func (a *App) UpdateInOrgInsuranceLics(finsurance *model.Org_insurance_lics) (*model.Org_insurance_lics, *model.AppError) {
	if result := <-a.Srv.Store.Insurance().Update(finsurance); result.Err != nil {
		return nil, result.Err
	} else {
		return finsurance, nil
	}
}

// Обновление страховой компании в Org_emails
func (a *App) UpdateInOrgEmails(frecord *model.Org_emails) (*model.Org_emails, *model.AppError) {
	if result := <-a.Srv.Store.OrgEmails().Update(frecord); result.Err != nil {
		return nil, result.Err
	} else {
		return frecord, nil
	}
}

// Обновление страховой компании в Org_faxs
func (a *App) UpdateInOrgFaxs(frecord *model.Org_faxs) (*model.Org_faxs, *model.AppError) {
	if result := <-a.Srv.Store.OrgFaxs().Update(frecord); result.Err != nil {
		return nil, result.Err
	} else {
		return frecord, nil
	}
}

// Обновление страховой компании в Org_lics
func (a *App) UpdateInOrgLics(frecord *model.Org_lics) (*model.Org_lics, *model.AppError) {
	if result := <-a.Srv.Store.OrgLics().Update(frecord); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't update the insurance err=%v", result.Err))
		return nil, result.Err
	} else {
		return frecord, nil
	}
}

// Обновление страховой компании в Org_phones
func (a *App) UpdateInOrgPhones(frecord *model.Org_phones) (*model.Org_phones, *model.AppError) {
	if result := <-a.Srv.Store.OrgPhones().Update(frecord); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't update the phones err=%v", result.Err))
		return nil, result.Err
	} else {
		return frecord, nil
	}
}

// Обновление страховой компании в Org_websites
func (a *App) UpdateInOrgWebsites(frecord *model.Org_websites) (*model.Org_websites, *model.AppError) {
	if result := <-a.Srv.Store.OrgWebsites().Update(frecord); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't update the websites err=%v", result.Err))
		return nil, result.Err
	} else {
		return frecord, nil
	}
}

// Удаление в Org_phones
func (a *App) ClearAllPhones() (bool, *model.AppError) {
	if result := <-a.Srv.Store.OrgPhones().Clear(); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't clear the phones err=%v", result.Err))
		return false, result.Err
	} else {
		return true, nil
	}
}

// Удаление в Org_phones
func (a *App) ClearAllLics() (bool, *model.AppError) {
	if result := <-a.Srv.Store.OrgLics().Clear(); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't clear the lics err=%v", result.Err))
		return false, result.Err
	} else {
		return true, nil
	}
}

// Удаление в Org_phones
func (a *App) ClearAllFaxs() (bool, *model.AppError) {
	if result := <-a.Srv.Store.OrgFaxs().Clear(); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't clear the faxs err=%v", result.Err))
		return false, result.Err
	} else {
		return true, nil
	}
}

// Удаление в Org_phones
func (a *App) ClearAllEmails() (bool, *model.AppError) {
	if result := <-a.Srv.Store.OrgEmails().Clear(); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't clear the emails err=%v", result.Err))
		return false, result.Err
	} else {
		return true, nil
	}
}

// Удаление в Org_phones
func (a *App) ClearAllWebsites() (bool, *model.AppError) {
	if result := <-a.Srv.Store.OrgWebsites().Clear(); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't clear the websites err=%v", result.Err))
		return false, result.Err
	} else {
		return true, nil
	}
}

// Удаление в Org_phones
func (a *App) ClearAllInsurance() (bool, *model.AppError) {
	if result := <-a.Srv.Store.Insurance().Clear(); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't clear the insurance err=%v", result.Err))
		return false, result.Err
	} else {
		return true, nil
	}
}

// Удаление в Popmail
func (a *App) ClearAllPopmail() (bool, *model.AppError) {
	if result := <-a.Srv.Store.Popmail().Clear(); result.Err != nil {
		// l4g.Error(fmt.Sprintf("Couldn't clear the popmail err=%v", result.Err))
		return false, result.Err
	} else {
		return true, nil
	}
}

// Получить все письма
func (a *App) GetAllMails() ([]*model.Popmail, *model.AppError) {
	if result := <-a.Srv.Store.Popmail().GetAllPopmail(); result.Err != nil && result.Err.Id == "store.sql_popmail.missing.const" {
		result.Err.StatusCode = http.StatusNotFound
		return result.Data.([]*model.Popmail), result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return result.Data.([]*model.Popmail), result.Err
	} else {
		return result.Data.([]*model.Popmail), nil
	}
}
// Получить письма чанком
func (a *App) GetChunkMails(offset int, limit int) ([]*model.Popmail, *model.AppError) {
	if result := <-a.Srv.Store.Popmail().GetChunk(offset, limit); result.Err != nil && result.Err.Id == "store.sql_popmail.missing.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.([]*model.Popmail), nil
	}
}
// Получить количество сообщений
func (a *App) GetPopmailCount() (int64, *model.AppError) {
	if result := <-a.Srv.Store.Popmail().GetPopmailCount(); result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return result.Data.(int64), result.Err
	} else {
		return result.Data.(int64), nil
	}
}

// Получить сообщение по Uidl
func (a *App) GetPopmailByUidl(uidl string) (*model.Popmail, *model.AppError) {
	if result := <-a.Srv.Store.Popmail().GetByUidl(uidl); result.Err != nil && result.Err.Id == "store.sql_finhelp.missing_inn.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Popmail), nil
	}
}

// Получить микрофин. по ИНН
func (a *App) GetOneByInn(Inn string) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().GetByInn(Inn); result.Err != nil && result.Err.Id == "store.sql_finhelp.missing_inn.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.FinHelp), nil
	}
}
// Получить микрофин. по BVKEY
func (a *App) GetOneByBvkey(Bvkey string) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().GetByBvkey(Bvkey); result.Err != nil && result.Err.Id == "store.sql_finhelp.missing_inn.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.FinHelp), nil
	}
}

// Получить банк по Regn
func (a *App) GetOneByRegn(Regn string) (*model.FinHelp, *model.AppError) {

	if result := <-a.Srv.Store.FinHelp().GetByRegn(Regn); result.Err != nil && result.Err.Id == "store.sql_finhelp.missing_inn.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.FinHelp), nil
	}
}
// Получить банк по Num
func (a *App) GetOneByNum(Num string) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().GetByNum(Num); result.Err != nil && result.Err.Id == "store.sql_finhelp.missing_inn.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.FinHelp), nil
	}
}

// Получить insurance lics по Org_id
func (a *App) GetInsuranceLicsByOrgId(org_id string) (*model.Org_insurance_lics, *model.AppError) {
	if result := <-a.Srv.Store.Insurance().GetByOrgId(org_id); result.Err != nil && result.Err.Id == "store.sql_org_insurance.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Org_insurance_lics), nil
	}
}
// Получить insurance lics по Org_id and Lic_type_code
func (a *App) GetInsuranceLicsByOrgIdAndType(org_id string, lic_type_code string) (*model.Org_insurance_lics, *model.AppError) {
	if result := <-a.Srv.Store.Insurance().GetByOrgIdAndType(org_id, lic_type_code); result.Err != nil && result.Err.Id == "store.sql_org_insurance.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Org_insurance_lics), nil
	}
}

// Получить emails по Org_id
func (a *App) GetEmailsByOrgId(org_id string) (*model.Org_emails, *model.AppError) {
	if result := <-a.Srv.Store.OrgEmails().GetByOrgId(org_id); result.Err != nil && result.Err.Id == "store.sql_org_emails.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Org_emails), nil
	}
}

// Получить faxs по Org_id
func (a *App) GetFaxsByOrgId(org_id string) (*model.Org_faxs, *model.AppError) {
	if result := <-a.Srv.Store.OrgFaxs().GetByOrgId(org_id); result.Err != nil && result.Err.Id == "store.sql_org_faxs.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {

		return result.Data.(*model.Org_faxs), nil
	}
}

// Получить lics по Org_id
func (a *App) GetLicsByOrgId(org_id string) (*model.Org_lics, *model.AppError) {
	if result := <-a.Srv.Store.OrgLics().GetByOrgId(org_id); result.Err != nil && result.Err.Id == "store.sql_org_lics.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Org_lics), nil
	}
}

// Получить lics по Org_id and Type
func (a *App) GetLicsByOrgIdAndType(org_id string, rtype string) (*model.Org_lics, *model.AppError) {
	if result := <-a.Srv.Store.OrgLics().GetByOrgIdAndType(org_id, rtype); result.Err != nil && result.Err.Id == "store.sql_org_lics.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Org_lics), nil
	}
}

// Получить phones по Org_id
func (a *App) GetPhonesByOrgId(org_id string) (*model.Org_phones, *model.AppError) {
	if result := <-a.Srv.Store.OrgPhones().GetByOrgId(org_id); result.Err != nil && result.Err.Id == "store.sql_org_phones.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Org_phones), nil
	}
}

// Получить phones по Org_websites
func (a *App) GetWebsitesByOrgId(org_id string) (*model.Org_websites, *model.AppError) {
	if result := <-a.Srv.Store.OrgWebsites().GetByOrgId(org_id); result.Err != nil && result.Err.Id == "store.sql_org_websites.missing_org_id.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.Org_websites), nil
	}
}

// Получить страховую компанию по ИНН
func (a *App) GetInsuranceByInn(Inn string) (*model.FinHelp, *model.AppError) {
	if result := <-a.Srv.Store.FinHelp().GetByInn(Inn); result.Err != nil && result.Err.Id == "store.sql_finhelp.missing_inn.const" {
		result.Err.StatusCode = http.StatusNotFound
		return nil, result.Err
	} else if result.Err != nil {
		result.Err.StatusCode = http.StatusBadRequest
		return nil, result.Err
	} else {
		return result.Data.(*model.FinHelp), nil
	}
}

// Разархивация, arj
func (a *App) Unarj(dst string, src string) ([]string, error) {
	var files []string
	target := utils.FindExeFile("arj.exe")

	c := exec.Command(target, `x`, src, dst)
	if err := c.Run(); err != nil {
		a.CLog(25103, "Не удалось распаковать arj архив")
		return files, err
	}

	filesInfo, err := ioutil.ReadDir(dst)
	if err != nil {
		a.CLog(25103, "Не удалось получить файлы из директории после распаковки")
		return files, err
	}

	if len(filesInfo) > 0 {
		for _, f := range filesInfo {
			path := filepath.Join(dst, f.Name())
			if path != src {
				files = append(files, path)
			}
		}
	} else {
		return files, nil
	}

	return files, nil
}

// Разархивация, tar
func (a *App) Untar(dst string, fileName string) ([]string, error) {
	var files []string

	r, err := os.Open(fileName)

	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	if err != nil {
		a.CLog(25103, "Не удалось открыть zip архив")
		return files, err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		switch {
		// if no more files are found return
		case err == io.EOF:
			return files, nil
			// return any other error
		case err != nil:
			a.CLog(25103, "Возникли ошибки во время чтения tar архива")
			return files, err
			// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)
		files = append(files, target)

		// check the file type
		switch header.Typeflag {
		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					a.CLogErr(25103, "Не удалось создать папку для tar архива")
					return files, err
				}
			}
			// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				a.CLog(25103, "Не удалось открыть файл из tar архива")
				return files, err
			}
			defer f.Close()
			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				a.CLog(25103, "Не удалось скопировать файл из tar архива")
				return files, err
			}
		}
	}
}

// Разархивация, zip
func (a *App) Unzip(dest, src string) ([]string, error) {
	var files []string
	r, err := zip.OpenReader(src)
	if err != nil {
		a.CLog(25103, "Не удалось открыть zip архив")
		return files, err
	}
	defer func() {
		if err := r.Close(); err != nil {
			a.CLog(25103, "Не удалось открыть zip архив")
			return
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			a.CLog(25103, "Не удалось открыть файл из zip архива")
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				a.CLog(25103, "Не удалось открыть файл из zip архива")
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)
		files = append(files, path)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				a.CLog(25103, "Не удалось создать файл из архива")
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					a.CLog(25103, "Не удалось создать файл из архива")
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				a.CLog(25103, "Не удалось скопировать файл из zip архива")
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return files, err
		}
	}

	return files, nil
}

func (a *App) getNormalize(v interface{}, d interface{}) (interface{})  {

	dataType := reflect.TypeOf(v)
	dataValue := reflect.ValueOf(v)

	dataTyped := reflect.TypeOf(d)
	newDatad := reflect.New(dataTyped).Elem()
	for i := 0; i < dataType.NumField(); i++ {
		f := dataType.Field(i)
		fieldValue := dataValue.Field(i)
		fieldName := f.Tag.Get("csv")

		for j := 0; j < dataTyped.NumField(); j++ {
			fd := dataTyped.Field(j)
			fieldNamed := fd.Tag.Get("csv")
			newField := newDatad.FieldByName(fd.Name)
			if fieldName == fieldNamed {
				if newField.IsValid() {
					if newField.CanSet() {
						value := reflect.ValueOf(fieldValue.String())
						newField.Set(value)
					}
				}
			}
		}
	}

	return newDatad.Interface()
}


func (a *App) GetNormalizeFieldFromFinHelp(d string) (string) {
	v := model.FinHelp{}
	return a.GetNormalizeField(v, d)
}

func (a *App) GetNormalizeField(v interface{}, d string) (string)  {
	dataType := reflect.TypeOf(v)

	for i := 0; i < dataType.NumField(); i++ {
		f := dataType.Field(i)
		fieldName := f.Tag.Get("csv")

		if strings.Contains(d, fieldName) {
			return fieldName
		}
	}
	return d
}