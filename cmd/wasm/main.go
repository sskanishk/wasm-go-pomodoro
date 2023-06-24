package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var sessionTimeDuration time.Duration
var timerRunning bool
var breakTimeDuration time.Duration

var htmlString = `<h4>Hello, I'm an HTML snippet from Go!</h4>`

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
		fmt.Println("starttimer ", sessionTimeDuration, timerRunning)

		inputDuration := js.Global().Get("document").Call("getElementById", "duration").Get("innerText").String()
		durationInt, err := strconv.Atoi(inputDuration)
		if err != nil {
			fmt.Println("Invalid duration:", inputDuration)
			return nil
		}

		if timerRunning {
			return nil
		}

		sessionTimeDuration = time.Duration(durationInt) * time.Minute
		timerRunning = true

		// Start the timer in a separate goroutine
		go RunTimer()

		return nil

	})
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
		sessionTimeDuration = 0
		SendRemainingTime()
		return nil
	})
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func RunTimer() {
	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	// for range ticker.C {
	// 	if !timerRunning {
	// 		// Timer paused or stopped, exit the goroutine
	// 		return
	// 	}

	// 	// Send the remaining time to the UI
	// 	SendRemainingTime()

	// 	fmt.Println("-----------------------------------------------")
	// 	fmt.Println("sessionTimeDuration 1", sessionTimeDuration)
	// 	fmt.Println("time.Second", time.Second)
	// 	sessionTimeDuration -= 1 * time.Second
	// 	fmt.Println("sessionTimeDuration 2", sessionTimeDuration)
	// 	fmt.Println("-----------------------------------------------")

	// 	if sessionTimeDuration <= 0 {
	// 		// Timer completed, stop the timer
	// 		timerRunning = false
	// 		SendRemainingTime()
	// 		return
	// 	}
	// }

	for timerRunning && sessionTimeDuration > 0 {
		SendRemainingTime()

		// fmt.Println("-----------------------------------------------")
		// fmt.Println("sessionTimeDuration 1:", sessionTimeDuration)
		// fmt.Println("-----------------------------------------------")

		sessionTimeDuration -= 1 * time.Second

		// fmt.Println("-----------------------------------------------")
		// fmt.Println("sessionTimeDuration 2:", sessionTimeDuration)
		// fmt.Println("-----------------------------------------------")

		time.Sleep(1 * time.Second)
	}

	if sessionTimeDuration == 0 {
		timerRunning = false
		SendRemainingTime()
	}
}

func SendRemainingTime() {
	js.Global().Call("updateRemainingTime", formatDuration(sessionTimeDuration))
	js.Global().Call("updateBreakRemainingTime", formatDuration(breakTimeDuration))
}

func SendMinutesTime(operation string) {
	if operation == "session" {
		js.Global().Call("updateMinutesTime", sessionTimeDuration.Seconds())
	} else if operation == "break" {
		js.Global().Call("updateBreakMinutesTime", breakTimeDuration.Seconds())
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
