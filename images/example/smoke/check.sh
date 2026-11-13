set -euo pipefail
ready() { [ "$(valkey-cli -p 56379 ping 2> /dev/null)" = PONG ]; }
for _ in $(seq 60); do ready && break; sleep 1; done
ready
databasify ctl is-active --quiet server.service
[ "$(databasify ctl show -p NRestarts --value server)" = 0 ]
databasify ctl logs -u server -n 3
echo "example $DATABASIFY_RELEASE runs under the databasify supervisor"
