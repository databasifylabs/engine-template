set -euo pipefail
ready() { [ "$(timeout 2 socat -u -T 1 TCP:127.0.0.1:25060 STDOUT 2> /dev/null)" = 'example ready' ]; }
for _ in $(seq 60); do ready && break; sleep 1; done
ready
[ "$(setpriv --reuid example --regid example --clear-groups sqlite3 /srv/example/data/smoke.db 'SELECT 1')" = 1 ]
databasify ctl is-active --quiet server.service
[ "$(databasify ctl show -p NRestarts --value server)" = 0 ]
databasify ctl logs -u server -n 3
echo "example $DATABASIFY_RELEASE runs under the databasify supervisor"
