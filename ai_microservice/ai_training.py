import sqlite3  # To connect to our .db file
import pandas as pd  # To manage our data easily like a table
import numpy as np  # For numerical operations
import matplotlib.pyplot as plt  # To make pretty graphs
from sklearn.preprocessing import MinMaxScaler  # To scale data for LSTM
from tensorflow.keras.models import Sequential  # To build the LSTM model
from tensorflow.keras.layers import LSTM, Dense  # LSTM layers

# Step 1: Connect to the .db file and fetch data from the table
conn = sqlite3.connect(r"D:\Coding\LIQUIDITY_TRACKER\backend\market_data.db")  # Connecting to our .db file
query = "SELECT * FROM records"  # SQL to fetch everything from our table
data = pd.read_sql_query(query, conn)  # Fetch and store the data as a table in Python
conn.close()  # We are done with the database, so let's close the connection

# Step 2: Scale the entire dataset
scaler = MinMaxScaler()
data_scaled = scaler.fit_transform(data[['timestamp', 'bid_ask_spread', 'volume', 'bid_price']])
data_scaled = pd.DataFrame(data_scaled, columns=['timestamp', 'bid_ask_spread', 'volume', 'bid_price'])
data_scaled['asset_type'] = data['asset_type']

# Step 3: Create sequences for LSTM
def create_sequences(data, sequence_length):
    sequences = []
    for i in range(len(data) - sequence_length):
        seq = data[i:i + sequence_length]
        label = data[i + sequence_length]
        sequences.append((seq, label))
    return np.array(sequences, dtype=object)

# Step 4: Prepare the training data by splitting into sequences based on asset type
sequence_length = 10
all_sequences = []
for asset_type in data_scaled['asset_type'].unique():
    asset_data = data_scaled[data_scaled['asset_type'] == asset_type][['timestamp', 'bid_ask_spread', 'volume', 'bid_price']].values
    sequences = create_sequences(asset_data, sequence_length)
    all_sequences.extend(sequences)

# Split sequences into inputs (X) and labels (Y)
X, Y = zip(*all_sequences)
X = np.array(X)
Y = np.array(Y)

# Step 5: Build a single LSTM model
model = Sequential()
model.add(LSTM(50, return_sequences=True, input_shape=(X.shape[1], X.shape[2])))
model.add(LSTM(50))
model.add(Dense(3))  # Predicting bid_ask_spread, volume, and bid_price
model.compile(optimizer='adam', loss='mean_squared_error')

# Step 6: Train the model
model.fit(X, Y, epochs=10, batch_size=32, verbose=1)

# Step 7: Define a function to predict future values and visualize them
def predict_and_plot(db_file, time_to_predict, time_interval):
    # Step 7.1: Connect to the provided .db file and fetch the historical data
    conn = sqlite3.connect(db_file)
    query = "SELECT * FROM records"
    new_data = pd.read_sql_query(query, conn)
    conn.close()

    # Step 7.2: Scale the new data using the previously fitted scaler
    new_data_scaled = scaler.transform(new_data[['timestamp', 'bid_ask_spread', 'volume', 'bid_price']])

    # Step 7.3: Prepare the last known data as input for prediction
    last_known_data = new_data_scaled[-sequence_length:]
    future_data = [last_known_data]

    # Step 7.4: Generate future timestamps based on the given time interval
    start_time = new_data['timestamp'].iloc[-1]
    future_timestamps = [start_time + i * time_interval for i in range(1, int(time_to_predict / time_interval) + 1)]

    # Step 7.5: Predict future values
    predictions = []
    for _ in range(len(future_timestamps)):
        pred = model.predict(np.array(future_data))
        predictions.append(pred[0])
        future_data.append(pred)
        future_data = future_data[1:]

    # Step 7.6: Rescale predictions back to original values
    predictions = scaler.inverse_transform(predictions)

    # Step 7.7: Convert predictions to a DataFrame for easier plotting
    predictions_df = pd.DataFrame(predictions, columns=['bid_ask_spread', 'volume', 'bid_price'])
    predictions_df['timestamp'] = future_timestamps

    # Step 7.8: Plot the predictions
    plt.figure(figsize=(10, 6))
    plt.plot(new_data['timestamp'], new_data['bid_ask_spread'], label='Historical Bid-Ask Spread')
    plt.plot(predictions_df['timestamp'], predictions_df['bid_ask_spread'], label='Predicted Bid-Ask Spread')
    plt.xlabel('Timestamp')
    plt.ylabel('Bid-Ask Spread')
    plt.title('Bid-Ask Spread Prediction')
    plt.legend()
    plt.show()

    plt.figure(figsize=(10, 6))
    plt.plot(new_data['timestamp'], new_data['volume'], label='Historical Volume')
    plt.plot(predictions_df['timestamp'], predictions_df['volume'], label='Predicted Volume')
    plt.xlabel('Timestamp')
    plt.ylabel('Volume')
    plt.title('Volume Prediction')
    plt.legend()
    plt.show()

    plt.figure(figsize=(10, 6))
    plt.plot(new_data['timestamp'], new_data['bid_price'], label='Historical Bid Price')
    plt.plot(predictions_df['timestamp'], predictions_df['bid_price'], label='Predicted Bid Price')
    plt.xlabel('Timestamp')
    plt.ylabel('Bid Price')
    plt.title('Bid Price Prediction')
    plt.legend()
    plt.show()

# Example usage
# predict_and_plot("path_to_unknown_asset.db", 3600, 300)  # Predict for 1 hour with 5-minute intervals
