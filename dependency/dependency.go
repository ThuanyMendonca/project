package dependency

import (
	"net/http"
	"time"

	dbInit "github.com/ThuanyMendonca/project/config/db"
	"github.com/ThuanyMendonca/project/config/dbTransaction"
	"github.com/ThuanyMendonca/project/config/env"
	db "github.com/ThuanyMendonca/project/config/migrations"
	"github.com/ThuanyMendonca/project/repository"
	"github.com/ThuanyMendonca/project/service/authorization"
)

// Repositories
var (
	UserRepository        repository.IUserRepository
	TypeRepository        repository.ITypeRepository
	BalanceRepository     repository.IBalanceRepository
	TransactionRepository repository.ITransactionRepository
)

// Db transaction
var (
	DbTransaction dbTransaction.IDbTransaction
)

// Services
var (
	AuthorizationService authorization.IAuthorizationService
)

func Load() error {
	// Db
	projectDb, err := dbInit.InitMySqlDb(env.ProjectDb.Host, env.ProjectDb.Port, env.ProjectDb.User, env.ProjectDb.Name, env.ProjectDb.Password, env.ProjectDb.TimeZone)
	if err != nil {
		return err
	}

	db.Load(projectDb)

	// Repositories
	TypeRepository = repository.NewTypeRepository(projectDb)
	UserRepository = repository.NewUserRepository(projectDb)
	TransactionRepository = repository.NewTransactionRepository(projectDb)
	BalanceRepository = repository.NewBalanceRepository(projectDb)
	DbTransaction = dbTransaction.NewDbTransaction(projectDb)

	client := http.Client{
		Timeout: time.Duration(2) * time.Second,
	}

	// Services
	AuthorizationService = authorization.NewAuthorizationService(client)

	return nil
}
