# scaffold
[![Build Status](https://travis-ci.org/izumin5210/scaffold.svg?branch=master)](https://travis-ci.org/izumin5210/scaffold)
[![Coverage Status](https://coveralls.io/repos/github/izumin5210/scaffold/badge.svg)](https://coveralls.io/github/izumin5210/scaffold)
[![MIT License](https://img.shields.io/github/license/izumin5210/scaffold.svg)][license]
[![Version](https://img.shields.io/github/release/izumin5210/scaffold.svg)](./releases)

[![https://gyazo.com/756d165b512a3d93c08e2094eedecd86](https://i.gyazo.com/756d165b512a3d93c08e2094eedecd86.gif)](https://gyazo.com/756d165b512a3d93c08e2094eedecd86)

## Usage
### Example

[`.scaffold`](./.scaffold)

```
$ tree .scaffold
.scaffold
├── command
│   ├── app
│   │   └── cmd
│   │       ├── {{name}}.go
│   │       └── {{name}}_test.go
│   └── meta.toml
└── usecase
    ├── app
    │   └── usecase
    │       ├── {{name}}.go
    │       └── {{name}}_test.go
    └── meta.toml

$ scaffold g command destroy
       exist  .
       exist  app/cmd
      create  app/cmd/destroy.go
      create  app/cmd/destroy_test.go

$ ls app/cmd/destroy*
app/cmd/destroy.go
app/cmd/destroy_test.go
```


### Available filters

- `toUpper`
- `toLower`
- `camelizer`
- `pascalize`
- `underscore`
- `dasherize`


## Development

```
# install dependencies
$ make deps

# build
$ make build

# run lint and test
$ make test
```

## License
Licensed under [MIT License][license].

[license]: ./LICENSE
