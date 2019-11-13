package controllers

import (
	"github.com/go-chi/chi"
	"html/template"
	"io"
	"net/http"
	"obas/config"
	addressIO "obas/io/address"
	demograhpyIO "obas/io/demographics"
	documentIO "obas/io/documents"
	loginIO "obas/io/login"
	storageIO "obas/io/storage"
	usersIO "obas/io/users"
	"os"
	"strings"
	"time"
)

const (
	layoutOBAS        = "2006-01-02"
	dangerAlertStyle  = "alert-danger"
	successAlertStyle = "alert-success"
)

type AddressPlaceHolder struct {
	AddressName string
	Address     string
	PostalCode  string
}

type PageToast struct {
	AlertType string
	AlertInfo string
}

type ContactPlaceHolder struct {
	ContactName   string
	ContactDetail string
}

type DistrictData struct {
	Student          usersIO.User
	Provinces        []demograhpyIO.Province
	Districts        []demograhpyIO.District
	Towns            []demograhpyIO.Town
	SelectedProvince demograhpyIO.Province
	SelectedDistrict demograhpyIO.District
	SelectedTown     demograhpyIO.Town
	Alert            PageToast
}

func Users(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", UsersHandler(app))
	r.Get("/admin", AdminHandler(app))
	r.Get("/student", StudentHandler(app))

	r.Get("/processingStatus", ProcessingStatusTypeHandler(app))
	r.Get("/student/application", StudentApplicationStatusHandler(app))
	r.Get("/studentContact", StudentContactsHandler(app))
	r.Get("/student/documents", StudentDocumentsHandler(app))
	r.Get("/studentResults", StudentResultsHandler(app))

	r.Get("/student/profile/personal", StudentProfilePersonalHandler(app))
	r.Get("/student/profile/demography", StudentProfileDemographyHandler(app))
	r.Get("/student/profile/address", StudentProfileAddressHandler(app))
	r.Get("/student/profile/relative", StudentProfileRelativeHandler(app))
	r.Get("/student/profile/settings", StudentProfileSettingsHandler(app))
	r.Get("/student/profile/courses", StudentProfileCourseHandler(app))
	r.Get("/student/profile/subjects", StudentProfileSubjectHandler(app))
	r.Get("/student/profile/districts", StudentProfileDistrictHandler(app))
	r.Get("/student/profile/contacts", StudentProfileContactsHandler(app))
	r.Get("/student/profile/application", StudentProfileApplicationHandler(app))
	r.Get("/student/profile/application_process", StudentProfileApplicationProcessHandler(app))
	r.Get("/student/profile/documents", StudentProfileDocumentsHandler(app))

	r.Post("/student/profile/personal/update", UpdateStudentProfilePersonalHandler(app))
	r.Post("/student/profile/address/addresstype", StudentProfileAddressTypeHandler(app))
	r.Post("/student/profile/address/update", StudentProfileAddressUpdateHandler(app))
	r.Post("/student-profile-relative-upate", StudentProfileRelativeUpdateHandler(app))
	r.Post("/student-profile-demography-update", StudentProfileDemographyUpdateHandler(app))
	r.Post("/student-profile-password-update", StudentProfilePasswordUpdate(app))
	r.Post("/student/profile/contact/contacttype", StudentProfileContactTypeHandler(app))
	r.Post("/student-profile-contact-update", StudentProfileContactUpdateHandler(app))
	r.Post("/student/profile/demographics/districts", StudentProfileDistrictsHandler(app))
	r.Post("/student/profile/demographics/towns", StudentProfileTownsHandler(app))
	r.Post("/student-profile-town-update", StudentProfileTownUpdateHandler(app))
	r.Post("/student-document-file-upload", StudentDocumentsUploadHandler(app))

	return r
}

func StudentProfileDocumentsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}

		type PageData struct {
			Student usersIO.User
		}

		data := PageData{user}
		files := []string{
			app.Path + "content/student/application/student_documents.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileApplicationProcessHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}

		type PageData struct {
			Student usersIO.User
		}

		data := PageData{user}
		files := []string{
			app.Path + "content/student/application/student_Application_process.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileApplicationHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}

		type PageData struct {
			Student usersIO.User
		}

		data := PageData{user}
		files := []string{
			app.Path + "content/student/application/student_Application.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}

}

func StudentDocumentsUploadHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseMultipartForm(10 << 20)
		documentTypeId := r.PostFormValue("documenttype")
		file, handler, err := r.FormFile("file")
		if err != nil {
			app.ErrorLog.Println(err.Error())
			setSessionMessage(app, r, dangerAlertStyle, "Internal Server Error!")
		} else {
			defer file.Close()
			tempFilePath := "temp_file/" + handler.Filename
			tempFile, err := os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				app.ErrorLog.Println(err.Error())
				setSessionMessage(app, r, dangerAlertStyle, "Internal Server Error!")
			} else {
				defer tempFile.Close()
				io.Copy(tempFile, file)
				fileName := handler.Filename
				successMessage := "File: " + fileName + " updated!"
				failureMessage := "File: " + fileName + " NOT updated!"
				app.InfoLog.Println("File Name: ", fileName)
				app.InfoLog.Println("File Size: ", handler.Size)
				app.InfoLog.Println("File Header: ", handler.Header)

				fileData, err := storageIO.UploadFile(tempFilePath, token)
				os.Remove(tempFilePath)
				if err != nil {
					app.ErrorLog.Println(err.Error())
					setSessionMessage(app, r, dangerAlertStyle, failureMessage)
				} else {
					if strings.Contains(fileData.Id, "error") {
						app.ErrorLog.Println(fileData.Id)
						setSessionMessage(app, r, dangerAlertStyle, "Internal Server Error!")
					} else {
						document := documentIO.Document{
							fileData.Id,
							documentTypeId,
							fileName,
							fileData.Url,
							handler.Header.Get("Content-Type"),
							time.Now(),
							"", ""}
						app.InfoLog.Println("Document to create: ", document)
						documentSaved, err := documentIO.CreateDocument(document, token)
						if err != nil {
							app.ErrorLog.Println(err.Error())
							setSessionMessage(app, r, dangerAlertStyle, "Document NOT saved!")
						} else {
							if documentSaved {
								userDocument := usersIO.UserDocument{email, fileData.Id}
								app.InfoLog.Println("User document to save: ", userDocument)
								userDocumentSaved, err := usersIO.UpdateUserDocument(userDocument, token)
								if err != nil {
									app.ErrorLog.Println(err.Error())
									setSessionMessage(app, r, dangerAlertStyle, "User document saved!")
								} else {
									if userDocumentSaved {
										setSessionMessage(app, r, successAlertStyle, successMessage)
									} else {
										setSessionMessage(app, r, dangerAlertStyle, failureMessage)
									}
								}
							} else {
								setSessionMessage(app, r, dangerAlertStyle, "User document NOT saved!")
							}
						}
					}
				}
			}
		}
		http.Redirect(w, r, "/users/student/documents", 301)
	}
}

func StudentDocumentsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		type DocumentData struct {
			Document            documentIO.Document
			DocumentType        string
			DocumentDate        string
			DocumentStatusBadge string
		}
		var alert PageToast
		var documentTypes []documentIO.DocumentType
		var userDocuments []DocumentData
		documentTypes, err = documentIO.GetDocumentTypes()
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve document types!"}
		}

		userDocumentsObj, err := usersIO.GetUserDocuments(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve user documents!"}
		} else {
			for _, userDocument := range userDocumentsObj {
				documentId := userDocument.DocumentId
				document, err := documentIO.GetDocument(documentId)
				if err != nil {
					app.ErrorLog.Println(err.Error() + " - Could not retrieve document for id: " + documentId)
				} else {
					documentType, err := documentIO.GetDocumentType(document.DocumentTypeId)
					if err != nil {
						app.ErrorLog.Println(err.Error() + " - Could not retrieve document type for document!")
					} else {
						date := getDate_YYYYMMDD(document.Date.String())
						var progressBadge string
						documentStatus := document.DocumentStatus
						if documentStatus == "Approved" {
							progressBadge = "badge-success"
						} else if documentStatus == "Not Approved" {
							progressBadge = "badge-danger"
						} else {
							progressBadge = "badge-warning"
						}
						documentData := DocumentData{document, documentType.DocumentTypename, date, progressBadge}
						userDocuments = append(userDocuments, documentData)
					}
				}
			}
			message := app.Session.GetString(r.Context(), "message")
			messageType := app.Session.GetString(r.Context(), "message-type")
			if message != "" && messageType != "" {
				alert = PageToast{messageType, message}
				app.Session.Remove(r.Context(), "message")
				app.Session.Remove(r.Context(), "message-type")
			}
		}

		type PageData struct {
			Student       usersIO.User
			DocumentTypes []documentIO.DocumentType
			UserDocuments []DocumentData
			Alert         PageToast
		}
		data := PageData{
			user,
			documentTypes,
			userDocuments,
			alert,
		}
		files := []string{
			app.Path + "content/student/student_documents.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileTownUpdateHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseForm()
		townCode := r.PostFormValue("town")
		userTown := usersIO.UserTown{email, townCode}
		app.InfoLog.Println("UserTown to update: ", userTown)
		updated, err := usersIO.UpdateUserTown(userTown, token)
		successMessage := "User town updated!"
		failureMessage := "User town NOT updated!"
		if err != nil {
			app.ErrorLog.Println(err.Error())
			setSessionMessage(app, r, dangerAlertStyle, failureMessage)
		} else {
			if updated {
				setSessionMessage(app, r, successAlertStyle, successMessage)
			} else {
				setSessionMessage(app, r, dangerAlertStyle, failureMessage)
			}
		}
		app.InfoLog.Println("Update response is ", updated)
		http.Redirect(w, r, "/users/student/profile/districts", 301)
	}
}

func StudentProfileTownsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		var alert PageToast
		var districts []demograhpyIO.District
		var selectedProvince demograhpyIO.Province
		var towns []demograhpyIO.Town
		var selectedDistrict demograhpyIO.District
		var selectedTown demograhpyIO.Town

		provinces, err := demograhpyIO.GetProvinces()
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve provinces!"}
		} else {
			r.ParseForm()
			provinceCode := r.PostFormValue("province")
			districtCode := r.PostFormValue("district")
			districtsInProvince, errorAlert := getDistrictsInProvince(app, provinceCode)
			if errorAlert.AlertInfo != "" {
				alert = errorAlert
			} else {
				districts = districtsInProvince
				townsInDistrict, errorAlert := getTownsInDistrict(app, districtCode)
				if errorAlert.AlertInfo != "" {
					alert = errorAlert
				} else {
					towns = townsInDistrict
				}
			}
			selectedProvince = filterProvince(provinces, provinceCode)
			selectedDistrict = filterDistrict(districts, districtCode)
		}

		data := DistrictData{
			user,
			provinces,
			districts,
			towns,
			selectedProvince,
			selectedDistrict,
			selectedTown,
			alert}

		files := []string{
			app.Path + "content/student/profile/district_and_municipality.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func filterDistrict(districts []demograhpyIO.District, districtCode string) demograhpyIO.District {
	var retdata demograhpyIO.District
	for _, district := range districts {
		if district.DistrictCode == districtCode {
			retdata = district
			break
		}
	}
	return retdata
}

func getTownsInDistrict(app *config.Env, districtCode string) ([]demograhpyIO.Town, PageToast) {
	var towns []demograhpyIO.Town
	var alert PageToast
	townsInDistrict, err := demograhpyIO.GetTownsInDistrict(districtCode)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		alert = PageToast{dangerAlertStyle, "Could not retrieve towns in district!"}
		return towns, alert
	}
	for _, districtTown := range townsInDistrict {
		town, err := demograhpyIO.GetTown(districtTown.TownCode)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		} else {
			towns = append(towns, town)
		}
	}
	return towns, alert
}

func StudentProfileDistrictsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		var alert PageToast
		var districts []demograhpyIO.District
		var selectedProvince demograhpyIO.Province
		var towns []demograhpyIO.Town
		var selectedDistrict demograhpyIO.District
		var selectedTown demograhpyIO.Town

		provinces, err := demograhpyIO.GetProvinces()
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve provinces!"}
		} else {
			r.ParseForm()
			provinceCode := r.PostFormValue("province")
			districtsInProvince, errorAlert := getDistrictsInProvince(app, provinceCode)
			if errorAlert.AlertInfo != "" {
				alert = errorAlert
			} else {
				districts = districtsInProvince
			}
			selectedProvince = filterProvince(provinces, provinceCode)
		}

		data := DistrictData{
			user,
			provinces,
			districts,
			towns,
			selectedProvince,
			selectedDistrict,
			selectedTown,
			alert}
		files := []string{
			app.Path + "content/student/profile/district_and_municipality.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func filterProvince(provinces []demograhpyIO.Province, provinceCode string) demograhpyIO.Province {
	var retdata demograhpyIO.Province
	for _, province := range provinces {
		if province.ProvinceCode == provinceCode {
			retdata = province
			break
		}
	}
	return retdata
}

