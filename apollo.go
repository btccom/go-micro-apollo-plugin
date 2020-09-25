package apollo

import (
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/zouyx/agollo/v3"
	"github.com/zouyx/agollo/v3/env/config"
	"github.com/zouyx/agollo/v3/storage"
	"strings"
	"time"
)

type apolloSource struct {
	namespaceName  string
	customConfig   *CustomConfig
	isCustomConfig bool
	opts           source.Options
}

type CustomConfig struct {
	isBackupConfig   bool
	backupConfigPath string
	appID            string
	cluster          string
	ip               string
}

func (a *apolloSource) String() string {
	return "apollo"
}

func (a *apolloSource) Read() (*source.ChangeSet, error) {
	if a.isCustomConfig {
		readyConfig := &config.AppConfig{
			IsBackupConfig:   a.customConfig.isBackupConfig,
			BackupConfigPath: a.customConfig.backupConfigPath,
			AppID:            a.customConfig.appID,
			Cluster:          a.customConfig.cluster,
			NamespaceName:    a.namespaceName,
			IP:               a.customConfig.ip,
		}
		agollo.InitCustomConfig(func() (*config.AppConfig, error) {
			return readyConfig, nil
		})
	}

	if err := agollo.Start(); err != nil {
		log.Error(err)
		return nil, err
	}
	c := agollo.GetConfig(a.namespaceName)

	var format string
	var content []byte

	namespaceParts := strings.Split(a.namespaceName, ".")
	if len(namespaceParts) > 1 {
		content = []byte(c.GetValue("content"))
		format = namespaceParts[len(namespaceParts)-1]
	} else {
		changes := make(map[string]interface{})
		content = []byte(c.GetContent("properties"))
		list := strings.Split(string(content), "\n")
		for _, env := range list {
			pair := strings.SplitN(env, "=", 2)
			if len(pair) < 2 {
				continue
			}
			changes[pair[0]] = interface{}(pair[1])
		}
		b, err := a.opts.Encoder.Encode(changes)
		if err != nil {
			return nil, err
		}
		content = b
		format = a.opts.Encoder.String()
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    format,
		Source:    a.String(),
		Data:      content,
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (a *apolloSource) Watch() (source.Watcher, error) {
	watcher, err := newWatcher(a.String(), a.namespaceName)
	storage.AddChangeListener(watcher)
	return watcher, err
}

func (a *apolloSource) Write(cs *source.ChangeSet) error {
	return nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	var nName string
	if p, ok := options.Context.Value(namespaceName{}).(string); ok {
		nName = p
	}

	var isCustom bool
	if p, ok := options.Context.Value(isCustomConfig{}).(bool); ok {
		isCustom = p
	}

	var cConfig *CustomConfig
	if p, ok := options.Context.Value(customConfig{}).(*CustomConfig); ok {
		cConfig = p
	}

	return &apolloSource{opts: options, namespaceName: nName, isCustomConfig: isCustom, customConfig: cConfig}
}
