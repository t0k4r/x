package slogx

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

type DayFile struct {
	day  time.Time
	file *os.File
}

func OpenDayFile() (df *DayFile, err error) {
	day := time.Now()
	file, err := os.OpenFile(fmt.Sprintf("%v-%d-%v.log", day.Day(), day.Month(), day.Year()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	return &DayFile{file: file, day: day}, err
}
func (df *DayFile) Close() error {
	return df.file.Close()
}

func (df *DayFile) Write(b []byte) (n int, err error) {
	now := time.Now()
	if df.day.Sub(now).Hours() > 24 || df.day.Day() != now.Day() {
		df.Close()
		df, err = OpenDayFile()
		if err != nil {
			return n, err
		}
	}
	return df.file.Write(b)
}

func NewJson(w io.Writer, opts *slog.HandlerOptions) *slog.JSONHandler {
	df, err := OpenDayFile()
	if err != nil {
		panic(err)
	}
	return slog.NewJSONHandler(io.MultiWriter(w, df), opts)
}

func NewText(w io.Writer, opts *slog.HandlerOptions) *slog.TextHandler {
	df, err := OpenDayFile()
	if err != nil {
		panic(err)
	}
	return slog.NewTextHandler(io.MultiWriter(w, df), opts)
}
