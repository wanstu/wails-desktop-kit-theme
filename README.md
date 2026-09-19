# Wails Desktop Kit Theme

`wails-desktop-kit-theme` 是 [Wails Desktop Kit](https://github.com/wanstu/wails-desktop-kit) 的可选通用主题包仓库。

边界保持简单：

- **Kit**：主题协议、默认 light/dark、组件、Runtime。
- **Theme**：可选通用配色，只覆盖 `--dk-*` token。
- **Application**：产品品牌主题和业务 UI。

Theme 不重新定义 `.dk-button`、`.dk-panel`、`.dk-nav-link` 等组件，因此新增主题不需要升级 Kit 本体。

## 要求

需要 Wails Desktop Kit **v0.3.0+**，因为 Theme Pack 使用：

- `data-dk-theme-pack`
- `desktopKitTheme.setPack()`
- `desktopKitTheme.getPack()`
- `desktopKitTheme.clearPack()`

## 安装

~~~powershell
go get github.com/wanstu/wails-desktop-kit@v0.3.0
go get github.com/wanstu/wails-desktop-kit-theme@v0.1.0
~~~

## 挂载

~~~go
//go:embed all:frontend/dist
var frontend embed.FS

func assets() (fs.FS, error) {
    app, err := fs.Sub(frontend, "frontend/dist")
    if err != nil {
        return nil, err
    }
    return theme.MountWithKit(app), nil
}
~~~

需要：

~~~go
import theme "github.com/wanstu/wails-desktop-kit-theme"
~~~

最终资源路径：

~~~text
/desktopkit/tokens.css
/desktopkit/theme.js
/desktopkit/components.css
/desktopkit/navigation.css

/desktopkit-theme/aurora.css
/desktopkit-theme/ocean.css
/desktopkit-theme/forest.css
/desktopkit-theme/sunset.css
~~~

## HTML

主题包 CSS 必须在 `tokens.css` 之后加载：

~~~html
<link rel="stylesheet" href="/desktopkit/tokens.css">
<link rel="stylesheet" href="/desktopkit-theme/aurora.css">
<link rel="stylesheet" href="/desktopkit/base.css">
<link rel="stylesheet" href="/desktopkit/components.css">
<link rel="stylesheet" href="/desktopkit/navigation.css">
<link rel="stylesheet" href="/app.css">

<script src="/desktopkit/theme.js"></script>
<script>
  desktopKitTheme.apply("system");
  desktopKitTheme.setPack("aurora");
</script>
~~~

明暗模式与 Theme Pack 正交：

~~~text
system -> light  + aurora
system -> dark   + aurora

light  + forest
dark   + forest
~~~

应用负责保存用户自己的 mode / pack 选择，Theme 和 Kit 都不碰 localStorage 或业务配置。

## 内置主题

| 稳定 ID | 显示名 | 风格 |
| --- | --- | --- |
| `aurora` | 极光 | 蓝紫冷色，默认推荐 |
| `ocean` | 海洋 | 蓝青清爽 |
| `forest` | 森林 | 柔和自然 |
| `sunset` | 落日 | 橙粉暖色 |

Go 侧也可以读取 manifest：

~~~go
for _, pack := range theme.Packs() {
    fmt.Println(pack.Name, pack.DisplayName)
}
~~~

## 主题贡献契约

新增主题时：

1. 只新增 `assets/<name>.css`。
2. 只覆盖现有 `--dk-*` token。
3. 同时提供 light 默认和 `:root[data-dk-theme="dark"]` 变体。
4. 不复制 Kit 组件 CSS。
5. 不加入具体产品名主题，例如 `frp-blue`、`adm-dark`。

主题需要改变圆角、阴影等时，也应通过 `--dk-radius-*`、`--dk-shadow-*` 等 token 完成。

## 版本关系

Theme 和 Kit 独立发布：

~~~text
wails-desktop-kit        v0.3.0
wails-desktop-kit-theme  v0.1.x
~~~

新增主题通常只发布 Theme，不要求 Kit 发版。
