const goWasm = new Go()

function updateRemainingTime(remaining) {
    console.log("remainhg ", remaining)
    const remainingTimeElement = document.getElementById("remainingTime");
    remainingTimeElement.textContent = remaining;
}

function updateBreakRemainingTime(remaining) {
    console.log("remainhg ", remaining)
    const remainingTimeElement = document.getElementById("breakRemainingTime");
    remainingTimeElement.textContent = remaining;
}

function updateBreakMinutesTime(minutes) {
    console.log("Script.js minutes ", minutes)
    if (minutes >= 1 && minutes <= 60) {
        const minutesElement = document.getElementById("breakDuration");
        minutesElement.innerText = minutes
    }
}

function updateMinutesTime(minutes) {
    console.log("Script.js minutes ", minutes)
    if (minutes >= 1 && minutes <= 60) {
        const minutesElement = document.getElementById("duration");
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