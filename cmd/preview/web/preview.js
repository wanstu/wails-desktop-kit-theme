(() => {
  "use strict";

  const themePicker = document.getElementById("theme-picker");
  const packDescription = document.getElementById("pack-description");
  const packSummary = document.getElementById("pack-summary");
  const packPosition = document.getElementById("pack-position");
  const prevPackButton = document.getElementById("prev-pack");
  const nextPackButton = document.getElementById("next-pack");
  const modeTabs = Array.from(document.querySelectorAll("[data-mode]"));
  const loadedStyles = new Set();

  let packs = [];
  let activePackIndex = 0;

  function ensurePackStylesheet(pack) {
    if (loadedStyles.has(pack.name)) return;

    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = pack.stylesheetPath;
    link.dataset.previewPack = pack.name;

    const base = document.getElementById("kit-base");
    document.head.insertBefore(link, base);
    loadedStyles.add(pack.name);
  }

  function selectMode(mode) {
    desktopKitTheme.apply(mode);
    localStorage.setItem("dk-preview-mode", mode);
    modeTabs.forEach((button) => {
      const active = button.dataset.mode === mode;
      button.classList.toggle("is-active", active);
      button.setAttribute("aria-selected", active ? "true" : "false");
    });
  }

  function updatePackButtons() {
    Array.from(themePicker.children).forEach((button, index) => {
      const active = index === activePackIndex;
      button.classList.toggle("is-active", active);
      button.setAttribute("aria-selected", active ? "true" : "false");
      button.tabIndex = active ? 0 : -1;
    });
  }

  function selectPackByIndex(index, focusButton = false) {
    if (packs.length === 0) return;

    activePackIndex = (index + packs.length) % packs.length;
    const pack = packs[activePackIndex];

    ensurePackStylesheet(pack);
    desktopKitTheme.setPack(pack.name);
    packDescription.textContent = pack.description;
    packSummary.textContent = pack.displayName + " · " + pack.name;
    packPosition.textContent = (activePackIndex + 1) + " / " + packs.length;
    localStorage.setItem("dk-preview-pack", pack.name);

    updatePackButtons();

    if (focusButton) {
      themePicker.children[activePackIndex]?.focus({ preventScroll: true });
    }
  }

  function renderThemePicker() {
    themePicker.replaceChildren(...packs.map((pack, index) => {
      const button = document.createElement("button");
      button.type = "button";
      button.className = "preview-theme-option";
      button.setAttribute("role", "option");
      button.title = pack.description;
      button.innerHTML = "<strong></strong><small></small>";
      button.querySelector("strong").textContent = pack.displayName;
      button.querySelector("small").textContent = pack.name;
      button.addEventListener("click", () => selectPackByIndex(index));
      return button;
    }));
  }

  function shouldIgnoreShortcut(event) {
    if (event.altKey || event.ctrlKey || event.metaKey) return true;
    const target = event.target;
    if (!(target instanceof HTMLElement)) return false;
    return ["INPUT", "SELECT", "TEXTAREA"].includes(target.tagName) || target.isContentEditable;
  }

  async function init() {
    const response = await fetch("/api/packs", { cache: "no-store" });
    if (!response.ok) {
      throw new Error("主题列表加载失败: HTTP " + response.status);
    }

    packs = await response.json();
    if (!Array.isArray(packs) || packs.length === 0) {
      throw new Error("没有可预览的主题");
    }

    renderThemePicker();

    const savedPack = localStorage.getItem("dk-preview-pack");
    const savedPackIndex = packs.findIndex((pack) => pack.name === savedPack);
    selectPackByIndex(savedPackIndex >= 0 ? savedPackIndex : 0);

    const savedMode = localStorage.getItem("dk-preview-mode");
    const initialMode = ["light", "dark", "system"].includes(savedMode) ? savedMode : "light";
    selectMode(initialMode);

    prevPackButton.addEventListener("click", () => selectPackByIndex(activePackIndex - 1, true));
    nextPackButton.addEventListener("click", () => selectPackByIndex(activePackIndex + 1, true));

    modeTabs.forEach((button) => {
      button.addEventListener("click", () => selectMode(button.dataset.mode));
    });

    document.addEventListener("keydown", (event) => {
      if (shouldIgnoreShortcut(event)) return;

      if (event.key === "ArrowLeft") {
        event.preventDefault();
        selectPackByIndex(activePackIndex - 1);
      } else if (event.key === "ArrowRight") {
        event.preventDefault();
        selectPackByIndex(activePackIndex + 1);
      } else if (event.key.toLowerCase() === "l") {
        selectMode("light");
      } else if (event.key.toLowerCase() === "d") {
        selectMode("dark");
      } else if (event.key.toLowerCase() === "s") {
        selectMode("system");
      }
    });
  }

  init().catch((error) => {
    packDescription.textContent = error.message;
    console.error(error);
  });
})();
