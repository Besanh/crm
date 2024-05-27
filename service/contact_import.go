package service

import (
	"contactcenter-api/common/constants"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

func (s *Contact) PostFileImportContacts(ctx context.Context, domainUuid, userUuid, filePath, extension string, isUpdateContact bool) (int, any) {
	file, err := os.Stat(filePath)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	size := file.Size()
	if size < constants.FILE_SIZE_THRESHOLD {
		return s.UploadContactSmall(ctx, domainUuid, userUuid, filePath, extension, isUpdateContact)
	} else {
		return s.UploadContactLarge(ctx, domainUuid, userUuid, filePath, extension, isUpdateContact)
	}
}

func (s *Contact) UploadContactSmall(ctx context.Context, domainUuid, userUuid, filePath, extension string, isUpdateContact bool) (int, any) {
	userExist, err := repository.UserCrmRepo.GetUserCrmById(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(userExist.UserUuid) < 1 {
		return response.BadRequestMsg("User not found")
	}
	unitUuid := userExist.UnitUuid
	var records [][]string
	if extension == "csv" {
		f, err := os.Open(filePath)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		defer f.Close()
		csvReader := csv.NewReader(f)
		records, err = csvReader.ReadAll()
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	} else if extension == "xlsx" {
		file, err := excelize.OpenFile(filePath)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		defer file.Close()
		list := file.GetSheetList()
		records, err = file.GetRows(list[0], excelize.Options{
			RawCellValue: true,
		})
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	}
	headers := make([]string, 0)
	// create headers
	if len(records) > 1 {
		headers = s.handleHeader(extension, records)
	} else {
		return response.BadRequestMsg("No records processing, please upload contact >= 1 row")
	}

	maxWorker := 20
	records = records[2:]
	recordsLength := len(records)
	patchLength := 0
	var mutex = &sync.Mutex{}
	contactsChan := make(chan *[]model.Contact)
	contactsPhoneChan := make(chan *[]model.ContactPhone)
	contactsEmailChan := make(chan *[]model.ContactEmail)
	contacts := []model.Contact{}
	contactPhones := []model.ContactPhone{}
	contactEmails := []model.ContactEmail{}

	go func(dataChan <-chan *[]model.Contact) {
		defer RecoveryApp()
		for c := range contactsChan {
			if err := repository.ContactRepo.InsertContactsTransaction(ctx, *c, contactPhones, contactEmails); err != nil {
				log.Error(err)
			}
		}
	}(contactsChan)

	go func(dataChan <-chan *[]model.ContactPhone) {
		defer RecoveryApp()
		for p := range contactsPhoneChan {
			if err := repository.ContactRepo.InsertContactsTransaction(ctx, contacts, *p, contactEmails); err != nil {
				log.Error(err)
			}
		}
	}(contactsPhoneChan)

	go func(dataChan <-chan *[]model.ContactEmail) {
		defer RecoveryApp()
		for e := range contactsEmailChan {
			if err := repository.ContactRepo.InsertContactsTransaction(ctx, contacts, contactPhones, *e); err != nil {
				log.Error(err)
			}
		}
	}(contactsEmailChan)

	total := 0
	totalSuccess := 0
	totalFail := 0

	var wg sync.WaitGroup
	for i := 0; i < maxWorker; i++ {
		wg.Add(1)
		splitedPart := [][]string{}
		if i == 0 {
			patchLength = recordsLength / (maxWorker - 1)
			if patchLength <= 20 {
				patchLength = 20
			}
			if len(records) < patchLength {
				patchLength = len(records)
			}
			splitedPart = records[0:patchLength]
		} else {
			if patchLength >= len(records) {
				patchLength = len(records)
			}
			splitedPart = records[:patchLength]
		}
		if len(records) >= patchLength {
			records = records[patchLength:]
		}
		total += len(splitedPart)
		totalSuccessChan := make(chan int, 1)
		totalFailChan := make(chan int, 1)
		go func(partsRecord [][]string, mutex *sync.Mutex, wg *sync.WaitGroup, successChan, failChan chan int) {
			defer wg.Done()
			totalSuccess := 0
			totalFail := 0
			contacts := make([]model.Contact, 0)
			contactPhones := make([]model.ContactPhone, 0)
			contactEmails := make([]model.ContactEmail, 0)
			for _, record := range partsRecord {
				if record == nil {
					continue
				}
				// use mutex for atomic mutate the value of  totalFail , avoid data-racing when run concurrent
				if extension == "csv" {
					if strings.Contains(record[0], ";") {
						tmp := strings.Split(record[0], ";")
						if len(tmp) > 0 {
							record = tmp
						}
					}
				}
				contactRow, ok := s.handleRowRecord(headers, record)
				if ok && contactRow != nil {
					totalSuccess += 1
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					contact, contactPhoneArr, contactEmailArr, isNew, err := handleRowContact(ctx, domainUuid, userUuid, unitUuid, contactRow, isUpdateContact)
					if err != nil {
						totalFail += 1
					} else if isNew {
						contacts = append(contacts, contact)
						contactPhones = append(contactPhones, contactPhoneArr...)
						contactEmails = append(contactEmails, contactEmailArr...)
					}
				}
			}
			failChan <- totalFail
			successChan <- totalSuccess
			batchSize := 1000
			{
				totalRows := len(contacts)
				for totalRows > 0 {
					if totalRows > batchSize {
						arr := contacts[:batchSize]
						contactsChan <- &arr
						contacts = contacts[batchSize:]
						totalRows -= batchSize
					} else {
						contactsChan <- &contacts
						totalRows = 0
					}
				}
			}
			{
				totalRows := len(contactPhones)
				for totalRows > 0 {
					if totalRows > batchSize {
						arr := contactPhones[:batchSize]
						contactsPhoneChan <- &arr
						contactPhones = contactPhones[batchSize:]
						totalRows -= batchSize
					} else {
						contactsPhoneChan <- &contactPhones
						totalRows = 0
					}
				}
			}
			{
				totalRows := len(contactEmails)
				for totalRows > 0 {
					if totalRows > batchSize {
						arr := contactEmails[:batchSize]
						contactsEmailChan <- &arr
						contactEmails = contactEmails[batchSize:]
						totalRows -= batchSize
					} else {
						contactsEmailChan <- &contactEmails
						totalRows = 0
					}
				}
			}
		}(splitedPart, mutex, &wg, totalSuccessChan, totalFailChan)
		totalSuccessTmp := <-totalSuccessChan
		totalFailTmp := <-totalFailChan
		totalSuccess += totalSuccessTmp
		totalFail += totalFailTmp
	}
	wg.Wait()

	return response.OK(map[string]any{
		"message":       "file upload is processing",
		"total_success": totalSuccess,
		"total_fail":    totalFail,
	})
}

func (s *Contact) UploadContactLarge(ctx context.Context, domainUuid, userUuid, filePath, extension string, isUpdateContact bool) (int, any) {
	userExist, err := repository.UserCrmRepo.GetUserCrmById(ctx, userUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(userExist.UserUuid) < 1 {
		log.Error("User not found")
		return response.ServiceUnavailableMsg("User not found")
	}
	unitUuid := userExist.UnitUuid
	contactsChan := make(chan *[]model.Contact)
	contactsPhoneChan := make(chan *[]model.ContactPhone)
	contactsEmailChan := make(chan *[]model.ContactEmail)
	contacts := []model.Contact{}
	contactPhones := []model.ContactPhone{}
	contactEmails := []model.ContactEmail{}

	go func(dataChan <-chan *[]model.Contact) {
		defer RecoveryApp()
		for c := range contactsChan {
			if err := repository.ContactRepo.InsertContactsTransaction(ctx, *c, contactPhones, contactEmails); err != nil {
				log.Error(err)
			}
		}
	}(contactsChan)

	go func(dataChan <-chan *[]model.ContactPhone) {
		defer RecoveryApp()
		for p := range contactsPhoneChan {
			if err := repository.ContactRepo.InsertContactsTransaction(ctx, contacts, *p, contactEmails); err != nil {
				log.Error(err)
			}
		}
	}(contactsPhoneChan)

	go func(dataChan <-chan *[]model.ContactEmail) {
		defer RecoveryApp()
		for e := range contactsEmailChan {
			if err := repository.ContactRepo.InsertContactsTransaction(ctx, contacts, contactPhones, *e); err != nil {
				log.Error(err)
			}
		}
	}(contactsEmailChan)

	go func(domainUuid, userUuid, unitUuid, filePath, extension string, isUpdateContact bool, contactsCh chan<- *[]model.Contact, contactsPhoneCh chan<- *[]model.ContactPhone, contactsEmailCh chan<- *[]model.ContactEmail) {
		defer RecoveryApp()
		s.handleFileDataChannel(domainUuid, userUuid, unitUuid, filePath, extension, isUpdateContact, contactsCh, contactsPhoneCh, contactsEmailCh)
	}(domainUuid, userUuid, unitUuid, filePath, extension, isUpdateContact, contactsChan, contactsPhoneChan, contactsEmailChan)

	return response.OK(map[string]any{
		"message": "file upload is processing",
	})
}

func (s *Contact) handleFileDataChannel(domainUuid, userUuid, unitUuid, filePath, extension string, isUpdateContact bool, contactsChan chan<- *[]model.Contact, contactsPhoneChan chan<- *[]model.ContactPhone, contactsEmailChan chan<- *[]model.ContactEmail) {
	index := 0
	headers := make([]string, 0)
	totalSuccess := 0
	totalFail := 0
	contacts := []model.Contact{}
	contactPhones := []model.ContactPhone{}
	contactEmails := []model.ContactEmail{}
	if extension == "csv" {
		f, err := os.Open(filePath)
		if err != nil {
			log.Error(err)
			return
		}
		defer f.Close()
		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
		if err != nil {
			log.Error(err)
		}
		headers := s.handleHeader(extension, records)

		for k, val := range records {
			record := []string{}
			if strings.Contains(val[0], ";") {
				tmp := strings.Split(val[0], ";")
				if len(tmp) > 0 {
					record = tmp
				}
			}
			if k >= 2 {
				contactRow, ok := s.handleRowRecord(headers, record)
				if !ok {
					totalFail += 1
				} else if ok && contactRow != nil {
					totalSuccess += 1
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					contact, contactPhoneArr, contactEmailArr, isNew, err := handleRowContact(ctx, domainUuid, userUuid, unitUuid, contactRow, isUpdateContact)
					if err != nil {
						totalFail += 1
					} else if isNew {
						contacts = append(contacts, contact)
						contactPhones = append(contactPhones, contactPhoneArr...)
						contactEmails = append(contactEmails, contactEmailArr...)
					}
				}
			}
		}

		if len(contacts) > 0 {
			contactsChan <- &contacts
		}
		if len(contactPhones) > 0 {
			contactsPhoneChan <- &contactPhones
		}
		if len(contactEmails) > 0 {
			contactsEmailChan <- &contactEmails
		}
	} else if extension == "xlsx" {
		file, err := excelize.OpenFile(filePath)
		if err != nil {
			log.Error(err)
			return
		}
		defer file.Close()
		sheets := file.GetSheetList()
		rows, err := file.Rows(sheets[0])
		if err != nil {
			log.Error(err)
		}
		// Get 2 row header
		headersRecord, err := file.GetRows(sheets[0], excelize.Options{
			RawCellValue: true,
		})
		if err != nil {
			log.Error(err)
		} else {
			if len(headersRecord) > 1 {
				for k, col := range headersRecord {
					if k < 2 {
						if k == 0 {
							for _, v := range col {
								if len(v) > 0 && strings.ToLower(v) != "phone" && strings.ToLower(v) != "email" {
									header := strings.Replace(v, "\uFEFF", "", -1)
									headers = append(headers, header)
								}
							}
						} else {
							for _, v := range col {
								if len(v) > 0 {
									header := strings.Replace(v, "\uFEFF", "", -1)
									headers = append(headers, header)
								}
							}
						}
					}
				}
			} else {
				log.Error("please upload contact >= 1 row")
			}
		}
		for rows.Next() {
			record, err := rows.Columns(excelize.Options{
				RawCellValue: true,
			})
			if err != nil {
				log.Error(err)
				continue
			}
			index++
			if index >= 2 {
				contactRow, ok := s.handleRowRecord(headers, record)
				if !ok {
					totalFail += 1
				} else if ok && contactRow != nil {
					totalSuccess += 1
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					contact, contactPhoneArr, contactEmailArr, isNew, err := handleRowContact(ctx, domainUuid, userUuid, unitUuid, contactRow, isUpdateContact)
					if err != nil {
						totalFail += 1
					} else if isNew {
						contacts = append(contacts, contact)
						contactPhones = append(contactPhones, contactPhoneArr...)
						contactEmails = append(contactEmails, contactEmailArr...)
					}
				}
			}
		}
		if len(contacts) > 0 {
			contactsChan <- &contacts
		}
		if len(contactPhones) > 0 {
			contactsPhoneChan <- &contactPhones
		}
		if len(contactEmails) > 0 {
			contactsEmailChan <- &contactEmails
		}
	}
}

// func (s *Contact) handleHeaderRecord(record []string) []string {
// 	headers := make([]string, 0)
// 	for _, col := range record {
// 		col = strings.Replace(col, "\uFEFF", "", -1)
// 		headers = append(headers, col)
// 	}
// 	return headers
// }

func (s *Contact) handleRowRecord(headers []string, record []string) (map[string]any, bool) {
	row := make(map[string]any)
	phoneNumber := ""
	for index, header := range headers {
		colValue := ""
		if index < len(record) {
			colValue = strings.TrimSpace(record[index])
		}
		var value any
		if strings.ToLower(header) == "phone_main" || strings.ToLower(header) == "phone_work" || strings.ToLower(header) == "phone_personal" || strings.ToLower(header) == "phone_home" {
			phoneNumber = util.ParsePhoneNumber(colValue)
			value = phoneNumber
		} else if util.CheckPatternTelephone(colValue) {
			value = util.ParseTelToTelStr(colValue)
		} else if util.CheckPatternNumeric(colValue) {
			value, _ = strconv.Atoi(colValue)
		} else if strings.ToLower(header) == "birthday" && len(record[index]) > 0 {
			birthdayFloat, _ := strconv.ParseFloat(record[index], 64)
			birthday, err := excelize.ExcelDateToTime(birthdayFloat, false)
			if err != nil {
				value = ""
			} else {
				value = birthday.Format("2006-01-02")
			}
		} else {
			value = colValue
		}

		row[header] = value
	}
	if len(phoneNumber) < 1 {
		return nil, false
	}
	return row, true
}

func handleRowContact(ctx context.Context, domainUuid, userUuid, unitUuid string, row map[string]any, isUpdateContact bool) (model.Contact, []model.ContactPhone, []model.ContactEmail, bool, error) {
	/*isNew := true
	contactUuid := uuid.NewString()
	contact := model.Contact{
		ContactUuid: contactUuid,
		DomainUuid:  domainUuid,
		UnitUuid:    unitUuid,
		CreatedBy:   userUuid,
		CreatedAt:   time.Now(),
	}
	contactPhones := make([]model.ContactPhone, 0)
	contactEmails := make([]model.ContactEmail, 0)

	contact.PhoneNumber, _ = row["phone_main"].(string)
	if len(contact.PhoneNumber) <= 1 {
		return model.Contact{}, []model.ContactPhone{}, []model.ContactEmail{}, false, errors.New("phone number " + contact.PhoneNumber + " is invalid")
	}
	_, err := repository.ContactRepo.GetContactByPhoneNumber(ctx, domainUuid, contact.PhoneNumber)
	if err != nil {
		log.Error(err)
		return model.Contact{}, []model.ContactPhone{}, []model.ContactEmail{}, false, errors.New("contact " + contact.PhoneNumber + " is existed")
	}
	phonePersonal, _ := row["phone_personal"].(string)
	phoneWork, _ := row["phone_work"].(string)
	phoneHome, _ := row["phone_home"].(string)
	phoneNumbers := make([]string, 0)
	if len(contact.PhoneNumber) > 0 {
		phoneNumbers = append(phoneNumbers, contact.PhoneNumber)
	}
	if len(phonePersonal) > 0 && phonePersonal != "0" {
		phoneNumbers = append(phoneNumbers, phonePersonal)
	}
	if len(phoneWork) > 0 && phoneWork != "0" {
		phoneNumbers = append(phoneNumbers, phoneWork)
	}
	if len(phoneHome) > 0 && phoneHome != "0" {
		phoneNumbers = append(phoneNumbers, phoneHome)
	}
	if !isUpdateContact {
		if ct, err := repository.ContactRepo.GetContactByPhoneNumber(ctx, domainUuid, phoneNumbers...); err != nil {
			log.Error(err)
			return model.Contact{}, []model.ContactPhone{}, []model.ContactEmail{}, false, err
		} else if ct != nil {
			return model.Contact{}, []model.ContactPhone{}, []model.ContactEmail{}, false, errors.New("contact " + ct.PhoneNumber + " is existed")
		}
	}

	if len(phonePersonal) > 0 && phonePersonal != "0" {
		contactPhones = append(contactPhones,
			model.ContactPhone{
				ContactUuid: contactUuid,
				DomainUuid:  domainUuid,
				Data:        phonePersonal,
				Type:        "personal",
			})
	}
	if len(phoneWork) > 0 && phoneWork != "0" {
		contactPhones = append(contactPhones, model.ContactPhone{
			ContactUuid: contactUuid,
			DomainUuid:  domainUuid,
			Data:        phoneWork,
			Type:        "work",
		})
	}
	if len(phoneHome) > 0 && phoneHome != "0" {
		contactPhones = append(contactPhones, model.ContactPhone{
			ContactUuid: contactUuid,
			DomainUuid:  domainUuid,
			Data:        phoneHome,
			Type:        "home",
		})
	}

	contact.Fullname, _ = row["fullname"].(string)

	// Email
	contact.Email, _ = row["email_main"].(string)
	emailPersonal, _ := row["email_personal"].(string)
	emailWork, _ := row["email_work"].(string)
	emailHome, _ := row["email_home"].(string)
	if len(emailPersonal) > 0 {
		contactEmails = append(contactEmails,
			model.ContactEmail{
				ContactUuid: contactUuid,
				DomainUuid:  domainUuid,
				Data:        emailPersonal,
				EmailType:   "personal",
			})
	}
	if len(emailWork) > 0 {
		contactEmails = append(contactEmails, model.ContactEmail{
			ContactUuid: contactUuid,
			DomainUuid:  domainUuid,
			Data:        emailWork,
			EmailType:   "work",
		})
	}
	if len(emailHome) > 0 {
		contactEmails = append(contactEmails, model.ContactEmail{
			ContactUuid: contact.ContactUuid,
			DomainUuid:  domainUuid,
			Data:        emailHome,
			EmailType:   "home",
		})
	}
	contact.Birthday, _ = row["birthday"].(string)

	// Contact tag
	tagTmp, _ := row["tags"].(string)
	if len(tagTmp) > 0 {
		tagStr := strings.Split(tagTmp, constants.DELIMITER)
		tags := []string{}
		for _, val := range tagStr {
			contactTag, err := repository.ContactTagRepo.GetContacTagByTagName(ctx, domainUuid, val)
			if err != nil {
				log.Error(err)
				continue
			} else if contactTag != nil {
				tags = append(tags, contactTag.ContactTagUuid)
			} else {
				// TODO: create new contact tag
			}
		}
		contact.Tags = tags
	}

	// Contact group
	contactGroupTmp, _ := row["contact_groups"].(string)
	if len(contactGroupTmp) > 0 {
		groupStr := strings.Split(contactGroupTmp, constants.DELIMITER)
		contactGroups := []string{}
		for _, val := range groupStr {
			contactGroup, err := repository.ContactGroupRepo.GetContactGroupByGroupName(ctx, domainUuid, val)
			if err != nil {
				log.Error(err)
				continue
			} else if contactGroup != nil {
				contactGroups = append(contactGroups, contactGroup.ContactGroupUuid)
			} else {
				// TODO: create new contact group
			}
		}
		contact.ContactGroups = contactGroups
	}

	// Contact career
	contactCareerTmp, _ := row["contact_careers"].(string)
	if len(contactCareerTmp) > 0 {
		careerStr := strings.Split(contactCareerTmp, constants.DELIMITER)
		contactCareers := []string{}
		for _, val := range careerStr {
			contactCareer, err := repository.ContactCareerRepo.GetContacCareerByCareerName(ctx, domainUuid, val)
			if err != nil {
				log.Error(err)
				continue
			} else if contactCareer != nil {
				contactCareers = append(contactCareers, contactCareer.ContactCareerUuid)
			} else {
				// TODO: create new contact group
			}
		}
		contact.ContactCareers = contactCareers
	}
	contact.Description, _ = row["description"].(string)
	contact.Address, _ = row["address"].(string)
	contact.Passport, _ = row["passport"].(string)
	contact.JobTitle, _ = row["job_title"].(string)
	gender, _ := row["gender"].(string)
	if gender == "male" || gender == "female" || gender == "other" {
		contact.Gender = gender
	} else {
		contact.Gender = "other"
	}
	contact.RefId, _ = row["ref_id"].(string)
	contact.RefCode, _ = row["ref_code"].(string)
	contact.MoreInformation, _ = row["more_information"].(string)
	country, _ := row["country"].(string)
	if len(country) > 0 {
		contact.Country = country
	} else {
		contact.Country = "VietNam"
	}
	contact.Province, _ = row["province"].(string)
	contact.District, _ = row["district"].(string)
	contact.Ward, _ = row["ward"].(string)
	// contact.Status, _ = row["status"].(bool)
	contact.Status = true

	return contact, contactPhones, contactEmails, isNew, nil*/
	return model.Contact{}, []model.ContactPhone{}, []model.ContactEmail{}, false, nil
}

func (s *Contact) handleHeader(extension string, records [][]string) []string {
	headers := make([]string, 0)
	for k, col := range records {
		if k < 2 {
			if k == 0 {
				for _, v := range col {
					if len(v) > 0 {
						if extension == "csv" {
							if strings.Contains(v, ";") {
								tmp := strings.Split(v, ";")
								if len(tmp) > 0 {
									for _, item := range tmp {
										if len(item) > 0 && strings.ToLower(item) != "phone" && strings.ToLower(item) != "email" {
											item = strings.Replace(item, "\uFEFF", "", -1)
											headers = append(headers, item)
										}
									}
								}
							}
						} else if extension == "xlsx" {
							if strings.ToLower(v) != "phone" && strings.ToLower(v) != "email" {
								header := strings.Replace(v, "\uFEFF", "", -1)
								headers = append(headers, header)
							}
						}
					}
				}
			} else {
				for _, v := range col {
					if len(v) > 0 {
						if extension == "csv" {
							if strings.Contains(v, ";") {
								tmp := strings.Split(v, ";")
								if len(tmp) > 0 {
									for _, item := range tmp {
										if len(item) > 0 && strings.ToLower(item) != "phone" && strings.ToLower(item) != "email" {
											item = strings.Replace(item, "\uFEFF", "", -1)
											headers = append(headers, item)
										}
									}
								}
							}
						} else if extension == "xlsx" {
							header := strings.Replace(v, "\uFEFF", "", -1)
							headers = append(headers, header)
						}
					}
				}
			}
		}
	}

	return headers
}
