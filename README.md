# hugo-patch

我从 [n0x1m/hugo-latex-patch](https://github.com/n0x1m/hugo-latex-patch) 抄了一份过来，支持 GitHub Action，方便从其他项目取用。

## 用法

我对 [kingreatwill/goldmark-katex](https://github.dev/kingreatwill/goldmark-katex) 进行了一些修改，能够解析 LaTeX 公式但什么都不做，原样交付 HTML，方便 KaTeX 渲染。

如需启用，请配置：

```yaml
markup:
  goldmark:
    extensions:
      LaTeX: true # default false
```

## 本地测试

```bash
$ # export patch
$ git format-patch -1
$ # apply patch
$ git apply ../../patch/20221224-update-add-goldmark-math.patch
```
