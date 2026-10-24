const arr = {
  default: {
    md5: "&zSMvVTK4gbDvCcU",
    aes: "M53B?xKnTVtVnjM^",
    iv: "3579653693526433",
  },
  tpower: {
    md5: "&W7e1vRWRH02jZIB",
    aes: "M53B?gQhkRtSeQZ^",
    iv: "4339794120398193",
  },
};

let env = "default";
let current = arr[env];

function updateCustomConfig() {
  current = {
    md5: document.getElementById("custom_md5").value,
    aes: document.getElementById("custom_aes").value,
    iv: document.getElementById("custom_iv").value,
  };
}

function encrypt(data) {
  if (env === "custom") {
    updateCustomConfig();
  }

  const key_hash = CryptoJS.MD5(current.aes).toString();
  const key = CryptoJS.enc.Utf8.parse(key_hash);
  const iv = CryptoJS.enc.Utf8.parse(current.iv);
  const encrypted = CryptoJS.AES.encrypt(data, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });
  return encrypted.toString();
}

function decrypt(encrypted) {
  if (env === "custom") {
    updateCustomConfig();
  }

  try {
    const key_hash = CryptoJS.MD5(current.aes).toString();
    const key = CryptoJS.enc.Utf8.parse(key_hash);
    const iv = CryptoJS.enc.Utf8.parse(current.iv);
    const decrypted = CryptoJS.AES.decrypt(encrypted, key, {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    });
    return decrypted.toString(CryptoJS.enc.Utf8);
  } catch (error) {
    console.error("解密失敗：", error);
    return "解密失敗";
  }
}

let disableAutoCopy = false;

require.config({
  paths: { vs: "https://example.invalid/monaco-editor@0.34.1/min/vs" },
});
require(["vs/editor/editor.main"], function () {
  const editor = monaco.editor.create(document.getElementById("editor"), {
    value: `{
  "domain": "example.invalid"
}`,
    language: "json",
    theme: "vs-dark",
    automaticLayout: true,
  });

  const resultOutput = document.getElementById("resultOutput");
  const copyButton = document.getElementById("copyButton");
  const autoCopyToggle = document.getElementById("autoCopyToggle");
  const statusDisplay = document.getElementById("status");

  // 複製按鈕功能
  copyButton.addEventListener("click", () => {
    navigator.clipboard
      .writeText(resultOutput.value)
      .then(() => {
        statusDisplay.innerText = "已複製到剪貼簿";
      })
      .catch((err) => {
        console.error("複製失敗：", err);
        statusDisplay.innerText = "複製失敗";
      });
  });

  document
    .getElementById("configSelector")
    .addEventListener("change", (event) => {
      env = event.target.value;
      current = arr[env];
      document.getElementById("customConfig").style.display =
        env === "custom" ? "block" : "none";
      statusDisplay.innerText = `已切換到配置：${env}`;
      makeAndRefreshEncrypt();
    });

  editor.onDidChangeModelContent(() => {
    makeAndRefreshEncrypt();
  });

  function makeAndRefreshEncrypt() {
    if (disableAutoCopy) {
      disableAutoCopy = false;
      return;
    }

    statusDisplay.innerText = "";
    const data = editor.getValue();
    const params = encrypt(data);
    const key = CryptoJS.MD5(params + current.md5).toString();
    const payload = JSON.stringify({ params, key });

    resultOutput.value = payload;

    if (autoCopyToggle.checked) {
      navigator.clipboard
        .writeText(payload)
        .then(() => {
          statusDisplay.innerText = "已複製到剪貼簿";
        })
        .catch((err) => {
          console.error("自動複製失敗：", err);
          statusDisplay.innerText = "自動複製失敗";
        });
    }
  }

  resultOutput.addEventListener("input", () => {
    statusDisplay.innerText = "";
    try {
      const input = resultOutput.value;
      const parsed = JSON.parse(input);
      const decrypted = decrypt(parsed.params);
      disableAutoCopy = true;
      editor.setValue(decrypted);
      disableAutoCopy = true;
      editor.getAction("editor.action.formatDocument").run();
      statusDisplay.innerText = "解密成功";
    } catch (error) {
      console.error("解析或解密失敗：", error);
      statusDisplay.innerText = "解密失敗";
    }
  });
});
