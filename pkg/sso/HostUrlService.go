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

package sso

import (
	"github.com/argoproj/argo-cd/util/session"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"time"
)

type HostUrlService interface {
	CreateHostUrl(request *HostUrlDto) (*HostUrlDto, error)
	UpdateHostUrl(request *HostUrlDto) (*HostUrlDto, error)
	GetById(id int) (*HostUrlDto, error)
	GetActiveList() ([]*HostUrlDto, error)
	GetByKey(key string) (*HostUrlDto, error)
}

type HostUrlDto struct {
	Id     int    `json:"id"`
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Active bool   `json:"active"`
	UserId int32  `json:"-"`
}

type HostUrlServiceImpl struct {
	sessionManager    *session.SessionManager
	hostUrlRepository repository.HostUrlRepository
	logger            *zap.SugaredLogger
}

func NewHostUrlServiceImpl(hostUrlRepository repository.HostUrlRepository, sessionManager *session.SessionManager,
	logger *zap.SugaredLogger) *HostUrlServiceImpl {
	serviceImpl := &HostUrlServiceImpl{
		hostUrlRepository: hostUrlRepository,
		sessionManager:    sessionManager,
		logger:            logger,
	}
	return serviceImpl
}

func (impl HostUrlServiceImpl) CreateHostUrl(request *HostUrlDto) (*HostUrlDto, error) {
	dbConnection := impl.hostUrlRepository.GetConnection()
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback tx on error.
	defer tx.Rollback()

	existingModel, err := impl.hostUrlRepository.FindByKey(request.Key)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error in update new host url", "error", err)
		return nil, err
	}
	if existingModel != nil && existingModel.Id > 0 {
		existingModel.Active = false
		err = impl.hostUrlRepository.Update(existingModel, tx)
		if err != nil {
			impl.logger.Errorw("error in creating new host url", "error", err)
			return nil, err
		}
	}

	model := &repository.HostUrl{
		Key:   request.Key,
		Value: request.Value,
	}
	model.Active = true
	model.CreatedBy = request.UserId
	model.UpdatedBy = request.UserId
	model.CreatedOn = time.Now()
	model.UpdatedOn = time.Now()
	_, err = impl.hostUrlRepository.Save(model, tx)
	if err != nil {
		impl.logger.Errorw("error in creating new host url", "error", err)
		return nil, err
	}
	request.Id = model.Id
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (impl HostUrlServiceImpl) UpdateHostUrl(request *HostUrlDto) (*HostUrlDto, error) {
	dbConnection := impl.hostUrlRepository.GetConnection()
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback tx on error.
	defer tx.Rollback()

	model, err := impl.hostUrlRepository.FindById(request.Id)
	if err != nil {
		impl.logger.Errorw("error in update new host url", "error", err)
		return nil, err
	}

	model = &repository.HostUrl{
		Key:   request.Key,
		Value: request.Value,
	}
	model.Active = true
	model.CreatedBy = request.UserId
	model.UpdatedBy = request.UserId
	model.CreatedOn = time.Now()
	model.UpdatedOn = time.Now()
	_, err = impl.hostUrlRepository.Save(model, tx)
	if err != nil {
		impl.logger.Errorw("error in creating new host url", "error", err)
		return nil, err
	}
	request.Id = model.Id
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (impl HostUrlServiceImpl) GetById(id int) (*HostUrlDto, error) {
	model, err := impl.hostUrlRepository.FindById(id)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error in update new host url", "error", err)
		return nil, err
	}
	if err == pg.ErrNoRows {
		return nil, nil
	}
	ssoLoginDto := &HostUrlDto{
		Id:     model.Id,
		Active: model.Active,
		Key:    model.Key,
		Value:  model.Value,
	}
	return ssoLoginDto, nil
}

func (impl HostUrlServiceImpl) GetActiveList() ([]*HostUrlDto, error) {
	models, err := impl.hostUrlRepository.FindActive()
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error", "error", err)
		return nil, err
	}
	if err == pg.ErrNoRows {
		return nil, nil
	}
	results := make([]*HostUrlDto, 0)
	for _, model := range models {
		dto := &HostUrlDto{
			Id:     model.Id,
			Active: model.Active,
			Key:    model.Key,
			Value:  model.Value,
		}
		results = append(results, dto)
	}
	return results, nil
}

func (impl HostUrlServiceImpl) GetByKey(key string) (*HostUrlDto, error) {
	model, err := impl.hostUrlRepository.FindByKey(key)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error", "error", err)
		return nil, err
	}
	if err == pg.ErrNoRows {
		return nil, nil
	}
	dto := &HostUrlDto{
		Id:     model.Id,
		Active: model.Active,
		Key:    model.Key,
		Value:  model.Value,
	}
	return dto, nil
}
