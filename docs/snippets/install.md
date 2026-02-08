=== ":simple-homebrew: brew"

    ```shell
    brew install cssnr/tap/install-release
    ```

=== ":simple-go: source"

    ```shell
    curl https://i.jpillora.com/smashedr/install-release! | bash  # (1)!
    ```

    1.  Windows users can download the [Windows Installer](https://github.com/smashedr/install-release/releases/latest/download/ir_Windows_Installer.exe).

        Alternatively, you can manually [download a release](https://github.com/smashedr/install-release/releases).

=== ":simple-go: source"

    ```shell
    go install github.com/smashedr/install-release@latest
    ```

=== ":simple-docker: docker"

    ```shell
    docker run --rm -v ~/bin:/out ghcr.io/smashedr/ir:latest -b /out owner/repo  # (1)!
    ```

    1. _Note: Docker requires you to mount your bin directory and specify the path._

:fontawesome-brands-windows: Windows users can download the [Windows&nbsp;Installer](https://github.com/smashedr/install-release/releases/latest/download/ir_Windows_Installer.exe).  
:lucide-download: Alternatively, you can manually [download&nbsp;a&nbsp;release](https://github.com/smashedr/install-release/releases).
