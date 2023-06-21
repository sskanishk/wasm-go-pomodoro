const goWasm = new Go()

function updateRemainingTime(remaining) {
    console.log("remainhg ", remaining)
    const remainingTimeElement = document.getElementById("remainingTime");
    remainingTimeElement.textContent = remaining;
}

function updateMinutes(minutes) {
    if(minutes >= 1 && minutes <= 60 ) {
        const minutesElement = document.getElementById("duration");
        minutesElement.innerText = minutes
    }
}

WebAssembly.instantiateStreaming(fetch("main.wasm"), goWasm.importObject).then((result) => {
    goWasm.run(result.instance)

    document.getElementById("get-html").addEventListener("click", () => {
        document.body.innerHTML += getHtml()
    })
})


    // const go = new Go();
    // WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    //     go.run(result.instance);
    // });