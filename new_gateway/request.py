periodic_requests = []


class PeriodicRequest:
    def __init__(self, device_type, device_id, datapoint_name, periodic_time=30, prio=208):
        self.device_type = device_type
        self.device_id = device_id
        self.datapoint_name = datapoint_name
        self.periodic_time = periodic_time
        self.prio = prio

    def __str__(self):
        return str.format("Periodic request with device type: {}, device id: {}, datapoint_name: {}, periodic_time: {}"
                          ", prio: {}", self.device_type, self.device_id, self.datapoint_name, self.periodic_time,
                          self.prio)


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
                    periodic_request["prio"] if "prio" in periodic_request else PRIO
                )
            )
