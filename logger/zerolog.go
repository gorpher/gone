package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const IsoZonedDateTime = "2006-01-02 15:04:05"

func init() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    true,
		TimeFormat: IsoZonedDateTime,
	}).With().Caller().Logger()
}

type LogConfig struct {
	LogLevel       string `json:"log_level" yaml:"log_level"`               // 日志级别
	LogDir         string `json:"log_dir" yaml:"log_dir"`                   // 日志文件路径
	LogWithDate    bool   `json:"log_with_date" yaml:"log_with_date"`       // 按日期存放日志
	LogWithLevel   bool   `json:"log_with_level" yaml:"log_with_level"`     // 按级别存放日志
	LogWithColor   bool   `json:"log_with_color" yaml:"log_with_color"`     // 是否输出彩色日志
	LogWithConsole bool   `json:"log_with_console" yaml:"log_with_console"` // 是否输出到控制台
	LogSaveDay     int    `json:"log_save_day" yaml:"log_save_day"`         // 保存时间，单位天
}
type OptFunc func(*MultiWriter) *MultiWriter

func WithLogConfig(cfg *LogConfig) OptFunc {
	return func(s *MultiWriter) *MultiWriter {
		s.cfg = cfg
		return s
	}
}
func WithFileName(name string) OptFunc {
	return func(s *MultiWriter) *MultiWriter {
		s.name = name
		return s
	}
}

// SetZerolog
func SetZerolog(opts ...OptFunc) {
	// 设置时间格式
	zerolog.TimeFieldFormat = IsoZonedDateTime
	writer := NewMultiWriter(opts...)
	switch writer.cfg.LogLevel {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "empty":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
	if strings.ToLower(writer.cfg.LogLevel) == "debug" {
		log.Logger = log.Logger.With().Caller().Logger()
	}
}

type syncFiles struct {
	mutex sync.Mutex
	files map[string]*os.File
}

func newSyncFiles() *syncFiles {
	return &syncFiles{
		files: make(map[string]*os.File, 7),
	}
}
func (s *syncFiles) get(filename string) *os.File {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	fd := s.files[filename]
	return fd
}

func (s *syncFiles) set(filename string, fd *os.File) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.files[filename] = fd
}

func (s *syncFiles) removeAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var keys = make([]string, len(s.files))
	var i int
	for name, _ := range s.files {
		keys[i] = name
		i += 1
	}
	for _, key := range keys {
		fd := s.files[key]
		fd.Close() //nolint
		delete(s.files, key)
	}
}

type MultiWriter struct {
	name  string
	cfg   *LogConfig
	cw    *zerolog.ConsoleWriter
	date  string
	files *syncFiles
}

func (m *MultiWriter) filename(level string) string {
	var sb strings.Builder
	if m.name != "" {
		sb.WriteString(m.name + "-")
	}
	if m.cfg.LogWithDate {
		sb.WriteString(time.Now().Format("20060102") + "-")
	}
	if m.cfg.LogWithLevel {
		sb.WriteString(level)
	}
	if sb.Len() == 0 {
		sb.WriteString("output")
	}
	sb.WriteString(".log")
	return sb.String()
}

func NewMultiWriter(opts ...OptFunc) *MultiWriter {
	m := &MultiWriter{
		name: "",
		cfg: &LogConfig{
			LogLevel:       "info",
			LogDir:         "./logs",
			LogWithColor:   true,
			LogWithConsole: true,
			LogSaveDay:     7,
		},
		date: time.Now().Format("20060102"),
		cw: &zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    true,
			TimeFormat: IsoZonedDateTime,
		},
		files: newSyncFiles(),
	}
	for _, opt := range opts {
		m = opt(m)
	}
	m.cw.NoColor = !m.cfg.LogWithColor
	return m
}

func (m *MultiWriter) checkDate() {
	nowDate := time.Now().Format("20060102")
	if m.date == nowDate {
		return
	}
	files := m.files
	m.files = newSyncFiles()
	m.date = nowDate
	files.removeAll()
}

func (m *MultiWriter) get(level string) io.Writer {
	filename := m.filename(level)
	fd := m.files.get(filename)
	if fd != nil {
		if m.cfg.LogWithConsole {
			return io.MultiWriter(fd, m.cw)
		}
		return fd
	}
	var err error
	if m.cfg.LogDir != "" {
		filename = filepath.Join(m.cfg.LogDir, filename)
	}
	fd, err = initLoggerFile(filename)
	if err != nil {
		return m.cw
	}
	m.files.set(filename, fd)
	if m.cfg.LogWithConsole {
		return io.MultiWriter(fd, m.cw)
	}
	return fd
}
func (m *MultiWriter) Write(p []byte) (n int, err error) {
	m.checkDate()
	w := m.get("")
	return w.Write(p)
}
func (m *MultiWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	m.checkDate()
	w := m.get(level.String())
	return w.Write(p)
}

var _ zerolog.LevelWriter = (*MultiWriter)(nil)

// initLoggerFile 初始化日志文件
func initLoggerFile(logPath string) (logFile *os.File, err error) {
	var location string
	if !filepath.IsAbs(logPath) {
		location, err = os.Executable()
		if err != nil {
			log.Error().Err(err).Msg("os path error")
			return nil, err
		}
		logPath = filepath.Join(filepath.Dir(location), logPath)
	}
	//log.Debug().Str("log_path", logPath).Msg("Setup Logger")
	return os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, os.ModePerm) // nolint
}
