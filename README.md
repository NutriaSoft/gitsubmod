### Requirements

```bash
brew install protobuf
```

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

#### Init submodule based on .gitmodules
```bash
subop -init
```

#### Init submodule based on choice and edit .gitmodules
```bash
subop -choose
```

#### Init submodule based on choice and edit .gitmodules and nuke previous module
```bash
subop -choose -nuke
```

#### Delete git submodule
```bash
subop -nuke
```

```bash
subop -repo <repo-url> -branch <branch> -name <name> -init
```
