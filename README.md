[![GitHub Release Version](https://img.shields.io/github/v/release/smashedr/install-release?logo=github)](https://github.com/smashedr/install-release/releases)
[![GitHub Downloads](https://img.shields.io/github/downloads/smashedr/install-release/total?logo=rolldown&logoColor=white)](https://github.com/smashedr/install-release/releases/latest)
[![Image Size](https://badges.cssnr.com/ghcr/size/smashedr/install-release)](https://github.com/smashedr/install-release/pkgs/container/install-release)
[![Go Version](https://img.shields.io/github/go-mod/go-version/smashedr/install-release?logo=go&logoColor=white&label=go)](https://github.com/smashedr/install-release/blob/master/go.mod)
[![Deployment Docs](https://img.shields.io/github/deployments/smashedr/install-release/docs?logo=materialformkdocs&logoColor=white&label=docs)](https://github.com/smashedr/install-release/deployments/docs)
[![Deployment Preview](https://img.shields.io/github/deployments/smashedr/install-release/preview?logo=materialformkdocs&logoColor=white&label=preview)](https://github.com/smashedr/install-release/deployments/preview)
[![Workflow Release](https://img.shields.io/github/actions/workflow/status/smashedr/install-release/release.yaml?logo=testcafe&logoColor=white&label=release)](https://github.com/smashedr/install-release/actions/workflows/release.yaml)
[![Workflow Lint](https://img.shields.io/github/actions/workflow/status/smashedr/install-release/lint.yaml?logo=testcafe&logoColor=white&label=lint)](https://github.com/smashedr/install-release/actions/workflows/lint.yaml)
[![GitHub Last Commit](https://img.shields.io/github/last-commit/smashedr/install-release?logo=listenhub&label=updated)](https://github.com/smashedr/install-release/pulse)
[![GitHub Repo Size](https://img.shields.io/github/repo-size/smashedr/install-release?logo=buffer&label=repo%20size)](https://github.com/smashedr/install-release?tab=readme-ov-file#readme)
[![GitHub Top Language](https://img.shields.io/github/languages/top/smashedr/install-release?logo=devbox)](https://github.com/smashedr/install-release?tab=readme-ov-file#readme)
[![GitHub Contributors](https://img.shields.io/github/contributors-anon/smashedr/install-release?logo=southwestairlines)](https://github.com/smashedr/install-release/graphs/contributors)
[![GitHub Issues](https://img.shields.io/github/issues/smashedr/install-release?logo=codeforces&logoColor=white)](https://github.com/smashedr/install-release/issues)
[![GitHub Discussions](https://img.shields.io/github/discussions/smashedr/install-release?logo=theconversation)](https://github.com/smashedr/install-release/discussions)
[![GitHub Forks](https://img.shields.io/github/forks/smashedr/install-release?style=flat&logo=forgejo&logoColor=white)](https://github.com/smashedr/install-release/forks)
[![GitHub Repo Stars](https://img.shields.io/github/stars/smashedr/install-release?style=flat&logo=gleam&logoColor=white)](https://github.com/smashedr/install-release/stargazers)
[![GitHub Org Stars](https://img.shields.io/github/stars/cssnr?style=flat&logo=apachespark&logoColor=white&label=org%20stars)](https://cssnr.github.io/)
[![Discord](https://img.shields.io/discord/899171661457293343?logo=discord&logoColor=white&label=discord&color=7289da)](https://discord.gg/wXy6m2X8wY)
[![Ko-fi](https://img.shields.io/badge/Ko--fi-72a5f2?logo=kofi&label=support)](https://ko-fi.com/cssnr)

# Install Release

<a title="Install Release" href="https://smashedr.github.io/install-release" target="_blank">
<img alt="Install Release" align="right" width="128" height="auto" src="https://raw.githubusercontent.com/smashedr/install-release/refs/heads/master/docs/assets/images/logo.svg"></a>

- [Install](#install)
- [Usage](#usage)
- [Development](#development)
- [Support](#Support)
- [Contributing](#contributing)

CLI to Install a GitHub Release.

Easily Install GitHub Release binaries with Windows support.

[![VHS Tape](https://cssnr.s3.amazonaws.com/install-release/demo.gif)](https://smashedr.github.io/install-release/)

> [!IMPORTANT]  
> This project is in development.
> It is functional but may have bugs.

## Install

#### Homebrew

```shell
brew install smashedr/test/install-release
```

#### GitHub

```shell
curl https://i.jpillora.com/smashedr/install-release! | bash
```

Windows users can download the [Windows Installer](https://github.com/smashedr/install-release/releases/latest/download/ir_Windows_Installer.exe).  
Alternatively, you can manually [download a release](https://github.com/smashedr/install-release/releases).

#### Source

```shell
go install github.com/smashedr/install-release@latest
```

#### Docker

```shell
docker run --rm -v ~/bin:/out ghcr.io/smashedr/ir:latest --bin=/out owner/repo
```

_Note: Docker requires you to mount your bin directory and specify the path._

[![View Documentation](https://img.shields.io/badge/view_documentation-blue?style=for-the-badge&logo=googledocs&logoColor=white)](https://smashedr.github.io/install-release/)

## Usage

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

List installed apps.

```shell
ir list
```

Remove an app.

```shell
ir remove name
```

[![View Documentation](https://img.shields.io/badge/view_documentation-blue?style=for-the-badge&logo=googledocs&logoColor=white)](https://smashedr.github.io/install-release/)

# Development

Go: <https://go.dev/doc/install>

```shell
go run main.go
```

Task: <https://taskfile.dev/docs/installation>

```shell
task build
task lint
```

Docs: <https://zensical.org/docs/get-started>

```shell
task docs
```

Inno Setup: <https://jrsoftware.org/isdl.php>

```shell
task pathmgr
task inno
```

# Support

If you run into any issues or need help getting started, please do one of the following:

- Report an Issue: <https://github.com/smashedr/install-release/issues>
- Q&A Discussion: <https://github.com/smashedr/install-release/discussions/categories/q-a>
- Request a Feature: <https://github.com/smashedr/install-release/issues/new?template=1-feature.yaml>
- Chat with us on Discord: <https://discord.gg/wXy6m2X8wY>

[![Features](https://img.shields.io/badge/features-brightgreen?style=for-the-badge&logo=rocket&logoColor=white)](https://github.com/smashedr/install-release/issues/new?template=1-feature.yaml)
[![Issues](https://img.shields.io/badge/issues-red?style=for-the-badge&logo=southwestairlines&logoColor=white)](https://github.com/smashedr/install-release/issues)
[![Discussions](https://img.shields.io/badge/discussions-blue?style=for-the-badge&logo=theconversation&logoColor=white)](https://github.com/smashedr/install-release/discussions)
[![Discord](https://img.shields.io/badge/discord-5865F2?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/wXy6m2X8wY)

# Contributing

If you would like to submit a PR, please review the [CONTRIBUTING.md](#contributing-ov-file).

Please consider making a donation to support the development of this project
and [additional](https://cssnr.com/) open source projects.

[![Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/cssnr)

For a full list of current projects visit: [https://cssnr.github.io/](https://cssnr.github.io/)
