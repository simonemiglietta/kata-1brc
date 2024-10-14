# Go Shared Tools

This module is named `lvciot/shared` and contains all the resources shared between the several Go projects in this Kata.

In order to use this module must be installed inside the consumer module by the commands:

```shell
go mod edit -replace lvciot/shared=../go-shared
go mod tidy
```

Remember to use the right relative path
