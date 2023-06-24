package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var timerDuration time.Duration
var timerRunning bool
var breakTimeDuration time.Duration

var htmlString = `<h4>Hello, I'm an HTML snippet from Go!</h4>`

func GetHtml() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return htmlString
	})
}

func AddMinutes() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("timer 1", timerDuration)
		if timerDuration < time.Duration(60)*time.Second {
			timerDuration = timerDuration + time.Second
			SendMinutesTime()
		}
		return nil
	})
}

func SubMinutes() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("good Sub", breakTimeDuration)
		if timerDuration > time.Duration(1)*time.Second {
			timerDuration = timerDuration - time.Second
			SendMinutesTime()
		}
		return nil
	})
}

func BreakSubMinutes() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("good ", breakTimeDuration)
		if breakTimeDuration > time.Duration(1)*time.Second {
			breakTimeDuration = breakTimeDuration - time.Second
			js.Global().Call("updateBreakMinutesTime", breakTimeDuration.Seconds())
		}
		return nil
	})
}

func BreakAddMinutes() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("B Add", breakTimeDuration)
		if breakTimeDuration < time.Duration(60)*time.Second {
			breakTimeDuration = breakTimeDuration + time.Second
			js.Global().Call("updateBreakMinutesTime", breakTimeDuration.Seconds())
			// SendMinutesTime()
		}
		return nil
	})
}

func StartTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("starttimer ", timerDuration, timerRunning)

		inputDuration := js.Global().Get("document").Call("getElementById", "duration").Get("innerText").String()
		durationInt, err := strconv.Atoi(inputDuration)
		if err != nil {
			fmt.Println("Invalid duration:", inputDuration)
			return nil
		}

		if timerRunning {
			return nil
		}

		timerDuration = time.Duration(durationInt) * time.Minute
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
		timerDuration = 0
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
	// 	fmt.Println("timerDuration 1", timerDuration)
	// 	fmt.Println("time.Second", time.Second)
	// 	timerDuration -= 1 * time.Second
	// 	fmt.Println("timerDuration 2", timerDuration)
	// 	fmt.Println("-----------------------------------------------")

	// 	if timerDuration <= 0 {
	// 		// Timer completed, stop the timer
	// 		timerRunning = false
	// 		SendRemainingTime()
	// 		return
	// 	}
	// }

	for timerRunning && timerDuration > 0 {
		SendRemainingTime()

		// fmt.Println("-----------------------------------------------")
		// fmt.Println("timerDuration 1:", timerDuration)
		// fmt.Println("-----------------------------------------------")

		timerDuration -= 1 * time.Second

		// fmt.Println("-----------------------------------------------")
		// fmt.Println("timerDuration 2:", timerDuration)
		// fmt.Println("-----------------------------------------------")

		time.Sleep(1 * time.Second)
	}

	if timerDuration == 0 {
		timerRunning = false
		SendRemainingTime()
	}
}

func SendRemainingTime() {
	js.Global().Call("updateRemainingTime", formatDuration(timerDuration))
	js.Global().Call("updateBreakRemainingTime", formatDuration(breakTimeDuration))
}

func SendMinutesTime() {
	js.Global().Call("updateMinutesTime", timerDuration.Seconds())
}

func SendResetMinutes() {
	js.Global().Call("resetMinutes", timerDuration.Seconds())
}

func RegisterCallbacks() {
	js.Global().Set("startTimer", StartTimer())
	js.Global().Set("pauseTimer", PauseTimer())
	js.Global().Set("resumeTimer", ResumeTimer())
	js.Global().Set("resetTimer", ResetTimer())
	js.Global().Set("addMinutes", AddMinutes())
	js.Global().Set("subMinutes", SubMinutes())
	js.Global().Set("breakAddMinutes", BreakAddMinutes())
	js.Global().Set("breakSubMinutes", BreakSubMinutes())
	js.Global().Set("getHtml", GetHtml())
}

func main() {
	c := make(chan struct{}, 0)
	RegisterCallbacks()
	<-c
}
