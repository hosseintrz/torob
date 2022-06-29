package db

type DBType string

const (
	MONGODB  DBType = "mongodb"
	DYNAMODB DBType = "dynamodb"
)

type DataBase interface {
	Connect() error
}

func NewDatabase(dbType DBType, conn string) DataBase {
	switch dbType {
	case MONGODB:
		return NewMongoDB(conn)
	case DYNAMODB:
		//not implemented
	}
	return nil
}
