package std

import (
	"context"
	"errors"

	"github.com/janstoon/toolbox/bricks"
	"github.com/spf13/viper"

	"github.com/janstoon/toolbox/kareless"
)

type MapSettingSource map[string]any

func (ss MapSettingSource) Get(_ context.Context, key string) (any, error) {
	v, ok := ss[key]
	if ok {
		return v, nil
	}

	return nil, bricks.ErrNotFound
}

type localEarlySettings struct {
	v *viper.Viper
}

func LocalEarlyLoadedSettingSource(name string, paths ...string) kareless.SettingSource {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigName(name)
	for _, p := range paths {
		v.AddConfigPath(p)
	}

	if len(paths) > 0 {
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				panic(errors.Join(bricks.ErrNotFound, err))
			}
		}
	}

	return localEarlySettings{
		v: v,
	}
}

func (ss localEarlySettings) Get(_ context.Context, key string) (any, error) {
	if ss.v.IsSet(key) {
		return ss.v.Get(key), nil
	}

	return nil, bricks.ErrNotFound
}
