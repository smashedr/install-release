Output dist/demo.gif
#Output dist/demo.mp4
#Output dist/demo.webm

Require ir

Set Theme "MaterialDark"
Set MarginFill "#00ADD8"
Set Margin 8
Set WindowBar Colorful
Set BorderRadius 12
Set Padding 20
Set Width 1200
Set Height 720
Set FontSize 34

Env TERM "xterm-256color"
Env COLORTERM "truecolor"

Type "ir smashedr/bup"
Sleep 0.5s
Enter

Wait / bup_/
Sleep 1.5s
Up
Sleep 0.2s
Up
Sleep 0.5s
Down
Sleep 0.2s
Down
Sleep 0.5s
Enter
Wait
Sleep 2s

Type "which bup"
Sleep 0.5s
Enter
Wait
Sleep 1.5s

Type "bup -V"
Sleep 0.5s
Enter
Wait
Sleep 1.5s

Type "ir list"
Sleep 0.5s
Enter
Wait
Sleep 2s

Type "clear"
Sleep 0.3s
Enter
