package logger

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/transerver/commons/utils"
	"strings"
	"sync"
)

type DebugFormatter struct {
	Config         Config
	colors         map[Level]*color.Color
	colorOnce      sync.Once
	underLineColor *color.Color
}

func (f *DebugFormatter) init() {
	f.colorOnce.Do(func() {
		f.underLineColor = color.New(color.Underline)
		f.colors = make(map[Level]*color.Color)
		if len(f.Config.TimestampFormat) == 0 {
			f.Config.TimestampFormat = "2006-01-02 15:04:05.000000"
		}

		for _, l := range AllLevels {
			f.initLevelColor(l)
		}
	})
}

func (f *DebugFormatter) Format(entry *Entry) ([]byte, error) {
	f.init()
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	cf := f.colors[entry.Level]
	_, _ = cf.Fprintf(b, "[%s]", f.upperLevel(entry.Level))
	_, _ = cf.Fprintf(b, " [%s]", entry.Time.Format(f.Config.TimestampFormat))

	var prefixLen int
	for k, v := range entry.Data {
		prefix := v == nil
		if !prefix {
			if s, ok := v.(string); ok && len(s) == 0 {
				prefix = true
			}
		}

		if prefix && utils.NotBlank(k) {
			_, _ = cf.Fprintf(b, " %s:", f.underLineColor.Sprintf("[%s]", k))
			prefixLen = len(k) + 4
			delete(entry.Data, k)
			break
		}
	}

	if len(entry.Message)+prefixLen > 50 || len(entry.Data) == 0 {
		_, _ = cf.Fprintf(b, " %s", entry.Message)
	} else {
		_, _ = cf.Fprintf(b, " %-*s", 50-prefixLen, entry.Message)
	}

	for k, v := range entry.Data {
		_, _ = cf.Fprintf(b, " %s", k)
		_, _ = fmt.Fprintf(b, "=")
		f.appendValue(b, v)
	}

	return append(b.Bytes(), '\n'), nil
}

func (f *DebugFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *DebugFormatter) needsQuoting(text string) bool {
	if f.Config.ForceQuote {
		return true
	}
	if f.Config.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.Config.DisableQuote {
		return false
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *DebugFormatter) upperLevel(level Level) string {
	switch level {
	case DebugLevel:
		return "DBUG"
	default:
		return strings.ToUpper(level.String()[:4])
	}
}

func (f *DebugFormatter) initLevelColor(level Level) {
	var lc color.Attribute
	switch level {
	case DebugLevel, TraceLevel:
		lc = color.FgWhite
	case WarnLevel:
		lc = color.FgYellow
	case ErrorLevel, FatalLevel, PanicLevel:
		lc = color.FgRed
	default:
		lc = color.FgCyan
	}
	c := color.New(lc, color.Bold)
	f.colors[level] = c
}
