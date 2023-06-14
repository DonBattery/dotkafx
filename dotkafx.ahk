; spin up the DotkaFX Server when Ctrl + F1 is pressed
^F1::
Run, dotkafx
return

; starts the DotkaFX Scheduler when Ctrl + F2 is pressed
^F2::
Run, dotkafx start,, Min
return

; roll backward the DotkaFX Scheduler by 5 seconds when Ctrl + F3 is pressed
^F3::
Run, dotkafx back5,, Min
return

; roll forward the DotkaFX Scheduler by 5 seconds when Ctrl + F4 is pressed
^F4::
Run, dotkafx forward5,, Min
return

; pauses/unpause the DotkaFX Scheduler when Ctrl + F5 is pressed
^F5::
Run, dotkafx pause,, Min
return

; shut down the DotkaFX Server when Ctrl + F6 is pressed
^F6::
Run, dotkafx shutdown,, Min
return
