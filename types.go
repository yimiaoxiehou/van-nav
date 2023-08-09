package main

import (
	"database/sql"
	"encoding/json"
)

type loginDto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// 默认是 0
type Setting struct {
	Id      int    `json:"id"`
	Favicon string `json:"favicon"`
	Title   string `json:"title"`
	Logo192 string `json:"logo192"`
	Logo512 string `json:"logo512"`
}

type Token struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Disabled int    `json:"disabled"`
}
type AddTokenDto struct {
	Name string `json:"name"`
}
type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
type Img struct {
	Id    int    `json:"id"`
	Url   string `json:"url"`
	Value string `json:"value"`
}

type resUserDto struct {
	Name string `json:"name"`
}

type updateUserDto struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Tool struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Logo     string `json:"logo"`
	Catelog  string `json:"catelog"`
	Desc     string `json:"desc"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type addToolDto struct {
	Name     string         `json:"name"`
	Url      string         `json:"url"`
	Logo     string         `json:"logo"`
	Catelog  string         `json:"catelog"`
	Desc     string         `json:"desc"`
	Username sql.NullString `json:"username"`
	Password sql.NullString `json:"password"`
}

func (d addToolDto) UnmarshalJSON(data []byte) error {

	var t struct {
		Name     string `json:"name"`
		Url      string `json:"url"`
		Logo     string `json:"logo"`
		Catelog  string `json:"catelog"`
		Desc     string `json:"desc"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	d.Name = t.Name
	d.Url = t.Url
	d.Logo = t.Logo
	d.Catelog = t.Catelog
	d.Desc = t.Desc
	if t.Username != "" {
		d.Username = sql.NullString{
			String: t.Username,
			Valid:  true,
		}
	} else {
		d.Username = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	if t.Password != "" {
		d.Password = sql.NullString{
			String: t.Password,
			Valid:  true,
		}
	} else {
		d.Password = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return nil
}

type updateToolDto struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Logo     string `json:"logo"`
	Catelog  string `json:"catelog"`
	Desc     string `json:"desc"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type updateCatelogDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type addCatelogDto struct {
	Name string `json:"name"`
}

type Catelog struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
