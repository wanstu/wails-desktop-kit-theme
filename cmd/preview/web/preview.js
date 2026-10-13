(() => {
  "use strict";

  const themePicker = document.getElementById("theme-picker");
  const packDescription = document.getElementById("pack-description");
  const packSummary = document.getElementById("pack-summary");
  const packPosition = document.getElementById("pack-position");
  const prevPackButton = document.getElementById("prev-pack");
  const nextPackButton = document.getElementById("next-pack");
  const modeTabs = Array.from(document.querySelectorAll("[data-mode]"));
  const styleTabs = Array.from(document.querySelectorAll("[data-style]"));
  const builder = document.getElementById("gradient-builder");
  const gradientOrder = document.getElementById("gradient-order");
  const gradientMessage = document.getElementById("gradient-message");
  const showSignature = document.getElementById("show-signature");
  const showAll = document.getElementById("show-all");
  const compareGrid = document.getElementById("theme-compare-grid");
  const loadedStyles = new Set();
  const paletteCache = new Map();

  let packs = [];
  let activePackIndex = 0;
  let basePalettePromise = null;
  let themeStyle = "solid";
  let gradientSelection = ["aurora", "sakura"];
  let gradientRequest = 0;
  let signatureOnly = true;

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

  function findRuleBody(css, selector) {
    const selectorIndex = css.indexOf(selector);
    if (selectorIndex < 0) return "";

    const openIndex = css.indexOf("{", selectorIndex + selector.length);
    if (openIndex < 0) return "";

    const closeIndex = css.indexOf("}", openIndex + 1);
    if (closeIndex < 0) return "";

    return css.slice(openIndex + 1, closeIndex);
  }

  function readTokens(ruleBody) {
    const result = {};
    const tokenPattern = /(--dk-[a-z0-9-]+)\s*:\s*([^;]+);/gi;
    let match;

    while ((match = tokenPattern.exec(ruleBody)) !== null) {
      result[match[1]] = match[2].trim();
    }

    return result;
  }

  async function loadBasePalette() {
    if (basePalettePromise) return basePalettePromise;

    basePalettePromise = fetch("/desktopkit/tokens.css", { cache: "no-store" })
      .then((response) => {
        if (!response.ok) throw new Error("Kit tokens 加载失败");
        return response.text();
      })
      .then((css) => ({
        light: readTokens(findRuleBody(css, ":root")),
        dark: readTokens(findRuleBody(css, ':root[data-dk-theme="dark"]')),
      }));

    return basePalettePromise;
  }

  async function loadPackPalette(pack) {
    if (paletteCache.has(pack.name)) return paletteCache.get(pack.name);

    const [base, response] = await Promise.all([
      loadBasePalette(),
      fetch(pack.stylesheetPath, { cache: "no-store" }),
    ]);

    if (!response.ok) throw new Error(pack.name + " 主题 CSS 加载失败");

    const css = await response.text();
    const lightSelector = ':root[data-dk-theme-pack="' + pack.name + '"]';
    const darkSelector = ':root[data-dk-theme="dark"][data-dk-theme-pack="' + pack.name + '"]';

    const palette = {
      light: { ...base.light, ...readTokens(findRuleBody(css, lightSelector)) },
      dark: { ...base.dark, ...readTokens(findRuleBody(css, darkSelector)) },
    };

    paletteCache.set(pack.name, palette);
    return palette;
  }

  async function hydrateThemeSwatches(pack, button) {
    try {
      const palette = await loadPackPalette(pack);
      const colors = [
        palette.light["--dk-bg-app"],
        palette.light["--dk-primary"],
        palette.dark["--dk-bg-app"],
        palette.dark["--dk-primary"],
      ];
      const swatches = button.querySelectorAll(".preview-theme-swatch");

      colors.forEach((color, index) => {
        if (swatches[index] && color) {
          swatches[index].style.setProperty("--preview-swatch", color);
        }
      });
    } catch (error) {
      console.error(error);
    }
  }

  function selectMode(mode) {
    desktopKitTheme.apply(mode);
    localStorage.setItem("dk-preview-mode", mode);
    modeTabs.forEach((button) => {
      const active = button.dataset.mode === mode;
      button.classList.toggle("is-active", active);
      button.setAttribute("aria-selected", active ? "true" : "false");
    });
    if (themeStyle === "gradient") renderGradientOrder();
  }

  function updatePackButtons() {
    Array.from(themePicker.children).forEach((button, index) => {
      const pack = packs[index];
      const order = gradientSelection.indexOf(pack.name);
      button.hidden = signatureOnly && !pack.signature && !(themeStyle === "gradient" && order >= 0);
      const active = themeStyle === "gradient" ? order >= 0 : index === activePackIndex;
      button.classList.toggle("is-active", themeStyle === "solid" && active);
      button.classList.toggle("is-gradient-selected", themeStyle === "gradient" && active);
      button.dataset.gradientOrder = themeStyle === "gradient" && active ? String(order + 1) : "";
      button.setAttribute("aria-selected", String(active));
      button.tabIndex = 0;
    });
  }


  function visiblePackIndices() {
    return packs.flatMap((pack, index) => signatureOnly && !pack.signature ? [] : [index]);
  }

  function navigatePack(direction, focus = false) {
    const visible = visiblePackIndices();
    if (!visible.length) return;
    const index = visible.indexOf(activePackIndex);
    selectPackByIndex(visible[(index + direction + visible.length) % visible.length], focus);
  }

  function setVisibility(onlySignature) {
    signatureOnly = onlySignature;
    showSignature.classList.toggle("is-active", signatureOnly);
    showAll.classList.toggle("is-active", !signatureOnly);
    showSignature.setAttribute("aria-pressed", String(signatureOnly));
    showAll.setAttribute("aria-pressed", String(!signatureOnly));
    localStorage.setItem("dk-preview-signature-only", String(signatureOnly));
    if (themeStyle === "solid" && signatureOnly && !packs[activePackIndex].signature) {
      selectPackByIndex(visiblePackIndices()[0]);
    }
    updatePackButtons();
    if (themeStyle === "solid") {
      const visible = visiblePackIndices();
      packPosition.textContent = (visible.indexOf(activePackIndex) + 1) + " / " + visible.length;
    }
  }

  function mockupForTheme(palette, mode) {
    const mock = document.createElement("div");
    mock.className = "preview-compare-mock";
    const vars = {
      "bg": "--dk-bg-app", "text": "--dk-text", "subtle": "--dk-bg-subtle",
      "muted": "--dk-bg-muted", "elevated": "--dk-bg-elevated",
      "border": "--dk-border", "primary": "--dk-primary",
      "background": "--dk-app-gradient", "sidebar-gradient": "--dk-sidebar-gradient",
      "panel-gradient": "--dk-panel-gradient", "primary-gradient": "--dk-primary-gradient",
      "radius": "--dk-radius-lg", "small-radius": "--dk-radius-sm",
      "shadow": "--dk-shadow-panel"
    };
    Object.entries(vars).forEach(([name, token]) => {
      const value = palette[token];
      if (value) mock.style.setProperty("--mini-" + name, value);
    });
    const modeLabel = document.createElement("small");
    modeLabel.textContent = mode === "dark" ? "DARK" : "LIGHT";
    const content = document.createElement("div");
    content.className = "preview-compare-mock-content";
    const sidebar = document.createElement("div");
    sidebar.className = "preview-compare-mock-sidebar";
    const panel = document.createElement("div");
    panel.className = "preview-compare-mock-panel";
    ["preview-compare-mock-line", "preview-compare-mock-line short", "preview-compare-mock-action"]
      .forEach((name) => {
        const shape = document.createElement("div");
        shape.className = name;
        panel.appendChild(shape);
      });
    content.append(sidebar, panel);
    mock.append(modeLabel, content);
    return mock;
  }

  async function renderComparisons() {
    const signatures = packs.filter((pack) => pack.signature);
    const cards = await Promise.all(signatures.map(async (pack) => {
      const palette = await loadPackPalette(pack);
      const button = document.createElement("button");
      button.type = "button";
      button.className = "preview-compare-card";
      const title = document.createElement("div");
      title.className = "preview-compare-card-header";
      const heading = document.createElement("strong");
      heading.textContent = pack.displayName;
      const subtitle = document.createElement("small");
      subtitle.textContent = pack.name;
      title.append(heading, subtitle);
      const mockups = document.createElement("div");
      mockups.className = "preview-compare-modes";
      mockups.append(mockupForTheme(palette.light, "light"), mockupForTheme(palette.dark, "dark"));
      button.append(title, mockups);
      button.title = pack.description;
      button.addEventListener("click", () => {
        if (themeStyle !== "solid") setStyle("solid");
        selectPackByIndex(packs.indexOf(pack));
        document.getElementById("overview").scrollIntoView({ behavior: "smooth" });
      });
      return button;
    }));
    compareGrid.replaceChildren(...cards);
  }

  function selectPackByIndex(index, focusButton = false) {
    if (packs.length === 0) return;

    activePackIndex = (index + packs.length) % packs.length;
    const pack = packs[activePackIndex];

    ensurePackStylesheet(pack);
    desktopKitTheme.setPack(pack.name);
    packDescription.textContent = pack.description;
    packSummary.textContent = pack.displayName + " · " + pack.name;
    const visible = visiblePackIndices();
    packPosition.textContent = (visible.indexOf(activePackIndex) + 1) + " / " + visible.length;
    localStorage.setItem("dk-preview-pack", pack.name);

    updatePackButtons();

    const activeButton = themePicker.children[activePackIndex];
    if (activeButton) {
      activeButton.scrollIntoView({ block: "nearest" });
      if (focusButton) activeButton.focus({ preventScroll: true });
    }
  }

  function renderThemePicker() {
    const buttons = packs.map((pack, index) => {
      const button = document.createElement("button");
      button.type = "button";
      button.className = "preview-theme-option";
      button.setAttribute("role", "option");
      button.title = pack.description;
      button.innerHTML = [
        '<span class="preview-theme-option-copy">',
        "<strong></strong>",
        "<small></small>",
        "</span>",
        '<span class="preview-theme-swatches" aria-hidden="true">',
        '<i class="preview-theme-swatch"></i>',
        '<i class="preview-theme-swatch"></i>',
        '<i class="preview-theme-swatch"></i>',
        '<i class="preview-theme-swatch"></i>',
        "</span>",
      ].join("");

      button.querySelector("strong").textContent = pack.displayName;
      button.querySelector("small").textContent = pack.name;
      button.addEventListener("click", () => {
        if (themeStyle !== "gradient") { selectPackByIndex(index); return; }
        const name = packs[index].name;
        const next = gradientSelection.includes(name)
          ? gradientSelection.filter((item) => item !== name)
          : [...gradientSelection, name];
        updateGradient(next);
      });
      hydrateThemeSwatches(pack, button);
      return button;
    });

    themePicker.replaceChildren(...buttons);
  }


  function renderGradientOrder() {
    const dark = desktopKitTheme.getResolvedTheme() === "dark";
    gradientOrder.replaceChildren(...gradientSelection.map((name, index) => {
      const pack = packs.find((item) => item.name === name);
      const item = document.createElement("div");
      item.className = "preview-gradient-stop";
      const color = document.createElement("span");
      color.className = "preview-gradient-stop-color";
      color.style.setProperty("--stop-color", dark ? pack.gradientDark : pack.gradientLight);
      const label = document.createElement("span");
      label.className = "preview-gradient-stop-label";
      label.textContent = (index + 1) + ". " + pack.displayName;
      item.append(color, label);
      const controls = [
        ["←", "前移", index > 0, () => reorderGradient(index, -1)],
        ["→", "后移", index < gradientSelection.length - 1, () => reorderGradient(index, 1)],
        ["×", "移除", gradientSelection.length > 2, () =>
          updateGradient(gradientSelection.filter((entry) => entry !== name))]
      ];
      controls.forEach(([text, action, enabled, callback]) => {
        const button = document.createElement("button");
        button.type = "button";
        button.textContent = text;
        button.title = action + " " + pack.displayName;
        button.setAttribute("aria-label", button.title);
        button.disabled = !enabled;
        button.addEventListener("click", callback);
        item.appendChild(button);
      });
      return item;
    }));
  }

  function reorderGradient(index, direction) {
    const next = gradientSelection.slice();
    const other = index + direction;
    [next[index], next[other]] = [next[other], next[index]];
    updateGradient(next);
  }

  async function updateGradient(names) {
    if (names.length < 2 || names.length > 4) {
      gradientMessage.textContent = "请选择 2～4 个主题。";
      return;
    }
    const request = ++gradientRequest;
    const previous = gradientSelection.slice();
    gradientSelection = names.slice();
    renderGradientOrder();
    updatePackButtons();
    try {
      await desktopKitTheme.applyGradient(names);
      if (request !== gradientRequest || themeStyle !== "gradient") return;
      localStorage.setItem("dk-preview-gradient", JSON.stringify(gradientSelection));
      gradientMessage.textContent = "已应用：" + names.map((name) =>
        packs.find((pack) => pack.name === name).displayName).join(" → ");
      packDescription.textContent = "渐变背景与操作强调；正文和卡片保持清晰。";
      packSummary.textContent = names.map((name) =>
        packs.find((pack) => pack.name === name).displayName).join(" → ");
      packPosition.textContent = names.length + " 种配色";
      renderGradientOrder();
      updatePackButtons();
    } catch (error) {
      if (request !== gradientRequest || themeStyle !== "gradient") return;
      gradientSelection = previous;
      renderGradientOrder();
      updatePackButtons();
      gradientMessage.textContent = error.message;
      console.error(error);
    }
  }

  function setStyle(nextStyle) {
    themeStyle = nextStyle;
    builder.hidden = nextStyle !== "gradient";
    styleTabs.forEach((button) => {
      const active = button.dataset.style === nextStyle;
      button.classList.toggle("is-active", active);
      button.setAttribute("aria-selected", String(active));
    });
    prevPackButton.hidden = nextStyle === "gradient";
    nextPackButton.hidden = nextStyle === "gradient";
    localStorage.setItem("dk-preview-style", nextStyle);
    if (nextStyle === "gradient") updateGradient(gradientSelection);
    else {
      gradientRequest++;
      desktopKitTheme.clearGradient();
      selectPackByIndex(activePackIndex);
    }
    updatePackButtons();
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
    signatureOnly = localStorage.getItem("dk-preview-signature-only") !== "false";
    if (savedPackIndex >= 0 && !packs[savedPackIndex].signature) signatureOnly = false;
    setVisibility(signatureOnly);
    showSignature.addEventListener("click", () => setVisibility(true));
    showAll.addEventListener("click", () => setVisibility(false));
    renderComparisons().catch((error) => console.error("主题对比加载失败", error));

    const savedMode = localStorage.getItem("dk-preview-mode");
    const initialMode = ["light", "dark", "system"].includes(savedMode) ? savedMode : "light";
    selectMode(initialMode);

    prevPackButton.addEventListener("click", () => navigatePack(-1, true));
    nextPackButton.addEventListener("click", () => navigatePack(1, true));

    modeTabs.forEach((button) => {
      button.addEventListener("click", () => selectMode(button.dataset.mode));
    });
    styleTabs.forEach((button) => {
      button.addEventListener("click", () => setStyle(button.dataset.style));
    });
    document.querySelectorAll("[data-gradient-preset]").forEach((button) => {
      button.addEventListener("click", () => updateGradient(button.dataset.gradientPreset.split(",")));
    });
    try {
      const saved = JSON.parse(localStorage.getItem("dk-preview-gradient") || "null");
      if (Array.isArray(saved) && saved.length >= 2 && saved.length <= 4 &&
          new Set(saved).size === saved.length &&
          saved.every((name) => packs.some((pack) => pack.name === name))) {
        gradientSelection = saved;
      }
    } catch (_) {
      // Ignore malformed gallery preferences.
    }
    setStyle(localStorage.getItem("dk-preview-style") === "gradient" ? "gradient" : "solid");

    document.addEventListener("keydown", (event) => {
      if (shouldIgnoreShortcut(event)) return;

      if (themeStyle === "solid" && event.key === "ArrowLeft") {
        event.preventDefault();
        navigatePack(-1);
      } else if (themeStyle === "solid" && event.key === "ArrowRight") {
        event.preventDefault();
        navigatePack(1);
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
