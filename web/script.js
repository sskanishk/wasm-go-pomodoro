const goWasm = new Go()

function updateRunningTime(remaining, remainingTimeInMilliSec) {
    const runningTimeElement = document.getElementById("runningTime");
    runningTimeElement.textContent = remaining;    
}

function updateBgColor(className) {
    const html = document.getElementById('pomodoro')
    if(html.classList.length > 0) {
        html.classList.remove(html.classList?.item(0))
    }
    html.classList.add(className)
    const currentOperation = document.getElementById("runningOperation")
    currentOperation.innerText = className
}

function updateDurationTime(minutes, operation) {
    console.log("Script.js minutes ", minutes, operation)
    if (minutes >= 1 && minutes <= 60) {
        if(operation == "break") {
            const minutesElement = document.getElementById("breakDuration");
            minutesElement.innerText = minutes
        } else if(operation == "session") {
            const minutesElement = document.getElementById("sessionDuration");
            minutesElement.innerText = minutes        
        }
    }
}


function resetMinutes(minutes) {
    const minutesElement = document.getElementById("duration");
    console.log("minutes ", minutes)
    minutesElement.innerText = minutes
}

WebAssembly.instantiateStreaming(fetch("main.wasm"), goWasm.importObject).then((result) => {
    goWasm.run(result.instance)
    document.getElementById("get_html").addEventListener("click", () => {
        const element = document.getElementById("test_container")
        element.innerHTML += getHtml()
    })
})