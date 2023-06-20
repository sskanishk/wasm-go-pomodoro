const goWasm = new Go()

function updateRemainingTime(remaining) {
    const remainingTimeElement = document.getElementById("remainingTime");
    remainingTimeElement.textContent = remaining + " seconds";
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