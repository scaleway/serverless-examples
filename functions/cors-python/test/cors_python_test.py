import requests

OFFLINE_TESTING_SERVER="http://localhost:8080"

def test_handler_offline():
    response = requests.get(OFFLINE_TESTING_SERVER)
    assert response.headers["Access-Control-Allow-Origin"] == "*"
    assert response.headers["Access-Control-Allow-Headers"] == "*"
    assert response.headers["Access-Control-Allow-Methods"] == "*"
    assert response.headers["Content-Type"] == "text/plain"
    assert response.text == "This is checking CORS"
