s#!/bin/sh

# configuration variables
SERVER_USERNAME="control-server"
LOCAL_MQTT_CONFIG_FILE="./config/mosquitto.conf"
LOCAL_CONFIG_FILE="config.toml"

MQTT_CONFIG_FILE="/etc/mosquitto/conf.d/custom-mosquitto.conf"
MQTT_PASSWD_FILE="/etc/mosquitto/passwd"

CERT_BITS=2048
CERT_DAYS=365
CN_SUBJ="/CN=control-server-broker"

# execute with root privileges
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# install mosquitto via snap package
apt install mosquitto -y

# copy config file to a directory where it is readable for mosquitto
cp $LOCAL_MQTT_CONFIG_FILE $MQTT_CONFIG_FILE

# set username/password, create random password first
echo "Creating user and password"
PASSWD=$(head -c30 /dev/urandom | base64)
cp config.toml.example $LOCAL_CONFIG_FILE
sed -i "s/password = \"\"/$PASSWD/g" $LOCAL_CONFIG_FILE
echo "" > $MQTT_PASSWD_FILE
mosquitto_passwd -b $MQTT_PASSWD_FILE $SERVER_USERNAME $PASSWD	# create passwd file
PASSWD=""

# generate CA certificate and server key (creates ca.key, ca.crt, server.key, server.crt)
openssl req -new -x509 -days $CERT_DAYS -extensions v3_ca -keyout ca.key -out ca.crt
openssl genrsa -out server.key $CERT_BITS							# generete server key
openssl req -out server.csr -key server.key -new -subj $CN_SUBJ 	# generate a certificate signing request to send to the CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days $CERT_DAYS # sign server key
rm server.csr
cp server.crt certs/
mv ca.* /etc/mosquitto/certs # copy created certificates
mv server.* /etc/mosquitto/certs

# restart service
echo "Restarting mosquitto service"
systemctl restart mosquitto