package service

import (
	"aino_document/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddApplicationRole(addAppRole models.AddApplicationRole) error {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	uuid := uuid.New()
	uuidString := uuid.String()
	app_role_id := currentTimestamp + int64(uniqueID)

	var roleID int64
	err := db.Get(&roleID, "SELECT role_id FROM role_ms WHERE role_uuid = $1", addAppRole.Role_UUID)
	if err != nil {
		log.Println("Error getting role_id:", err)
		return err
	}
	var applicationID int64
	err = db.Get(&applicationID, "SELECT application_id FROM application_ms WHERE application_uuid = $1", addAppRole.Application_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return err
	}
	_, errInsert := db.NamedExec("INSERT INTO application_role_ms(application_role_id, application_role_uuid, application_id, role_id, created_by) VALUES(:application_role_id, :application_role_uuid, :application_id, :role_id, :created_by)", map[string]interface{}{
		"application_role_id":   app_role_id,
		"application_role_uuid": uuidString,
		"application_id":        applicationID,
		"role_id":               roleID,
		"created_by":            addAppRole.Created_by,
	})

	if errInsert != nil {
		log.Print(errInsert)
		return errInsert
	}
	return nil
}
