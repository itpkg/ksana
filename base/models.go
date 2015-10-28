package base

import (
	"time"

	ks "github.com/itpkg/ksana"
	"github.com/jinzhu/gorm"
)

type Setting struct {
	ks.Model
	Key  string `sql:"size:255;not null;unique"`
	Flag bool   `sql:"not null"`
	Val  []byte `sql:"not null"`
}

func (p Setting) TableName() string {
	return "base_settings"
}

type Locale struct {
	ks.Model
	Lang string `sql:"size:5;not null;index"`
	Key  string `sql:"size:255;not null;index"`
	Val  string `sql:"type:TEXT;not null"`
	Type string `sql:"size:8;not null;index"`
}

func (p Locale) TableName() string {
	return "base_locales"
}

type User struct {
	ks.Model
	Username    string `sql:"size:255;not null;index"`
	Email       string `sql:"size:255;not null;unique"`
	Uid         string `sql:"size:255;not null;index"`
	Provider    string `sql:"size:8;not null;index"`
	Password    string `sql:"size:255"`
	Details     []byte
	Logs        []Log
	Permissions []Permission
}

func (p User) TableName() string {
	return "base_users"
}

type Log struct {
	ID        uint `gorm:"primary_key"`
	User      User
	UserID    uint      `sql:"not null"`
	Message   string    `sql:"size:255;not null"`
	CreatedAt time.Time `sql:"not null"`
}

func (p Log) TableName() string {
	return "base_logs"
}

type Role struct {
	ks.Model

	Name         string `sql:"size:255;not null;index"`
	ResourceType string `sql:"size:255;not null;index;default:''"`
	ResourceID   uint   `sql:"not null;index;default:0"`
	Permissions  []Permission
}

func (p Role) TableName() string {
	return "base_roles"
}

type Permission struct {
	ks.Model

	User     User
	UserID   uint `sql:"not null"`
	Role     Role
	RoleID   uint      `sql:"not null"`
	StartUp  time.Time `sql:"not null;type:DATE;default:'9999-12-31'"`
	ShutDown time.Time `sql:"not null;type:DATE;default:'2015-10-27'"`
}

func (p Permission) TableName() string {
	return "base_permissions"
}

//==============================================================================

func (p *BaseEngine) Migrate() error {
	db := p.Db

	db.AutoMigrate(&Setting{}, &Locale{}, &User{}, &Log{}, &Role{}, &Permission{})
	db.Model(&Locale{}).AddUniqueIndex("idx_base_locales_lang_type_key", "lang", "type", "key")
	db.Model(&User{}).AddUniqueIndex("idx_base_users_uid_provider", "uid", "provider")
	db.Model(&Role{}).AddUniqueIndex("idx_base_roles_name_resource_type_id", "name", "resource_type", "resource_id")
	db.Model(&Permission{}).AddUniqueIndex("idx_base_permissions_role_user", "role_id", "user_id")
	return nil
}

//==============================================================================

type Dao struct {
	Aes  *ks.Aes  `inject:""`
	Hmac *ks.Hmac `inject:""`
}

func (p *Dao) Get(db *gorm.DB, key string, val interface{}) error {
	s := Setting{}
	if err := db.Where(&Setting{Key: key}).First(&s).Error; err != nil {
		return err
	}

	if s.Flag {
		buf, err := p.Aes.Decrypt(s.Val)
		if err != nil {
			return err
		}
		return ks.FromJson(buf, val)
	} else {
		return ks.FromJson(s.Val, val)
	}
}
func (p *Dao) Set(db *gorm.DB, key string, val interface{}, encrypt bool) error {
	s := Setting{}
	vb, e := ks.ToJson(val)
	if e != nil {
		return e
	}
	if encrypt {
		if vb, e = p.Aes.Encrypt(vb); e != nil {
			return e
		}
	}

	if db.Where(&Setting{Key: key, Flag: encrypt}).First(&s).RecordNotFound() {
		s.Key = key
		s.Flag = encrypt
		s.Val = vb
		db.Create(&s)
	} else {
		s.Flag = encrypt
		s.Val = vb
		db.Save(&s)
	}
	return nil
}
func (p *Dao) IsEmailUserExist(db *gorm.DB, email string) bool {
	var cn int
	db.Model(User{}).Where(&User{Email: email, Provider: "email"}).Count(&cn)
	return cn > 0
}
func (p *Dao) CreateEmailUser(db *gorm.DB, username, email, password string) (*User, error) {
	pwd, err := ks.Ssha512([]byte(password), 8)
	if err != nil {
		return nil, err
	}
	user := User{
		Username: username,
		Email:    email,
		Uid:      ks.Uuid(),
		Password: pwd,
		Provider: "email",
	}
	db.Create(&user)
	return &user, nil
}
