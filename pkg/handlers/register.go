package handlers

import (
	"errors"
	"service-reg/pkg/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func Register(user *model.User) (*model.User, error) {
	pswErr := checkPassword(user.Password)
	if pswErr != nil {
		return user, pswErr
	}

	user.Password, _ = hashPassword(user.Password)
	user.Id = uuid.New().String()

	return storeUser(user)
}

func hashPassword(password string) (string, error) {
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// return err == nil

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPassword(psw string) error {
	if len(psw) < 10 {
		return errors.New("password is too short")
	}
	return nil
}

func getDb() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}

func storeUser(user *model.User) (*model.User, error) {
	userAttributes, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Panic().Msgf("Got error marshalling user: %s", err)
		return user, err
	}
	userUniqueLoginKeyAttributes, err := dynamodbattribute.MarshalMap(map[string]string{
		"id":    "login#" + user.Login,
		"login": "login#" + user.Login,
	})
	if err != nil {
		log.Panic().Msgf("Got error marshalling user: %s", err)
		return user, err
	}

	db := getDb()

	twii := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			&dynamodb.TransactWriteItem{
				Put: &dynamodb.Put{
					Item:                userUniqueLoginKeyAttributes,
					ConditionExpression: aws.String("attribute_not_exists(id)"),
					TableName:           aws.String("KoH_UsersTable"),
				},
			},
			&dynamodb.TransactWriteItem{
				Put: &dynamodb.Put{
					Item:                userAttributes,
					ConditionExpression: aws.String("attribute_not_exists(login)"),
					TableName:           aws.String("KoH_UsersTable"),
				},
			},
		},
	}

	_, err = db.TransactWriteItems(twii)

	if err != nil {
		log.Warn().Msgf("Error during new user creation: %s", err.Error())
		return user, errors.New("user with the same login already exists")
	}

	user.Password = ""

	return user, nil
}
