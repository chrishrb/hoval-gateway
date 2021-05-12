from unittest import TestCase
import new_gateway.datapoint as datapoint
from new_gateway.exceptions import NoDatapointFoundError


class TestDatapoint(TestCase):
    def test_get_datapoint_type(self):
        self.fail()


class Test(TestCase):
    def setUp(self):
        datapoint.datapoints_by_id = {}
        datapoint.datapoints_by_name = {}

        self.element = [{
            "name": "example_datapoint",
            "function_group": 1,
            "function_number": 2,
            "datapoint_id": 3,
            "type": "LIST"
        }]

    def test_parse_datapoints(self):
        datapoint.parse_datapoints(self.element)

        self.assertIn((1, 2, 3), datapoint.datapoints_by_id)
        self.assertIn("example_datapoint", datapoint.datapoints_by_name)

    def test_get_datapoint_by_name(self):
        datapoint.datapoints_by_name["example_name"] = 0
        name = datapoint.get_datapoint_by_name("example_name")

        self.assertEqual(0, name)
        self.assertRaises(NoDatapointFoundError, lambda: datapoint.get_datapoint_by_name("not_found",))

    def test_get_datapoint_by_id(self):
        datapoint.datapoints_by_id[(1, 2, 3)] = 0
        by_id = datapoint.get_datapoint_by_id(1, 2, 3)

        self.assertEqual(0, by_id)
        self.assertRaises(NoDatapointFoundError, lambda: datapoint.get_datapoint_by_id(2, 3, 7))
