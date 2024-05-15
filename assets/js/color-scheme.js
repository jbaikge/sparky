// Adapted from the shoelace.style docs.js
(function() {
  function getThemePreference() {
    return localStorage.getItem("theme") || "auto";
  }

  function setThemePreference(newTheme) {
    localStorage.setItem("theme", newTheme);
    updateThemeMenu(newTheme);
    document.documentElement.classList.toggle("sl-theme-dark", prefersDark(newTheme));
  }

  function prefersDark(theme) {
    if (theme == "auto") {
      return window.matchMedia("(prefers-color-scheme: dark)").matches;
    }
    return theme == "dark";
  }

  function updateThemeMenu(theme) {
    const menu = document.querySelector("#theme-selector sl-menu");
    if (menu == null) return;
    Array.from(menu.querySelectorAll("sl-menu-item")).map(function(item) {
      item.checked = (item.getAttribute("value") == theme);
    });
  }

  // Listen for opening the dropdown
  document.addEventListener("sl-show", function(event) {
    const themeSelector = event.target.closest("#theme-selector");
    if (themeSelector == null) return;
    updateThemeMenu(getThemePreference());
  });

  // Listen for selections
  document.addEventListener("sl-select", function(event) {
    const menu = event.target.closest("#theme-selector sl-menu");
    if (menu == null) return;
    setThemePreference(event.detail.item.value);
  });

  // Listen for system theme changes
  window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", function() {
    setThemePreference(getThemePreference());
  });

  setThemePreference(getThemePreference());
})();
