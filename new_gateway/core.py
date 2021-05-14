import logging

from new_gateway import datapoint
from new_gateway.exceptions import UnknownDatatypeError, NoValidMessageException, NoDatapointFoundError
from new_gateway.message import get_message_header, get_message_len, get_message_id, \
    get_operation_id, Operation, build_arbitration_id, ReceiveMessage, SendMessage
from new_gateway.request import periodic_requests

_pending_msg = {}


async def read(can0):
    """Read data from CAN and export to mqtt"""
    can0.open()

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

                if operation == Operation.RESPONSE:
                    # todo: add mqtt handler
                    pass
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


async def send(can0):
    # todo: send can message
    pass


async def send_periodic(can0):
    payload = 0
    operation = Operation.GET_REQUEST

    for request in periodic_requests:
        logging.debug(request)

        # build message
        arbitration_id = build_arbitration_id(request.priority, request.device_type, request.device_id)
        datapoint_of_message = datapoint.get_datapoint_by_name(request.datapoint_name)
        message = SendMessage(arbitration_id, operation.value, datapoint_of_message)

        # set data to 0
        message.put_data([payload])
        logging.debug("Send message: %s", message)

        # to can message
        try:
            can_message = message.to_can_message()
            logging.debug("Build message: %s", can_message)
            can0.send_periodic(can_message, request.periodic_time)
        except NoValidMessageException as e:
            logging.error(e)
