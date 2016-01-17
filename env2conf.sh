cd kdaemon
echo "[default]" > /kdaemon/config/app.conf
echo "bind_ip = $BIND_IP" >> /kdaemon/config/app.conf
echo "api_port = $API_PORT" >> /kdaemon/config/app.conf
echo "ui_port = $UI_PORT" >> /kdaemon/config/app.conf
echo "mysql_host = $MYSQL_HOST" >> /kdaemon/config/app.conf
echo "mysql_user = $MYSQL_USER" >> /kdaemon/config/app.conf
echo "mysql_password = $MYSQL_PASSWORD" >> /kdaemon/config/app.conf
echo "mysql_dbname = $MYSQL_DBNAME" >> /kdaemon/config/app.conf
chmod a+x /kdaemon/kdaemon
./kdaemon/kdaemon
