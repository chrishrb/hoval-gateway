from unittest import TestCase

import new_gateway.datatypes as datatypes


class TestUnsigned(TestCase):
    number_8 = int.to_bytes(150, byteorder='big', signed=False, length=8)
    number_16 = int.to_bytes(150, byteorder='big', signed=False, length=16)

    def test_convert_from_bytes(self):
        datatype = datatypes.Unsigned(0)

        converted_number_8 = datatype.convert_from_bytes(self.number_8)
        converted_number_16 = datatype.convert_from_bytes(self.number_16)

        self.assertEqual(150, converted_number_8)
        self.assertEqual(150, converted_number_16)

    def test_convert_from_bytes_decimal(self):
        datatype = datatypes.Unsigned(1)

        converted_number_8 = datatype.convert_from_bytes(self.number_8)
        converted_number_16 = datatype.convert_from_bytes(self.number_16)

        self.assertEqual(15.0, converted_number_8)
        self.assertEqual(15.0, converted_number_16)

    def test_convert_to_bytes(self):
        # todo: add test
        pass


class TestSigned(TestCase):
    def test_convert_from_bytes(self):
        # todo: add test
        self.fail()

    def test_convert_to_bytes(self):
        # todo: add test
        pass


class TestList(TestCase):
    number_8 = int.to_bytes(5, byteorder='big', signed=False, length=8)

    def test_convert_from_bytes(self):
        datatype = datatypes.List()

        converted_number_8 = datatype.convert_from_bytes(self.number_8)

        self.assertEqual(5, converted_number_8)

    def test_convert_to_bytes(self):
        # todo: add test
        pass


class TestString(TestCase):
    string = b'\x68\x65\x6c\x6c\x6f\x20\x77\x6f\x72\x6c\x64'

    def test_convert_from_bytes(self):
        datatype = datatypes.String()

        converted_string = datatype.convert_from_bytes(self.string)

        self.assertEqual("hello world", converted_string)

    def test_convert_to_bytes(self):
        # todo: add test
        pass
