# log — довідник

Повний довідник пакета `log`: ментальна модель, виходи та їх конфігурація, рівні
й макети, усі методи логування, міст до `slog` і практичні рецепти.

Англійська версія: **[DOC.md](DOC.md)**.

## Зміст

- [Ментальна модель](#ментальна-модель)
- [Логери й дефолтний логер](#логери-й-дефолтний-логер)
- [Методи логування](#методи-логування)
- [Рівні](#рівні)
- [Виходи (Outputs)](#виходи-outputs)
- [Макети й форматування](#макети-й-форматування)
- [Текстовий і JSON-вивід](#текстовий-і-json-вивід)
- [Керування виходами](#керування-виходами)
- [Ad-hoc writer'и](#ad-hoc-writerи)
- [Структуроване логування зі slog](#структуроване-логування-зі-slog)
- [Умовне логування й обробка помилок](#умовне-логування-й-обробка-помилок)
- [Кадри стеку й префікси](#кадри-стеку-й-префікси)
- [Рецепти й поради](#рецепти-й-поради)

## Ментальна модель

`log` — це рівневий логер із багатьма виходами. Один виклик логування
розгалужується на кожен налаштований **вихід**; кожен вихід самостійно вирішує,
які рівні приймати, рендерити текст чи JSON, використовувати колір і як
викладати префікс повідомлення.

API формують дві ідеї:

1. **Одне повідомлення, багато виходів.** Ви не створюєте окремий логер під
   кожен пункт призначення. Ви приєднуєте кілька значень `Output` до одного
   логера — скажімо, кольорову консоль для `Info+` і JSON-файл для `Error+` — і
   кожен виклик маршрутизується до тих, кому він потрібен.
2. **Робота пропускається, коли нікому не потрібна.** Фільтрація за рівнем
   відбувається біля джерела, а захоплення кадру стеку й форматування робляться,
   лише коли їх вимагає хоча б один вихід, тож «тихі» рівні дешеві.

Логер безпечний для конкурентного використання, а гарячий шлях пулить буфери,
щоб тримати алокації низькими.

```go
import (
    "github.com/goloop/log/v2"
    "github.com/goloop/log/v2/level"
    "github.com/goloop/log/v2/layout"
)
```

## Логери й дефолтний логер

```go
func New(prefixes ...string) *Logger
func (logger *Logger) Copy() *Logger

func Log() *Logger
func SetDefault(logger *Logger)
```

`New` створює логер; будь-які префікси обрізаються й з'єднуються дефісами
(`New("APP", "API")` → `APP-API`) та додаються до кожного повідомлення. `Copy`
клонує логер із його виходами, тож його можна підправити, не чіпаючи оригінал.

Пакет також тримає дефолтний логер, який використовують пакетні функції
(`log.Info`, `log.Errorf`, …). `Log` повертає його; `SetDefault` замінює його
(зручно в тестах чи щоб встановити наперед налаштований логер на старті).

```go
logger := log.New("APP")
logger.Info("started")

log.Info("використовує дефолтний логер")
```

## Методи логування

Кожен рівень має три форми, і на `*Logger`, і на пакеті:

| Форма | Приклад | Поведінка |
|-------|---------|-----------|
| проста | `Info(a ...any)`             | як `fmt.Print` |
| `…f`   | `Infof(format string, a ...any)` | як `fmt.Printf` |
| `…ln`  | `Infoln(a ...any)`           | як `fmt.Println` |

Рівні від найсуворішого до найм'якшого: `Panic`, `Fatal`, `Error`, `Warn`,
`Info`, `Debug`, `Trace`:

- `Panic*` логує, тоді викликає `panic()`.
- `Fatal*` логує, тоді викликає `os.Exit(1)`.
- Решта логують і повертають.

```go
logger.Info("Application started")
logger.Infof("User %s logged in", user)
logger.Errorln("Failed to connect to database")
```

Є також родина з префіксом `F` (`Finfo`, `Ferrorf`, …) — див.
[Ad-hoc writer'и](#ad-hoc-writerи).

## Рівні

`level.Level` — це набір бітових прапорів із підпакета `level`:

```go
const (
    Panic level.Level = 1 << iota
    Fatal
    Error
    Warn
    Info
    Debug
    Trace
)
var Default = Panic | Fatal | Error | Warn | Info | Debug | Trace
```

Поєднуйте їх через `|`, щоб задати, які рівні приймає вихід:

```go
Levels: level.Info | level.Warn | level.Error
```

`Enabled(l)` повідомляє, чи прийняв би якийсь вихід рівень `l`, — використовуйте
його, щоб захистити дорогі аргументи (див.
[Умовне логування](#умовне-логування-й-обробка-помилок)).

## Виходи (Outputs)

`Output` — це пункт призначення плюс його правила рендерингу. `Name` і `Writer`
обов'язкові; решта мають розумні дефолти.

| Поле | Тип | Призначення |
|------|-----|-------------|
| `Name`            | `string`         | унікальний ідентифікатор для `Outputs`/`EditOutputs`/`DeleteOutputs` |
| `Writer`          | `io.Writer`      | куди йдуть байти (`os.Stdout`, файл, власний writer) |
| `Levels`          | `level.Level`    | які рівні приймає цей вихід |
| `Layouts`         | `layout.Layout`  | які блоки контексту виклику включати в префікс |
| `Space`           | `string`         | роздільник між блоками префікса |
| `WithPrefix`      | `trit.Trit`      | показувати префікс логера (за замовч. увімкнено) |
| `WithColor`       | `trit.Trit`      | ANSI-колір за рівнем (лише текст, UNIX-подібні; за замовч. вимкнено) |
| `Enabled`         | `trit.Trit`      | увімкнути/вимкнути вихід (за замовч. увімкнено) |
| `TextStyle`       | `trit.Trit`      | текст (`true`) vs JSON (`false`); за замовч. текст |
| `TimestampFormat` | `string`         | макет `time.Format` для мітки часу |
| `LevelFormat`     | `string`         | обгортка навколо назви рівня, напр. `"[%s]"` |

Поля `trit.Trit` використовують трійкову логіку (з пакета `trit`): значення
`> 0` — істина, `< 0` — хиба, а `0` означає «дефолт» (або «не змінювати» в
режимі редагування). Можна передавати сирі `1`/`-1` чи `trit.True`/`trit.False`.

Є два готові виходи: `log.Stdout` і `log.Stderr`.

```go
log.SetOutputs(
    log.Output{
        Name:      "console",
        Writer:    os.Stdout,
        Levels:    level.Info | level.Warn | level.Error,
        Layouts:   layout.Default,
        WithColor: 1,
        TextStyle: 1,
    },
    log.Output{
        Name:       "file",
        Writer:     file,
        Levels:     level.Error | level.Fatal,
        TextStyle:  -1, // JSON
        WithPrefix: 1,
    },
)
```

## Макети й форматування

`layout.Layout` — це набір бітових прапорів, що керує тим, які блоки контексту
виклику з'являються в префіксі повідомлення:

```go
const (
    FullFilePath layout.Layout = 1 << iota
    ShortFilePath
    FuncName
    FuncAddress
    LineNumber
)
var Default = ShortFilePath | FuncName | LineNumber
```

```go
Layouts: layout.FullFilePath | layout.FuncName | layout.LineNumber
```

`TimestampFormat` і `LevelFormat` додатково формують префікс; `Space` задає
роздільник між блоками.

## Текстовий і JSON-вивід

Вихід рендерить текст, коли `TextStyle` істинний (дефолт), і JSON, коли хибний.

Текст:

```
APP: 2023/12/02 15:04:05 INFO main.go:42 Starting application
```

JSON (порожні поля опускаються):

```json
{
    "prefix": "APP",
    "level": "INFO",
    "timestamp": "2023/12/02 15:04:05",
    "message": "Starting application",
    "filePath": "/home/user/app/main.go",
    "lineNumber": 42,
    "funcName": "main"
}
```

Ключі JSON: `prefix`, `level`, `timestamp`, `message`, `filePath`,
`lineNumber`, `funcName` і `funcAddress`.

## Керування виходами

```go
func (logger *Logger) SetOutputs(outputs ...Output) error
func (logger *Logger) EditOutputs(outputs ...Output) error
func (logger *Logger) DeleteOutputs(names ...string)
func (logger *Logger) Outputs(names ...string) []Output
```

`SetOutputs` замінює весь набір. `EditOutputs` змінює названі виходи на місці —
застосовуються лише задані вами поля (поле `trit`, лишене на `0`, не чіпається),
тож можна перемкнути колір чи рівні, не перевизначаючи writer. `DeleteOutputs`
видаляє виходи за іменем; `Outputs` повертає всі або лише названі.

```go
logger.EditOutputs(log.Output{Name: "console", Levels: level.Error | level.Fatal})
logger.EditOutputs(log.Output{Name: "console", Enabled: -1}) // вимкнути
logger.DeleteOutputs("file")
```

## Ad-hoc writer'и

Методи з префіксом `F` (`Finfo`, `Ferrorf`, `Fdebugln`, …) пишуть у налаштовані
виходи **й додатково** у writer, переданий першим аргументом, не змінюючи
конфігурацію логера:

```go
var buf bytes.Buffer
logger.Finfo(&buf, "захоплено тут і в налаштованих виходах")
```

Це корисно для захоплення конкретного повідомлення в тесті чи буфері в межах
запиту, зберігаючи звичайне логування.

## Структуроване логування зі slog

Логер може бути бекендом для `log/slog` зі стандартної бібліотеки:

```go
func NewSlog(prefixes ...string) *slog.Logger
func (logger *Logger) Handler() slog.Handler
```

```go
slogger := log.NewSlog("APP")
slogger.Info("user logged in", "user", "bob", "id", 42)

// Або приєднати handler до наявного логера.
logger := log.New("APP")
slogger = slog.New(logger.Handler())
```

Рівні slog відображаються на рівні логера (Debug, Info, Warn, Error). Атрибути
запису — включно з доданими через `With`/`WithGroup` — стають типізованими
полями JSON у JSON-виходах і парами `key=value` в текстових виходах.

## Умовне логування й обробка помилок

```go
func (logger *Logger) Enabled(l level.Level) bool
func (logger *Logger) SetErrorHandler(handler func(o Output, n int, err error))
```

Захищайте дорогі аргументи через `Enabled`, щоб нічого не обчислювалось для
рівня, який жодному виходу не потрібен:

```go
if logger.Enabled(level.Debug) {
    logger.Debug(expensiveDump())
}
```

За замовчуванням запис — best-effort, а помилки запису ігноруються. Зареєструйте
handler, щоб їх спостерігати, — наприклад, щоб сповіщати про несправний файловий
чи мережевий вихід:

```go
logger.SetErrorHandler(func(o log.Output, n int, err error) {
    fmt.Fprintf(os.Stderr, "log output %q failed: %v\n", o.Name, err)
})
```

## Кадри стеку й префікси

```go
func (logger *Logger) SetSkipStackFrames(skip int) int
func (logger *Logger) SkipStackFrames() int
func (logger *Logger) SetPrefix(prefix string) string
func (logger *Logger) Prefix() string
```

Коли ви обгортаєте логер власним помічником, зафіксований файл/рядок вказує на
обгортку. `SetSkipStackFrames` каже логеру, скільки кадрів пропустити, щоб
натомість фіксувалося місце виклику. `SetPrefix`/`Prefix` читають чи змінюють
префікс після створення.

```go
logger.SetSkipStackFrames(2) // пропустити функції-обгортки
```

## Рецепти й поради

**Розділення консоль + файл.** Приєднайте кольорову текстову консоль для `Info+`
і JSON-файл для `Error+`; один виклик живить обидва, кожен фільтрується
незалежно.

**Продакшн-рівні.** Тримайте `Debug`/`Trace` поза продакшн-виходами, щоб
форматування й захоплення стеку для цих рівнів пропускалися повністю; поєднуйте
з `Enabled`, щоб не будувати дорогі дебаг-навантаження.

**Захоплення в межах запиту.** Використовуйте родину `F`, щоб роздвоїти
повідомлення в буфер у межах запиту, поки звичайне логування триває.

**Перемикання без перевизначення.** `EditOutputs` застосовує лише задані поля —
перемкніть `WithColor`, змініть `Levels` чи вимкніть через `Enabled: -1`, не
перезадаючи writer.

**Міст до наявного slog-коду.** Якщо застосунок уже логує через `slog`,
встановіть `logger.Handler()`, щоб його записи текли крізь ті самі виходи,
кольори й формати.
