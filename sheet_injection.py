import gspread
from oauth2client.service_account import ServiceAccountCredentials
import csv
import os

# Define scope and credentials
scope = ["https://spreadsheets.google.com/feeds", "https://www.googleapis.com/auth/drive"]
creds = ServiceAccountCredentials.from_json_keyfile_name('credentials.json', scope)
client = gspread.authorize(creds)

# Open the main Google Sheet
spreadsheet = client.open("NOTION.Shops and Evaluations")

# Get all CSV files in the current directory
csv_files = [f for f in os.listdir('.') if f.endswith('.csv')]

for csv_file in csv_files:
    sheet_name = os.path.splitext(csv_file)[0]  # Remove .csv extension
    print(f"Uploading {csv_file} to Google Sheet: {sheet_name}")
    
    # Try to open the specific worksheet, create if it doesn't exist
    try:
        sheet = spreadsheet.worksheet(sheet_name)
    except gspread.exceptions.WorksheetNotFound:
        print(f"Worksheet '{sheet_name}' not found. Creating a new one.")
        sheet = spreadsheet.add_worksheet(title=sheet_name, rows="1000", cols="26")
    
    # Read the CSV file
    with open(csv_file, 'r', encoding='utf-8') as csvfile:
        reader = csv.reader(csvfile)
        data = list(reader)

    # Resize the sheet to fit the data
    num_rows = len(data)
    num_cols = max(len(row) for row in data) if data else 0
    sheet.resize(rows=num_rows, cols=num_cols)
    
    # Clear the existing sheet content
    sheet.clear()
    
    # Convert data into batch update format
    requests = []
    for row_idx, row in enumerate(data, start=1):
        for col_idx, cell in enumerate(row, start=1):
            cell_range = gspread.utils.rowcol_to_a1(row_idx, col_idx)
            requests.append({
                "range": cell_range,
                "values": [[cell]]  # Ensures formulas are recognized properly
            })
    
    # Batch update the sheet to ensure formulas work
    sheet.batch_update(requests, value_input_option="USER_ENTERED")

    print(f"{csv_file} synced to Google Sheet '{sheet_name}' successfully.")
