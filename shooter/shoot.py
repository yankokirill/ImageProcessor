import requests
import json
import base64
import sys
import uuid
import time

BASE_URL = "http://localhost:8000"

def get_headers():
    username = f'user_{uuid.uuid4()}'
    password = 'password228'
    user_data = {'username': username, 'password': password}

    register_url = f"{BASE_URL}/register"
    login_url = f"{BASE_URL}/login"

    requests.post(register_url, json=user_data)
    response = requests.post(login_url, json=user_data)

    data = response.json()
    auth_token = data['token']
    return {'Content-Type': 'application/json',
            'Authorization': f'Bearer {auth_token}'}


def get_image():
    with open('../tests/static/sigma.png', 'rb') as image_file:
        return base64.b64encode(image_file.read()).decode('utf-8')


if len(sys.argv) < 2 or len(sys.argv) > 3:
    print(f"Usage: python {sys.argv[0]} FilterName Parameters")
    print("Parameters should be in JSON format")
    sys.exit(1)

filter_name = sys.argv[1]
payload = {
    "filter": {"name": filter_name},
    "image": get_image()
}

if len(sys.argv) == 3:
    parameters = json.loads(sys.argv[2])
    payload["filter"]["parameters"] = parameters

headers = get_headers()

response = requests.post(f"{BASE_URL}/task", headers=headers, json=payload)
data = response.json()

task_id = data['task_id']
status_url = f"{BASE_URL}/status/{task_id}"
result_url = f"{BASE_URL}/result/{task_id}"
while True:
    data = requests.get(status_url, headers=headers).json()
    if data['status'] == 'ready':
        break
    time.sleep(0.1)

data = requests.get(result_url, headers=headers).json()
with open(f'results/{filter_name}Sigma.png', 'wb') as file:
    file.write(base64.b64decode(data['result']))