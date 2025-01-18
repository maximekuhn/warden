// By default, HTMX only swap elements on success (2xx status code)
// However, we want HTMX to also swap elements when a know error occurs.
// This script fixes this :)
document.body.addEventListener('htmx:beforeSwap', function(evt) {
    if (evt.detail.xhr.status === 409) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    }
});
