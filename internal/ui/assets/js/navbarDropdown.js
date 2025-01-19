const dropdownButton = document.getElementById("navbar-dropdown-button");
const dropdownMenu = document.getElementById("navbar-dropdown-menu");

dropdownButton.addEventListener("click", function() {
    dropdownMenu.classList.toggle("hidden");
});

window.addEventListener("click", function(e) {
    if (!dropdownButton.contains(e.target) && !dropdownMenu.contains(e.target)) {
        dropdownMenu.classList.add("hidden");
    }
});
