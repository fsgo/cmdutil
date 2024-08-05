// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/12

package cmdutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// OSEnv 系统环境变量
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

// Get 读取环境变量的值，若 key 不存在则返回空字符串
//
//	PATH =  /home/work/bin:/opt/bin
//	Get("PATH")  ---> "/home/work/bin:/opt/bin"
func (oe *OSEnv) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	key = envKey(key)
	vs := oe.Environ()
	for _, kv := range vs {
		k, v, found := strings.Cut(kv, "=")
		if !found {
			continue
		}
		if envKey(k) == key {
			return v
		}
	}
	return ""
}

// GetValues 读取环境变量的值，并采用 ListSeparator 分割
//
//	PATH =  /home/work/bin:/opt/bin
//	GetValues("PATH")  ---> ["/home/work/bin", "/opt/bin"]
func (oe *OSEnv) GetValues(key string) []string {
	str := oe.Get(key)
	if str == "" {
		return nil
	}
	return strings.Split(str, string(filepath.ListSeparator))
}

var errOEEmptyKey = errors.New("empty key")

// Set 设置环境变量，若有相同的 key，则整体覆盖掉
//
//	PATH = /home/work/bin:/opt/bin
//	Set("PATH","/home/root/bin")
//	--->
//	PATH = /home/root/bin
func (oe *OSEnv) Set(key string, value string) error {
	if len(key) == 0 {
		return errOEEmptyKey
	}
	key = envKey(key)
	vs := oe.Environ()
	var found bool
	for i, kv := range vs {
		k, _, has := strings.Cut(kv, "=")
		if has {
			found = envKey(k) == key
		} else {
			//  k1=v1;k2;k3=v3 中的 k2
			found = envKey(kv) == key
		}
		if found {
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

// TrySet Set 的别名，忽略错误
func (oe *OSEnv) TrySet(key string, value string) {
	_ = oe.Set(key, value)
}

// MustSet Set 的别名，若有错误会 panic
func (oe *OSEnv) MustSet(key string, value string) {
	if err := oe.Set(key, value); err != nil {
		panic(fmt.Errorf("oe.Set(%q,%q): %w", key, value, err))
	}
}

// Append 设置环境变量，若有相同的 key，则将 value 补充在最后
//
//	PATH = /home/work/bin:/opt/bin
//	Append("PATH","/home/root/bin")
//	--->
//	PATH =  /home/work/bin:/opt/bin:/home/root/bin
func (oe *OSEnv) Append(key string, value string) error {
	if len(key) == 0 {
		return errOEEmptyKey
	}
	key = envKey(key)
	vs := oe.Environ()
	var found bool
	for i, kv := range vs {
		k, _, has := strings.Cut(kv, "=")
		if has {
			if envKey(k) == key {
				found = true
				vs[i] = kv + string(filepath.ListSeparator) + value
				break
			}
		} else {
			if envKey(kv) == key {
				found = true
				vs[i] = key + "=" + value
				break
			}
		}
	}
	if !found {
		vs = append(vs, key+"="+value)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

// TryAppend Append 的别名，忽略错误
func (oe *OSEnv) TryAppend(key string, value string) {
	_ = oe.Append(key, value)
}

// MustAppend Append 的别名，若有错误会 panic
func (oe *OSEnv) MustAppend(key string, value string) {
	if err := oe.Append(key, value); err != nil {
		panic(fmt.Errorf("oe.Append(%q,%q): %w", key, value, err))
	}
}

// Insert 设置环境变量，若有相同的 key，则将 value 插入到最前面
//
//	PATH = /home/work/bin:/opt/bin
//	Insert("PATH","/home/root/bin")
//	--->
//	PATH =  /home/root/bin:/home/work/bin:/opt/bin
func (oe *OSEnv) Insert(key string, value string) error {
	if len(key) == 0 {
		return errOEEmptyKey
	}
	key = envKey(key)
	vs := oe.Environ()
	var found bool
	for i, kv := range vs {
		k, v, has := strings.Cut(kv, "=")
		if has {
			if envKey(k) == key {
				found = true
				vs[i] = k + "=" + value + string(filepath.ListSeparator) + v
				break
			}
		} else {
			if envKey(kv) == key {
				found = true
				vs[i] = key + "=" + value
				break
			}
		}
	}
	if !found {
		vs = append(vs, key+"="+value)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

// TryInsert Insert 的别名，忽略错误
func (oe *OSEnv) TryInsert(key string, value string) {
	_ = oe.Insert(key, value)
}

// MustInsert Insert 的别名，若有错误会 panic
func (oe *OSEnv) MustInsert(key string, value string) {
	if err := oe.Insert(key, value); err != nil {
		panic(fmt.Errorf("oe.Insert(%q,%q): %w", key, value, err))
	}
}

// Delete 删除指定 key 的环境变量
//
//	PATH =  /home/work/bin:/opt/bin
//	Delete("PATH")
//	Get("PATH")  ---> ""
func (oe *OSEnv) Delete(key string) error {
	if len(key) == 0 {
		return errOEEmptyKey
	}
	key = envKey(key)
	vs := oe.Environ()
	found := -1
	for i, kv := range vs {
		k, _, has := strings.Cut(kv, "=")
		if has {
			if envKey(k) == key {
				found = i
				break
			}
		} else {
			if envKey(kv) == key {
				found = i
				break
			}
		}
	}

	if found > -1 {
		vs = append(vs[:found], vs[found+1:]...)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

// TryDelete Delete 的别名，忽略错误
func (oe *OSEnv) TryDelete(key string) {
	_ = oe.Delete(key)
}

// MustDelete Delete 的别名，若有错误会 panic
func (oe *OSEnv) MustDelete(key string) {
	if err := oe.Delete(key); err != nil {
		panic(fmt.Errorf("oe.Delete(%q): %w", key, err))
	}
}

// DeleteValue 删除环境变量中 key 的 指定的值
//
//	PATH =  /home/work/bin:/opt/bin
//	DeleteValue("PATH","/opt/bin")
//	Get("PATH")  ---> "/home/work/bin"
func (oe *OSEnv) DeleteValue(key string, value string) error {
	if len(key) == 0 {
		return errOEEmptyKey
	}
	key = envKey(key)
	vs := oe.Environ()
	deleteIndex := -1 // 需要删除的索引 id
	for i, kv := range vs {
		k, v, has := strings.Cut(kv, "=")
		if !has || envKey(k) != key {
			continue
		}
		arr := strings.Split(v, string(filepath.ListSeparator))
		arrNew := make([]string, 0, len(arr))
		for _, item := range arr {
			if item != value {
				arrNew = append(arrNew, item)
			}
		}
		if len(arrNew) == 0 {
			deleteIndex = i
		} else {
			vs[i] = k + "=" + strings.Join(arrNew, string(filepath.ListSeparator))
		}
		break
	}
	if deleteIndex > -1 {
		vs = append(vs[:deleteIndex], vs[deleteIndex+1:]...)
	}
	oe.values = oe.valuesUnique(vs)
	return nil
}

// TryDeleteValue DeleteValue 的别名，忽略错误
func (oe *OSEnv) TryDeleteValue(key string, value string) {
	_ = oe.DeleteValue(key, value)
}

// MustDeleteValue DeleteValue 的别名，若有错误会 panic
func (oe *OSEnv) MustDeleteValue(key string, value string) {
	if err := oe.DeleteValue(key, value); err != nil {
		panic(fmt.Errorf("oe.DeleteValue(%q,%q): %w", key, value, err))
	}
}

// WithEnviron 设置初始化的环境变量信息
//
// 若不调用该方法，会自动采用 os.Environ()作为默认值
func (oe *OSEnv) WithEnviron(es []string) {
	oe.values = es
}

// Environ 读取最终的环境变量信息
func (oe *OSEnv) Environ() []string {
	if len(oe.values) == 0 {
		return os.Environ()
	}
	return oe.values
}

func (oe *OSEnv) valuesUnique(vs []string) []string {
	for i, kv := range vs {
		k, v, found := strings.Cut(kv, "=")
		if !found {
			continue
		}
		vs[i] = k + "=" + oe.unique(v)
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
