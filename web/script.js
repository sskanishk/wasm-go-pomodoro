const goWasm = new Go()

function updateRunningTime(remaining) {
    // console.log("remainhg ", remaining)
    const runningTimeElement = document.getElementById("runningTime");
    runningTimeElement.textContent = remaining;
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

function updateBreakMinutesTime(minutes) {
    if (minutes >= 1 && minutes <= 60) {
        const minutesElement = document.getElementById("breakDuration");
        minutesElement.innerText = minutes
    }
}

function updateMinutesTime(minutes) {
    console.log("Script.js minutes ", minutes)
    if (minutes >= 1 && minutes <= 60) {
        const minutesElement = document.getElementById("sessionDuration");
        minutesElement.innerText = minutes
    }
}

function resetMinutes(minutes) {
    const minutesElement = document.getElementById("duration");
    console.log("minutes ", minutes)
    minutesElement.innerText = minutes
}

WebAssembly.instantiateStreaming(fetch("main.wasm"), goWasm.importObject).then((result) => {
    goWasm.run(result.instance)
    document.getElementById("get-html").addEventListener("click", () => {
        document.body.innerHTML += getHtml()
    })
})