# QOR example application

This is an example application to show and explain features of [QOR](http://getqor.com).

Chat Room: [![Join the chat at https://gitter.im/qor/qor](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/qor/qor?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Quick Started

### Go version: 1.8+

```shell
# Get example app
$ go get -u github.com/TheDoctor028/bazar

# Setup database
# Or if you prefer postgreSQL see the postgres_env file for its configuration.
$ mysql -uroot -p
mysql> CREATE DATABASE qor_example;

# Run Application
$ cd $GOPATH/src/github.com/TheDoctor028/bazar
$ go run main.go
```

### Generate sample data

```go
$ go run config/db/seeds/main.go config/db/seeds/seeds.go
```

### Run tests (Pending)

```
$ go test $(go list ./... | grep -v /vendor/ | grep  -v /db/)
```

## Admin Management Interface

[Qor Example admin configuration](https://github.com/TheDoctor028/bazar/blob/master/config/admin/admin.go)

Online Demo Website: [demo.getqor.com/admin](http://demo.getqor.com/admin)

## RESTful API

[Qor Example API configuration](https://github.com/TheDoctor028/bazar/blob/master/config/api/api.go)

Online Example APIs:

* Users: [http://demo.getqor.com/api/users.json](http://demo.getqor.com/api/users.json)
* User 1: [http://demo.getqor.com/api/users/1.json](http://demo.getqor.com/api/users/1.json)
* User 1's Orders [http://demo.getqor.com/api/users/1/orders.json](http://demo.getqor.com/api/users/1/orders.json)
* User 1's Order 1 [http://demo.getqor.com/api/users/1/orders/1.json](http://demo.getqor.com/api/users/1/orders/1.json)
* User 1's Orders 1's Items [http://demo.getqor.com/api/users/1/orders.json](http://demo.getqor.com/api/users/1/orders/1/items.json)
* Orders: [http://demo.getqor.com/api/orders.json](http://demo.getqor.com/api/orders.json)
* Products: [http://demo.getqor.com/api/products.json](http://demo.getqor.com/api/products.json)

## License

Released under the MIT License.

[@QORSDK](https://twitter.com/qorsdk)