const goWasm = new Go()

WebAssembly.instantiateStreaming(fetch("main.wasm"), goWasm.importObject)
    .then((result) => {
        goWasm.run(result.instance)
        console.log("Result ", result)
        document.getElementById("get-html").addEventListener("click", () => {
            document.body.innerHTML += getHtml()
        })
    })