import unittest

from new_gateway import core


class CoreTest(unittest.TestCase):
    def test_something(self):
        add = core.add(1, 1)
        self.assertEqual(2, add)


if __name__ == '__main__':
    unittest.main()
