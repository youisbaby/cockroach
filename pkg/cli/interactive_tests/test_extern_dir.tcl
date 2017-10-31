#! /usr/bin/env expect -f
#
source [file join [file dirname $argv0] common.tcl]

set storedir "mystore"
set externdir "/some/extern"

spawn /bin/bash
send "PS1=':''/# '\r"
eexpect ":/# "

start_test "Check that non-absolute external-io-dir rejected"

send "$argv start --insecure --store=$storedir --external-io-dir=blah\r"
eexpect "external-io-dir path must be absolute"

end_test

start_test "Check disabling external IO explicitly"

send "$argv start --insecure --store=$storedir --external-io-dir=disabled\r"
eexpect "external I/O path:   <disabled>"
interrupt
eexpect "shutdown completed"

end_test

start_test "Check setting external IO explicitly"

send "$argv start --insecure --store=$storedir --external-io-dir=$externdir\r"
eexpect "external I/O path:   $externdir"
interrupt
eexpect "shutdown completed"

end_test

start_test "Check implicit external I/O dir under store dir"

send "$argv start --insecure --store=$storedir\r"
eexpect "external I/O path:   /mystore/extern"
interrupt
eexpect "shutdown completed"

end_test

