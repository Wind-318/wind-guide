// @Author: Wind
// @file discovery.go
// @brief: discovery controller.

// Copyright (c) 2023 Wind. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package controllers

import (
	"wind-guide/config"
	"wind-guide/protobuf_data"
	"wind-guide/response"
	"wind-guide/services"

	"github.com/Wind-318/wind-chimes/logger"
	"github.com/lesismal/arpc"
	"google.golang.org/protobuf/proto"
)

// DiscoveryService is used to discover a service address by service name and unique id.
func DiscoveryService(ctx *arpc.Context) {
	data := &protobuf_data.DiscoveryRequest{}
	bin := ctx.Body()

	// Unmarshal request data.
	if err := proto.Unmarshal(bin, data); err != nil {
		services.SendErrorResponse(ctx, response.BadRequestCode, response.BadRequestMessage, err)
		return
	}
	// Check request data.
	if len(data.ServiceName) == 0 || len(data.UniqueId) == 0 {
		services.SendErrorResponse(ctx, response.BadRequestCode, response.BadRequestMessage, nil)
		return
	}

	// Record log.
	logger.Logger.Info().
		Str(response.RequestID, data.UniqueId).
		Str(response.ServiceID, data.CallerUniqueId).
		Str(response.ServiceName, data.CallerServiceName).
		Str(response.ServiceAddr, data.CallerServiceAddr).
		Str(response.ServicePort, data.CallerServicePort).
		Str(response.ServiceVersion, data.CallerServiceVersion).
		Msg(response.SuccessMessage)

	// Get service info from local cache.
	var minn int64 = response.MaxNumber
	result := &protobuf_data.RegisterRequest{}
	config.LocalCache.Mutex.Lock()
	defer config.LocalCache.Mutex.Unlock()

	// Check if service name exists.
	if _, ok := config.LocalCache.Infos[data.ServiceName]; !ok {
		services.SendErrorResponse(ctx, response.BadRequestCode, response.BadRequestMessage, nil)
		return
	} else {
		// Choose a service address by least connections.
		for _, info := range config.LocalCache.Infos[data.ServiceName] {
			if info.UsageCount < minn && info.ServiceVersion == data.Version {
				minn = info.UsageCount
				result = info
			}
		}
		// Add used times.
		result.UsageCount++
	}

	// Marshal response data.
	bin, err := proto.Marshal(result)

	// Send response.
	if err != nil {
		services.SendErrorResponse(ctx, response.MarshalFailedCode, response.MarshalFailedMessage, err)
		return
	}
	// Record service address.
	logger.Logger.Info().Msgf("Caller service %+v, call service %+v", data, result)

	ctx.Write(bin)
}
