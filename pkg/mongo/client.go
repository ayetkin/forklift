package mongo

import (
	"context"
	"github.com/globalsign/mgo"
	"time"
)

type client struct {
	session *mgo.Session
}

type Client interface {
	NewSession() *mgo.Session
	NewSessionWithSecondaryPreferred() *mgo.Session
	Ping(ctx context.Context) bool
	EnsureIndex(fields []string, isUnique bool, indexName string, database string, collection string) error
}

func NewClient(connectionString string, timeout time.Duration) (Client, error) {

	session, err := mgo.DialWithTimeout(connectionString, timeout)

	if err != nil {
		return nil, err
		//return nil, errors.NewWithCause(ConnectionError, err, connectionString)
	}

	return &client{session: session}, nil
}

func (c *client) NewSession() *mgo.Session {
	newSession := c.session.Copy()
	newSession.SetMode(mgo.Strong, true)
	return newSession
}

func (c *client) NewSessionWithSecondaryPreferred() *mgo.Session {
	newSession := c.session.Copy()
	newSession.SetMode(mgo.SecondaryPreferred, true)
	return newSession
}

func (c *client) Ping(ctx context.Context) bool {
	localSession := c.NewSession()
	defer localSession.Close()

	if err := localSession.Ping(); err != nil {
		//c.logger.Exception(ctx, errors.NewWithCause(PingError, err).Error(), err)
		return false
	}

	return true
}

func (c *client) EnsureIndex(fields []string, isUnique bool, indexName string, database string, collection string) error {
	localSession := c.NewSession()
	defer localSession.Close()

	index := mgo.Index{
		Key:        fields,
		Unique:     isUnique,
		Name:       indexName,
		Background: true,
	}

	col := localSession.DB(database).C(collection)

	if err := col.EnsureIndex(index); err != nil {
		return err
		//return errors.NewWithCause(IndexCreationError, err, strings.Join(fields, ","), database, collection)
	}
	return nil
}
