---
icon: lucide/rocket
---

# :lucide-rocket: Get Started

[![Install Release](assets/images/logo.svg){ align=right width=96 }](https://github.com/smashedr/install-release?tab=readme-ov-file#readme)

[![GitHub Release Version](https://img.shields.io/github/v/release/smashedr/install-release?logo=github)](https://github.com/smashedr/install-release/releases)
[![GitHub Downloads](https://img.shields.io/github/downloads/smashedr/install-release/total?logo=rolldown&logoColor=white)](https://github.com/smashedr/install-release/releases/latest)
[![Image Size](https://badges.cssnr.com/ghcr/size/smashedr/install-release)](https://github.com/smashedr/install-release/pkgs/container/install-release)
[![Go Version](https://img.shields.io/github/go-mod/go-version/smashedr/install-release?logo=go&logoColor=white&label=go)](https://github.com/smashedr/install-release/blob/master/go.mod)
[![GitHub Last Commit](https://img.shields.io/github/last-commit/smashedr/install-release?logo=listenhub&label=updated)](https://github.com/smashedr/install-release/pulse)
[![GitHub Repo Size](https://img.shields.io/github/repo-size/smashedr/install-release?logo=buffer&label=repo%20size)](https://github.com/smashedr/install-release?tab=readme-ov-file#readme)
[![GitHub Top Language](https://img.shields.io/github/languages/top/smashedr/install-release?logo=devbox)](https://github.com/smashedr/install-release?tab=readme-ov-file#readme)
[![GitHub Contributors](https://img.shields.io/github/contributors-anon/smashedr/install-release?logo=southwestairlines)](https://github.com/smashedr/install-release/graphs/contributors)
[![GitHub Issues](https://img.shields.io/github/issues/smashedr/install-release?logo=codeforces&logoColor=white)](https://github.com/smashedr/install-release/issues)
[![GitHub Discussions](https://img.shields.io/github/discussions/smashedr/install-release?logo=theconversation&logoColor=white)](https://github.com/smashedr/install-release/discussions)
[![GitHub Forks](https://img.shields.io/github/forks/smashedr/install-release?style=flat&logo=forgejo&logoColor=white)](https://github.com/smashedr/install-release/forks)
[![GitHub Repo Stars](https://img.shields.io/github/stars/smashedr/install-release?style=flat&logo=gleam&logoColor=white)](https://github.com/smashedr/install-release/stargazers)
[![GitHub Org Stars](https://img.shields.io/github/stars/cssnr?style=flat&logo=apachespark&logoColor=white&label=org%20stars)](https://cssnr.github.io/)
[![Discord](https://img.shields.io/discord/899171661457293343?logo=discord&logoColor=white&label=discord&color=7289da)](https://discord.gg/wXy6m2X8wY)
[![Ko-fi](https://img.shields.io/badge/Ko--fi-72a5f2?logo=kofi&label=support)](https://ko-fi.com/cssnr)

CLI to Install a GitHub Release.

Easily Install GitHub Release binaries with Windows, Linux and macOS Support.

--8<-- "docs/snippets/install.md"

If you run into any issues or have any questions, [support](support.md) is available.

## :lucide-terminal: Demo

[![VHS Tape](https://cssnr.s3.amazonaws.com/install-release/demo.gif)](#install)

:lucide-videotape: This demo was generated with [charmbracelet/vhs](https://github.com/charmbracelet/vhs).

## :lucide-sparkles: Features

- Supports Windows
- Custom `bin` Path
- Automatic Release Detection
- Select Asset and Name Interactively
- Set Asset and Name Programmatically
- List and Remove Installed Apps

## :lucide-plane-takeoff: Install

--8<-- "docs/snippets/install.md"

[![Latest Release](https://img.shields.io/github/v/release/smashedr/install-release?style=for-the-badge&logo=github&label=latest%20version)](https://github.com/smashedr/install-release/releases/latest)

## :lucide-terminal-square: Usage

Install the latest release.

```shell
ir owner/repo
```

Install a specific version.

```shell
ir owner/repo v4.2.0
```

Skip the asset and name prompts.

```shell
ir owner/repo -y
```

Set the name and asset programmatically.

```shell
ir owner/repo -n name -a name_asset.zip
```

Install to a different bin directory.

```shell
ir owner/repo -b /usr/local/bin
```

Get package information.

```shell
ir info owner/repo
```

List installed apps.

```shell
ir list
```

Remove an app.

```shell
ir remove name
```

&nbsp;

!!! question

    If you need **help** getting started or run into any issues, [support](support.md) is available!