func getDistrictsInProvince(app *config.Env, provinceCode string) ([]demograhpyIO.District, PageToast) {
	var districts []demograhpyIO.District
	var alert PageToast
	districtsInProvince, err := demograhpyIO.GetDistrictsInProvince(provinceCode)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		alert = PageToast{dangerAlertStyle, "Could not retrieve districts in province!"}
		return districts, alert
	}
	for _, provinceDistrict := range districtsInProvince {
		district, err := demograhpyIO.GetDistrict(provinceDistrict.DistrictCode)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		} else {
			districts = append(districts, district)
		}
	}
	return districts, alert

}

func StudentProfileDistrictHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}

		var alert PageToast
		var selectedProvince demograhpyIO.Province
		var districts []demograhpyIO.District
		var towns []demograhpyIO.Town
		var selectedDistrict demograhpyIO.District
		var selectedTown demograhpyIO.Town

		provinces, err := demograhpyIO.GetProvinces()
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve provinces!"}
		} else {
			userTown, err := usersIO.ReadUserTown(email)
			if err != nil {
				app.ErrorLog.Println(err.Error())
				alert = PageToast{dangerAlertStyle, "Could not retrieve user town!"}
			} else {
				districtTown, err := demograhpyIO.GetDistrictForTown(userTown.TownCode)
				if err != nil {
					app.ErrorLog.Println(err.Error())
					alert = PageToast{dangerAlertStyle, "Could not retrieve district for town!"}
				} else {
					provinceDistrict, err := demograhpyIO.GetProvinceForDistrict(districtTown.DistrictCode)
					if err != nil {
						app.ErrorLog.Println(err.Error())
						alert = PageToast{dangerAlertStyle, "Could not retrieve province for district!"}
					} else {
						selectedProvince = filterProvince(provinces, provinceDistrict.ProvinceCode)
						districtsInProvince, errorAlert := getDistrictsInProvince(app, provinceDistrict.ProvinceCode)
						if errorAlert.AlertInfo != "" {
							alert = errorAlert
						} else {
							districts = districtsInProvince
							selectedDistrict = filterDistrict(districts, districtTown.DistrictCode)
							townsInDistrict, errorAlert := getTownsInDistrict(app, districtTown.DistrictCode)
							if errorAlert.AlertInfo != "" {
								alert = errorAlert
							} else {
								towns = townsInDistrict
								for _, town := range towns {
									if town.TownCode == districtTown.TownCode {
										selectedTown = town
										break
									}
								}
								message := app.Session.GetString(r.Context(), "message")
								messageType := app.Session.GetString(r.Context(), "message-type")
								if message != "" && messageType != "" {
									alert = PageToast{messageType, message}
									app.Session.Remove(r.Context(), "message")
									app.Session.Remove(r.Context(), "message-type")
								}
							}
						}
					}
				}
			}
		}

		data := DistrictData{
			user,
			provinces,
			districts,
			towns,
			selectedProvince,
			selectedDistrict,
			selectedTown,
			alert}

		files := []string{
			app.Path + "content/student/profile/district_and_municipality.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileContactUpdateHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseForm()
		contactTypeId := r.PostFormValue("contactTypeId")
		contact := r.PostFormValue("contact")
		userContact := usersIO.UserContact{email, contactTypeId, contact}
		app.InfoLog.Println("UserContact to update: ", userContact)
		updated, err := usersIO.UpdateUserContact(userContact, token)
		successMessage := "User contact updated!"
		failureMessage := "User contact NOT updated!"
		if err != nil {
			app.ErrorLog.Println(err.Error())
			setSessionMessage(app, r, dangerAlertStyle, failureMessage)
		} else {
			if updated {
				setSessionMessage(app, r, successAlertStyle, successMessage)
			} else {
				setSessionMessage(app, r, dangerAlertStyle, failureMessage)
			}
		}
		app.InfoLog.Println("Update response is ", updated)
		http.Redirect(w, r, "/users/student/profile/contacts", 301)
	}
}

func StudentProfileContactTypeHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		var alert PageToast
		var userContact usersIO.UserContact
		var contactName string
		var contactTypes []addressIO.ContactType
		var contacts []ContactPlaceHolder
		r.ParseForm()
		contactTypeId := r.PostFormValue("contacttypes")
		if contactTypeId == "" {
			errMsg := "No contact type selected!"
			app.ErrorLog.Println(errMsg)
			alert = PageToast{dangerAlertStyle, errMsg}
		} else {
			contactTypes, err := addressIO.GetContactTypes()
			if err != nil {
				errMsg := "Could not retrieve contact types!"
				app.ErrorLog.Println(err.Error() + " - " + errMsg)
				alert = PageToast{dangerAlertStyle, errMsg}
			} else {
				for _, contactType := range contactTypes {
					if contactTypeId == contactType.ContactTypeId {
						contactName = contactType.Name
					}
					userContact, err := usersIO.GetUserContact(email, contactType.ContactTypeId)
					if err != nil {
						errMsg := "Could not retrieve user contact for " + contactType.Name
						app.ErrorLog.Println(err.Error() + " - " + errMsg)
					} else {
						contacts = append(contacts, ContactPlaceHolder{contactType.Name, userContact.Contact})
					}
				}
				userContact, err = usersIO.GetUserContact(email, contactTypeId)
				if err != nil {
					errMsg := "Could not retrieve user contact for " + contactName
					app.ErrorLog.Println(err.Error() + " - " + errMsg)
					alert = PageToast{dangerAlertStyle, errMsg}
				}
			}
		}

		type PageData struct {
			Student       usersIO.User
			ContactTypes  []addressIO.ContactType
			Contacts      []ContactPlaceHolder
			Contact       usersIO.UserContact
			ContactTypeId string
			ContactName   string
			Alert         PageToast
		}
		data := PageData{user, contactTypes, contacts, userContact, contactTypeId, contactName, alert}
		files := []string{
			app.Path + "content/student/profile/contacts.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileContactsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		var alert PageToast
		var contacts []ContactPlaceHolder
		contactTypes, err := addressIO.GetContactTypes()
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve contact types!"}
		} else {
			for _, contactType := range contactTypes {
				userContact, err := usersIO.GetUserContact(email, contactType.ContactTypeId)
				if err != nil {
					errMsg := "Could not retrieve user contact for " + contactType.Name
					app.ErrorLog.Println(err.Error() + " - " + errMsg)
				} else {
					contacts = append(contacts, ContactPlaceHolder{contactType.Name, userContact.Contact})
				}
			}
			message := app.Session.GetString(r.Context(), "message")
			messageType := app.Session.GetString(r.Context(), "message-type")
			if message != "" && messageType != "" {
				alert = PageToast{messageType, message}
				app.Session.Remove(r.Context(), "message")
				app.Session.Remove(r.Context(), "message-type")
			}
		}

		type PageData struct {
			Student       usersIO.User
			ContactTypes  []addressIO.ContactType
			Contacts      []ContactPlaceHolder
			Contact       usersIO.UserContact
			ContactTypeId string
			ContactName   string
			Alert         PageToast
		}

		data := PageData{user, contactTypes, contacts, usersIO.UserContact{}, "", "", alert}
		files := []string{
			app.Path + "content/student/profile/contacts.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfilePasswordUpdate(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseForm()
		currentPassword := r.PostFormValue("current_password")
		newPasswordOne := r.PostFormValue("new_password_one")
		newPasswordTwo := r.PostFormValue("new_password_two")

		successMessage := "User password updated!"
		failureMessage := "User password NOT Updated!"

		if newPasswordOne != newPasswordTwo {
			errMsg := "New password mismatch."
			app.ErrorLog.Println(errMsg)
			setSessionMessage(app, r, dangerAlertStyle, failureMessage+" - "+errMsg)
		} else {
			userChangePassword := loginIO.ChangePassword{email, currentPassword, newPasswordOne, time.Now()}
			app.InfoLog.Println("User password to update: ", userChangePassword)
			loginToken, err := loginIO.DoChangePassword(userChangePassword, token)

			if err != nil {
				app.ErrorLog.Println(err.Error())
				setSessionMessage(app, r, dangerAlertStyle, failureMessage)
			}
			app.Session.Put(r.Context(), "userId", loginToken.Email)
			app.Session.Put(r.Context(), "token", loginToken.Token)
			setSessionMessage(app, r, successAlertStyle, successMessage)
		}
		http.Redirect(w, r, "/users/student/profile/settings", 301)
	}
}

func StudentProfileSettingsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		var alert PageToast
		message := app.Session.GetString(r.Context(), "message")
		messageType := app.Session.GetString(r.Context(), "message-type")
		if message != "" && messageType != "" {
			alert = PageToast{messageType, message}
			app.Session.Remove(r.Context(), "message")
			app.Session.Remove(r.Context(), "message-type")
		}

		type PageData struct {
			Student usersIO.User
			Alert   PageToast
		}

		data := PageData{user, alert}
		files := []string{
			app.Path + "content/student/profile/settings.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileDemographyUpdateHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseForm()
		title := r.PostFormValue("title")
		gender := r.PostFormValue("gender")
		race := r.PostFormValue("race")
		userDemograpgy := usersIO.UserDemography{email, title, gender, race}
		app.InfoLog.Println("userDemography to update: ", userDemograpgy)

		updated, err := usersIO.UpdateUserDemographics(userDemograpgy, token)
		successMessage := "User demography updated!"
		failureMessage := "User demography NOT Updated!"

		if err != nil {
			app.ErrorLog.Println(err.Error())
			setSessionMessage(app, r, dangerAlertStyle, failureMessage)
		} else {
			if updated {
				setSessionMessage(app, r, successAlertStyle, successMessage)
			} else {
				setSessionMessage(app, r, dangerAlertStyle, failureMessage)
			}
		}
		app.InfoLog.Println("UserDemography update response is ", updated)
		http.Redirect(w, r, "/users/student/profile/demography", 301)
	}
}
func StudentProfileDemographyHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		var alert PageToast
		var genders []demograhpyIO.Gender
		var races []demograhpyIO.Race
		var gender demograhpyIO.Gender
		var race demograhpyIO.Race
		var title demograhpyIO.Title
		titles, err := demograhpyIO.GetTitles()
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve titles!"}
		} else {
			genders, err = demograhpyIO.GetGenders()
			if err != nil {
				app.ErrorLog.Println(err.Error())
				alert = PageToast{dangerAlertStyle, "Could not retrieve genders!"}
			} else {
				races, err = demograhpyIO.GetRaces()
				if err != nil {
					app.ErrorLog.Println(err.Error())
					alert = PageToast{dangerAlertStyle, "Could not retrieve races!"}
				} else {
					userDemography, err := usersIO.GetUserDemographic(email)
					if err != nil {
						app.ErrorLog.Println(err.Error())
						alert = PageToast{dangerAlertStyle, "Could not retrieve student demography!"}
					} else {
						message := app.Session.GetString(r.Context(), "message")
						messageType := app.Session.GetString(r.Context(), "message-type")
						if message != "" && messageType != "" {
							alert = PageToast{messageType, message}
							app.Session.Remove(r.Context(), "message")
							app.Session.Remove(r.Context(), "message-type")
						}
						title = getUserTitle(userDemography, titles)
						gender = getUserGender(userDemography, genders)
						race = getUserRace(userDemography, races)
					}
				}
			}
		}

		type PageData struct {
			Student       usersIO.User
			Titles        []demograhpyIO.Title
			Genders       []demograhpyIO.Gender
			Races         []demograhpyIO.Race
			Alert         PageToast
			StudentTitle  demograhpyIO.Title
			StudentGender demograhpyIO.Gender
			StudentRace   demograhpyIO.Race
		}

		data := PageData{user, titles, genders, races, alert, title, gender, race}
		app.InfoLog.Println("DistrictData: ", data.Alert)
		files := []string{
			app.Path + "content/student/profile/demography.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func getUserRace(demography usersIO.UserDemography, races []demograhpyIO.Race) demograhpyIO.Race {
	for _, race := range races {
		if demography.RaceId == race.RaceId {
			return race
		}
	}
	return demograhpyIO.Race{}
}

func getUserTitle(demography usersIO.UserDemography, titles []demograhpyIO.Title) demograhpyIO.Title {
	for _, title := range titles {
		if demography.TitleId == title.TitleId {
			return title
		}
	}
	return demograhpyIO.Title{}
}

func getUserGender(demography usersIO.UserDemography, genders []demograhpyIO.Gender) demograhpyIO.Gender {
	for _, gender := range genders {
		if demography.GenderId == gender.GenderId {
			return gender
		}
	}
	return demograhpyIO.Gender{}
}

func StudentProfileRelativeUpdateHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseForm()
		relativeName := r.PostFormValue("relative_name")
		relationship := r.PostFormValue("relationship")
		cellphone := r.PostFormValue("relative_cellphone")
		relativeEmail := r.PostFormValue("relative_email")
		userRelative := usersIO.UserRelative{email, relativeName, cellphone, relativeEmail, relationship}
		app.InfoLog.Println("UserRelative to update: ", userRelative)
		updated, err := usersIO.UpdateUserRelative(userRelative, token)

		successMessage := "User relative updated!"
		failureMessage := "User relative NOT Updated!"

		if err != nil {
			app.ErrorLog.Println(err.Error())
			setSessionMessage(app, r, dangerAlertStyle, failureMessage)
		} else {
			if updated {
				setSessionMessage(app, r, successAlertStyle, successMessage)
			} else {
				setSessionMessage(app, r, dangerAlertStyle, failureMessage)
			}
		}
		app.InfoLog.Println("UserRelative update response is ", updated)
		http.Redirect(w, r, "/users/student/profile/relative", 301)

	}
}

