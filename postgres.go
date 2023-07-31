package postgres

import (
	"log"

	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var conf Config
	if err := o.UnmarshalConfig(&conf); err != nil {
		return err
	}

	return container.Container.Provide(func() (*gorm.DB, error) {
		dbConfig := postgres.Config{
			DSN: conf.DSN(), // DSN data source name
		}
		log.Println("PostgreSQL DSN: ", dbConfig.DSN)

		gormConfig := gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   conf.Prefix,
				SingularTable: conf.Singular,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   &Logger{Level: logger.Info},
		}

		db, err := gorm.Open(postgres.New(dbConfig), &gormConfig)
		if err != nil {
			return nil, err
		}

		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

		return db, err
	}, o.DiOptions()...)
}
