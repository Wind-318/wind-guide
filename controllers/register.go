// @Author: Wind
// @file register.go
// @brief: register controller.

// Copyright (c) 2023 Wind. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Package controllers is a controller package.
package controllers

import (
	"errors"
	"wind-guide/config"
	"wind-guide/protobuf_data"
	"wind-guide/response"
	"wind-guide/services"

	"github.com/Wind-318/wind-chimes/logger"
	"github.com/lesismal/arpc"
	"google.golang.org/protobuf/proto"
)

// RegisterService is a controller for registering service.
func RegisterService(ctx *arpc.Context) {
	registerData := &protobuf_data.RegisterRequest{}
	bin := ctx.Body()

	if err := proto.Unmarshal(bin, registerData); err != nil {
		services.SendErrorResponse(ctx, response.BadRequestCode, response.BadRequestMessage, err)
		return
	}

	// Check request data.
	if len(registerData.ServiceName) == 0 || len(registerData.ServiceAddr) == 0 || len(registerData.UniqueId) == 0 ||
		len(registerData.HealthCheckUrl) == 0 || len(registerData.ServicePort) == 0 || len(registerData.ServiceId) == 0 {
		services.SendErrorResponse(ctx, response.BadRequestCode, response.BadRequestMessage, errors.New(response.BadRequestMessage))
		return
	}

	// Check service id exists.
	config.LocalCache.Mutex.Lock()
	defer config.LocalCache.Mutex.Unlock()
	if _, ok := config.LocalCache.Infos[registerData.ServiceName]; !ok {
		config.LocalCache.Infos[registerData.ServiceName] = make(map[string]*protobuf_data.RegisterRequest)
	} else if _, ok := config.LocalCache.Infos[registerData.ServiceName][registerData.ServiceId]; ok {
		services.SendErrorResponse(ctx, response.ServiceIDEixstsCode, response.ServiceIDEixstsMessage, errors.New(response.ServiceIDEixstsMessage))
		return
	}

	// Add service info.
	config.LocalCache.Infos[registerData.ServiceName][registerData.ServiceId] = &protobuf_data.RegisterRequest{
		ServiceId:      registerData.ServiceId,
		ServiceName:    registerData.ServiceName,
		ServiceAddr:    registerData.ServiceAddr,
		UniqueId:       registerData.UniqueId,
		HealthCheckUrl: registerData.HealthCheckUrl,
		ServicePort:    registerData.ServicePort,
		ServiceVersion: registerData.ServiceVersion,
		UsageCount:     0,
	}

	// Record log.
	logger.Logger.Info().
		Str(response.RequestID, registerData.UniqueId).
		Str(response.ServiceID, registerData.ServiceId).
		Str(response.ServiceName, registerData.ServiceName).
		Str(response.ServiceAddr, registerData.ServiceAddr).
		Str(response.HealthCheckUrl, registerData.HealthCheckUrl).
		Str(response.ServicePort, registerData.ServicePort).
		Str(response.ServiceVersion, registerData.ServiceVersion).
		Msg(response.SuccessMessage)

	// Send response.
	res, err := proto.Marshal(&protobuf_data.RegisterResponse{
		Code:    response.SuccessCode,
		Message: response.SuccessMessage,
	})
	if err != nil {
		services.SendErrorResponse(ctx, response.MarshalFailedCode, response.MarshalFailedMessage, err)
		return
	}

	ctx.Write(res)
}