func setSessionMessage(app *config.Env, r *http.Request, messageType string, message string) {
	app.Session.Put(r.Context(), "message-type", messageType)
	app.Session.Put(r.Context(), "message", message)
}

func StudentProfileSubjectHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}

		type PageData struct {
			Student usersIO.User
		}

		data := PageData{user}
		files := []string{
			app.Path + "content/student/profile/subjects.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileCourseHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}

		type PageData struct {
			Student usersIO.User
		}

		data := PageData{user}
		files := []string{
			app.Path + "content/student/profile/courses.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}

}

func StudentProfileRelativeHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}

		var alert PageToast

		userRelative, err := usersIO.GetUserRelative(user.Email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			alert = PageToast{dangerAlertStyle, "Could not retrieve student relative!"}
		} else {
			message := app.Session.GetString(r.Context(), "message")
			messageType := app.Session.GetString(r.Context(), "message-type")
			if message != "" && messageType != "" {
				alert = PageToast{messageType, message}
				app.Session.Remove(r.Context(), "message")
				app.Session.Remove(r.Context(), "message-type")
			}
		}

		type PageData struct {
			Student         usersIO.User
			StudentRelative usersIO.UserRelative
			Alert           PageToast
		}

		data := PageData{user, userRelative, alert}
		files := []string{
			app.Path + "content/student/profile/relative.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileAddressUpdateHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseForm()
		addressTypeId := r.PostFormValue("addressTypeId")
		address := r.PostFormValue("address")
		postalCode := r.PostFormValue("postalCode")
		userAddress := usersIO.UserAddress{email, addressTypeId, address, postalCode}
		app.InfoLog.Println("UserAddress to update: ", userAddress)
		updated, err := usersIO.UpdateUserAddress(userAddress, token)

		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		app.InfoLog.Println("Update response is ", updated)
		http.Redirect(w, r, "/users/student/profile/address", 301)
	}
}

func StudentProfileAddressHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		addressTypes, err := addressIO.GetAddressTypes()
		if err != nil {
			app.ErrorLog.Println(err.Error(), addressTypes)
		}

		addresses := []AddressPlaceHolder{}

		for _, addressType := range addressTypes {
			userAddress, err := usersIO.GetUserAddress(email, addressType.AddressTypeID)
			if err != nil {
				app.ErrorLog.Println(err.Error())
			} else {
				addresses = append(addresses, AddressPlaceHolder{addressType.AddressName, userAddress.Address, userAddress.PostalCode})
			}
		}

		type PageData struct {
			Student       usersIO.User
			AddressTypes  []addressIO.AddressType
			Addresses     []AddressPlaceHolder
			Address       usersIO.UserAddress
			AddressTypeID string
			AddressName   string
		}

		data := PageData{user, addressTypes, addresses, usersIO.UserAddress{}, "", ""}
		files := []string{
			app.Path + "content/student/profile/address.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentProfileAddressTypeHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		r.ParseForm()
		addressTypeId := r.PostFormValue("addresstypes")
		userAddress, err := usersIO.GetUserAddress(email, addressTypeId)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}

		addressTypes, err := addressIO.GetAddressTypes()
		if err != nil {
			app.ErrorLog.Println(err.Error(), addressTypes)
		}

		addresses := []AddressPlaceHolder{}
		var addressName string

		for _, addressType := range addressTypes {
			if addressTypeId == addressType.AddressTypeID {
				addressName = addressType.AddressName
			}
			userAddress, err := usersIO.GetUserAddress(email, addressType.AddressTypeID)
			if err != nil {
				app.ErrorLog.Println(err.Error())
			} else {
				addresses = append(addresses, AddressPlaceHolder{addressType.AddressName, userAddress.Address, userAddress.PostalCode})
			}
		}

		type PageData struct {
			Student       usersIO.User
			AddressTypes  []addressIO.AddressType
			Addresses     []AddressPlaceHolder
			Address       usersIO.UserAddress
			AddressTypeID string
			AddressName   string
		}

		data := PageData{user, addressTypes, addresses, userAddress, addressTypeId, addressName}
		files := []string{
			app.Path + "content/student/profile/address.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func StudentHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" || len(email) <= 0 {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		type PageData struct {
			Student usersIO.User
		}
		data := PageData{user}
		files := []string{
			app.Path + "content/student/student_dashboard.page.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}

}

func StudentProfilePersonalHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userId")
		if email == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		user, err := usersIO.GetUser(email)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Redirect(w, r, "/login", 301)
			return
		}
		dobString := getDate_YYYYMMDD(user.DateOfBirth.String()) // split date and get in format: yyy-mm-dd

		type PageData struct {
			Student     usersIO.User
			DateOfBirth string
		}

		data := PageData{user, dobString}
		files := []string{
			app.Path + "content/student/profile/personal.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func getDate_YYYYMMDD(dateString string) string {
	return strings.Split(dateString, " ")[0]
}

func UpdateStudentProfilePersonalHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := app.Session.GetString(r.Context(), "userId")
		token := app.Session.GetString(r.Context(), "token")
		if email == "" || token == "" {
			http.Redirect(w, r, "/login", 301)
			return
		}
		idNumber := r.PostFormValue("id_number")
		firstName := r.PostFormValue("first_name")
		lastName := r.PostFormValue("last_name")
		dateOfBirthStr := r.PostFormValue("dateOfBirth")
		dateOfBirth, _ := time.Parse(layoutOBAS, dateOfBirthStr)
		user := usersIO.User{email, idNumber, firstName, "", lastName, dateOfBirth}
		app.InfoLog.Println("User to update: ", user)
		updated, err := usersIO.UpdateUser(user, token)

		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		app.InfoLog.Println("Update response is ", updated)
		http.Redirect(w, r, "/users/student/profile/personal", 301)
	}
}

func UsersHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type PageData struct {
			//courses []io.Users
			name string
		}
		data := PageData{""}

		files := []string{
			app.Path + "base/register/register.page.html",
			/**app.Path + "/users/users.page.html",
			app.Path + "/base/base.page.html",
			app.Path + "/base/navbar.page.html",
			app.Path + "/base/sidebarOld.page.html",
			app.Path + "/base/footer.page.html",*/
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}

	}
}

func AdminHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//allAdmin, err := io.GetAdmins()
		//
		//if err != nil {
		//	app.ServerError(w, err)
		//}

		type PageData struct {
			//courses []io.Admin
			name string
		}
		data := PageData{""}

		files := []string{
			app.Path + "/users/users.page.html",
			app.Path + "/base/base.page.html",
			app.Path + "/base/navbar.page.html",
			app.Path + "/base/sidebarOld.page.html",
			app.Path + "/base/footer.page.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}

	}
}

func ProcessingStatusTypeHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//allProcess, err := io.GetProcessingStatusTypes()
		//
		//if err != nil {
		//	app.ServerError(w, err)
		//}

		type PageData struct {
			//subjects []io.ProcessingStatusType
			name string
		}
		data := PageData{""}

		files := []string{
			app.Path + "/users/users.page.html",
			app.Path + "/base/base.page.html",
			app.Path + "/base/navbar.page.html",
			app.Path + "/base/sidebarOld.page.html",
			app.Path + "/base/footer.page.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}

	}
}

func StudentApplicationStatusHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := []string{
			app.Path + "content/student/Student_Application.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())

		}
	}
}

func StudentContactsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//allStudentContacts, err := io.GetStudentContacts()
		//
		//if err != nil {
		//	app.ServerError(w, err)
		//}

		type PageData struct {
			//subjects []io.StudentContacts
			name string
		}
		data := PageData{""}

		files := []string{
			app.Path + "/users/users.page.html",
			app.Path + "/base/base.page.html",
			app.Path + "/base/navbar.page.html",
			app.Path + "/base/sidebarOld.page.html",
			app.Path + "/base/footer.page.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}

	}
}

func StudentResultsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//allStudentResults, err := io.GetStudentResults()
		//
		//if err != nil {
		//	app.ServerError(w, err)
		//}

		type PageData struct {
			//subjects []io.StudentResults
			name string
		}
		data := PageData{""}

		files := []string{
			app.Path + "/users/users.page.html",
			app.Path + "/base/base.page.html",
			app.Path + "/base/navbar.page.html",
			app.Path + "/base/sidebarOld.page.html",
			app.Path + "/base/footer.page.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}

	}
}
