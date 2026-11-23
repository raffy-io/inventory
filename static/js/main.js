document.addEventListener("htmx:beforeSwap", (event) => {
  const status = event.detail.xhr.status;
  if (status === 400 || status === 422) {
    event.detail.shouldSwap = true;
    event.detail.isError = false;
  }
});
