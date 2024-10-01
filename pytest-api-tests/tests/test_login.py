# tests/test_login.py
import pytest
import requests
import uuid

def test_login_success(base_url, create_test_user):
    """Test successful login with correct email and password."""
    url = f"{base_url}/api/login"
    login_data = {
        "email": "lanes@example.com",
        "password": "04234"
    }
    response = requests.post(url, json=login_data)
    assert response.status_code == 200, f"Expected 200 OK, got {response.status_code}"
    user = response.json()
    assert 'id' in user, "Response does not contain 'id'"
    assert user['email'] == "lanes@example.com", "Email does not match"
    assert 'password' not in user, "Response should not contain 'password'"

def test_login_failure_wrong_password(base_url, create_test_user):
    """Test login failure with incorrect password."""
    url = f"{base_url}/api/login"
    login_data = {
        "email": "lanes@example.com",
        "password": "wrongpassword"
    }
    response = requests.post(url, json=login_data)
    assert response.status_code == 401, f"Expected 401 Unauthorized, got {response.status_code}"
    error_message = response.json().get('error')
    assert error_message == "Incorrect email or password", f"Unexpected error message: {error_message}"

def test_login_failure_nonexistent_email(base_url):
    """Test login failure with non-existent email."""
    url = f"{base_url}/api/login"
    login_data = {
        "email": "nonexistent@example.com",
        "password": "somepassword"
    }
    response = requests.post(url, json=login_data)
    assert response.status_code == 401, f"Expected 401 Unauthorized, got {response.status_code}"
    error_message = response.json().get('error')
    assert error_message == "Incorrect email or password", f"Unexpected error message: {error_message}"

