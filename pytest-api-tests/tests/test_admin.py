# tests/test_admin.py
import pytest
import requests

def test_reset(client ,base_url):
    """Test to reset the system."""
    url = f"{base_url}/admin/reset"
    response = client.post(url)
    assert response.status_code == 200  # Adjust based on actual expected status code

