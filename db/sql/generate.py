import csv
import datetime
import uuid
from faker import Faker
import bcrypt
import random

fake = Faker()

# Number of workers
rows = 1000

# Ваш пароль
password = 'test123'.encode('utf-8')

# Хеширование пароля, соль генерируется автоматически
hashed = bcrypt.hashpw(password, bcrypt.gensalt())

# Set to keep track of generated emails
used_emails = set()

# Load worker IDs from workers_data.csv
worker_ids = []
with open('workers_data.csv', mode="r") as worker_file:
    worker_reader = csv.reader(worker_file, delimiter=';')
    next(worker_reader)  # Skip header row
    for row in worker_reader:
        if row:  # Check if the row is not empty
            worker_ids.append(row[0])  # Assuming worker_id is in the first column

# Load user IDs from users_data.csv
user_ids = []
with open('users_data.csv', mode="r") as user_file:
    user_reader = csv.reader(user_file, delimiter=';')
    next(user_reader)  # Skip header row
    for row in user_reader:
        user_ids.append(row[0])  # Assuming user_id is in the first column

order_ids = []
with open('orders_data.csv', mode="r") as order_file:
    order_reader = csv.reader(order_file, delimiter=';')
    next(order_reader)  # Skip header row
    for row in order_reader:
        if row:  # Check if the row is not empty
            order_ids.append(row[0])  # Assuming worker_id is in the first column

# Load user IDs from users_data.csv
task_ids = []
with open('tasks_data.csv', mode="r") as task_file:
    task_reader = csv.reader(task_file, delimiter=';')
    next(task_reader)  # Skip header row
    for row in task_reader:
        task_ids.append(row[0])


def add_random_weeks(date):
    weeks_to_add = random.randint(1, 5)
    new_date = date + datetime.timedelta(weeks=weeks_to_add)
    return new_date


# Generate orders data and write to a CSV file
def generate_orders_data():
    with open("orders_data.csv", mode="w", newline="") as file:
        writer = csv.writer(file, delimiter=';')
        writer.writerow(["id", "worker_id", "user_id", "status", "deadline", "address", "creation_date", "rate"])

        for _ in range(rows):
            order_id = str(uuid.uuid4())  # Generate a UUID
            worker_id = fake.random_element(elements=worker_ids)
            user_id = fake.random_element(elements=user_ids)
            address = fake.address().replace("\n", " ")  # Replace line breaks with spaces
            creation_date = fake.date_this_year()
            deadline = add_random_weeks(creation_date)
            status = fake.random_element(elements=("3", "4")) if datetime.datetime.today().date() > deadline else fake.random_element(elements=("1", "2", "3", "4"))  # ожидает, выполняется, выполнено

            rate = fake.random_int(min=0, max=5) if status == "3" else 0

            writer.writerow([order_id, worker_id, user_id, status, deadline, address, creation_date, rate])


def generate_phone_number():
    return f"+7{''.join([str(random.randint(0, 9)) for _ in range(10)])}"


# Generate workers data and write to a CSV file
def generate_workers_data():
    with open("workers_data.csv", mode="w", newline="") as file:
        writer = csv.writer(file, delimiter=';')
        writer.writerow(["id", "name", "surname", "email", "phone_number", "address", "password", "role"])

        for _ in range(rows):
            worker_id = str(uuid.uuid4())  # Generate a UUID
            name = fake.first_name()
            surname = fake.last_name()

            # Generate a unique email (not previously used)
            while True:
                email = fake.email()
                if email not in used_emails:
                    used_emails.add(email)
                    break

            phone_number = generate_phone_number()
            address = fake.address().replace("\n", " ")
            password = hashed
            role = 2

            writer.writerow([worker_id, name, surname, email, phone_number, address, password, role])


# Generate users data and write to a CSV file
def generate_users_data():
    with open("users_data.csv", mode="w", newline="") as file:
        writer = csv.writer(file, delimiter=';')
        writer.writerow(["id", "name", "surname", "email", "phone_number", "address", "password"])

        for _ in range(rows):
            user_id = str(uuid.uuid4())  # Generate a UUID
            name = fake.first_name()
            surname = fake.last_name()

            # Generate a unique email (not previously used)
            while True:
                email = fake.unique.email()
                if email not in used_emails:
                    used_emails.add(email)
                    break

            phone_number = generate_phone_number()  # Generate a random phone number
            address = fake.address().replace("\n", " ")  # Replace line breaks with spaces
            password = hashed

            writer.writerow([user_id, name, surname, email, phone_number, address, password])


def generate_order_contains_data():
    with open("order_contains_data.csv", mode="w", newline="") as file:
        writer = csv.writer(file, delimiter=';')
        writer.writerow(["id", "order_id", "task_id", "quantity"])

        for _ in range(rows):
            id = str(uuid.uuid4())  # Generate a UUID
            order_id = fake.random_element(elements=order_ids)
            task_id = fake.random_element(elements=task_ids)
            quantity = fake.random_int(min=1, max=10)

            writer.writerow([id, order_id, task_id, quantity])


if __name__ == '__main__':
    # generate_workers_data()
    # generate_users_data()
    # generate_orders_data()
    # generate_order_contains_data()

    # print(hashed)  # $2b$12$TbfG11CRR9OSEsNX.Awije1.DmStMp.Erq1nJ/xIYwu.ilYjSbwOm
