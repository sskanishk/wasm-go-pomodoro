// package main

// import (
// 	"fmt"
// 	"syscall/js"
// )

// var htmlString = `<h4>Hello, I'm an HTML snippet from Go!</h4>`
// func GetHtml() js.Func {
// 	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		return htmlString
// 	})
// }

// func main() {

// 	ch := make(chan struct{}, 0)
// 	fmt.Printf("Hello Web Assembly from Go!\n")

// 	js.Global().Set("getHtml", GetHtml())
// 	<-ch
// }

package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var timerDuration time.Duration
var timerRunning bool

var htmlString = `<h4>Hello, I'm an HTML snippet from Go!</h4>`

func GetHtml() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return htmlString
	})
}

func AddMinutes() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// fmt.Println("timer 1", timerDuration)
		if timerDuration < time.Duration(60)*time.Second {
			timerDuration = timerDuration + time.Second
			// fmt.Println("timer 2", timerDuration)
			js.Global().Call("updateMinutes", timerDuration.Seconds())
		}
		return nil
	})
}

func SubMinutes() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if timerDuration > time.Duration(1)*time.Second {
			timerDuration = timerDuration - time.Second
			js.Global().Call("updateMinutes", timerDuration.Seconds())
		}
		return nil
	})
}

func StartTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		inputDuration := js.Global().Get("document").Call("getElementById", "duration").Get("innerText").String()
		durationInt, err := strconv.Atoi(inputDuration)
		if err != nil {
			fmt.Println("Invalid duration:", inputDuration)
			return nil
		}

		if timerRunning {
			// Timer is already running, do nothing
			return nil
		}

		// Set the timer duration to 25 minutes
		timerDuration = time.Duration(durationInt) * time.Second

		// timerDuration = 25 * time.Minute
		timerRunning = true

		// Start the timer in a separate goroutine
		go RunTimer()

		return nil

	})
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

	for timerRunning && timerDuration >= 0 {
		// Send the remaining time to the UI
		SendRemainingTime()

		fmt.Println("-----------------------------------------------")
		fmt.Println("timerDuration 1:", timerDuration)
		fmt.Println("-----------------------------------------------")

		// Decrement the timer duration by 1 second
		timerDuration -= 1 * time.Second

		fmt.Println("-----------------------------------------------")
		fmt.Println("timerDuration 2:", timerDuration)
		fmt.Println("-----------------------------------------------")

		// js.Global().Call("updateMinutes", timerDuration.Seconds())
		time.Sleep(1 * time.Second)
	}

	if !timerRunning {
		// Timer paused or stopped
		fmt.Println("Its herer")
		SendRemainingTime()
	}

}

func PauseTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !timerRunning {
			// Timer is not running, do nothing
			return ""
		}

		timerRunning = false
		SendRemainingTime()
		return ""
	})
}

func ResumeTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if timerRunning {
			// Timer is already running, do nothing
			return ""
		}

		timerRunning = true
		go RunTimer()
		return ""
	})
}

func ResetTimer() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		timerRunning = false
		timerDuration = 0
		SendRemainingTime()
		return ""
	})
}

func SendRemainingTime() {
	js.Global().Call("updateRemainingTime", timerDuration.Seconds())
}

// func sendJSONResponse(data interface{}) string {
// 	return toJSON(data)
// }

// func toJSON(data interface{}) string {
// 	json, _ := json.Marshal(data)
// 	println("JSON in toJson ", json, string(json))
// 	return string(json)
// }

func RegisterCallbacks() {
	js.Global().Set("startTimer", StartTimer())
	js.Global().Set("pauseTimer", PauseTimer())
	js.Global().Set("resumeTimer", ResumeTimer())
	js.Global().Set("resetTimer", ResetTimer())
	js.Global().Set("addMinutes", AddMinutes())
	js.Global().Set("subMinutes", SubMinutes())
}

func main() {
	c := make(chan struct{}, 0)
	RegisterCallbacks()
	js.Global().Set("getHtml", GetHtml())
	<-c
}
