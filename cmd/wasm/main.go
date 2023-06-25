package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var sessionTimeDuration time.Duration
var timerRunning bool = false
var breakTimeDuration time.Duration

var sessionTimeDurationValue string
var breakTimeDurationValue string

var runningTimeDuration time.Duration

var isSessionOn bool = true
var isBreakOn bool = false

var htmlString = `<p>Hello, I'm an HTML snippet from Go!</p>`

func GetHtml() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return htmlString
	})
}

func ModifyTime() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		if len(args) < 2 {
			return nil
		}

		operation := args[0].String()
		operator := args[1].String()

		if operator == "+" {
			if operation == "session" {
				if sessionTimeDuration < time.Duration(60)*time.Second {
					sessionTimeDuration = sessionTimeDuration + time.Second
				}
			} else if operation == "break" {
				if breakTimeDuration < time.Duration(60)*time.Second {
					breakTimeDuration = breakTimeDuration + time.Second
				}
			}
		} else if operator == "-" {
			if operation == "session" {
				if sessionTimeDuration > time.Duration(1)*time.Second {
					sessionTimeDuration = sessionTimeDuration - time.Second
				}
			} else if operation == "break" {
				if breakTimeDuration > time.Duration(1)*time.Second {
					breakTimeDuration = breakTimeDuration - time.Second
				}
			}
		}
		SendMinutesTime(operation)
		return nil
	})
}

func StartTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// fmt.Println("starttimer ", sessionTimeDuration, timerRunning)
		var inputDuration string
		sessionTimeDurationValue = js.Global().Get("document").Call("getElementById", "sessionDuration").Get("innerText").String()
		breakTimeDurationValue = js.Global().Get("document").Call("getElementById", "breakDuration").Get("innerText").String()
		timerRunning = true

		if isSessionOn {
			inputDuration = sessionTimeDurationValue
		} else if isBreakOn {
			inputDuration = breakTimeDurationValue
		}
		StartTimerGo(inputDuration)
		return nil
	})
}

func StartTimerGo(inputDuration string) {

	fmt.Println("Timmer Running ", timerRunning)

	// in future will use as pause
	if !timerRunning {
		return
	}

	durationInt, err := strconv.Atoi(inputDuration)
	if err != nil {
		fmt.Println("Invalid duration:", inputDuration)
		return
	}

	runningTimeDuration = time.Duration(durationInt) * time.Second

	// Start the timer in a separate goroutine
	go RunTimer()

	return
}

func RunTimer() {
	fmt.Println("Inside issession of runtimer ", sessionTimeDuration, isBreakOn, isSessionOn)

	for timerRunning && runningTimeDuration > 0 {
		SendRemainingTime()
		runningTimeDuration -= 1 * time.Second
		time.Sleep(1 * time.Second)
	}

	if runningTimeDuration == 0 {
		isBreakOn = !isBreakOn
		isSessionOn = !isSessionOn
		SendRemainingTime()
		if isSessionOn {
			UpdateColor("session")
			StartTimerGo(sessionTimeDurationValue)
		} else if isBreakOn {
			UpdateColor("break")
			StartTimerGo(breakTimeDurationValue)
		}
	}
}

func PauseTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !timerRunning {
			return nil
		}
		timerRunning = false
		SendRemainingTime()
		return nil
	})
}

func ResumeTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if timerRunning {
			return nil
		}
		timerRunning = true
		go RunTimer()
		return nil
	})
}

func ResetTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		timerRunning = false
		runningTimeDuration = 0
		SendRemainingTime()
		return nil
	})
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func SendRemainingTime() {
	js.Global().Call("updateRunningTime", formatDuration(runningTimeDuration))
}

func UpdateColor(className string) {
	js.Global().Call("updateBgColor", className)
}

func SendMinutesTime(operation string) {
	if operation == "session" {
		js.Global().Call("updateDurationTime", sessionTimeDuration.Seconds(), "session")
	} else if operation == "break" {
		js.Global().Call("updateDurationTime", breakTimeDuration.Seconds(), "break")
	}
}

func SendResetMinutes() {
	js.Global().Call("resetMinutes", sessionTimeDuration.Seconds())
}

func RegisterCallbacks() {
	js.Global().Set("startTimer", StartTimer())
	js.Global().Set("pauseTimer", PauseTimer())
	js.Global().Set("resumeTimer", ResumeTimer())
	js.Global().Set("resetTimer", ResetTimer())
	js.Global().Set("getHtml", GetHtml())
	js.Global().Set("modifyTime", ModifyTime())
}

func main() {
	c := make(chan struct{}, 0)
	RegisterCallbacks()
	<-c
}
