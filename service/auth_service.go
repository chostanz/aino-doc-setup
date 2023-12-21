package service

import (
	"aino_document/models"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ValidationError struct {
	Message string
	Field   string
	Tag     string
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

func RegisterUser(userRegister models.Register) error {
	if len(userRegister.Password) < 8 {
		return &ValidationError{
			Message: "Password should be of 8 characters long",
			Field:   "password",
			Tag:     "strong_password",
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPasswordStr := base64.StdEncoding.EncodeToString(hashedPassword)
	fmt.Println(hashedPassword)
	fmt.Println(hashedPasswordStr)

	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	user_id := currentTimestamp + int64(uniqueID)
	uuid := uuid.New()
	uuidString := uuid.String()
	_, errInsert := db.NamedExec("INSERT INTO user_ms (user_id, user_uuid, user_name, user_email, user_password, created_by) VALUES (:user_id, :user_uuid, :user_name, :user_email, :user_password, :created_by)", map[string]interface{}{
		"user_id":       user_id,
		"user_uuid":     uuidString,
		"user_name":     userRegister.Username,
		"user_email":    userRegister.Email,
		"user_password": hashedPasswordStr,
		"created_by":    userRegister.Created_by,
	})

	if errInsert != nil {
		return errInsert
	}

	err = db.Get(&user_id, "SELECT user_id FROM user_ms WHERE user_name = $1", userRegister.Username)
	if err != nil {
		return err
	}
	var roleID int64
	err = db.Get(&roleID, "SELECT role_id FROM role_ms WHERE role_uuid = $1", userRegister.ApplicationRole.Role_UUID)
	if err != nil {
		log.Println("Error getting role_id:", err)
		return err
	}
	var applicationID int64
	err = db.Get(&applicationID, "SELECT application_id FROM application_ms WHERE application_uuid = $1", userRegister.ApplicationRole.Application_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return err
	}

	// Get division_id
	var divisionID int64
	err = db.Get(&divisionID, "SELECT division_id FROM division_ms WHERE division_uuid = $1", userRegister.ApplicationRole.Division_UUID)
	if err != nil {
		log.Println("Error fetching division_id:", err)
		return err
	}

	AppRoleId := currentTimestamp + int64(uniqueID)
	uudiNew := uuid.String()
	// Insert data ke application_role_ms
	_, err = db.Exec("INSERT INTO application_role_ms(application_role_id, application_role_uuid, application_id, role_id, created_by) VALUES ($1, $2, $3, $4, $5)",
		AppRoleId, uudiNew, applicationID, roleID, userRegister.Created_by)
	if err != nil {
		log.Println("Error inserting data into application_role_ms:", err)
		return err
	}
	log.Println("Data inserted into application_role_ms successfully")

	// Get application_role_id
	var applicationRoleID int64
	err = db.Get(&applicationRoleID, "SELECT application_role_id FROM application_role_ms WHERE application_id = $1 AND role_id = $2",
		applicationID, roleID)
	if err != nil {
		log.Println("Error getting application_role_id:", err)
		return err
	}
	log.Println("Application Role ID:", applicationRoleID)

	// Insert user_application_role_ms data
	_, err = db.Exec("INSERT INTO user_application_role_ms(user_application_role_uuid, user_id, application_role_id, division_id, created_by) VALUES ($1, $2, $3, $4, $5)", uuidString, user_id, applicationRoleID, divisionID, userRegister.Created_by)
	if err != nil {
		log.Println("Error inserting data into user_application_role_ms:", err)
		return err
	}

	return nil
}

func Login(userLogin models.Login) (int, bool, error) {
	var isAuthentication bool
	var user_id int
	// var application_role_id int
	// var division_id int

	rows, err := db.Query("SELECT CASE WHEN COUNT(*) > 0 THEN 'true' ELSE 'false' END FROM user_ms WHERE user_email = $1 AND user_password = $2", userLogin.Username, userLogin.Password)
	if err != nil {
		return 0, false, err
	}

	defer rows.Close()

	rows, err = db.Query("SELECT user_id, user_password from user_ms where user_name = $1", userLogin.Username)
	if err != nil {
		fmt.Println("Error querying users:", err)
		return 0, false, err
	}

	defer rows.Close()

	var dbPasswordBase64 string
	if rows.Next() {
		err = rows.Scan(&user_id, &dbPasswordBase64)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return 0, false, err
		}
		dbPassword, errBycript := base64.StdEncoding.DecodeString(dbPasswordBase64)

		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return 0, false, errBycript
		}
		errBycript = bcrypt.CompareHashAndPassword(dbPassword, []byte(userLogin.Password))
		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return 0, false, errBycript
		}
		isAuthentication = true
	}

	if isAuthentication {
		rows, err = db.Query("SELECT application_role_id, division_id FROM user_application_role_ms WHERE user_id = $1", user_id)
		if err != nil {
			fmt.Println("Error querying user roles:", err)
			return 0, false, err
		}
		defer rows.Close()

		// if rows.Next() {
		// 	err = rows.Scan(&application_role_id, &division_id)
		// 	if err != nil {
		// 		fmt.Println("Error scanning role row:", err)
		// 		return 0, false, 0, 0, err
		// 	}
		// }
		return user_id, isAuthentication, nil
	}
	return 0, false, nil // Jika tidak ada authentikasi yang berhasil

}
