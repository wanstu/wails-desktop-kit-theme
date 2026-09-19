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
...
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
| `aurora` | 极光 | 蓝紫冷色，通用默认 |
| `ocean` | 海洋 | 蓝青清爽 |
| `forest` | 森林 | 柔和自然 |
| `sunset` | 落日 | 橙粉暖色 |
| `nord` | 北境 | 低饱和蓝灰 |
| `mint` | 薄荷 | 清透青绿 |
| `rose` | 蔷薇 | 柔和玫红 |
| `amber` | 琥珀 | 金黄暖色 |
| `graphite` | 石墨 | 中性灰阶、紧凑圆角 |
| `mocha` | 摩卡 | 咖啡棕与暖米色 |
| `lavender` | 薰衣草 | 浅紫、圆润 |
| `midnight` | 午夜 | 深蓝靛色、高对比 |
| `cobalt` | 钴蓝 | 清晰强烈的蓝色工具风格 |
| `jade` | 翡翠 | 玉石绿、圆润沉稳 |
| `cherry` | 樱桃 | 鲜明红粉 |
| `sand` | 沙丘 | 米沙与棕褐暖色 |
| `plum` | 梅紫 | 成熟紫红 |
| `ice` | 冰川 | 冰蓝灰与清冷青 |
| `terminal` | 终端 | 低圆角绿色开发工具风格 |
| `cyber` | 赛博 | 紫色主调、青色强调 |

Go 侧也可以读取 manifest：

~~~go
for _, pack := range theme.Packs() {
    fmt.Println(pack.Name, pack.DisplayName)
}
~~~

## Theme Gallery

仓库自带真实 Kit 组件预览页，可以同时检查主题和 light / dark / system 模式：

~~~powershell
cd D:\projects\wails-desktop-kit-theme
go run ./cmd/preview
~~~

然后打开：

~~~text
http://127.0.0.1:8080
~~~

也可以修改监听地址：

~~~powershell
go run ./cmd/preview -addr 127.0.0.1:8090
~~~

Gallery 会从 `Packs()` 自动读取主题清单，因此新增并注册 Theme Pack 后无需再维护预览页面。

## 主题贡献契约

新增主题时：

1. 新增 `assets/<name>.css`。
2. 在 `pack.go` 中注册对应 `Pack` 元数据。
3. 只覆盖现有 `--dk-*` token。
4. 同时提供 light 默认和 `:root[data-dk-theme="dark"]` 变体。
5. 不复制 Kit 组件 CSS。
6. 不加入具体产品名主题，例如 `frp-blue`、`adm-dark`。

主题需要改变圆角、阴影等时，也应通过 `--dk-radius-*`、`--dk-shadow-*` 等 token 完成。

测试会同时校验：

- 注册的 Pack 必须存在对应 CSS。
- `assets/*.css` 必须全部注册。
- Theme CSS 不允许覆盖 `.dk-*` 组件选择器。
- Theme CSS 不允许使用 `!important`。
- 自定义属性声明只能使用 `--dk-*`。

## 版本关系

Theme 和 Kit 独立发布：

~~~text
wails-desktop-kit        v0.3.0
wails-desktop-kit-theme  v0.1.x
~~~

新增主题通常只发布 Theme，不要求 Kit 发版。
