=== "brew"

    ```shell
    brew install smashedr/test/install-release
    ```

=== "github"

    ```shell
    curl https://i.jpillora.com/smashedr/install-release! | bash  # (1)!
    ```

    1.  Windows users can download the [Windows Installer](https://github.com/smashedr/install-release/releases/latest/download/ir_Windows_Installer.exe).

        Alternatively, you can manually [download a release](https://github.com/smashedr/install-release/releases).

=== "source"

    ```shell
    go install github.com/smashedr/install-release@latest
    ```

=== "docker"

    ```shell
    docker run --rm -v ~/bin:/out ghcr.io/smashedr/ir:latest --bin=/out owner/repo  # (1)!
    ```

    1. _Note: Docker requires you to mount your bin directory and specify the path._

Windows users can download the [Windows Installer](https://github.com/smashedr/install-release/releases/latest/download/ir_Windows_Installer.exe).  
Alternatively, you can manually [download a release](https://github.com/smashedr/install-release/releases).
