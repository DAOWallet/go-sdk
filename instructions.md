# Публикация пакета
* В файле `go.mod` изменить имя модуля. Сейчас просто `daowallet`, обычно указывают что-то вроде `github.com/example-org/daowallet`.
* Изменить имя модуля в тестах. Сейчас модуль в файле `client_test.go` импортирован как:
* Изменить имя модуля в README.md
```go
import (
	...

	"daowallet"
)
```

нужно будет поменять на новое имя:
```go
import (
	...

	"github.com/example-org/daowallet"
)
```
* Публикация пакета. Публиковать нужно полностью исходные коды как есть, в GO все зависимости распостраняются в виде исходных кодов. Корнем импорта нужно считать директорию где расположен файл `go.mod`.
* После этого можно устанавливать пакет командой `go get github.com/example-org/daowallet` и использовать следующим образом:
```go
import (
	...

	"github.com/example-org/daowallet"
)

...

client := daowallet.NewDefaultClient(key, secret)
res, err := client.Withdraw(ctx, ...

```
