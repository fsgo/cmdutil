// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/12

package cmdutils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type OSEnv struct {
	values []string
}

var caseInsensitiveEnv = runtime.GOOS == "windows"

func envKey(k string) string {
	if caseInsensitiveEnv {
		return strings.ToLower(k)
	}
	return k
}

func (oe *OSEnv) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	key = envKey(key)
	vs := oe.Environ()
	for _, kv := range vs {
		idx := strings.Index(kv, "=")
		if idx < 1 {
			continue
		}
		k := envKey(kv[:idx])
		if k == key {
			return kv[idx+1:]
		}
	}
	return ""
}

func (oe *OSEnv) Set(key string, value string) error {
	if len(key) == 0 {
		return errors.New("empty key")
	}
	key = envKey(key)
	vs := oe.Environ()
	var found bool
	for i, kv := range vs {
		idx := strings.Index(kv, "=")
		if idx < 1 {
			if envKey(kv) == key {
				found = true
				vs[i] = key + "=" + value
				break
			}
			continue
		}
		k := envKey(kv[:idx])
		if k == key {
			found = true
			vs[i] = key + "=" + value
			break
		}
	}
	if !found {
		vs = append(vs, key+"="+value)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

func (oe *OSEnv) Append(key string, value string) error {
	if len(key) == 0 {
		return errors.New("empty key")
	}
	key = envKey(key)
	vs := oe.Environ()
	var found bool
	for i, kv := range vs {
		idx := strings.Index(kv, "=")
		if idx < 1 {
			if envKey(kv) == key {
				found = true
				vs[i] = key + "=" + value
				break
			}
			continue
		}
		k := envKey(kv[:idx])
		if k == key {
			found = true
			vs[i] = kv + string(filepath.ListSeparator) + value
			break
		}
	}
	if !found {
		vs = append(vs, key+"="+value)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

func (oe *OSEnv) Insert(key string, value string) error {
	if len(key) == 0 {
		return errors.New("empty key")
	}
	key = envKey(key)
	vs := oe.Environ()
	var found bool
	for i, kv := range vs {
		idx := strings.Index(kv, "=")
		if idx < 1 {
			if envKey(kv) == key {
				found = true
				vs[i] = key + "=" + value
				break
			}
			continue
		}
		k := envKey(kv[:idx])
		if k == key {
			found = true
			vs[i] = kv[:idx+1] + value + string(filepath.ListSeparator) + kv[idx+1:]
			break
		}
	}
	if !found {
		vs = append(vs, key+"="+value)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

func (oe *OSEnv) Delete(key string) error {
	if len(key) == 0 {
		return errors.New("empty key")
	}
	key = envKey(key)
	vs := oe.Environ()
	found := -1
	for i, kv := range vs {
		idx := strings.Index(kv, "=")
		if idx < 1 {
			if envKey(kv) == key {
				found = i
				break
			}
			continue
		}
		k := envKey(kv[:idx])
		if k == key {
			found = i
			break
		}
	}

	if found > -1 {
		vs = append(vs[:found], vs[found+1:]...)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

func (oe *OSEnv) WithEnviron(es []string) {
	oe.values = es
}

func (oe *OSEnv) Environ() []string {
	if len(oe.values) == 0 {
		return os.Environ()
	}
	return oe.values
}

func (oe *OSEnv) valuesUnique(vs []string) []string {
	for i, kv := range vs {
		idx := strings.Index(kv, "=")
		if idx < 1 {
			continue
		}
		v := oe.unique(kv[idx+1:])
		vs[i] = kv[:idx+1] + v
	}
	return vs
}

func (oe *OSEnv) unique(str string) string {
	arr := strings.Split(str, string(filepath.ListSeparator))
	result := make([]string, 0, len(arr))
	ks := make(map[string]bool, len(arr))
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if len(v) == 0 || ks[v] {
			continue
		}
		result = append(result, v)
		ks[v] = true
	}
	return strings.Join(result, string(filepath.ListSeparator))
}
