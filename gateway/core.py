import logging

from gateway import datapoint
from gateway.exceptions import UnknownDatatypeError, NoValidMessageException, NoDatapointFoundError
from gateway.message import get_message_header, get_message_len, get_message_id, \
    get_operation_id, Operation, build_arbitration_id, ReceiveMessage, SendMessage
from gateway.mqtt import Subscriber
from gateway.request import periodic_requests, subscribe_requests

_pending_msg = {}


async def read(can0, mqtt_client, topic):
    """Read data from CAN and export to mqtt"""
    # no mqtt client
    if not mqtt_client:
        logging.error("Publishing not possible without mqtt client")

    while True:
        try:
            msg = await can0.get_message()
        except EOFError:
            break

        # Add message to pending
        add_to_pending_msg(msg)

        # When no more messages are expected, parse data
        if get_message_header(msg.data) in _pending_msg and \
                _pending_msg[get_message_header(msg.data)].nb_remaining == 0:
            process_message = _pending_msg[get_message_header(msg.data)]
            logging.debug("Current message: %s", process_message)

            try:
                dp, operation, parsed_message = process_message.parse()
                logging.debug("Message for datapoint [function_group: %s, function_number: %s, function_datapoint: "
                              "%s, operation: %s] : %s",
                              dp.function_group,
                              dp.function_number,
                              dp.datapoint_id,
                              operation,
                              parsed_message)

                # Publish to mqtt broker
                if operation == Operation.RESPONSE and mqtt_client:
                    topic_message = "{}/{}/status".format(topic, dp.name)

                    logging.info("Publish message to mqtt server: [{} {}]".format(topic_message, parsed_message))
                    mqtt_client.publish(topic_message, parsed_message)
            except (UnknownDatatypeError, NoDatapointFoundError, NoValidMessageException) as e:
                logging.error(e)

            del _pending_msg[get_message_header(msg.data)]


def add_to_pending_msg(msg):
    """Add message to queue"""
    if get_message_id(msg.arbitration_id) == 0x1f:
        parsed_msg = ReceiveMessage(msg.arbitration_id, get_operation_id(msg.data), get_message_len(msg.data))
        parsed_msg.put_data(msg.data)
        _pending_msg[get_message_header(msg.data)] = parsed_msg
    elif get_message_header(msg.data) in _pending_msg:
        _pending_msg[get_message_header(msg.data)].put_extended_data(msg.data)


async def send(can0, mqtt_client, topic):
    """Retrieve messages from mqtt broker and send to can interface"""
    if mqtt_client is None:
        logging.error("Sending not possible without mqtt client")
        return

    # subscribe to topic
    for key, request in subscribe_requests.items():
        try:
            datapoint_of_message = datapoint.get_datapoint_by_name(key)
        except NoDatapointFoundError as e:
            logging.error(e)
            continue

        logging.debug("Subscribe to %s/%s/set", topic, datapoint_of_message.name)
        mqtt_client.subscribe("{}/{}/set".format(topic, datapoint_of_message.name))

    # Handle received messages from mqtt broker
    subscriber = Subscriber(can0, topic)
    mqtt_client.on_message = subscriber.on_message
    mqtt_client.loop_start()


async def send_periodic(can0):
    """Send periodic requests to can interface"""
    _payload = 0
    _message_id = 31
    _operation = Operation.GET_REQUEST

    for key, request in periodic_requests.items():
        logging.debug(request)

        # build message
        try:
            datapoint_of_message = datapoint.get_datapoint_by_name(request.datapoint_name)
        except NoDatapointFoundError as e:
            logging.error(e)
            continue

        _message_id += 1
        arbitration_id = build_arbitration_id(_message_id, request.priority, request.device_type,
                                              request.device_id)
        message = SendMessage(arbitration_id, _operation.value, datapoint_of_message)

        # set data to 0
        message.put_single_data(_payload)

        # to can message
        try:
            can_message = message.to_can_message()
            can0.send_periodic(can_message, request.periodic_time)
            logging.debug("Send message: %s", message)
        except (ValueError, NoValidMessageException) as e:
            logging.error(e)
