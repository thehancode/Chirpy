# tests/conftest.py
import pytest
import httpx
import logging

# Constants from env.sh
BASE_URL = "http://localhost:8080"
TOKEN = "your-auth-token"

logging.basicConfig(level=logging.INFO)
logging.getLogger('httpx').setLevel(logging.DEBUG)

def log_response(response):
    # Read the response content to ensure it's available
    content = response.read()
    # Decode the content to a string, handling any decoding errors
    try:
        content_str = content.decode('utf-8')
    except UnicodeDecodeError:
        content_str = repr(content)
    # Log the response body
    logger = logging.getLogger('httpx.response')
    logger.info('Response body: %s', content_str)


@pytest.fixture(scope='session')
def client():
    return httpx.Client(event_hooks={'response': [log_response]})
    


@pytest.fixture(scope='session')
def base_url():
    """Fixture to provide the base URL."""
    return BASE_URL

@pytest.fixture(scope='session')
def auth_headers():
    """Fixture to provide authorization headers."""
    return {
        'Authorization': f'Bearer {TOKEN}',
        'Content-Type': 'application/json'
    }

@pytest.fixture(scope='session')
def user(client, auth_headers, base_url):
    """Fixture to create a user and return the user object."""
    # Create a new user
    create_url = f"{base_url}/api/users"
    user_data = {
        "email": "lluser@4example.com",
        "password": "04234"
    }
    create_response = client.post(create_url, json=user_data, headers=auth_headers)
    assert create_response.status_code == 201
    created_user = create_response.json()
    assert 'id' in created_user, "User ID not found in the users list."
    yield created_user  # Provide the user object to tests

    url = f"{base_url}/admin/reset"
    response = client.post(url)
    assert response.status_code == 200  
