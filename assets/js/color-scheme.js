(function() {
  function setColorScheme(e) {
    if (e.matches) {
      document.documentElement.classList.add("dark", "sl-theme-dark");
    } else {
      document.documentElement.classList.remove("dark", "sl-theme-dark");
    }
  }
  const darkModePreference = window.matchMedia("(prefers-color-scheme: dark)");
  darkModePreference.addEventListener("change", setColorScheme);
  setColorScheme(darkModePreference);
})();
