Output dist/banner.gif
#Output dist/banner.mp4
#Output dist/banner.webm

Require ir

Set Theme "Catppuccin Frappe"
#Set MarginFill "#6B50FF"
#Set Margin 10
#Set WindowBar Colorful
#Set BorderRadius 10
Set Padding 84
Set Width 1280
Set Height 640
Set FontSize 24


Hide
Type "ir --help | head -n 16 && echo '  ...'"
Enter
Show
Sleep 1.5s

Type "ir smashedr/bup"
Sleep 0.5s
Enter
Wait
Sleep 2s

Type "which bup"
Sleep 0.5s
Enter
Wait
Sleep 1s

Type "bup -V"
Sleep 0.5s
Enter
Wait
Sleep 1s

Type "ir list"
Sleep 0.5s
Enter
Wait
Sleep 2s

Type "clear"
Sleep 0.3s
Enter
