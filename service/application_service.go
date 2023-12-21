package service

import (
	"aino_document/models"
	"time"

	"github.com/google/uuid"
)

func AddApplication(addApplication models.Application) error {

	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	app_id := currentTimestamp + int64(uniqueID)

	uuid := uuid.New()
	uuidString := uuid.String()
	_, err := db.NamedExec("INSERT INTO application_ms (application_id, application_uuid, application_code, application_title, application_description, created_by) VALUES (:application_id, :application_uuid, :application_code, :application_title, :application_description, :created_by)", map[string]interface{}{
		"application_id":          app_id,
		"application_uuid":        uuidString,
		"application_code":        addApplication.Code,
		"application_title":       addApplication.Title,
		"application_description": addApplication.Description,
		"created_by":              addApplication.Created_by,
	})
	if err != nil {
		return err
	}
	return nil
}
