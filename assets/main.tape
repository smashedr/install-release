
Type "ir smashedr/bup"
Sleep 0.5s
Enter

# select version
Wait / [0-9]\.[0-9]/
Sleep 1.5s
Down
Sleep 0.2s
Down
Sleep 0.5s
Up
Sleep 0.2s
Up
Sleep 0.5s
Enter
# select asset
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
# set name
Wait /\> /
Sleep 1s
Backspace
Sleep 0.2s
Backspace
Sleep 0.2s
Type u
Sleep 0.2s
Type p
Sleep 0.5s
Enter
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

Type "ir rm bup"
Sleep 0.5s
Enter
Wait /Remove\!/
Sleep 1.5s
Right
Sleep 0.4s
Left
Sleep 0.4s
Enter
Sleep 1.5s
