import logging


def on_connect(client, userdata, flags, rc):
    if rc == 0:
        logging.debug("Connected to broker")
        client.subscribe("hoval-gw/lueftungs_modulation/set")
    else:
        logging.error("Connection failed")


def on_message(client, userdata, message):
    print("Message received : " + str(message.payload) + " on " + message.topic)

