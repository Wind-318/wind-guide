package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"wind-guide/protobuf_data"
	"wind-guide/response"

	"github.com/Wind-318/wind-chimes/logger"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

var (
	ConfigSettings *Config
	LocalCache     *ServiceData
	configPath     = "config/config.json"
)

type ServiceData struct {
	Infos map[string]map[string]*protobuf_data.RegisterRequest `json:"infos"`
	Mutex *sync.RWMutex                                        `json:"mutex"`
}

type ServerConfig struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	HeartbeatInterval int64  `json:"heartbeat_interval"`
	HeartbeatTimeout  int64  `json:"heartbeat_timeout"`
	EnableHealthCheck bool   `json:"enable_health_check"`
}

type LoggingConfig struct {
	Level     string `json:"level"`
	FilePath  string `json:"file_path"`
	MaxSize   int    `json:"max_size"`
	MaxAge    int    `json:"max_age"`
	MaxBack   int    `json:"max_back"`
	Compress  bool   `json:"compress"`
	CallDepth int    `json:"call_depth"`
}

type RoutesConfig struct {
	Path    string `json:"path"`
	Handler string `json:"handler"`
}

type Config struct {
	Server  ServerConfig  `json:"server"`
	Logging LoggingConfig `json:"logging"`
	Routes  []RoutesConfig
}

// Read config file.
func ReadConfig() error {
	// Initialize local cache.
	LocalCache = &ServiceData{
		Infos: make(map[string]map[string]*protobuf_data.RegisterRequest),
		Mutex: &sync.RWMutex{},
	}
	ConfigSettings = &Config{}

	// Read config file.
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	err = json.Unmarshal(configFile, ConfigSettings)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	// Generate server id.
	ConfigSettings.Server.ID = uuid.New().String()

	// Initialize logger.
	logger.InitLogger(ConfigSettings.Logging.FilePath, ConfigSettings.Logging.MaxSize, ConfigSettings.Logging.MaxAge,
		ConfigSettings.Logging.MaxBack, ConfigSettings.Logging.CallDepth, true, map[string]interface{}{
			response.ServiceID:   ConfigSettings.Server.ID,
			response.ServiceName: ConfigSettings.Server.Name,
			response.ServiceAddr: ConfigSettings.Server.Host,
			response.ServicePort: ConfigSettings.Server.Port,
		})

	logger.Logger.Info().Msgf("%+v", ConfigSettings)

	if ConfigSettings.Server.EnableHealthCheck {
		// Start heartbeat.
		go healthCheck()
	}

	return nil
}

// Heartbeat every HeartbeatInterval seconds for each service.
func healthCheck() {
	for {
		// Sleep.
		time.Sleep(time.Duration(ConfigSettings.Server.HeartbeatInterval) * time.Second)

		// Record the service that needs to be deleted.
		var deleteIDs []string
		deleteNames := map[string]interface{}{}

		// Check the service status.
		LocalCache.Mutex.RLock()
		for name, serviceInfos := range LocalCache.Infos {
			for id, serviceInfo := range serviceInfos {
				// Use fasthttp to send heartbeat request.
				statusCode, _, err := fasthttp.GetTimeout(nil, serviceInfo.HealthCheckUrl, time.Duration(ConfigSettings.Server.HeartbeatTimeout)*time.Second)
				if err != nil || statusCode != fasthttp.StatusOK {
					logger.Logger.Error().Msgf("Service %s(%s) is not available, delete it", name, id)
					deleteIDs = append(deleteIDs, id)
					deleteNames[name] = nil
				}
			}
		}
		LocalCache.Mutex.RUnlock()

		// Delete the service that is not available.
		LocalCache.Mutex.Lock()
		for name := range deleteNames {
			for index := range deleteIDs {
				// Record the service that needs to be deleted.
				logger.Logger.Info().Msgf("Delete service %s(%s)", name, deleteIDs[index])
				// Delete the service from local cache.
				delete(LocalCache.Infos[name], deleteIDs[index])
			}
			if len(LocalCache.Infos[name]) == 0 {
				// Record the service that needs to be deleted.
				logger.Logger.Info().Msgf("Delete service %s", name)
				// Delete the service from local cache.
				delete(LocalCache.Infos, name)
			}
		}
		LocalCache.Mutex.Unlock()
	}
}
