/*
 * Copyright (c) 2020 Devtron Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package repository

import (
	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/go-pg/pg"
)

type HostUrl struct {
	tableName struct{} `sql:"host_url" pg:",discard_unknown_columns"`
	Id        int      `sql:"id,pk"`
	Url       string   `sql:"url,notnull"`
	Active    bool     `sql:"active, notnull"`
	models.AuditLog
}

type HostUrlRepository interface {
	Save(model *HostUrl, tx *pg.Tx) (*HostUrl, error)
	Update(model *HostUrl, tx *pg.Tx) error
	FindByUrl(url string) (*HostUrl, error)
	FindById(id int) (*HostUrl, error)
	FindActive() (*HostUrl, error)
	GetConnection() (dbConnection *pg.DB)
}

type HostUrlRepositoryImpl struct {
	dbConnection *pg.DB
}

func NewHostUrlRepositoryImpl(dbConnection *pg.DB) *HostUrlRepositoryImpl {
	return &HostUrlRepositoryImpl{dbConnection: dbConnection}
}

func (impl *HostUrlRepositoryImpl) GetConnection() (dbConnection *pg.DB) {
	return impl.dbConnection
}


func (repo HostUrlRepositoryImpl) Save(model *HostUrl, tx *pg.Tx) (*HostUrl, error) {
	err := tx.Insert(model)
	return model, err
}

func (repo HostUrlRepositoryImpl) Update(model *HostUrl, tx *pg.Tx) error {
	err := tx.Insert(model)
	return err
}

func (repo HostUrlRepositoryImpl) FindByUrl(url string) (*HostUrl, error) {
	model := &HostUrl{}
	err := repo.dbConnection.Model(model).Where("url = ?", url).Where("active = ?", true).
		Select()
	return model, err
}

func (repo HostUrlRepositoryImpl) FindById(id int) (*HostUrl, error) {
	model := &HostUrl{}
	err := repo.dbConnection.Model(model).Where("id = ?", id).Where("active = ?", true).
		Select()
	return model, err
}

func (repo HostUrlRepositoryImpl) FindActive() (*HostUrl, error) {
	model := &HostUrl{}
	err := repo.dbConnection.Model(model).Where("active = ?", true).
		Select()
	return model, err
}
