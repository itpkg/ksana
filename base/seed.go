package base

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	ks "github.com/itpkg/ksana"
)

func (p *BaseEngine) Seed() error {
	db := p.Db
	//--------------administrator-------------
	admin_e := "root@localhost.localdomain"

	if !p.Dao.IsEmailUserExist(db, admin_e) {
		admin_u, err := p.Dao.CreateEmailUser(db, "Admin", admin_e, "changeme")
		if err != nil {
			return err
		}
		role_a := Role{Name: "admin"}
		role_r := Role{Name: "root"}
		db.Create(&role_a)
		db.Create(&role_r)

		begin := time.Now()
		end := begin.AddDate(10, 0, 0)
		db.Create(&Permission{
			RoleID:   role_a.ID,
			UserID:   admin_u.ID,
			StartUp:  begin,
			ShutDown: end,
		})
		db.Create(&Permission{
			RoleID:   role_r.ID,
			UserID:   admin_u.ID,
			StartUp:  begin,
			ShutDown: end,
		})
		db.Create(&Log{
			UserID:  admin_u.ID,
			Message: "Init.",
		})
	}
	//--------------locales-------------------
	root := fmt.Sprintf("%s/locales", ks.PkgRoot(p))
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, f := range files {
		fn := fmt.Sprintf("%s/%s", root, f.Name())
		log.Printf("Find locale file %s", fn)
		ss := strings.Split(f.Name(), ".")
		if len(ss) != 3 {
			return errors.New(fmt.Sprintf("bad locale file name %s", f.Name))
		}
		items := make(map[string]string, 0)
		if _, err := toml.DecodeFile(fn, &items); err != nil {
			return err
		}
		for k, v := range items {
			var cn int
			db.Model(Locale{}).Where(&Locale{Type: ss[0], Lang: ss[1], Key: k}).Count(&cn)
			if cn == 0 {
				db.Create(&Locale{
					Type: ss[0],
					Lang: ss[1],
					Key:  k,
					Val:  v,
				})
			}
		}

	}

	return nil
}
