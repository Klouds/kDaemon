cd /kdaemon
echo "[default]" > /kdaemon/config/app.conf
echo "bind_ip = $BIND_IP" >> /kdaemon/config/app.conf
echo "api_port = $API_PORT" >> /kdaemon/config/app.conf
echo "ui_port = $UI_PORT" >> /kdaemon/config/app.conf
echo "rethinkdb_host = $RETHINKDB_HOST" >> /kdaemon/config/app.conf
echo "rethinkdb_port = $RETHINKDB_PORT" >> /kdaemon/config/app.conf
echo "rethinkdb_dbname = $RETHINKDB_DBNAME" >> /kdaemon/config/app.conf

kdaemon
