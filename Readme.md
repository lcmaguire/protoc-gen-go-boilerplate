# protoc-gen-go-boilerplate

plugin to generate go boilerplate code with the ability to support custom templates.

## installation

```sh

go install github.com/lcmaguire/protoc-gen-go-boilerplate@latest

```

## usage

Heavily recommended to use buf cli v2

## gRPC go gen

currently supports generating boilerplate code for the following.

go gRPC

|                | unary | streaming | streaming |
|----------------|-------|-----------|-----------| 
| method gen     | ✅     | ✅         | ✅         |
| service struct | ✅     | ✅         | ✅         |
| server         | 🚧    | 🚧        | 🚧        |

## connect rpc go gen 🚧

|                | unary | streaming | streaming |
|----------------|-------|-----------|-----------| 
| method gen     | ✅     | ✅         | ✅         |
| service struct | ✅     | ✅         | ✅         |
| server         | 🚧    | 🚧        | 🚧        |

## 🚧🚧🚧 In progress 🚧🚧🚧

- templates for generating message related functions
- server generation
- connect rpc support

## Potential future features

- dockerfile generation
- test generation