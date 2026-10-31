(() => {
  if (!window.hljs) return;
  document.querySelectorAll("pre.lab-code code[class^='language-']").forEach((code) => {
    const language = code.className.replace("language-", "");
    if (!window.hljs.getLanguage(language)) return;
    code.innerHTML = window.hljs.highlight(code.textContent, { language }).value;
  });
})();
