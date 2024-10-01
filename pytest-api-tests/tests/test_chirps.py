# tests/test_chirps.py
import pytest
import requests
import uuid

def test_create_chirp(auth_headers, base_url, user_id):
    """Test to create a chirp using the retrieved user ID."""
    url = f"{base_url}/api/chirps"
    chirp_data = {
        "body": "Hello2, world!",
        "user_id": user_id
    }
    response = requests.post(url, json=chirp_data, headers=auth_headers)
    assert response.status_code == 201
    chirp = response.json()
    assert 'id' in chirp

def test_get_all_chirps(base_url):
    """Test to retrieve all chirps."""
    url = f"{base_url}/api/chirps"
    response = requests.get(url)
    assert response.status_code == 200
    chirps = response.json()
    assert isinstance(chirps, list)


def test_get_chirp_by_id_success(auth_headers, base_url, user_id):
    """Test retrieving a single chirp by its ID successfully."""
    # First, create a chirp to ensure it exists
    create_url = f"{base_url}/api/chirps"
    chirp_data = {
        "body": "Test chirp for retrieval",
        "user_id": user_id
    }
    create_response = requests.post(create_url, json=chirp_data, headers=auth_headers)
    assert create_response.status_code == 201, f"Failed to create chirp, status code {create_response.status_code}"
    created_chirp = create_response.json()
    chirp_id = created_chirp['id']

    # Now, retrieve the chirp by ID
    get_url = f"{base_url}/api/chirps/{chirp_id}"
    get_response = requests.get(get_url, headers=auth_headers)
    assert get_response.status_code == 200, f"Expected status code 200, got {get_response.status_code}"
    chirp = get_response.json()
    assert chirp['id'] == chirp_id, "Chirp ID does not match"
    assert chirp['body'] == "Test chirp for retrieval", "Chirp body does not match"
    assert chirp['user_id'] == user_id, "Chirp user_id does not match"
    # Optionally, verify timestamps if applicable
    assert 'created_at' in chirp, "Chirp does not contain 'created_at'"
    assert 'updated_at' in chirp, "Chirp does not contain 'updated_at'"

def test_get_chirp_by_id_not_found(auth_headers, base_url):
    """Test retrieving a chirp with a non-existent ID returns 404."""
    non_existent_id = str(uuid.uuid4())  # Generates a random UUID
    get_url = f"{base_url}/api/chirps/{non_existent_id}"
    get_response = requests.get(get_url, headers=auth_headers)
    assert get_response.status_code == 404, f"Expected status code 404, got {get_response.status_code}"
    # Optionally, verify the error message
    error_response = get_response.json()
    assert "error" in error_response, "Error response does not contain 'error' key"
    assert error_response["error"] == "Chirp not found", "Error message does not match expected"

def test_get_chirp_by_id_invalid_uuid(auth_headers, base_url):
    """Test retrieving a chirp with an invalid UUID format returns 400 or appropriate error."""
    invalid_id = "invalid-uuid-format"
    get_url = f"{base_url}/api/chirps/{invalid_id}"
    get_response = requests.get(get_url, headers=auth_headers)
    assert get_response.status_code in [400, 422], f"Expected status code 400 or 422, got {get_response.status_code}"
    # Optionally, verify the error message
    error_response = get_response.json()
    assert "error" in error_response, "Error response does not contain 'error' key"
    assert "invalid UUID" .lower() in error_response["error"].lower(), "Error message does not indicate invalid UUID"
