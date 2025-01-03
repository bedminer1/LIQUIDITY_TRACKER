import sqlite3

conn = sqlite3.connect(r"..\backend\market_data.db")
cursor = conn.cursor()

query = "select distinct asset_type from records"
cursor.execute(query)

rows = cursor.fetchall()  # Get all rows
for i in rows:
    print(i)

    # Extract the asset_type value from the tuple (i[0] if it's the first element)
    asset_type = i[0]

    # Use parameterized queries to avoid SQL injection and ensure proper formatting
    newquery = "select * from records where asset_type = ? limit 10"
    cursor.execute(newquery, (asset_type,))  # Pass the value as a tuple
    newrows = cursor.fetchall()

    # Print the new rows retrieved
    for x in newrows:
        print(x)

conn.close()