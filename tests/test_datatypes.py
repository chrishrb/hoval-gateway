from unittest import TestCase

import gateway.datatypes as datatypes
from gateway.exceptions import NoValidMessageException


class TestUnsigned(TestCase):
    def test_convert_from_bytes(self):
        datatype = datatypes.Unsigned(8)
        converted_number = datatype.convert_from_bytes(b'\x14')
        self.assertEqual(20, converted_number)

    def test_convert_from_bytes_decimal(self):
        datatype = datatypes.Unsigned(16, 2)
        converted_number = datatype.convert_from_bytes(b'\x07\xe4')
        self.assertEqual(20.20, converted_number)

    def test_convert_to_bytes(self):
        datatype = datatypes.Unsigned(16)
        converted_number = datatype.convert_to_bytes("20")
        self.assertEqual(b'\x00\x14', converted_number)

    def test_convert_to_bytes_decimal(self):
        datatype = datatypes.Unsigned(16, 2)
        converted_number = datatype.convert_to_bytes("20.20")
        self.assertEqual(b'\x07\xe4', converted_number)

    def test_convert_wrong_type(self):
        datatype = datatypes.Unsigned(16, 2)
        with self.assertRaises(NoValidMessageException):
            datatype.convert_to_bytes("example123")


class TestSigned(TestCase):
    def test_convert_from_bytes(self):
        datatype = datatypes.Signed(16)
        converted_number = datatype.convert_from_bytes(b'\xff\xec')
        self.assertEqual(-20, converted_number)

    def test_convert_from_bytes_decimal(self):
        datatype = datatypes.Signed(16, 2)
        converted_number = datatype.convert_from_bytes(b'\xf8\x19')
        self.assertEqual(-20.23, converted_number)

    def test_convert_to_bytes(self):
        datatype = datatypes.Signed(16)
        converted_number = datatype.convert_to_bytes(-20)
        self.assertEqual(b'\xff\xec', converted_number)

    def test_convert_to_bytes_decimal(self):
        datatype = datatypes.Signed(16, 2)
        converted_number = datatype.convert_to_bytes(-20.23)
        self.assertEqual(b'\xf8\x19', converted_number)

    def test_convert_wrong_type(self):
        datatype = datatypes.Signed(16, 2)
        with self.assertRaises(NoValidMessageException):
            datatype.convert_to_bytes("example123")


class TestList(TestCase):
    def test_convert_from_bytes(self):
        datatype = datatypes.List()
        converted_number = datatype.convert_from_bytes(b'\x05')
        self.assertEqual(5, converted_number)

    def test_convert_to_bytes(self):
        datatype = datatypes.List()
        converted_number = datatype.convert_to_bytes(5)
        self.assertEqual(b'\x05', converted_number)


class TestString(TestCase):
    string = 'hello world'
    byte_string = b'\x68\x65\x6c\x6c\x6f\x20\x77\x6f\x72\x6c\x64'

    def test_convert_from_bytes(self):
        datatype = datatypes.String()
        converted_string = datatype.convert_from_bytes(self.byte_string)
        self.assertEqual(self.string, converted_string)

    def test_convert_to_bytes(self):
        datatype = datatypes.String()
        converted_string = datatype.convert_to_bytes(self.string)
        self.assertEqual(self.byte_string, converted_string)
