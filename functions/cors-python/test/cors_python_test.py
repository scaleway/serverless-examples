import requests

OFFLINE_TESTING_SERVER="http://localhost:8080"

def test_handler_offline():
    response = requests.get(OFFLINE_TESTING_SERVER)
    assert response.text == "This is checking CORS"