=== ":simple-homebrew: brew"

    ```shell
    brew install cssnr/tap/install-release
    ```

=== ":lucide-terminal-square: bash"

    ```shell
    curl 'https://raw.githubusercontent.com/smashedr/install-release/master/scripts/install.sh' | bash  # (1)!
    ```

    1.  Alternatively, you can manually [download a release](https://github.com/smashedr/install-release/releases).

=== ":lucide-terminal: powershell"

    ```powershell
    iex (iwr -useb 'https://raw.githubusercontent.com/smashedr/install-release/master/scripts/install.ps1').Content  # (1)!
    ```

    1.  Windows users can download the [Windows Installer](https://github.com/smashedr/install-release/releases/latest/download/ir_Windows_Installer.exe).

=== ":simple-go: source"

    ```shell
    go install github.com/smashedr/install-release@latest  # (1)!
    ```

    1. Requires Go: <https://go.dev/doc/install>

=== ":simple-docker: docker"

    ```shell
    docker run --rm -itv ~/bin:/out ghcr.io/smashedr/ir:latest -b /out smashedr/install-release  # (1)!
    ```

    1. _Note: Docker requires you to mount the target bin directory._

:fontawesome-brands-windows: Windows users can download the [Windows&nbsp;Installer](https://github.com/smashedr/install-release/releases/latest/download/ir_Windows_Installer.exe).  
:lucide-download: Alternatively, you can manually [download&nbsp;a&nbsp;release](https://github.com/smashedr/install-release/releases).
