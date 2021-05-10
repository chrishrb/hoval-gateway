periodic_requests = []

PERIODIC_TIME = 30
PRIORITY = 208


class PeriodicRequest:
    """Periodic request"""
    def __init__(self, device_type, device_id, datapoint_name, periodic_time, priority):
        self.device_type = device_type
        self.device_id = device_id
        self.datapoint_name = datapoint_name
        self.periodic_time = periodic_time
        self.priority = priority

    def __str__(self):
        return str.format("Periodic request with device type: {}, device id: {}, datapoint_name: {}, periodic_time: {}"
                          ", priority: {}", self.device_type, self.device_id, self.datapoint_name, self.periodic_time,
                          self.priority)


def parse_periodic_requests(element):
    for periodic_request in element:
        for datapoint_element in periodic_request["datapoints"]:
            # add periodic request to queue
            periodic_requests.append(
                PeriodicRequest(
                    periodic_request["device_type"],
                    periodic_request["device_id"],
                    datapoint_element,
                    periodic_request["periodic_time"] if "periodic_time" in periodic_request else PERIODIC_TIME,
                    periodic_request["priority"] if "priority" in periodic_request else PRIORITY
                )
            )
