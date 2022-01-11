package domain

import (
	"fmt"
	"forklift/internal/domain/consts"
	"forklift/internal/domain/entity"
	"forklift/pkg/enums"
	"forklift/pkg/mongo"
	"github.com/globalsign/mgo/bson"
)

type MigrationTaskRepository interface {
	GetMigrations() ([]*entity.MigrationTask, error)
	GetMigrationByMessageId(messageId string) (*entity.MigrationTask, error)
	ControlMigrationByMessageId(messageId string) error
	UpsertMigrationTask(task *entity.MigrationTask) error
	UpsertMigrationDeletedTask(task *entity.MigrationTask) error
	UpdateMigrationTask(task *entity.MigrationTask) error
	RemoveMigrationTaskByMessageId(messageId string) error
}

type migrationTaskRepository struct {
	mongoClient  mongo.Client
	databaseName string
}

func (m *migrationTaskRepository) UpdateMigrationTask(task *entity.MigrationTask) error {

	var session = m.mongoClient.NewSession()
	defer session.Close()

	err := session.
		DB(m.databaseName).
		C(consts.CollectionTasks).
		Update(bson.M{"MessageId": task.MessageId}, task)

	if err != nil {
		return err
	}

	return nil
}

func (m *migrationTaskRepository) UpsertMigrationTask(task *entity.MigrationTask) error {

	var session = m.mongoClient.NewSession()
	defer session.Close()

	_, err := session.
		DB(m.databaseName).
		C(consts.CollectionTasks).
		Upsert(bson.M{"MessageId": task.MessageId}, task)

	if err != nil {
		return err
	}

	return nil
}

func (m *migrationTaskRepository) UpsertMigrationDeletedTask(task *entity.MigrationTask) error {

	var session = m.mongoClient.NewSession()
	defer session.Close()

	task.Status = enums.Deleted

	_, err := session.
		DB(m.databaseName).
		C(consts.CollectionDeletedTasks).
		Upsert(bson.M{"MessageId": task.MessageId}, task)

	if err != nil {
		return err
	}

	return nil
}

func (m *migrationTaskRepository) GetMigrationByMessageId(messageId string) (*entity.MigrationTask, error) {
	var session = m.mongoClient.NewSession()
	defer session.Close()

	var task *entity.MigrationTask

	err := session.DB(m.databaseName).C(consts.CollectionTasks).Find(bson.M{"MessageId": messageId}).One(&task)
	if err != nil {
		return nil, err
	}

	return task, err
}

func (m *migrationTaskRepository) ControlMigrationByMessageId(messageId string) error {
	var session = m.mongoClient.NewSession()
	defer session.Close()

	var task *entity.MigrationTask

	err := session.DB(m.databaseName).C(consts.CollectionTasks).Find(bson.M{"MessageId": messageId}).One(&task)
	if err != nil {
		return err
	}

	return nil
}

func (m *migrationTaskRepository) GetMigrations() ([]*entity.MigrationTask, error) {
	var session = m.mongoClient.NewSession()
	defer session.Close()

	var records []*entity.MigrationTask

	err := session.DB(m.databaseName).C(consts.CollectionTasks).Find(nil).All(&records)

	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return records, nil
}

func (m *migrationTaskRepository) RemoveMigrationTaskByMessageId(messageId string) error {

	var session = m.mongoClient.NewSession()
	defer session.Close()

	_, err := session.
		DB(m.databaseName).
		C(consts.CollectionTasks).
		RemoveAll(bson.M{"MessageId": messageId})

	if err != nil {
		return err
	}

	return nil
}

func NewMigrationRepository(client mongo.Client, databaseName string) MigrationTaskRepository {

	if err := client.EnsureIndex([]string{"MessageId"}, true, "MessageId", databaseName, consts.CollectionTasks); err != nil {
		panic(fmt.Sprintln("could not create a new client", err))
	}

	return &migrationTaskRepository{mongoClient: client, databaseName: databaseName}
}
