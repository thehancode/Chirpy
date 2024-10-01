# tests/test_users.py
import pytest
import requests

def test_get_users(base_url):
    """Test to retrieve all users."""
    url = f"{base_url}/api/users"
    response = requests.get(url)
    assert response.status_code == 200
    users = response.json()
    assert isinstance(users, list)

