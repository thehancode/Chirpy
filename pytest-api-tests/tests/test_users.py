# tests/test_users.py
import pytest
import requests

def test_get_users(client, base_url):
    """Test to retrieve all users."""
    url = f"{base_url}/api/users"
    response = client.get(url)
    assert response.status_code == 200, f"Expected status code 200, got {response.status_code}"
    users = response.json()
    assert isinstance(users, list), "Response is not a list of users"

