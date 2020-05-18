# daowallet - This is Golang library for cryptocurrency payments api of daowallet.com ...

Add your description ...

## Installation

Import the package and use it:

```go
import (
	github.com/example-org/daowallet
)

```

### Older versions of Go (<= 1.12)

To install the package, use `go get`:

```
go get github.com/example-org/daowallet

```

Or import the package into your project and then do `go get .`

## Usage

### Obtaining crypto-address

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	foreignID := "user-1250"
	currency := "BTC"

	adr, err := client.Addresses(context.Background(), foreignID, currency)
	if err != nil {
		panic(err)
	}

	fmt.Println("crypto address: %v", adr)
}
```

### Withdrawing cryptocurrency to crypto address

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	foreignID := "user-1250"
	amount := 0.01
	currency := "BTC"
	address := "1MDY9GwakAUYKPXEnFxrTEufzZW3kTE7Rx"

	wtl, err := client.Withdraw(context.Background(), foreignID, amount, currency, address)
	if err != nil {
		panic(err)
	}

	fmt.Println("withdrawal: %v", wtl)
}
```

### Invoice issuing

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	amount := 1250,
	fiatCurrency := "USD"

	inv, err := client.InvoiceNew(context.Background(), amount, fiatCurrency)
	if err != nil {
		panic(err)
	}

	fmt.Println("invoice: %v", inv)
}
```

### Invoice status checking

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	id := "eif0Z2bfnkY6WU5mg7gIqTUQBgDs5zWI"

	inv, err := client.InvoiceStatus(context.Background(), id)
	if err != nil {
		panic(err)
	}

	fmt.Println("invoice: %v", inv)
}
```


# daowallet - Golang библиотека для сервиса криптовалютных платежей компании daowallet.com ...

Здесь какое-то описание ...

## Установка

Для установки пакета используйте команду `go get`:

```
go get github.com/example-org/daowallet

```

Или импортируйте пакет в ваш проект и затем выполните команду `go get .`

## Использование

### Получение крипто-адреса

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	foreignID := "user-1250"
	currency := "BTC"

	adr, err := client.Addresses(context.Background(), foreignID, currency)
	if err != nil {
		panic(err)
	}

	fmt.Println("crypto address: %v", adr)
}
```

### Вывод криптовалюты на крипто-адрес

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	foreignID := "user-1250"
	amount := 0.01
	currency := "BTC"
	address := "1MDY9GwakAUYKPXEnFxrTEufzZW3kTE7Rx"

	wtl, err := client.Withdraw(context.Background(), foreignID, amount, currency, address)
	if err != nil {
		panic(err)
	}

	fmt.Println("withdrawal: %v", wtl)
}
```

### Выставление инвойса

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	amount := 1250,
	fiatCurrency := "USD"

	inv, err := client.InvoiceNew(context.Background(), amount, fiatCurrency)
	if err != nil {
		panic(err)
	}

	fmt.Println("invoice: %v", inv)
}
```

### Получение статуса инвойса

```go
package main

import (
	"context"
	"fmt"

	"daowallet"
)

func main() {
	key := "sA4kH4BX4IAKBzm5DpOFoHL6XoUNJ0sP"
	secret := "GAD0DcpFiS2dSAZFucjScSuUhS9yQNEtHT2es4Fq"

	client := daowallet.NewDefaultClient(key, secret)

	id := "eif0Z2bfnkY6WU5mg7gIqTUQBgDs5zWI"

	inv, err := client.InvoiceStatus(context.Background(), id)
	if err != nil {
		panic(err)
	}

	fmt.Println("invoice: %v", inv)
}
```
