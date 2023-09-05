import unittest
import requests
from requests.adapters import HTTPAdapter
from urllib3.util import Retry

LOCAL_TESTING_URL = "http://localhost:8080"

FACTORIALS = {
    0: 1,
    1: 1,
    2: 2,
    3: 6,
    4: 24,
    5: 120,
}


class TestHandler(unittest.TestCase):
    """Sample test cases to run against the handler."""

    @classmethod
    def setUpClass(cls):
        cls.session = requests.Session()
        retries = Retry(total=5, backoff_factor=1, status_forcelist=[502, 503, 504])
        cls.session.mount("http://", HTTPAdapter(max_retries=retries))
        # Required for node offline testing
        cls.session.headers["Content-Type"] = "text/plain"

    def test_factorial(self):
        for factorial, expected in FACTORIALS.items():
            with self.subTest(f"factorial({factorial})"):
                resp = self.session.post(
                    LOCAL_TESTING_URL,
                    data=str(factorial),
                )
                self.assertEqual(resp.status_code, 200)
                self.assertEqual(resp.text, str(expected))

    def test_405_on_get(self):
        resp = self.session.get(LOCAL_TESTING_URL)
        self.assertEqual(resp.status_code, 405)


if __name__ == "__main__":
    unittest.main()
