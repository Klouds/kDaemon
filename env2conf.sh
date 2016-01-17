echo "[default]" > $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
echo "bind_ip = $BIND_IP" >> $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
echo "api_port = $API_PORT" >> $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
echo "ui_port = $UI_PORT" >> $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
echo "mysql_host = $MYSQL_HOST" >> $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
echo "mysql_user = $MYSQL_USER" >> $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
echo "mysql_password = $MYSQL_PASSWORD" >> $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
echo "mysql_dbname = $MYSQL_DBNAME" >> $GOPATH/src/github.com/superordinate/kDaemon/config/app.conf
