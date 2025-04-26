# import csv
#
# def csv_to_sql_insert(csv_file_path, output_file_path, table_name):
#     """
#     Converts a CSV file to SQL INSERT statements and saves them to an output file.
#
#     Args:
#         csv_file_path (str): Path to the input CSV file
#         output_file_path (str): Path to save the SQL output file
#         table_name (str): Name of the table to insert into
#     """
#     try:
#         with open(csv_file_path, mode='r', newline='', encoding='utf-8') as csv_file, \
#                 open(output_file_path, mode='w', encoding='utf-8') as sql_file:
#
#             csv_reader = csv.DictReader(csv_file, delimiter=';')
#             # Debug: Print column names
#             print(f"Available columns: {csv_reader.fieldnames}")
#
#             # Write the beginning of the SQL statement
#             sql_file.write(f"INSERT INTO {table_name} (id, worker_id, user_id, status, deadline, address, creation_date, rate) VALUES\n")
#
#             # Process each row
#             rows = []
#             for row in csv_reader:
#                 # Escape single quotes in address by doubling them (SQL standard)
#                 address = row['address'].replace("'", "''")
#
#                 values = (
#                     f"('{row['id']}', "
#                     f"'{row['worker_id']}', "
#                     f"'{row['user_id']}', "
#                     f"{row['status']}, "
#                     f"'{row['deadline']}', "
#                     f"'{address}', "
#                     f"'{row['creation_date']}', "
#                     f"{row['rate']})"
#                 )
#                 rows.append(values)
#
#             # Write all rows separated by commas
#             sql_file.write(",\n".join(rows))
#
#             # End the SQL statement
#             sql_file.write(";")
#
#         print(f"Successfully converted CSV to SQL INSERT statements. Output saved to {output_file_path}")
#
#     except FileNotFoundError:
#         print(f"Error: File not found - {csv_file_path}")
#     except Exception as e:
#         print(f"An error occurred: {str(e)}")
#
# # Example usage
# if __name__ == "__main__":
#     input_csv = "orders_data.csv"  # Replace with your input CSV file path
#     output_sql = "orders_data.sql"  # Replace with your desired output file path
#     table_name = "orders"  # Replace with your table name
#
#     csv_to_sql_insert(input_csv, output_sql, table_name)

import csv
import os

def csv_to_sql_insert(csv_file_path, output_file_path, table_name):
    try:
        with open(csv_file_path, mode='r', newline='', encoding='utf-8') as csv_file, \
                open(output_file_path, mode='w', encoding='utf-8') as sql_file:

            csv_reader = csv.DictReader(csv_file, delimiter=';')
            columns = csv_reader.fieldnames

            if not columns:
                raise ValueError("CSV file has no headers")

            sql_file.write(f"INSERT INTO {table_name} ({', '.join(columns)}) VALUES\n")

            rows = []
            for row in csv_reader:
                values = []
                for col in columns:
                    value = row[col]
                    if col in ['status', 'rate']:
                        values.append(str(value))
                    else:
                        escaped_value = value.replace("'", "''")
                        values.append(f"'{escaped_value}'")

                rows.append(f"({', '.join(values)})")

            sql_file.write(",\n".join(rows))
            sql_file.write(";")

        print(f"Successfully converted CSV to SQL INSERT statements. Output saved to {output_file_path}")

    except FileNotFoundError:
        print(f"Error: File not found - {csv_file_path}")
    except Exception as e:
        print(f"An error occurred: {str(e)}")

if __name__ == "__main__":
    input_csv = input("Enter input CSV file path: ") # orders_data.csv
    output_sql = input("Enter output SQL file path: ") #orders_data.sql
    table_name = input("Enter table name: ") #orders

    if not os.path.exists(input_csv):
        print(f"Error: Input file '{input_csv}' not found")
    else:
        csv_to_sql_insert(input_csv, output_sql, table_name)