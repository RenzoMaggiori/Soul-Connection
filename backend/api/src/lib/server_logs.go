package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

var progress = mpb.New(mpb.WithWidth(64))
var progressBars = make(map[string]*mpb.Bar)
var LoggerEnabled = true

func DisableLogger() {
	LoggerEnabled = false
}

func EnableLogger() {
	LoggerEnabled = true
}

func ServerLog(msgType string, message interface{}) {
	if !LoggerEnabled {
		return
	}
	color := map[string]string{
		"INFO":     Green,
		"WARNING":  Yellow,
		"ERROR":    Red,
		"DEBUG":    Magenta,
		"PROGRESS": Cyan,
	}

	colorCode, ok := color[msgType]
	if !ok {
		colorCode = White
	}

	var formattedMessage string
	switch msg := message.(type) {
	case string:
		formattedMessage = msg
	case error:
		formattedMessage = msg.Error()
	default:
		formattedMessage = fmt.Sprintf("%v", msg)
	}

	if msgType == "PROGRESS" {
		updateProgressBar(formattedMessage)
		return
	}

	fmt.Printf("%s Server [%s%s%s] %s%s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		colorCode,
		msgType,
		Reset,
		formattedMessage,
		Reset,
	)
}

func updateProgressBar(message string) {
	parts := strings.Split(message, ":")
	if len(parts) < 2 {
		return
	}

	taskName := parts[0]
	action := parts[1]

	switch action {
	case "START":
		if len(parts) != 3 {
			return
		}
		total := parts[2]
		totalSteps, err := strconv.Atoi(total)
		if err != nil {
			return
		}
		fmt.Printf("%s Server [%sPROGRESS%s] %s%s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			Cyan,
			Reset,
			taskName,
			Reset,
		)

		bar := progress.AddBar(int64(totalSteps),
			mpb.BarRemoveOnComplete(),
			mpb.PrependDecorators(
				decor.CountersNoUnit(fmt.Sprintf("%%d / %%d")),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WC{W: 5, C: decor.DSyncSpace}),
			),
		)
		progressBars[taskName] = bar

	case "INCREMENT":
		bar, exists := progressBars[taskName]
		if exists {
			bar.Increment()
		}

	case "COMPLETE":
		bar, exists := progressBars[taskName]
		if exists {
			bar.SetTotal(0, true)
			delete(progressBars, taskName)
		}
	}
}
