MicoPrompt
# micoprompt

> Собирает файлы проекта в один текстовый промпт для ИИ.

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Зачем

Скормить весь проект ChatGPT/Claude/DeepSeek одним сообщением —
больно. `micoprompt` обходит папку, отфильтровывает мусор и склеивает
всё в один аккуратный файл.

## Установка

```bash
go install github.com/Mico/micoprompt/cmd/micoprompt@latest
```

Или скачайте готовый бинарник из [Releases](../../releases).

## Использование

```bash
micoprompt -dir ./my-project -out prompt.txt
```

### Примеры

Собрать только Go и Markdown:
```bash
micoprompt -dir ./src -ext .go,.md -out prompt.txt
```

Вывести в stdout и скопировать в буфер:
```bash
micoprompt -dir ./src -out - | pbcopy       # macOS
micoprompt -dir ./src -out - | xclip -sel c # Linux
```

Формат XML (для Claude):
```bash
micoprompt -dir ./src -format xml -out prompt.xml
```

## Флаги

| Флаг            | По умолчанию                              | Описание                        |
|-----------------|-------------------------------------------|---------------------------------|
| `-dir`          | `.`                                       | Папка для обхода                |
| `-out`          | `prompt.txt`                              | Выходной файл (`-` = stdout)    |
| `-ext`          | `.go,.md,.txt,.mod,.yaml,.yml,.json`      | Расширения                      |
| `-ignore`       | `.git,node_modules,vendor,dist,build,...` | Игнорируемые папки              |
| `-ignore-files` | `.env,.DS_Store`                          | Игнорируемые имена файлов       |
| `-max-size`     | `102400`                                  | Макс. размер файла (байт)       |
| `-max-total`    | `0`                                       | Макс. общий размер (0 = без)    |
| `-format`       | `markdown`                                | `plain` \| `markdown` \| `xml`  |
| `-header`       | —                                         | Заголовок промпта               |
| `-footer`       | —                                         | Футер промпта                   |

## Форматы вывода

**markdown** (по умолчанию):
```markdown
## main.go

```go
package main
...
```
```

**xml** (для Claude):
```xml
<file path="main.go">
package main
...
</file>
```

**plain**:
```
===== main.go =====

package main
...
```

## Лицензия

MIT © Mico