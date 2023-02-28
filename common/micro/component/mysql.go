package component

import (
	"context"

	"go-micro/common/micro"
	"go-micro/common/mysql"
	"github.com/jinzhu/gorm"
)

// MysqlComponent is Component for mysql
type MysqlComponent struct {
	micro.EmptyComponent
	db *gorm.DB
}

// Name of the component
func (c *MysqlComponent) Name() string {
	return "Mysql"
}

// PreInit called before Init()
func (c *MysqlComponent) PreInit(ctx context.Context) error {
	// load config
	mysql.SetDefaultMysqlConfig()
	return nil
}

// Init the component
func (c *MysqlComponent) Init(server *micro.Server) error {
	// init
	var err error
	mysqlConf := mysql.GetMysqlConfig()
	// spew.Dump(logConf)
	c.db, err = mysql.CreateDB(*mysqlConf)
	if err != nil {
		return err
	}
	server.RegisterElement(&micro.MysqlElementKey, c.db)
	return nil
}

// PostStop called after Stop()
func (c *MysqlComponent) PostStop(ctx context.Context) error {
	// post stop
	return c.db.Close()
}
