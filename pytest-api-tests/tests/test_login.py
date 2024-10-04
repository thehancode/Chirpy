# tests/test_login.py
import pytest
import uuid

def test_login_success(client, base_url, user):
    """test successful login with correct email and password."""
    url = f"{base_url}/api/login"
    login_data = {
        "email": "lluser@4example.com",
        "password": "04234"
    }
    response = client.post(url, json=login_data)
    assert response.status_code == 200, f"expected 200 ok, got {response.status_code}"
    user = response.json()
    assert 'id' in user, "response does not contain 'id'"
    assert user['email'] == "lluser@4example.com", "email does not match"
    assert 'password' not in user, "response should not contain 'password'"

def test_login_failure_wrong_password(client, base_url, user):
    """test login failure with incorrect password."""
    url = f"{base_url}/api/login"
    login_data = {
        "email": "lluser@4example.com",
        "password": "wrongpassword"
    }
    response = client.post(url, json=login_data)
    assert response.status_code == 401, f"expected 401 unauthorized, got {response.status_code}"
    error_message = response.json().get('error')
    assert error_message.lower() == "incorrect email or password", f"unexpected error message: {error_message}"

def test_login_failure_nonexistent_email(client, base_url):
    """test login failure with non-existent email."""
    url = f"{base_url}/api/login" 
    login_data = {
        "email": "nonexistent@example.com",
        "password": "somepassword"
    }
    response = client.post(url, json=login_data)
    assert response.status_code == 401, f"Expected 401 Unauthorized, got {response.status_code}"
    error_message = response.json().get('error')
    assert error_message == "Incorrect email or password", f"Unexpected error message: {error_message}"

