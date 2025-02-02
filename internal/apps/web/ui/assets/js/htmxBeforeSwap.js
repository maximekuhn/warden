// By default, HTMX only swap elements on success (2xx status code)
// However, we want HTMX to also swap elements when a know error occurs.
// This script fixes this :)
document.body.addEventListener('htmx:beforeSwap', function(evt) {
    const status = evt.detail.xhr.status;
    const isUserError = status >= 400 && status < 500;
    if (isUserError) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    }
});
