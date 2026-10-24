import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Main from './Main';
import * as serviceWorker from './serviceWorker';

if (typeof WebAssembly === "object" && typeof WebAssembly.instantiate === "function") {
  if (!WebAssembly.instantiateStreaming) { // polyfill

      WebAssembly.instantiateStreaming = async (resp, importObject) => {
          const source = await (await resp).arrayBuffer();
          return await WebAssembly.instantiate(source, importObject);
      };
  }

  // WebAssembly.instantiateStreaming(fetch("wasm.wasm"), window.go.importObject).then((result) => {
  //     window.go.run(result.instance);
  // });

  fetch('wasm.wasm').then(response => 
  	response.arrayBuffer()
  ).then(bytes => 
    WebAssembly.instantiate(bytes, window.go.importObject)
  ).then(result => {
    window.go.run(result.instance);
  });

} else {
  alert('!WebAssembly')
}

ReactDOM.render(
    <React.StrictMode>
      <Main />
    </React.StrictMode>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.register();
