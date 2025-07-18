package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"path"

	"github.com/pdcgo/shared/yenstream"
)

var CACHE_DIR = "/tmp/stream_cache"

func init() {
	os.MkdirAll(CACHE_DIR, os.ModeDir)
}

func get_fpath(fname string) string {
	return path.Join(CACHE_DIR, fname)
}

func NewFileCache[T any](
	ctx *yenstream.RunnerContext,
	fname string,
	createItem func() T,
	handler func(ctx *yenstream.RunnerContext) yenstream.Pipeline,
) yenstream.Pipeline {
	fpath := get_fpath(fname)
	_, err := os.Stat(fpath)

	var pipe yenstream.Pipeline

	if os.IsNotExist(err) {
		source := handler(ctx)

		result := source.
			Via(fname, yenstream.NewMap(ctx, func(data T) (T, error) {
				return data, nil
			}))

		f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}

		go func() {
			<-ctx.Done()
			f.Close()
		}()

		writefile := source.
			Via(fname+"writing to file", yenstream.NewMap(ctx, func(data T) (T, error) {
				raw, err := json.Marshal(data)
				if err != nil {
					return data, err
				}

				_, err = f.WriteString(string(raw) + "\n")
				if err != nil {
					return data, err
				}

				return data, nil
			})).
			Via(fname+"filter", yenstream.NewFilter(ctx, func(data T) (bool, error) {
				return false, nil
			}))

		pipe = yenstream.NewFlatten(ctx, fname+"flatten", result, writefile)
	} else {
		in := make(chan any, 1)

		go func() {
			defer close(in)
			file, err := os.Open(fpath)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Bytes()
				data := createItem()

				err := json.Unmarshal(line, data)
				if err != nil {
					panic(err)
				}
				in <- data
			}
		}()

		pipe = yenstream.NewChannelSource(ctx, in).
			Via(fname, yenstream.NewMap(ctx, func(data T) (T, error) {
				return data, nil
			}))
	}

	return pipe

}
