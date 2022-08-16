package dependency

import (
	dbInit "github.com/ThuanyMendonca/project/config/db"
	"github.com/ThuanyMendonca/project/config/env"
	db "github.com/ThuanyMendonca/project/config/migrations"
	"github.com/ThuanyMendonca/project/repository"
)

var (
	UserRepository repository.IUserRepository
)

var (
	typeRep repository.TypeRepository
)

func Load() error {
	projectDb, err := dbInit.InitMySqlDb(env.ProjectDb.Host, env.ProjectDb.Port, env.ProjectDb.User, env.ProjectDb.Name, env.ProjectDb.Password, env.ProjectDb.TimeZone)
	if err != nil {
		return err
	}

	db.Load(projectDb)

	UserRepository = repository.NewUserRepository(projectDb, typeRep)

	return nil
}
