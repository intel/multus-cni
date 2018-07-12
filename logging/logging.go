// Copyright (c) 2018 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Level type
type Level uint32

const (
	PanicLevel Level = iota
	ErrorLevel
	DebugLevel
	MaxLevel
	UnknownLevel
)

var loggingStderr bool
var loggingFp	*os.File
var loggingLevel Level
const defaultTimestampFormat = time.RFC3339

func (l Level) String() string {
	switch l {
	case PanicLevel:
		return "panic"
	case ErrorLevel:
		return "error"
	case DebugLevel:
		return "debug"
	}
	return "unknown"
}

func Printf(level Level, format string, a ...interface{}) {
	header := "%s [%s] "
	t := time.Now()
	if level > loggingLevel {
		return
	}

	if loggingStderr {
		fmt.Fprintf(os.Stderr, header, t.Format(defaultTimestampFormat), level)
		fmt.Fprintf(os.Stderr, format, a...)
		fmt.Fprintf(os.Stderr, "\n")
	}

	if loggingFp != nil {
		fmt.Fprintf(loggingFp, header, t.Format(defaultTimestampFormat), level)
		fmt.Fprintf(loggingFp, format, a...)
		fmt.Fprintf(loggingFp, "\n")
	}
}

func Debugf(format string, a ...interface{}) {
	Printf(DebugLevel, format, a...)
}

func Errorf(format string, a ...interface{}) {
	Printf(ErrorLevel, format, a...)
}

func Panicf(format string, a ...interface{}) {
	Printf(PanicLevel, format, a...)
	Printf(PanicLevel, "========= Stack trace output ========")
	Printf(PanicLevel, "%+v", errors.New("Multus Panic"))
	Printf(PanicLevel, "========= Stack trace output end ========")
}

func GetLoggingLevel (levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DebugLevel
	case "error":
		return ErrorLevel
	case "panic":
		return PanicLevel
	}
	fmt.Fprintf(os.Stderr, "multus logging: cannot set logging level to %s\n", levelStr)
	return UnknownLevel
}

func SetLogLevel (levelStr string) {
	level := GetLoggingLevel(levelStr)
	if level < MaxLevel {
		loggingLevel = level
	}
}

func SetLogStderr (enable bool) {
	loggingStderr = enable
}

func SetLogFile (filename string) {
	if filename == "" {
		return
	}

	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		loggingFp = nil
		fmt.Fprintf(os.Stderr, "multus logging: cannot open %s", filename)
	}
	loggingFp = fp
}

func init() {
	loggingStderr = true
	loggingFp = nil
	loggingLevel = PanicLevel
}
