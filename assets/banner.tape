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

Env TERM "xterm-256color"
Env COLORTERM "truecolor"

Hide
Type "ir --help | head -n 16 && echo '  ...'"
Enter
Show
Sleep 1.5s

Type "clear"
Sleep 0.3s
Enter

Source assets/main.tape

Type "clear"
Sleep 0.3s
Enter
