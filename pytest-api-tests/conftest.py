# tests/conftest.py
import pytest
import requests
import logging

# Constants from env.sh
BASE_URL = "http://localhost:8080"
TOKEN = "your-auth-token"

logging.basicConfig(level=logging.DEBUG)
logging.getLogger('urllib3').setLevel(logging.DEBUG)

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
def user_id(auth_headers, base_url):
    """Fixture to create a user and return the user ID."""
    # Create a new user
    create_url = f"{base_url}/api/users"
    user_data = {
        "email": "lluser@4example.com"
    }
    create_response = requests.post(create_url, json=user_data, headers=auth_headers)
    assert create_response.status_code == 201
    created_user = create_response.json()
    created_email = created_user.get('email')

    # Get the list of users and find the user ID
    get_users_url = f"{base_url}/api/users"
    get_users_response = requests.get(get_users_url)
    assert get_users_response.status_code == 200
    users = get_users_response.json()
    user_id = None
    for user in users:
        if user.get('email') == created_email:
            user_id = user.get('id')
            break
    assert user_id is not None, "User ID not found in the users list."

    yield user_id  # Provide the user_id to tests

    # Teardown: Delete the user after all tests are done
    #delete_url = f"{base_url}/api/users/{user_id}"
    #delete_response = requests.delete(delete_url, headers=auth_headers)
    #assert delete_response.status_code == 204

