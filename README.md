# Bore

> [!WARNING]
> This project is still under development, use at your own risk, this is mostly functional enough for my use at work currently.

A clipboard implementation for headless (and non-headless) environments

# Install

You can download binaries for your platform from the [releases page](https://github.com/aosasona/bore/releases/latest) or install directly with the Go toolchain:

```sh
go install go.trulyao.dev/bore/cmd/bore@latest
```

# Usage

You can use `bore` the same way you would use `pbcopy` on Mac or `xclip` on Linux. For example, the command below will read from a file using `cat` and copy the result with `bore`.

```sh
cat path/to/file | bore copy
```

- Piping directly into bore via echo:

```sh
echo 'Hello world' | bore copy
```

> [!NOTE]
> If you run `bore copy` directly, it will prompt you to type in your text

- Pasting in a file

```sh
bore paste > /path/to/file
```

# Why...

### Why does this exist?

At work, I work on a server where I don't have a system clipboard or install permissions to add one, and asking my manager to install other clipboard tools would seemingly pull in other dependencies we don't need. I need a clipboard so I can copy across tabs and sessions in Neovim (while in and not in Zellij), so I made one (you can find the Neovim plugin I use [here](https://github.com/aosasona/bore.nvim))

### Why does it use a database?

For efficient search, querying, syncing etc (other future features) and... SQLite is just cool. It is also just easier than trying to re-implement a database on a plain text or JSON file anyway (which will happen eventually when you set out to make something like this). Also, it is your data, you can just grab that database file and query it as you wish from other applications.
