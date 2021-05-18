from gateway.exceptions import NoRequestFoundError

periodic_requests = {}
subscribe_requests = {}

PERIODIC_TIME = 30
PRIORITY = 208


class Request:
    """Abstract Request"""

    def __init__(self, device_type, device_id, datapoint_name, priority):
        self.device_type = device_type
        self.device_id = device_id
        self.datapoint_name = datapoint_name
        self.priority = priority


class PeriodicRequest(Request):
    """Periodic request"""

    def __init__(self, device_type, device_id, datapoint_name, priority, periodic_time):
        super().__init__(device_type, device_id, datapoint_name, priority)
        self.periodic_time = periodic_time

    def __str__(self):
        return str.format("Periodic request with device type: {}, device id: {}, datapoint_name: {}, periodic_time: {}"
                          ", priority: {}", self.device_type, self.device_id, self.datapoint_name, self.periodic_time,
                          self.priority)


class SubscribeRequest(Request):
    """Subscribe to mqtt topic"""

    def __init__(self, device_type, device_id, datapoint_name, priority):
        super().__init__(device_type, device_id, datapoint_name, priority)

    def __str__(self):
        return str.format("Subscribe request with device type: {}, device id: {}, datapoint_name: {}, priority: {}",
                          self.device_type, self.device_id, self.datapoint_name, self.priority)


def parse_requests(element):
    for request in element:
        # add periodic request to queue
        for periodic_request in request["periodic"]:
            periodic_requests[periodic_request] = PeriodicRequest(
                request["device_type"],
                request["device_id"],
                periodic_request,
                request["priority"] if "priority" in request else PRIORITY,
                request["periodic_time"] if "periodic_time" in request else PERIODIC_TIME
            )
        # add subscribe request to queue
        for subscribe_request in request["subscribe"]:
            subscribe_requests[subscribe_request] = SubscribeRequest(
                request["device_type"],
                request["device_id"],
                subscribe_request,
                request["priority"] if "priority" in request else PRIORITY,
            )


def get_subscribe_request_by_name(name):
    """Get datapoint by name (used for writing messages"""
    if name not in subscribe_requests:
        raise NoRequestFoundError(str.format("No subscribe request found with name [{}]", name))

    return subscribe_requests[name]
