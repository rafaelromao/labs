(() => {
  const title = document.querySelector("#source-title");
  const status = document.querySelector("#source-status");
  const code = document.querySelector("#source-code");
  const github = document.querySelector("#source-github");
  const requested = new URLSearchParams(window.location.search).get("file") || "";

  const fail = (message) => {
    title.textContent = "Source unavailable";
    status.textContent = message;
  };

  const isJava = /^(?:solutions\/)?java\/src\/(?:main|test)\/java\/[a-z0-9_./-]+\.java$/i.test(requested);
  const isGo = /^(?:solutions\/)?go\/[a-z0-9_./-]+\.go$/i.test(requested);
  if (!isJava && !isGo || requested.includes("..")) {
    fail("Only Java and Go files from this repository can be opened here.");
    return;
  }

  const filename = requested.split("/").pop();
  const language = isJava ? "java" : "go";
  title.textContent = filename;
  status.textContent = `Loading ${requested}…`;
  github.hidden = false;
  github.querySelector("a").href = `https://github.com/rafaelromao/labs/blob/main/${requested}`;

  fetch(requested, { credentials: "same-origin" })
    .then((response) => {
      if (!response.ok) throw new Error(`The server returned ${response.status}.`);
      return response.text();
    })
    .then((text) => {
      code.className = `language-${language}`;
      code.innerHTML = window.hljs.highlight(text, { language }).value;
      status.textContent = `${requested} · ${text.split("\n").length - 1} lines`;
      document.title = `${filename} · source`;
    })
    .catch((error) => fail(`Could not load ${filename}: ${error.message}`));
})();
