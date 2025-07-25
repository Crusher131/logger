package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type cfgopt func(*cfg)

type cfg struct {
	id            int
	file          *os.File
	writeFile     bool
	logfile       string
	writeTerminal bool
	debug         bool
	logpath       Logpath
	multiWriter   io.Writer
}
type Logpath struct {
	path string
	file string
}

func Init(f ...cfgopt) {
	c := defaultcfg()
	for _, fn := range f {
		fn(&c)
	}
	definePath(&c)

	if c.writeTerminal && c.writeFile {
		err := createLogFile(&c)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		c.multiWriter = io.MultiWriter(c.file, os.Stdout)
	} else if !c.writeTerminal && c.writeFile {
		err := createLogFile(&c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		c.multiWriter = io.MultiWriter(c.file)
	} else if c.writeTerminal && !c.writeFile {
		c.multiWriter = io.MultiWriter(os.Stdout)
	}

	_cfg = &c
}

var _cfg *cfg

func defaultcfg() cfg {
	cfg := cfg{
		id:            os.Getpid(),
		writeTerminal: true,
		writeFile:     false,
		debug:         false,
		logpath: Logpath{
			path: "log",
			file: "log.log",
		},
	}
	return cfg
}

func SetId(id int) func(c *cfg) {
	return func(c *cfg) {
		c.id = id
	}
}

func SetTerminal() func(c *cfg) {
	return func(c *cfg) {
		c.writeTerminal = true
	}
}

func SetDebug() func(c *cfg) {
	return func(c *cfg) {
		c.debug = true
	}
}
func SetLogFile(logfile string) func(c *cfg) {
	return func(c *cfg) {
		c.logfile = logfile
		c.writeFile = true
	}
}

func GetId() (id int) {
	return _cfg.id
}

func GetTerminal() (terminal bool) {
	return _cfg.writeTerminal
}

func GetLogFile() (logfile string) {
	return _cfg.logfile
}

func Close() error {
	return _cfg.close()
}

func (c *cfg) close() error {
	if c.file == nil {
		return nil
	}

	return c.file.Close()
}
func definePath(c *cfg) {
	_files := strings.Split(c.logfile, "/")
	_logfile := _files[len(_files)-1]
	if len(_files) > 0 {
		_files = _files[:len(_files)-1]
	}
	_filepath := strings.Join(_files, "/")
	c.logpath.file = _logfile
	c.logpath.path = _filepath
}

func createLogFile(c *cfg) error {
	if c.logpath.path != "" {

		if _, err := os.Stat(c.logpath.path); os.IsNotExist(err) {
			err := os.MkdirAll(c.logpath.path, 0755)
			if err != nil {
				return fmt.Errorf("falha ao criar a estrutura de diretorio: %w", err)
			}
		}
	}

	logFile, err := os.OpenFile(filepath.Join(c.logpath.path, filepath.Base(c.logpath.file)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("falha ao criar o arquivo de log: %w", err)
	}
	c.file = logFile
	return err
}

func Info(text ...any) {
	prefix := "INFO: "
	collor := ColorCyan
	fmtText := fmt.Sprint(text...)
	_cfg.writer(collor, prefix, fmtText)
}

func Warn(text ...any) {
	prefix := "WARN: "
	collor := ColorYellow
	fmtText := fmt.Sprint(text...)
	_cfg.writer(collor, prefix, fmtText)
}

func Debug(text ...any) {
	if _cfg.debug {

		prefix := "Debug: "
		collor := ColorYellow
		fmtText := fmt.Sprint(text...)
		_cfg.writer(collor, prefix, fmtText)

	}
}

func Error(err ...any) {
	prefix := "ERROR: "
	collor := ColorRed
	text := fmt.Sprint(err...)
	_cfg.writer(collor, prefix, text)
}

func Fatal(err ...any) {
	prefix := "FATAL: "

	collor := ColorRed
	text := fmt.Sprint(err...)
	_cfg.writer(collor, prefix, text)
	log.Fatal()
}

func (c *cfg) getDate() (lnow string) {
	lnow = fmt.Sprintf("%d/%d/%d %d:%d:%d", time.Now().Day(),
		time.Now().Month(),
		time.Now().Year(),
		time.Now().Hour(),
		time.Now().Hour(),
		time.Now().Second())
	return lnow
}
func (c *cfg) writer(collor string, prefix string, text string) {

	if c.writeTerminal {
		c.terminalWriter(collor, prefix, text, c.getDate())
	}
	textfile := fmt.Sprint(c.getDate(), " ", c.id, " ", prefix, text, "\n")
	c.file.Write([]byte(textfile))
}

func (c *cfg) terminalWriter(collor string, prefix string, text string, lnow string) {
	textterminal := fmt.Sprintln(lnow, c.id, collor, prefix, ColorReset, text)
	fmt.Print(textterminal)
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)
