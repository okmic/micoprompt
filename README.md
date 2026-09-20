# MicoPrompt

🗂️ Собирает файлы проекта в один промпт для ИИ.

## 📦 Установка

Готовый бинарь в корне репо. macOS ARM64, для других ОС собирайте сами

```bash
chmod +x micoprompt
./micoprompt
```

Собрать самому:

```bash
go build -o micoprompt ./cmd
```

## 🚀 Использование

```bash
micoprompt
```

Рядом появится `micoprompt.txt`.


# Не забудь игнорить в проекте файлы **/micoprompt*
.gitignore
**/micoprompt*

```bash
micoprompt -dir ./src
micoprompt -out - | pbcopy
micoprompt -format xml
micoprompt -ext .go,.md
```

## ⚙️ Флаги

| Флаг | По умолчанию |
|---|---|
| `-dir` | `.` |
| `-out` | `micoprompt.txt` |
| `-ext` | `.go,.md,.txt,.docx,...` |
| `-ignore` | `.git,node_modules,...` |
| `-format` | `markdown` |

MIT © [Mico](https://github.com/okmic)