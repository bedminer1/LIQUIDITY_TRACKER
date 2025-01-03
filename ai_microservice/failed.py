import sqlite3
import pandas as pd

# Step 1: Load data from the SQLite database based on asset_type
def load_data_by_asset_type(db_file):
    conn = sqlite3.connect(db_file)
    
    # Fetch distinct asset_types from the database
    query = "SELECT DISTINCT asset_type FROM records"
    asset_types = pd.read_sql(query, conn)['asset_type'].tolist()
    
    # Create a dictionary to hold DataFrames for each asset_type
    asset_data = {}
    
    for asset_type in asset_types:
        # Query data for each asset_type
        query = f"SELECT timestamp, bid_ask_spread, volume, bid_price FROM records WHERE asset_type = '{asset_type}'"
        df = pd.read_sql(query, conn)
        
        # Attempt to parse the timestamp using different formats
        df['timestamp'] = pd.to_datetime(df['timestamp'], format='%d.%m.%YT%H:%M:%SZ', errors='coerce')

        # Handle remaining NaT values using another attempt
        df['timestamp'].fillna(pd.to_datetime(df['timestamp'], errors='coerce'), inplace=True)
        
        # Set timestamp as index for time series
        df.set_index('timestamp', inplace=True)
        
        # Store the DataFrame for this asset_type in the dictionary
        asset_data[asset_type] = df
    
    conn.close()
    return asset_data

# Example usage
db_file = r"D:\Coding\LIQUIDITY_TRACKER\backend\market_data.db"
asset_data = load_data_by_asset_type(db_file)

# Print the data for a specific asset_type
print(asset_data.get('Crypto_BTC'))
