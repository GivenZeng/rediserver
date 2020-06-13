package rediserver

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

type flagSlice []string

func (i *flagSlice) String() string {
	if len(*i) == 0 {
		return ""
	}
	return fmt.Sprint(*i)
}
func (i *flagSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// 用于暂存从命令行解析到的数据
var (
	ms = make(map[string]interface{})
)

// Parse 从配置文件和命令行解析conf
func Parse(conf interface{}, defaultConfPath string) error {
	confPath := flag.String("conf", defaultConfPath, "path to config")
	register("", conf)
	flag.Parse()
	// 先从文件加载配置, 如有
	if *confPath != "" {
		err := ParseYAMLFile(*confPath, conf)
		if err != nil {
			return err
		}
	}
	// 使用命令行覆盖配置
	return get("", conf)
}

// register 将配置的字段注册到flag
func register(prefix string, conf interface{}) {
	prefix = strings.ToLower(prefix)
	if len(prefix) > 0 {
		prefix += "."
	}

	t := reflect.TypeOf(conf)
	v := reflect.ValueOf(conf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Int64, reflect.Int:
			ms[prefix+strings.ToLower(t.Field(i).Name)] = flag.Int64(prefix+strings.ToLower(t.Field(i).Name), 0, t.Field(i).Tag.Get("comment"))
		case reflect.Float64:
			ms[prefix+strings.ToLower(t.Field(i).Name)] = flag.Float64(prefix+strings.ToLower(t.Field(i).Name), 0, t.Field(i).Tag.Get("comment"))
		case reflect.String:
			ms[prefix+strings.ToLower(t.Field(i).Name)] = flag.String(prefix+strings.ToLower(t.Field(i).Name), "", t.Field(i).Tag.Get("comment"))
		case reflect.Bool:
			ms[prefix+strings.ToLower(t.Field(i).Name)] = flag.Bool(prefix+strings.ToLower(t.Field(i).Name), false, t.Field(i).Tag.Get("comment"))
		case reflect.Slice:
			s := flagSlice{}
			ms[prefix+strings.ToLower(t.Field(i).Name)] = &s
			flag.Var(&s, prefix+strings.ToLower(t.Field(i).Name), t.Field(i).Tag.Get("comment"))
		case reflect.Struct:
			register(prefix+t.Field(i).Name, v.Field(i).Addr().Interface())
		}
	}
}

// get 获取命令行暂存数据并写入conf
func get(prefix string, conf interface{}) error {
	prefix = strings.ToLower(prefix)
	if len(prefix) > 0 {
		prefix += "."
	}
	t := reflect.TypeOf(conf)
	v := reflect.ValueOf(conf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Int64, reflect.Int:
			if val, ok := ms[prefix+strings.ToLower(t.Field(i).Name)]; ok {
				if ival, ok := val.(*int64); ok && *ival > 0 {
					v.Field(i).SetInt(*ival)
					// fmt.Println("migo.Flag|", prefix+strings.ToLower(t.Field(i).Name), *ival)
				}
			} else {
				return fmt.Errorf("flag does not existed: %s", prefix+strings.ToLower(t.Field(i).Name))
			}
		case reflect.Float64:
			if val, ok := ms[prefix+strings.ToLower(t.Field(i).Name)]; ok {
				if fval, ok := val.(*float64); ok && *fval != 0 {
					v.Field(i).SetFloat(*fval)
					// fmt.Println("migo.Flag|", prefix+strings.ToLower(t.Field(i).Name), *fval)
				}
			} else {
				return fmt.Errorf("flag does not existed: %s", prefix+strings.ToLower(t.Field(i).Name))
			}
		case reflect.String:
			if s, ok := ms[prefix+strings.ToLower(t.Field(i).Name)]; ok {
				if str, ok := s.(*string); ok && len(*str) > 0 {
					v.Field(i).SetString(*str)
					// fmt.Println("migo.Flag|", prefix+strings.ToLower(t.Field(i).Name), *str)
				}
			} else {
				return fmt.Errorf("flag does not existed: %s", prefix+strings.ToLower(t.Field(i).Name))
			}
		case reflect.Bool:
			if b, ok := ms[prefix+strings.ToLower(t.Field(i).Name)]; ok {
				if bo, ok := b.(*bool); ok && *bo {
					v.Field(i).SetBool(*bo)
					// fmt.Println("migo.Flag|", prefix+strings.ToLower(t.Field(i).Name), *bo)
				}
			} else {
				return fmt.Errorf("flag does not existed: %s", prefix+strings.ToLower(t.Field(i).Name))
			}
		case reflect.Slice:
			if f, ok := ms[prefix+strings.ToLower(t.Field(i).Name)]; ok {
				if fs, ok := f.(*flagSlice); ok && len(*fs) > 0 {
					v.Field(i).Set(reflect.ValueOf(([]string)(*fs)))
					// fmt.Println("migo.Flag|", prefix+strings.ToLower(t.Field(i).Name), *fs)
				}
			} else {
				return fmt.Errorf("flag does not existed: %s", prefix+strings.ToLower(t.Field(i).Name))
			}
		case reflect.Struct:
			get(prefix+t.Field(i).Name, v.Field(i).Addr().Interface())
		}
	}
	return nil
}
